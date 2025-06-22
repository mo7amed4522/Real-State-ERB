import { Entity, PrimaryGeneratedColumn, Column, CreateDateColumn, ManyToOne, JoinColumn, Check } from 'typeorm';
import { User } from '../user/user.entity';
import { Company } from '../company/company.entity';
import { Developer } from '../company/developer.entity';
import { Building } from './building.entity';
import { ObjectType, Field, ID } from '@nestjs/graphql';

@ObjectType()
@Entity('likes')
@Check(`("company_id" IS NOT NULL AND "developer_id" IS NULL AND "building_id" IS NULL) OR ("company_id" IS NULL AND "developer_id" IS NOT NULL AND "building_id" IS NULL) OR ("company_id" IS NULL AND "developer_id" IS NULL AND "building_id" IS NOT NULL)`)
export class Like {
  @PrimaryGeneratedColumn('uuid')
  @Field(() => ID)
  id: string;

  @Column({ type: 'uuid' })
  user_id: string;

  @ManyToOne(() => User)
  @JoinColumn({ name: 'user_id' })
  user: User;

  @CreateDateColumn()
  created_at: Date;

  @Column({ type: 'uuid', nullable: true })
  @Field(() => ID, { nullable: true })
  company_id?: string;

  @ManyToOne(() => Company)
  @JoinColumn({ name: 'company_id' })
  company?: Company;

  @Column({ type: 'uuid', nullable: true })
  @Field(() => ID, { nullable: true })
  developer_id?: string;

  @ManyToOne(() => Developer)
  @JoinColumn({ name: 'developer_id' })
  developer?: Developer;
  
  @Column({ type: 'uuid', nullable: true })
  @Field(() => ID, { nullable: true })
  building_id?: string;

  @ManyToOne(() => Building)
  @JoinColumn({ name: 'building_id' })
  building?: Building;
} 