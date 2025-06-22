import { InputType, Field, ID } from '@nestjs/graphql';
import { IsString, IsEmail, IsUrl, IsEnum, IsDateString, IsOptional, IsUUID } from 'class-validator';
import { LegalStatus } from '../company.entity';

@InputType()
export class UpdateCompanyInput {
  @Field(() => ID)
  @IsUUID()
  id: string;

  @Field({ nullable: true })
  @IsString()
  @IsOptional()
  name?: string;

  @Field({ nullable: true })
  @IsString()
  @IsOptional()
  trade_license_number?: string;

  @Field(() => LegalStatus, { nullable: true })
  @IsEnum(LegalStatus)
  @IsOptional()
  legal_status?: LegalStatus;

  @Field({ nullable: true })
  @IsDateString()
  @IsOptional()
  registration_date?: string;

  @Field({ nullable: true })
  @IsEmail()
  @IsOptional()
  contact_email?: string;

  @Field({ nullable: true })
  @IsString()
  @IsOptional()
  contact_phone?: string;

  @Field({ nullable: true })
  @IsUrl()
  @IsOptional()
  website?: string;

  @Field({ nullable: true })
  @IsString()
  @IsOptional()
  address?: string;
} 