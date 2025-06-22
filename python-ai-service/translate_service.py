from fastapi import FastAPI
from pydantic import BaseModel
from typing import List, Dict
from transformers import pipeline
import threading

app = FastAPI()

# Supported language codes
LANG_CODE_MAP = {
    'en': 'English',
    'ar': 'Arabic',
    'fil': 'Filipino',
    'fr': 'French',
    'de': 'German',
    'hi': 'Hindi',
    'ru': 'Russian',
}

# Map (src, tgt) to Hugging Face model names
MODEL_MAP = {
    ('en', 'ar'): 'Helsinki-NLP/opus-mt-en-ar',
    ('en', 'fr'): 'Helsinki-NLP/opus-mt-en-fr',
    ('en', 'de'): 'Helsinki-NLP/opus-mt-en-de',
    ('en', 'ru'): 'Helsinki-NLP/opus-mt-en-ru',
    ('en', 'hi'): 'Helsinki-NLP/opus-mt-en-hi',
    # Filipino (Tagalog) not directly supported, fallback to echo
}

# Thread-safe model cache
model_cache = {}
model_lock = threading.Lock()

def get_translator(src: str, tgt: str):
    key = (src, tgt)
    with model_lock:
        if key not in model_cache:
            model_name = MODEL_MAP.get(key)
            if not model_name:
                return None
            model_cache[key] = pipeline('translation', model=model_name)
        return model_cache[key]

class TranslationRequest(BaseModel):
    text: str
    source_lang: str = 'en'
    target_langs: List[str]

@app.post('/translate')
def translate(req: TranslationRequest) -> Dict[str, str]:
    results = {}
    for lang in req.target_langs:
        if lang == req.source_lang:
            results[lang] = req.text
            continue
        translator = get_translator(req.source_lang, lang)
        if translator:
            try:
                results[lang] = translator(req.text)[0]['translation_text']
            except Exception as e:
                results[lang] = f"[Translation error: {e}]"
        else:
            # Filipino/tagalog or unsupported: fallback to echo
            results[lang] = f"{req.text} [{lang}]"
    return results 