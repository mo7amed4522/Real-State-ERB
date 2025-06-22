import { Injectable, OnModuleInit } from '@nestjs/common';
import * as crypto from 'crypto';

@Injectable()
export class EncryptionService implements OnModuleInit {
  private key: Buffer;
  private readonly algorithm = 'aes-256-gcm';
  private readonly ivLength = 16;
  private readonly tagLength = 16;

  onModuleInit() {
    const secret = process.env.ENCRYPTION_SECRET_KEY;
    if (!secret) {
      throw new Error('ENCRYPTION_SECRET_KEY environment variable not set.');
    }
    // Use SHA-256 to create a deterministic 32-byte key from the secret.
    this.key = crypto.createHash('sha256').update(secret).digest();
  }

  encrypt(text: string) {
    const iv = crypto.randomBytes(this.ivLength);
    const cipher = crypto.createCipheriv(this.algorithm, this.key, iv);
    const encrypted = Buffer.concat([cipher.update(text, 'utf8'), cipher.final()]);
    const tag = cipher.getAuthTag();
    return Buffer.concat([iv, tag, encrypted]).toString('hex');
  }

  decrypt(encryptedText: string) {
    const data = Buffer.from(encryptedText, 'hex');
    const iv = data.slice(0, this.ivLength);
    const tag = data.slice(this.ivLength, this.ivLength + this.tagLength);
    const encrypted = data.slice(this.ivLength + this.tagLength);
    const decipher = crypto.createDecipheriv(this.algorithm, this.key, iv);
    decipher.setAuthTag(tag);
    return Buffer.concat([decipher.update(encrypted), decipher.final()]).toString('utf8');
  }
} 