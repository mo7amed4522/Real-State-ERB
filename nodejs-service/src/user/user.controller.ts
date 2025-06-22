import { Controller, Get, UseGuards, Req, Res } from '@nestjs/common';
import { GqlAuthGuard } from '../auth/guards/gql-auth.guard'; // Using the GQL guard for session checking
import { EncryptionService } from '../common/encryption.service';
import { User } from './user.entity';
import { Response } from 'express';
import * as path from 'path';

@Controller('users')
export class UserController {
  constructor(private readonly encryptionService: EncryptionService) {}

  @Get('photo')
  @UseGuards(GqlAuthGuard)
  async getProfilePicture(@Req() req, @Res() res: Response) {
    const user: User = req.user;
    if (!user || !user.photoUrl || user.photoUrl.startsWith('http')) {
      return res.status(404).send('No profile picture found.');
    }

    try {
      const decryptedPath = this.encryptionService.decrypt(user.photoUrl);
      const filePath = path.resolve(`./storage/${decryptedPath}`);
      return res.sendFile(filePath);
    } catch (error) {
      return res.status(404).send('File not found or invalid.');
    }
  }
} 