import { Resolver, Mutation, Args, Context } from '@nestjs/graphql';
import { UseGuards } from '@nestjs/common';
import { User } from '../user/user.entity';
import { UserService } from '../user/user.service';
import { RegisterInput } from './dto/register.input';
import { LoginInput } from './dto/login.input';
import { LocalAuthGuard } from './guards/local-auth.guard';
import { AuthService } from './auth.service';
import { ForgotPasswordInput } from './dto/forgot-password.input';
import { ResetPasswordInput } from './dto/reset-password.input';

@Resolver()
export class AuthResolver {
  constructor(
    private readonly userService: UserService,
    private readonly authService: AuthService,
    ) {}

  @Mutation(() => User)
  async register(@Args('registerInput') registerInput: RegisterInput): Promise<User> {
    return this.userService.create(registerInput);
  }

  @Mutation(() => User)
  @UseGuards(LocalAuthGuard)
  async login(@Args('loginInput') loginInput: LoginInput, @Context() context): Promise<User> {
    // The LocalAuthGuard populates context.user
    return context.user;
  }

  @Mutation(() => String)
  async forgotPassword(@Args('forgotPasswordInput') forgotPasswordInput: ForgotPasswordInput): Promise<string> {
    const { resetToken } = await this.authService.forgotPassword(forgotPasswordInput.email);
    return resetToken;
  }

  @Mutation(() => User)
  async resetPassword(@Args('resetPasswordInput') resetPasswordInput: ResetPasswordInput): Promise<User> {
    return this.authService.resetPassword(resetPasswordInput.token, resetPasswordInput.password);
  }
} 