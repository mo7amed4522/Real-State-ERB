import { Injectable, UnauthorizedException } from '@nestjs/common';
import { UserService } from '../user/user.service';
import { User } from '../user/user.entity';
import * as crypto from 'crypto';

@Injectable()
export class AuthService {
  constructor(private readonly userService: UserService) {}

  async validateUser(email: string, pass: string): Promise<Omit<User, 'password'> | null> {
    return this.userService.validatePassword(email, pass);
  }

  async forgotPassword(email: string): Promise<{ resetToken: string }> {
    const user = await this.userService.findByEmail(email);
    if (!user) {
      // Don't reveal that the user does not exist
      return { resetToken: 'If a matching account exists, a token has been generated.' };
    }

    const resetToken = crypto.randomBytes(32).toString('hex');
    const hashedToken = crypto.createHash('sha256').update(resetToken).digest('hex');

    const expiration = new Date();
    expiration.setHours(expiration.getHours() + 1); // 1 hour expiry

    await this.userService.update(user.id, {
      resetPasswordToken: hashedToken,
      resetPasswordExpires: expiration,
    });
    
    // In a real app, you would email this resetToken to the user.
    // For development, we return it directly.
    return { resetToken };
  }

  async resetPassword(token: string, newPassword: string): Promise<User> {
    const hashedToken = crypto.createHash('sha256').update(token).digest('hex');
    
    const user = await this.userService.findByResetToken(hashedToken);

    if (!user || user.resetPasswordExpires < new Date()) {
      throw new UnauthorizedException('Password reset token is invalid or has expired.');
    }

    await this.userService.update(user.id, {
      password: newPassword, // The service will hash this automatically
      resetPasswordToken: null,
      resetPasswordExpires: null,
    });
    
    const { password, ...result } = user;
    return result as User;
  }
} 