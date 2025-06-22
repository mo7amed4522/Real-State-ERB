import { Entity, PrimaryGeneratedColumn, Column, CreateDateColumn, ManyToOne, JoinColumn } from 'typeorm';
import { Company } from '../company/company.entity';
import { Developer } from '../company/developer.entity';
import { Building } from './building.entity';
import { ObjectType, Field, ID, Float } from '@nestjs/graphql';

export enum OfferStatus {
  ACCEPTED = 'accepted',
  REJECTED = 'rejected',
  PENDING = 'pending',
}

@ObjectType()
@Entity('offers')
export class Offer {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ type: 'uuid' })
  company_id: string;

  @ManyToOne(() => Company)
  @JoinColumn({ name: 'company_id' })
  company: Company;

  @Column({ type: 'uuid' })
  developer_id: string;
  
  @ManyToOne(() => Developer)
  @JoinColumn({ name: 'developer_id' })
  developer: Developer;

  @Column()
  @Field(() => ID)
  building_id: string;

  @ManyToOne(() => Building, (building) => building.offers)
  @JoinColumn({ name: 'building_id' })
  building: Building;

  @Column('float')
  offer_price: number;

  @Column({ type: 'enum', enum: OfferStatus, default: OfferStatus.PENDING })
  status: OfferStatus;

  @Column()
  offer_date: Date;

  @Column({ nullable: true })
  accepted_date: Date;
} 