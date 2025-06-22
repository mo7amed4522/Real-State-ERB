import { Entity, PrimaryGeneratedColumn, Column, CreateDateColumn, ManyToOne, JoinColumn } from 'typeorm';
import { Company } from './company.entity';

export enum DocumentType {
  TRADE_LICENSE = 'TradeLicense',
  VAT = 'VAT',
  INSURANCE = 'Insurance',
  CHAMBER_OF_COMMERCE = 'ChamberOfCommerce',
}

@Entity('company_documents')
export class CompanyDocument {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ type: 'uuid' })
  company_id: string;

  @ManyToOne(() => Company, company => company.documents)
  @JoinColumn({ name: 'company_id' })
  company: Company;

  @Column({
    type: 'enum',
    enum: DocumentType,
  })
  document_type: DocumentType;

  @Column()
  file_url: string;

  @Column({ default: false })
  verified: boolean;
  
  @CreateDateColumn()
  uploaded_at: Date;
} 