import json
from kafka import KafkaConsumer, KafkaProducer
import psycopg2
import os
from transformers import pipeline
from encryption import encryption_service

# --- AI Model Initialization ---
# This will download the model on the first run.
# Using a conversational model as a base.
print("Loading conversational AI model...")
chatbot = pipeline('conversational', model='microsoft/DialoGPT-medium')
print("AI model loaded.")

# --- Kafka Initialization ---
print("Initializing Kafka consumer and producer...")
consumer = KafkaConsumer(
    'user_messages',
    bootstrap_servers='kafka:9092',
    value_deserializer=lambda v: v.decode('utf-8'),
    group_id='ai-worker-group'
)
producer = KafkaProducer(
    bootstrap_servers='kafka:9092',
    value_serializer=lambda v: v.encode('utf-8')
)
print("Kafka initialized.")

# --- Database Connection ---
def get_db_connection():
    try:
        conn = psycopg2.connect(os.environ.get('DATABASE_URL'))
        print("Database connection successful.")
        return conn
    except Exception as e:
        print(f"Database connection failed: {e}")
        return None

# --- Placeholder Functions ---
def query_property_data(message: str):
    """
    Placeholder function to query property data from the database.
    In a real application, you would parse the message to extract property IDs, etc.
    """
    conn = get_db_connection()
    if not conn:
        return "I can't access property information right now."
    
    # Example: Check if message contains a number (as a fake property ID)
    # This is a very basic example. You'd use NLP to extract entities.
    try:
        # cur = conn.cursor()
        # cur.execute("SELECT * FROM properties WHERE id = %s", (property_id,))
        # property_data = cur.fetchone()
        # cur.close()
        # conn.close()
        # if property_data:
        #     return f"Property {property_id} is..."
        # else:
        #     return f"I couldn't find any information on property {property_id}."
        return "Database query for properties would happen here."
    except Exception as e:
        return f"I had trouble querying the database: {e}"


# --- Main Processing Loop ---
print("AI Worker is running and waiting for messages...")
for message in consumer:
    try:
        decrypted_message = encryption_service.decrypt(message.value)
        data = json.loads(decrypted_message)
        room_id = data.get("room_id")
        user_text = data.get("text")
        print(f"Received encrypted message from room {room_id}")

        response_text = ""
        # Simple logic to decide if it's a property query
        if "property" in user_text.lower():
            response_text = query_property_data(user_text)
        else:
            # Use the conversational AI model
            conversation = chatbot(user_text)
            response_text = conversation.generated_responses[-1]

        print(f"Generated response for room {room_id}")
        
        # Encrypt the response before sending
        response_payload = {'room_id': room_id, 'text': response_text}
        encrypted_response = encryption_service.encrypt(json.dumps(response_payload))
        producer.send('bot_responses', encrypted_response)
    except Exception as e:
        print(f"Failed to process message: {e}") 