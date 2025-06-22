import { Resolver, Mutation, Args, Query, Parent, ResolveField } from '@nestjs/graphql';
import { UseGuards } from '@nestjs/common';
import { GqlAuthGuard } from '../auth/guards/gql-auth.guard';
import { InteractionService } from './interaction.service';
import { Like } from './like.entity';
import { User } from '../user/user.entity';
import { CurrentUser } from '../auth/decorators/current-user.decorator';
import { ToggleLikeInput, LikeableType } from './dto/toggle-like.input';
import { registerEnumType } from '@nestjs/graphql';

registerEnumType(LikeableType, {
  name: 'LikeableType',
});

@Resolver(() => Like)
export class LikeResolver {
  constructor(private readonly interactionService: InteractionService) {}

  @Mutation(() => Boolean)
  @UseGuards(GqlAuthGuard)
  toggleLike(
    @CurrentUser() user: User,
    @Args('toggleLikeInput') toggleLikeInput: ToggleLikeInput,
  ): Promise<boolean> {
    return this.interactionService.toggleLike(user.id, toggleLikeInput);
  }
} 