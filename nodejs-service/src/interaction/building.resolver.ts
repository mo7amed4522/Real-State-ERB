import { Resolver, Mutation, Args, Query, ID } from '@nestjs/graphql';
import { UseGuards } from '@nestjs/common';
import { GqlAuthGuard } from '../auth/guards/gql-auth.guard';
import { InteractionService } from './interaction.service';
import { Building } from './building.entity';
import { CreateBuildingInput } from './dto/create-building.input';
import { UpdateBuildingInput } from './dto/update-building.input';
import { GetBuildingsArgs } from './dto/get-buildings.args';
import { CurrentUser } from '../auth/decorators/current-user.decorator';
import { User } from '../user/user.entity';
import { OptionalGqlAuthGuard } from '../auth/guards/optional-gql-auth.guard';
import { LikeableType } from './dto/toggle-like.input';

@Resolver(() => Building)
export class BuildingResolver {
  constructor(private readonly interactionService: InteractionService) {}

  @Mutation(() => Building)
  @UseGuards(GqlAuthGuard)
  createBuilding(
    @Args('createBuildingInput') createBuildingInput: CreateBuildingInput,
  ) {
    return this.interactionService.createBuilding(createBuildingInput);
  }

  @Query(() => Building, { name: 'building' })
  @UseGuards(OptionalGqlAuthGuard)
  async findOneBuilding(
    @Args('id', { type: () => ID }) id: string,
    @CurrentUser() user?: User
  ) {
    const building = await this.interactionService.findOneBuilding(id);
    if (user) {
      building.is_liked = await this.interactionService.checkIfLiked(id, LikeableType.BUILDING, user.id);
    }
    return building;
  }

  @Query(() => [Building], { name: 'buildings' })
  @UseGuards(OptionalGqlAuthGuard)
  findAllBuildings(
    @Args() args: GetBuildingsArgs,
    @CurrentUser() user?: User
  ) {
    return this.interactionService.findAllBuildings(args, user?.id);
  }

  @Mutation(() => Building)
  @UseGuards(GqlAuthGuard)
  updateBuilding(
    @Args('updateBuildingInput') updateBuildingInput: UpdateBuildingInput,
  ) {
    return this.interactionService.updateBuilding(updateBuildingInput.id, updateBuildingInput);
  }

  @Mutation(() => String)
  @UseGuards(GqlAuthGuard)
  async removeBuilding(@Args('id', { type: () => ID }) id: string) {
    await this.interactionService.removeBuilding(id);
    return 'Building removed successfully';
  }
} 