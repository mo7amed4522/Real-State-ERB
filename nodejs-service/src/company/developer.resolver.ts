import { Resolver, Mutation, Args, Query, ID } from '@nestjs/graphql';
import { UseGuards } from '@nestjs/common';
import { GqlAuthGuard } from '../auth/guards/gql-auth.guard';
import { CompanyService } from './company.service';
import { Developer } from './developer.entity';
import { CreateDeveloperInput } from './dto/create-developer.input';
import { UpdateDeveloperInput } from './dto/update-developer.input';

@Resolver(() => Developer)
export class DeveloperResolver {
  constructor(private readonly companyService: CompanyService) {}

  @Mutation(() => Developer)
  @UseGuards(GqlAuthGuard)
  createDeveloper(
    @Args('createDeveloperInput') createDeveloperInput: CreateDeveloperInput,
  ) {
    return this.companyService.createDeveloper(createDeveloperInput);
  }
  
  @Query(() => Developer, { name: 'developer' })
  findOneDeveloper(@Args('id', { type: () => ID }) id: string) {
    return this.companyService.findOneDeveloper(id);
  }

  @Query(() => [Developer], { name: 'developers' })
  findAllDevelopers() {
    return this.companyService.findAllDevelopers();
  }
  
  @Mutation(() => Developer)
  @UseGuards(GqlAuthGuard)
  updateDeveloper(
    @Args('updateDeveloperInput') updateDeveloperInput: UpdateDeveloperInput,
  ) {
    return this.companyService.updateDeveloper(updateDeveloperInput.id, updateDeveloperInput);
  }

  @Mutation(() => String)
  @UseGuards(GqlAuthGuard)
  removeDeveloper(@Args('id', { type: () => ID }) id: string) {
    this.companyService.removeDeveloper(id);
    return 'Developer removed successfully';
  }
} 