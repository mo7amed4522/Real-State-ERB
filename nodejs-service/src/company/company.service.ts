import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Company } from './company.entity';
import { CreateCompanyInput } from './dto/create-company.input';
import { User } from '../user/user.entity';
import { CompanyUser, CompanyRole } from './company-user.entity';
import { FileService } from '../common/file.service';
import { EncryptionService } from '../common/encryption.service';
import { FileUpload } from 'graphql-upload';
import { v4 as uuidv4 } from 'uuid';
import { CompanyDocument, DocumentType } from './company-document.entity';
import { Developer } from './developer.entity';
import { CreateDeveloperInput } from './dto/create-developer.input';
import { UpdateDeveloperInput } from './dto/update-developer.input';

@Injectable()
export class CompanyService {
  constructor(
    @InjectRepository(Company)
    private companyRepository: Repository<Company>,
    @InjectRepository(CompanyUser)
    private companyUserRepository: Repository<CompanyUser>,
    @InjectRepository(CompanyDocument)
    private companyDocumentRepository: Repository<CompanyDocument>,
    @InjectRepository(Developer)
    private developerRepository: Repository<Developer>,
    private fileService: FileService,
    private encryptionService: EncryptionService,
  ) {}

  async createCompany(
    createCompanyInput: CreateCompanyInput,
    owner: User,
    logo?: FileUpload,
    documents?: FileUpload[],
  ): Promise<Company> {
    const company = this.companyRepository.create({
      ...createCompanyInput,
      registration_date: new Date(createCompanyInput.registration_date),
    });
    
    const savedCompany = await this.companyRepository.save(company);
    const companyFolder = `company/${savedCompany.id}`;

    if (logo) {
      const { createReadStream, filename } = await logo;
      const uniqueFilename = `${uuidv4()}-${filename}`;
      const logoPath = await this.fileService.savePrivateFile(
        createReadStream(),
        uniqueFilename,
        companyFolder,
      );
      savedCompany.logo_url = this.encryptionService.encrypt(logoPath);
    }
    
    await this.companyRepository.save(savedCompany);
    
    if (documents && documents.length > 0) {
      for (const doc of documents) {
        const { createReadStream, filename, mimetype } = await doc;
        // Basic mapping from filename to DocumentType, can be improved
        const docType = this.getDocTypeFromFilename(filename);
        const uniqueFilename = `${uuidv4()}-${filename}`;
        const docPath = await this.fileService.savePrivateFile(
          createReadStream(),
          uniqueFilename,
          `${companyFolder}/documents`,
        );
        const encryptedPath = this.encryptionService.encrypt(docPath);
        
        const companyDoc = this.companyDocumentRepository.create({
          company_id: savedCompany.id,
          document_type: docType,
          file_url: encryptedPath,
        });
        await this.companyDocumentRepository.save(companyDoc);
      }
    }

    const companyUser = this.companyUserRepository.create({
      user_id: owner.id,
      company_id: savedCompany.id,
      role: CompanyRole.OWNER,
    });
    await this.companyUserRepository.save(companyUser);

    return savedCompany;
  }

  async findOne(id: string): Promise<Company> {
    const company = await this.companyRepository.findOne({ where: { id } });
    if (!company) {
      throw new NotFoundException(`Company with ID ${id} not found`);
    }
    return company;
  }

  async findAll(): Promise<Company[]> {
    return this.companyRepository.find();
  }

  async update(id: string, updateCompanyInput: any): Promise<Company> {
    const company = await this.findOne(id);
    Object.assign(company, updateCompanyInput);
    return this.companyRepository.save(company);
  }

  async remove(id: string): Promise<{ id: string; message: string }> {
    const company = await this.findOne(id);
    await this.companyRepository.remove(company);
    return { id, message: 'Company removed successfully' };
  }

  private getDocTypeFromFilename(filename: string): DocumentType {
    if (filename.toLowerCase().includes('license')) return DocumentType.TRADE_LICENSE;
    if (filename.toLowerCase().includes('vat')) return DocumentType.VAT;
    if (filename.toLowerCase().includes('insurance')) return DocumentType.INSURANCE;
    return DocumentType.CHAMBER_OF_COMMERCE;
  }

  // Developer Methods
  async createDeveloper(createDeveloperInput: CreateDeveloperInput): Promise<Developer> {
    const developer = this.developerRepository.create(createDeveloperInput);
    return this.developerRepository.save(developer);
  }

  async findOneDeveloper(id: string): Promise<Developer> {
    const developer = await this.developerRepository.findOne({ where: { id } });
    if (!developer) {
      throw new NotFoundException(`Developer with ID ${id} not found`);
    }
    return developer;
  }

  async findAllDevelopers(): Promise<Developer[]> {
    return this.developerRepository.find();
  }

  async updateDeveloper(id: string, updateDeveloperInput: UpdateDeveloperInput): Promise<Developer> {
    const developer = await this.findOneDeveloper(id);
    Object.assign(developer, updateDeveloperInput);
    return this.developerRepository.save(developer);
  }

  async removeDeveloper(id: string): Promise<{ id: string; message: string }> {
    const developer = await this.findOneDeveloper(id);
    await this.developerRepository.remove(developer);
    return { id, message: 'Developer removed successfully' };
  }

  // Document Methods
  async addDocument(companyId: string, file: FileUpload, docType: DocumentType): Promise<CompanyDocument> {
    const { createReadStream, filename } = await file;
    const companyFolder = `company/${companyId}/documents`;
    const uniqueFilename = `${uuidv4()}-${filename}`;
    
    const docPath = await this.fileService.savePrivateFile(
      createReadStream(),
      uniqueFilename,
      companyFolder,
    );
    const encryptedPath = this.encryptionService.encrypt(docPath);
    
    const companyDoc = this.companyDocumentRepository.create({
      company_id: companyId,
      document_type: docType,
      file_url: encryptedPath,
    });
    
    return this.companyDocumentRepository.save(companyDoc);
  }

  async removeDocument(id: string): Promise<{ id: string; message: string }> {
    const doc = await this.companyDocumentRepository.findOne({ where: { id } });
    if (!doc) {
      throw new NotFoundException(`Document with ID ${id} not found`);
    }
    // Note: This only removes the DB record. 
    // You might want to also delete the file from storage.
    await this.companyDocumentRepository.remove(doc);
    return { id, message: 'Document removed successfully' };
  }
} 