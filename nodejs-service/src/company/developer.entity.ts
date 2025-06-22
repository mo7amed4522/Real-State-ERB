import { Entity, PrimaryGeneratedColumn, Column, CreateDateColumn, UpdateDateColumn, ManyToOne, JoinColumn } from 'typeorm';
import { Company } from './company.entity';

export enum DeveloperStatus {
  PENDING = 'Pending',
  VERIFIED = 'Verified',
  REJECTED = 'Rejected',
}

@Entity('developers')
export class Developer {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  full_name: string;

  @Column({ unique: true })
  email: string;

  @Column()
  phone: string;

  @Column({ type: 'uuid', nullable: true })
  company_id: string;

  @ManyToOne(() => Company, company => company.developers, { nullable: true })
  @JoinColumn({ name: 'company_id' })
  company: Company;

  @Column({ nullable: true })
  profile_photo_url: string;

  @Column({ unique: true })
  license_number: string;

  @Column()
  experience_years: number;

  @Column()
  specialization: string;

  @Column({
    type: 'enum',
    enum: DeveloperStatus,
    default: DeveloperStatus.PENDING,
  })
  status: DeveloperStatus;

  @Column({ type: 'int', default: 0 })
  total_buildings: number;

  @Column({ type: 'int', default: 0 })
  total_offers: number;

  @Column({ type: 'int', default: 0 })
  total_likes: number;

  @Column({ type: 'int', default: 0 })
  total_comments: number;

  @CreateDateColumn()
  created_at: Date;

  @UpdateDateColumn()
  updated_at: Date;
} 