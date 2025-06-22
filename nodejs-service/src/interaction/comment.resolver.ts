import { Resolver, Mutation, Args, ID, Query } from '@nestjs/graphql';
import { UseGuards } from '@nestjs/common';
import { GqlAuthGuard } from '../auth/guards/gql-auth.guard';
import { CurrentUser } from '../auth/decorators/current-user.decorator';
import { User } from '../user/user.entity';
import { InteractionService } from './interaction.service';
import { Comment } from './comment.entity';
import { CreateCommentInput } from './dto/create-comment.input';

@Resolver(() => Comment)
export class CommentResolver {
  constructor(private readonly interactionService: InteractionService) {}

  @Mutation(() => Comment)
  @UseGuards(GqlAuthGuard)
  createComment(
    @CurrentUser() user: User,
    @Args('createCommentInput') createCommentInput: CreateCommentInput,
  ): Promise<Comment> {
    return this.interactionService.createComment(user.id, createCommentInput);
  }

  @Mutation(() => String)
  @UseGuards(GqlAuthGuard)
  deleteComment(
    @CurrentUser() user: User,
    @Args('commentId', { type: () => ID }) commentId: string,
  ): Promise<string> {
    return this.interactionService.deleteComment(commentId, user.id).then(res => res.message);
  }
} 