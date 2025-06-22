import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Company } from './company.entity';
import { Developer } from './developer.entity';
import { CompanyDocument } from './company-document.entity';
import { CompanyUser } from './company-user.entity';
import { CompanyService } from './company.service';
import { CompanyResolver } from './company.resolver';
import { CompanyDocumentResolver } from './company-document.resolver';
import { DeveloperResolver } from './developer.resolver';
import { CommonModule } from '../common/common.module';
// Import Service and Resolver here later

@Module({
  imports: [
    TypeOrmModule.forFeature([
      Company,
      Developer,
      CompanyDocument,
      CompanyUser,
    ]),
    CommonModule,
  ],
  providers: [CompanyService, CompanyResolver, CompanyDocumentResolver, DeveloperResolver],
  // Add providers and exports later
})
export class CompanyModule {} 