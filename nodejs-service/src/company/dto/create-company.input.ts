import { InputType, Field } from '@nestjs/graphql';
import { IsString, IsEmail, IsUrl, IsEnum, IsDateString, IsNotEmpty } from 'class-validator';
import { LegalStatus } from '../company.entity';

@InputType()
export class CreateCompanyInput {
  @Field()
  @IsString()
  @IsNotEmpty()
  name: string;

  @Field()
  @IsString()
  @IsNotEmpty()
  trade_license_number: string;

  @Field(() => LegalStatus)
  @IsEnum(LegalStatus)
  legal_status: LegalStatus;

  @Field()
  @IsDateString()
  registration_date: string;

  @Field()
  @IsEmail()
  contact_email: string;

  @Field()
  @IsString()
  @IsNotEmpty()
  contact_phone: string;

  @Field({ nullable: true })
  @IsUrl()
  website?: string;

  @Field()
  @IsString()
  @IsNotEmpty()
  address: string;
} 