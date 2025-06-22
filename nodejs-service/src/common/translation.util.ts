import axios from 'axios';

export async function translateText(
  text: string,
  sourceLang: string,
  targetLangs: string[]
): Promise<Record<string, string>> {
  const reqBody = {
    text,
    source_lang: sourceLang,
    target_langs: targetLangs,
  };
  const response = await axios.post('http://python-translate-service:8010/translate', reqBody);
  return response.data;
} 