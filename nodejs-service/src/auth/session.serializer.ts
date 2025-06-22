import { PassportSerializer } from '@nestjs/passport';
import { Injectable } from '@nestjs/common';
import { UserService } from 'src/user/user.service';
import { User } from 'src/user/user.entity';

@Injectable()
export class SessionSerializer extends PassportSerializer {
  constructor(private readonly userService: UserService) {
    super();
  }

  serializeUser(user: User, done: (err: Error, user: { id: string }) => void): void {
    done(null, { id: user.id });
  }

  async deserializeUser(payload: { id: string }, done: (err: Error, user: User) => void): Promise<void> {
    const user = await this.userService.findById(payload.id); // We need to add findById to UserService
    done(null, user);
  }
} 