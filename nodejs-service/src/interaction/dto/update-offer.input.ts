import { InputType, Field, ID } from '@nestjs/graphql';
import { IsUUID, IsEnum, IsDateString, IsOptional } from 'class-validator';
import { OfferStatus } from '../offer.entity';

@InputType()
export class UpdateOfferInput {
  @Field(() => ID)
  @IsUUID()
  id: string;

  @Field(() => OfferStatus)
  @IsEnum(OfferStatus)
  status: OfferStatus;

  @Field({ nullable: true })
  @IsDateString()
  @IsOptional()
  accepted_date?: string;
} 