import { User } from './user.entity';
import { UserService } from './user.service';
export declare class UserResolver {
    private readonly userService;
    constructor(userService: UserService);
    me(context: any): Promise<User | undefined>;
}
