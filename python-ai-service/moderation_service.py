import os
import re
import base64
from typing import Dict, Any, Optional
from PIL import Image
import io
import numpy as np
from transformers import pipeline, AutoTokenizer, AutoModelForSequenceClassification
import torch
from fastapi import HTTPException
import logging

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class ModerationService:
    def __init__(self):
        self.text_classifier = None
        self.image_classifier = None
        self.toxic_words = set()
        self.spam_patterns = []
        
        # Initialize models
        self._load_models()
        self._load_filters()
    
    def _load_models(self):
        """Load AI models for text and image classification"""
        try:
            # Text classification model for toxicity detection
            model_name = "unitary/toxic-bert"
            self.text_classifier = pipeline(
                "text-classification",
                model=model_name,
                tokenizer=model_name,
                device=0 if torch.cuda.is_available() else -1
            )
            logger.info("Text moderation model loaded successfully")
        except Exception as e:
            logger.error(f"Failed to load text model: {e}")
            self.text_classifier = None
    
    def _load_filters(self):
        """Load rule-based filters"""
        # Toxic/inappropriate words (simplified list)
        self.toxic_words = {
            'badword', 'inappropriate', 'spam', 'scam', 'fake', 'fraud',
            'hate', 'discrimination', 'violence', 'threat', 'harassment'
        }
        
        # Spam patterns
        self.spam_patterns = [
            r'\b(?:buy|sell|discount|offer|limited|urgent|act now)\b',
            r'\b(?:click here|visit|subscribe|join now)\b',
            r'\b(?:free|money|cash|earn|profit|investment)\b',
            r'\b(?:lottery|winner|prize|jackpot)\b',
            r'\b(?:viagra|cialis|medication|prescription)\b',
            r'\b(?:casino|poker|betting|gambling)\b',
        ]
    
    def moderate_text(self, content: str) -> Dict[str, Any]:
        """Moderate text content for inappropriate language"""
        if not content or len(content.strip()) == 0:
            return {
                "allowed": True,
                "reason": "Empty content",
                "severity": "low",
                "flagged": False
            }
        
        # Rule-based checks
        content_lower = content.lower()
        
        # Check for toxic words
        found_toxic_words = [word for word in self.toxic_words if word in content_lower]
        if found_toxic_words:
            return {
                "allowed": False,
                "reason": f"Contains inappropriate language: {', '.join(found_toxic_words)}",
                "severity": "high",
                "flagged": True
            }
        
        # Check for spam patterns
        for pattern in self.spam_patterns:
            if re.search(pattern, content_lower):
                return {
                    "allowed": False,
                    "reason": "Detected spam content",
                    "severity": "medium",
                    "flagged": True
                }
        
        # Check for excessive caps (shouting)
        if len(content) > 10 and content.isupper():
            return {
                "allowed": False,
                "reason": "Excessive use of capital letters",
                "severity": "low",
                "flagged": True
            }
        
        # Check for repeated characters
        if re.search(r'(.)\1{4,}', content):
            return {
                "allowed": False,
                "reason": "Excessive character repetition",
                "severity": "low",
                "flagged": True
            }
        
        # AI-based classification if model is available
        if self.text_classifier:
            try:
                result = self.text_classifier(content)
                # Check if any toxic category has high confidence
                for item in result:
                    if item['score'] > 0.7:  # High confidence threshold
                        return {
                            "allowed": False,
                            "reason": f"AI detected {item['label']} content (confidence: {item['score']:.2f})",
                            "severity": "high",
                            "flagged": True
                        }
            except Exception as e:
                logger.error(f"AI text classification error: {e}")
        
        return {
            "allowed": True,
            "reason": "Content approved",
            "severity": "low",
            "flagged": False
        }
    
    def moderate_image(self, image_path: str) -> Dict[str, Any]:
        """Moderate image content for inappropriate content"""
        if not image_path or not os.path.exists(image_path):
            return {
                "allowed": True,
                "reason": "No image provided or file not found",
                "severity": "low",
                "flagged": False
            }
        
        try:
            # Basic image validation
            with Image.open(image_path) as img:
                # Check file size (prevent very large files)
                file_size = os.path.getsize(image_path)
                if file_size > 10 * 1024 * 1024:  # 10MB limit
                    return {
                        "allowed": False,
                        "reason": "Image file too large",
                        "severity": "medium",
                        "flagged": True
                    }
                
                # Check image dimensions
                width, height = img.size
                if width > 4000 or height > 4000:
                    return {
                        "allowed": False,
                        "reason": "Image dimensions too large",
                        "severity": "low",
                        "flagged": True
                    }
                
                # Check file format
                if img.format not in ['JPEG', 'PNG', 'GIF', 'WEBP']:
                    return {
                        "allowed": False,
                        "reason": "Unsupported image format",
                        "severity": "medium",
                        "flagged": True
                    }
                
                # TODO: Add AI-based image content analysis
                # This would require additional models like:
                # - NSFW detection
                # - Violence detection
                # - Hate symbol detection
                # - Spam/ads detection
                
                # For now, just basic validation
                return {
                    "allowed": True,
                    "reason": "Image passed basic validation",
                    "severity": "low",
                    "flagged": False
                }
                
        except Exception as e:
            logger.error(f"Image moderation error: {e}")
            return {
                "allowed": False,
                "reason": f"Image processing error: {str(e)}",
                "severity": "medium",
                "flagged": True
            }
    
    def moderate_content(self, content: str = "", image_path: str = "", 
                        user_id: int = 0, user_type: str = "", room_id: int = 0) -> Dict[str, Any]:
        """Main moderation function that handles both text and image content"""
        
        # Log moderation request
        logger.info(f"Moderating content for user {user_id} ({user_type}) in room {room_id}")
        
        # Moderate text content
        text_result = self.moderate_text(content)
        
        # Moderate image content
        image_result = self.moderate_image(image_path)
        
        # Combine results
        if not text_result["allowed"]:
            return text_result
        elif not image_result["allowed"]:
            return image_result
        elif text_result["flagged"] or image_result["flagged"]:
            # If either is flagged, return the higher severity one
            if text_result["severity"] == "high" or image_result["severity"] == "high":
                return text_result if text_result["severity"] == "high" else image_result
            else:
                return text_result if text_result["severity"] == "medium" else image_result
        
        # Both passed
        return {
            "allowed": True,
            "reason": "Content approved",
            "severity": "low",
            "flagged": False
        }

# Global moderation service instance
moderation_service = ModerationService() 