# AI Moderation System Guide

## Overview

The AI moderation system automatically checks chat messages and uploaded files for inappropriate content before they are sent to other users. This helps maintain a safe and professional environment in your chat application.

## How It Works

### 1. Text Moderation
- **Rule-based filtering**: Checks for toxic words, spam patterns, excessive caps, and character repetition
- **AI-based classification**: Uses the `unitary/toxic-bert` model to detect toxic content with high confidence
- **Real-time processing**: Messages are checked before being saved to the database

### 2. Image Moderation
- **File validation**: Checks file size, dimensions, and format
- **Basic security**: Prevents oversized or unsupported files
- **Future enhancement**: Will include AI-based content analysis (NSFW, violence, hate symbols)

### 3. Moderation Flow
1. User sends message or uploads file
2. Go service calls Python AI service for moderation
3. AI service analyzes content and returns result
4. If blocked: Message is rejected with explanation
5. If flagged: Message is marked but allowed through
6. If approved: Message is saved and broadcast normally

## Configuration

### Environment Variables

Add these to your `.env` file:

```bash
# AI Service URL
AI_SERVICE_URL=http://python-ai-service:8000

# Moderation settings
MODERATION_ENABLED=true
MODERATION_FAIL_OPEN=true  # Allow messages if AI service is down
```

### Moderation Thresholds

In `python-ai-service/moderation_service.py`:

```python
# Text classification confidence threshold
CONFIDENCE_THRESHOLD = 0.7

# File size limits
MAX_FILE_SIZE = 10 * 1024 * 1024  # 10MB
MAX_IMAGE_DIMENSIONS = 4000  # pixels
```

## API Endpoints

### Moderation Endpoint

**POST** `/moderate`

Request:
```json
{
  "content": "Message text to moderate",
  "image_path": "/path/to/image.jpg",
  "user_id": 123,
  "user_type": "user",
  "room_id": 456
}
```

Response:
```json
{
  "allowed": true,
  "reason": "Content approved",
  "severity": "low",
  "flagged": false
}
```

## Moderation Results

### Allowed Content
- `allowed: true`
- Content passes all checks
- Message is sent normally

### Flagged Content
- `allowed: true, flagged: true`
- Content has minor issues but is allowed
- Message is marked for review
- Admin notification may be sent

### Blocked Content
- `allowed: false`
- Content violates rules
- Message is rejected
- User receives error message

### Severity Levels
- **Low**: Minor issues (excessive caps, character repetition)
- **Medium**: Spam patterns, large files, unsupported formats
- **High**: Toxic language, inappropriate content

## Examples

### Text Moderation Examples

```graphql
# This message would be blocked
mutation {
  sendMessage(input: {
    roomId: "1",
    content: "This contains badword content",
    messageType: "text"
  }) {
    id
    content
  }
}
# Response: Error - "message blocked by AI moderation: Contains inappropriate language: badword"

# This message would be flagged
mutation {
  sendMessage(input: {
    roomId: "1",
    content: "BUY NOW LIMITED TIME OFFER!!!",
    messageType: "text"
  }) {
    id
    content
  }
}
# Response: Error - "message blocked by AI moderation: Detected spam content"
```

### File Upload Examples

```graphql
# This file would be blocked
mutation($file: Upload!) {
  uploadFile(input: {
    messageId: "1",
    file: $file  # 15MB file
  }) {
    id
    fileName
  }
}
# Response: Error - "file blocked by AI moderation: Image file too large"
```

## Monitoring and Logging

### Moderation Logs

All moderation events are logged in the `chat_moderation_logs` table:

```sql
SELECT * FROM chat_moderation_logs 
WHERE created_at > NOW() - INTERVAL '1 day'
ORDER BY created_at DESC;
```

### Key Metrics to Monitor

- **Block rate**: Percentage of messages blocked
- **Flag rate**: Percentage of messages flagged
- **AI service availability**: Uptime of moderation service
- **Response time**: How long moderation takes

## Customization

### Adding Custom Rules

Edit `python-ai-service/moderation_service.py`:

```python
def _load_filters(self):
    # Add custom toxic words
    self.toxic_words.update(['custom_word1', 'custom_word2'])
    
    # Add custom spam patterns
    self.spam_patterns.append(r'\b(?:custom_pattern)\b')
```

### Adjusting Sensitivity

```python
# Make moderation more strict
CONFIDENCE_THRESHOLD = 0.5  # Lower threshold = more sensitive

# Make moderation less strict
CONFIDENCE_THRESHOLD = 0.9  # Higher threshold = less sensitive
```

## Troubleshooting

### Common Issues

1. **AI service unavailable**
   - Check if Python AI service is running
   - Verify network connectivity between services
   - Check logs for connection errors

2. **High false positives**
   - Adjust confidence threshold
   - Review and update toxic words list
   - Fine-tune spam patterns

3. **Slow moderation**
   - Check AI model loading
   - Monitor system resources
   - Consider model optimization

### Debug Mode

Enable debug logging in `moderation_service.py`:

```python
logging.basicConfig(level=logging.DEBUG)
```

## Security Considerations

- All moderation requests are logged for audit purposes
- Failed moderation attempts don't expose sensitive information
- AI models are loaded securely with proper error handling
- File paths are validated to prevent path traversal attacks

## Future Enhancements

- **Advanced image analysis**: NSFW detection, violence detection
- **User reputation system**: Adjust moderation based on user history
- **Custom model training**: Train models on your specific content
- **Real-time learning**: Improve accuracy based on user feedback
- **Multi-language support**: Detect inappropriate content in multiple languages 