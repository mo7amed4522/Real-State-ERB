import os
import base64
from cryptography.hazmat.primitives.ciphers.aead import AESGCM

class EncryptionService:
    def __init__(self):
        key_hex = os.environ.get('ENCRYPTION_KEY')
        if not key_hex or len(key_hex) != 64:
            raise ValueError("ENCRYPTION_KEY environment variable must be a 64-character hex string (32 bytes).")
        self.key = bytes.fromhex(key_hex)
        self.aesgcm = AESGCM(self.key)

    def encrypt(self, plaintext: str) -> str:
        """Encrypts a plaintext string and returns a base64 encoded string containing nonce + ciphertext."""
        nonce = os.urandom(12)  # GCM nonce
        ciphertext = self.aesgcm.encrypt(nonce, plaintext.encode('utf-8'), None)
        encrypted_payload = base64.b64encode(nonce + ciphertext).decode('utf-8')
        return encrypted_payload

    def decrypt(self, encrypted_payload: str) -> str:
        """Decrypts a base64 encoded payload and returns the plaintext string."""
        data = base64.b64decode(encrypted_payload)
        nonce = data[:12]
        ciphertext = data[12:]
        plaintext = self.aesgcm.decrypt(nonce, ciphertext, None).decode('utf-8')
        return plaintext

encryption_service = EncryptionService() 