import { Module } from '@nestjs/common';
import { FileController } from './file.controller';
import { FileService } from '../common/file.service';
import { EncryptionService } from '../common/encryption.service';
import { AuthModule } from '../auth/auth.module';

@Module({
  imports: [AuthModule],
  controllers: [FileController],
  providers: [FileService, EncryptionService],
})
export class FileModule {} 