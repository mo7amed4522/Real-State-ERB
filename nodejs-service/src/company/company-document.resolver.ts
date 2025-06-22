import { Resolver, ResolveField, Parent, Mutation, Args, ID } from '@nestjs/graphql';
import { CompanyDocument, DocumentType } from './company-document.entity';
import { EncryptionService } from '../common/encryption.service';
import { CompanyService } from './company.service';
import { GraphQLUpload, FileUpload } from 'graphql-upload';
import { UseGuards } from '@nestjs/common';
import { GqlAuthGuard } from '../auth/guards/gql-auth.guard';

@Resolver(() => CompanyDocument)
export class CompanyDocumentResolver {
  constructor(
    private readonly encryptionService: EncryptionService,
    private readonly companyService: CompanyService,
    ) {}

  @ResolveField('file_url', () => String)
  getFileUrl(@Parent() document: CompanyDocument): string {
    // The file_url from DB is the encrypted path.
    // The client will use this to construct the final URL.
    return document.file_url;
  }

  @Mutation(() => CompanyDocument)
  @UseGuards(GqlAuthGuard)
  addCompanyDocument(
    @Args('companyId', { type: () => ID }) companyId: string,
    @Args('file', { type: () => GraphQLUpload }) file: FileUpload,
    @Args('documentType', { type: () => DocumentType }) docType: DocumentType,
  ) {
    return this.companyService.addDocument(companyId, file, docType);
  }

  @Mutation(() => String)
  @UseGuards(GqlAuthGuard)
  removeCompanyDocument(@Args('id', { type: () => ID }) id: string) {
    this.companyService.removeDocument(id);
    return 'Document removed successfully';
  }
} 