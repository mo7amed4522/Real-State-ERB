from fastapi import FastAPI, WebSocket, WebSocketDisconnect, HTTPException
from pydantic import BaseModel
from connection_manager import manager
from kafka import KafkaProducer, KafkaConsumer
from moderation_service import moderation_service
import json
import threading
import asyncio
from encryption import encryption_service

app = FastAPI()

# Pydantic models for moderation
class ModerationRequest(BaseModel):
    content: str = ""
    image_url: str = ""
    image_path: str = ""
    user_id: int = 0
    user_type: str = ""
    room_id: int = 0

class ModerationResponse(BaseModel):
    allowed: bool
    reason: str = ""
    severity: str = "low"
    flagged: bool = False

# Kafka Producer - value will be an encrypted string
producer = KafkaProducer(
    bootstrap_servers='kafka:9092',
    value_serializer=lambda v: v.encode('utf-8')
)

# Kafka Consumer in a background thread
def consume_bot_responses():
    consumer = KafkaConsumer(
        'bot_responses',
        bootstrap_servers='kafka:9092',
        value_deserializer=lambda v: v.decode('utf-8'),
        group_id='fastapi-group'
    )
    for message in consumer:
        try:
            decrypted_message = encryption_service.decrypt(message.value)
            response = json.loads(decrypted_message)
            room_id = response.get("room_id")
            text = response.get("text")
            if room_id and text:
                asyncio.run(manager.send_personal_message(text, room_id))
        except Exception as e:
            print(f"Failed to decrypt or process message: {e}")


consumer_thread = threading.Thread(target=consume_bot_responses)
consumer_thread.daemon = True
consumer_thread.start()


@app.get("/")
def read_root():
    return {"message": "Hello from FastAPI AI Service!"}

@app.post("/moderate", response_model=ModerationResponse)
async def moderate_content(request: ModerationRequest):
    """
    Moderate text and image content for inappropriate content
    """
    try:
        result = moderation_service.moderate_content(
            content=request.content,
            image_path=request.image_path,
            user_id=request.user_id,
            user_type=request.user_type,
            room_id=request.room_id
        )
        
        return ModerationResponse(**result)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Moderation error: {str(e)}")

@app.websocket("/ws/{room_id}")
async def websocket_endpoint(websocket: WebSocket, room_id: str):
    await manager.connect(room_id, websocket)
    try:
        while True:
            data = await websocket.receive_text()
            # Create payload, convert to JSON, then encrypt
            payload = {'room_id': room_id, 'text': data}
            encrypted_payload = encryption_service.encrypt(json.dumps(payload))
            producer.send('user_messages', encrypted_payload)
    except WebSocketDisconnect:
        manager.disconnect(room_id) 