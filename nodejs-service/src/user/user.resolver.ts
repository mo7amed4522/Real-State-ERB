import { Resolver, Query, Context, Mutation, Args, ResolveField, Parent } from '@nestjs/graphql';
import { UseGuards } from '@nestjs/common';
import { User, Role } from './user.entity';
import { UserService } from './user.service';
import { GqlAuthGuard } from '../auth/guards/gql-auth.guard';
import { UpdateUserInput } from './dto/update-user.input';
import { ChangePasswordInput } from './dto/change-password.input';
import { Roles } from '../auth/decorators/roles.decorator';
import { RolesGuard } from '../auth/guards/roles.guard';
import { FileUpload, GraphQLUpload } from 'graphql-upload';
import { FileService } from '../common/file.service';
import { EncryptionService } from '../common/encryption.service';

@Resolver(() => User)
@UseGuards(GqlAuthGuard)
export class UserResolver {
  constructor(
    private readonly userService: UserService,
    private readonly fileService: FileService,
    private readonly encryptionService: EncryptionService,
  ) {}

  @Query(() => User, { nullable: true })
  me(@Context() context): User {
    return context.req.user;
  }

  @Query(() => [User])
  @Roles(Role.ADMIN)
  @UseGuards(RolesGuard)
  allUsers(): Promise<User[]> {
    return this.userService.findAll();
  }

  @Mutation(() => User)
  async updateUser(
    @Args('updateUserInput') updateUserInput: UpdateUserInput,
    @Context() context,
  ): Promise<User> {
    const userId = context.req.user.id;
    return this.userService.update(userId, updateUserInput);
  }

  @Mutation(() => Boolean)
  async deleteUser(@Context() context): Promise<boolean> {
    const userId = context.req.user.id;
    return this.userService.delete(userId);
  }

  @Mutation(() => Boolean)
  @UseGuards(GqlAuthGuard)
  async changePassword(
    @Args('changePasswordInput') changePasswordInput: ChangePasswordInput,
    @Context() context,
  ): Promise<boolean> {
    const userId = context.req.user.id;
    return this.userService.changePassword(
      userId,
      changePasswordInput.oldPassword,
      changePasswordInput.newPassword,
    );
  }

  @Mutation(() => User)
  @UseGuards(GqlAuthGuard)
  async uploadProfilePicture(
    @Context() context,
    @Args({ name: 'file', type: () => GraphQLUpload })
    { createReadStream, filename }: FileUpload,
  ): Promise<User> {
    const userId = context.req.user.id;
    const stream = createReadStream();
    const filePath = await this.fileService.saveFile(stream, filename, userId);
    const encryptedPath = this.encryptionService.encrypt(filePath);

    return this.userService.update(userId, { photoUrl: encryptedPath });
  }

  @ResolveField(() => String, { nullable: true })
  photoUrl(@Parent() user: User): string | null {
    if (!user.photoUrl) {
      return null;
    }
    // If it's a google URL, return it directly
    if (user.photoUrl.startsWith('http')) {
      return user.photoUrl;
    }
    // Otherwise, build the secure URL to our serving endpoint
    return `http://localhost:3000/users/photo`;
  }
} 