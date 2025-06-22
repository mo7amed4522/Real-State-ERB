import { Module } from '@nestjs/common';
import { EncryptionService } from './encryption.service';
import { FileService } from './file.service';

@Module({
  providers: [EncryptionService, FileService],
  exports: [EncryptionService, FileService],
})
export class CommonModule {} 