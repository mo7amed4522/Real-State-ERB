import { Resolver, Mutation, Args, Query, ResolveField, Parent, ID } from '@nestjs/graphql';
import { UseGuards } from '@nestjs/common';
import { GqlAuthGuard } from '../auth/guards/gql-auth.guard';
import { CurrentUser } from '../auth/decorators/current-user.decorator';
import { User } from '../user/user.entity';
import { CompanyService } from './company.service';
import { Company } from './company.entity';
import { CreateCompanyInput } from './dto/create-company.input';
import { UpdateCompanyInput } from './dto/update-company.input';
import { GraphQLUpload, FileUpload } from 'graphql-upload';
import { EncryptionService } from '../common/encryption.service';

@Resolver(() => Company)
export class CompanyResolver {
  constructor(
    private readonly companyService: CompanyService,
    private readonly encryptionService: EncryptionService,
  ) {}

  @Mutation(() => Company)
  @UseGuards(GqlAuthGuard)
  createCompany(
    @Args('createCompanyInput') createCompanyInput: CreateCompanyInput,
    @CurrentUser() user: User,
    @Args({ name: 'logo', type: () => GraphQLUpload, nullable: true })
    logo?: FileUpload,
    @Args({ name: 'documents', type: () => [GraphQLUpload], nullable: true })
    documents?: FileUpload[],
  ) {
    return this.companyService.createCompany(createCompanyInput, user, logo, documents);
  }

  @ResolveField('logo_url', () => String, { nullable: true })
  getLogoUrl(@Parent() company: Company): string | null {
    if (!company.logo_url) {
      return null;
    }
    // The logo_url from DB is the encrypted path. We just return it.
    // The client will use this encrypted path to construct the final URL
    // e.g., https://your-api.com/files/private/<encrypted-path>
    return company.logo_url;
  }

  @Query(() => Company, { name: 'company' })
  findOne(@Args('id', { type: () => ID }) id: string) {
    return this.companyService.findOne(id);
  }

  @Query(() => [Company], { name: 'companies' })
  findAll() {
    return this.companyService.findAll();
  }

  @Mutation(() => Company)
  @UseGuards(GqlAuthGuard)
  updateCompany(
    @Args('updateCompanyInput') updateCompanyInput: UpdateCompanyInput
  ) {
    return this.companyService.update(updateCompanyInput.id, updateCompanyInput);
  }

  @Mutation(() => String)
  @UseGuards(GqlAuthGuard)
  removeCompany(@Args('id', { type: () => ID }) id: string) {
    this.companyService.remove(id);
    return 'Company removed successfully';
  }
} 