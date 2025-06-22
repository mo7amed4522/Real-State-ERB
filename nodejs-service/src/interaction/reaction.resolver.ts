import { Resolver, Mutation, Args, ID } from '@nestjs/graphql';
import { UseGuards } from '@nestjs/common';
import { GqlAuthGuard } from '../auth/guards/gql-auth.guard';
import { CurrentUser } from '../auth/decorators/current-user.decorator';
import { User } from '../user/user.entity';
import { InteractionService } from './interaction.service';
import { Emoji } from './reaction.entity';

@Resolver()
export class ReactionResolver {
  constructor(private readonly interactionService: InteractionService) {}

  @Mutation(() => Boolean)
  @UseGuards(GqlAuthGuard)
  async toggleReaction(
    @CurrentUser() user: User,
    @Args('commentId', { type: () => ID }) commentId: string,
    @Args('emoji', { type: () => Emoji }) emoji: Emoji,
  ): Promise<boolean> {
    const { reacted } = await this.interactionService.toggleReaction(
      user.id,
      commentId,
      emoji,
    );
    return reacted;
  }
} 