#!/bin/bash

# Start the AI worker in the background
echo "Starting AI worker..."
python ai_worker.py &

# Start the FastAPI server in the foreground
echo "Starting FastAPI server..."
uvicorn main:app --host 0.0.0.0 --port 8000 