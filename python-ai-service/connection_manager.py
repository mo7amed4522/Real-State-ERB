from fastapi import WebSocket

class ConnectionManager:
    def __init__(self):
        self.active_connections: dict[str, WebSocket] = {}

    async def connect(self, room_id: str, websocket: WebSocket):
        await websocket.accept()
        self.active_connections[room_id] = websocket

    def disconnect(self, room_id: str):
        if room_id in self.active_connections:
            del self.active_connections[room_id]

    async def send_personal_message(self, message: str, room_id: str):
        if room_id in self.active_connections:
            websocket = self.active_connections[room_id]
            await websocket.send_text(message)

manager = ConnectionManager() 