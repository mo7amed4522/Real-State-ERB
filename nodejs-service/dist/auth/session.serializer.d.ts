import { PassportSerializer } from '@nestjs/passport';
import { UserService } from 'src/user/user.service';
import { User } from 'src/user/user.entity';
export declare class SessionSerializer extends PassportSerializer {
    private readonly userService;
    constructor(userService: UserService);
    serializeUser(user: User, done: (err: Error, user: {
        id: string;
    }) => void): void;
    deserializeUser(payload: {
        id: string;
    }, done: (err: Error, user: User) => void): Promise<void>;
}
