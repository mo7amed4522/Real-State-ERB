import { InputType, Field, ID } from '@nestjs/graphql';
import { IsString, IsNotEmpty, IsUUID, IsNumber, Min, IsDateString } from 'class-validator';

@InputType()
export class CreateOfferInput {
  @Field(() => ID)
  @IsUUID()
  company_id: string;

  @Field(() => ID)
  @IsUUID()
  developer_id: string;

  @Field(() => ID)
  @IsUUID()
  building_id: string;

  @Field()
  @IsNumber()
  @Min(0)
  offer_price: number;
  
  @Field()
  @IsDateString()
  offer_date: string;
} 