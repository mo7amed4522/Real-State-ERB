import { Entity, PrimaryGeneratedColumn, Column, CreateDateColumn, ManyToOne, JoinColumn } from 'typeorm';
import { User } from '../user/user.entity';
import { Company } from './company.entity';

export enum CompanyRole {
  OWNER = 'Owner',
  HR = 'HR',
  ADMIN = 'Admin',
}

@Entity('company_users')
export class CompanyUser {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ type: 'uuid' })
  user_id: string;

  @ManyToOne(() => User)
  @JoinColumn({ name: 'user_id' })
  user: User;

  @Column({ type: 'uuid' })
  company_id: string;

  @ManyToOne(() => Company, company => company.users)
  @JoinColumn({ name: 'company_id' })
  company: Company;

  @Column({
    type: 'enum',
    enum: CompanyRole,
  })
  role: CompanyRole;

  @CreateDateColumn()
  created_at: Date;
} 