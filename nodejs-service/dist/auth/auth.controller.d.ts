import { Request, Response } from 'express';
export declare class AuthController {
    googleAuth(req: Request): Promise<void>;
    googleAuthRedirect(req: Request, res: Response): void;
}
