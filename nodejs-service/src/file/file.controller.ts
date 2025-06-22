import { Controller, Get, Param, Res, UseGuards } from '@nestjs/common';
import { Response } from 'express';
import { GqlAuthGuard } from '../auth/guards/gql-auth.guard';
import { FileService } from '../common/file.service';
import { EncryptionService } from '../common/encryption.service';

@Controller('files')
export class FileController {
  constructor(
    private readonly fileService: FileService,
    private readonly encryptionService: EncryptionService,
  ) {}

  @Get('private/:encryptedPath')
  @UseGuards(GqlAuthGuard) // Or your preferred auth guard
  async getPrivateFile(
    @Param('encryptedPath') encryptedPath: string,
    @Res() res: Response,
  ) {
    try {
      const decryptedPath = this.encryptionService.decrypt(encryptedPath);
      const fileStream = this.fileService.getPrivateFileStream(decryptedPath);
      
      // You might want to set more specific headers based on file type
      res.setHeader('Content-Type', 'application/octet-stream');
      fileStream.pipe(res);
    } catch (error) {
      res.status(404).send('File not found or access denied.');
    }
  }
} 