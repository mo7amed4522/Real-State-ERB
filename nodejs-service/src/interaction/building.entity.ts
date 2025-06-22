import { Entity, PrimaryGeneratedColumn, Column, CreateDateColumn, ManyToOne, JoinColumn, OneToMany } from 'typeorm';
import { Company } from '../company/company.entity';
import { Developer } from '../company/developer.entity';
import { Offer } from './offer.entity';
import { ObjectType, Field, ID, Float, Int } from '@nestjs/graphql';

export enum BuildingStatus {
  SOLD = 'sold',
  AVAILABLE = 'available',
}

@ObjectType()
@Entity('buildings')
export class Building {
  @PrimaryGeneratedColumn('uuid')
  @Field(() => ID)
  id: string;

  @Column()
  @Field()
  title: string;
  
  @Column({ type: 'text', nullable: true })
  @Field({ nullable: true })
  description: string;

  @Column({ type: 'float' })
  @Field(() => Float)
  latitude: number;

  @Column({ type: 'float' })
  @Field(() => Float)
  longitude: number;

  @Column()
  @Field()
  address: string;

  @Column()
  @Field()
  city: string;

  @Column()
  @Field()
  region: string;

  @Column({ type: 'decimal', precision: 10, scale: 2 })
  @Field(() => Float)
  price: number;

  @Column({
    type: 'enum',
    enum: BuildingStatus,
    default: BuildingStatus.AVAILABLE,
  })
  @Field(() => BuildingStatus)
  status: BuildingStatus;

  @Column({ type: 'timestamp', nullable: true })
  @Field({ nullable: true })
  sold_at: Date;

  @Column()
  @Field(() => ID)
  company_id: string;

  @ManyToOne(() => Company)
  @JoinColumn({ name: 'company_id' })
  company: Company;

  @Column()
  @Field(() => ID)
  developer_id: string;

  @ManyToOne(() => Developer)
  @JoinColumn({ name: 'developer_id' })
  developer: Developer;

  @OneToMany(() => Offer, (offer) => offer.building)
  offers: Offer[];

  @Column({ default: 0 })
  @Field(() => Int)
  total_likes: number;

  @Column({ default: 0 })
  @Field(() => Int)
  total_comments: number;

  @Column({ default: 0 })
  @Field(() => Int)
  total_views: number;

  @CreateDateColumn()
  @Field()
  created_at: Date;

  @Field(() => Boolean, { nullable: true })
  is_liked?: boolean;
} 