import { Entity, PrimaryGeneratedColumn, Column, CreateDateColumn, UpdateDateColumn, OneToMany } from 'typeorm';
import { Developer } from './developer.entity';
import { CompanyDocument } from './company-document.entity';
import { CompanyUser } from './company-user.entity';

export enum LegalStatus {
  LLC = 'LLC',
  JOINT_STOCK = 'JointStock',
  PARTNERSHIP = 'Partnership',
}

@Entity('companies')
export class Company {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  name: string;

  @Column({ unique: true })
  trade_license_number: string;

  @Column({ type: 'enum', enum: LegalStatus, default: LegalStatus.LLC })
  legal_status: LegalStatus;

  @Column()
  registration_date: Date;

  @Column({ unique: true })
  contact_email: string;

  @Column()
  contact_phone: string;

  @Column({ nullable: true })
  website: string;

  @Column()
  address: string;

  @Column({ nullable: true })
  logo_url: string;

  @Column({ default: false })
  verified: boolean;

  @Column({ type: 'int', default: 0 })
  total_buildings: number;

  @Column({ type: 'int', default: 0 })
  total_offers: number;

  @Column({ type: 'int', default: 0 })
  total_likes: number;

  @Column({ type: 'int', default: 0 })
  total_comments: number;
  
  @OneToMany(() => Developer, developer => developer.company)
  developers: Developer[];

  @OneToMany(() => CompanyDocument, document => document.company)
  documents: CompanyDocument[];

  @OneToMany(() => CompanyUser, companyUser => companyUser.company)
  users: CompanyUser[];

  @CreateDateColumn()
  created_at: Date;

  @UpdateDateColumn()
  updated_at: Date;
} 