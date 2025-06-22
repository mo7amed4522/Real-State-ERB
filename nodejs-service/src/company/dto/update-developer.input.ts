import { InputType, Field, ID } from '@nestjs/graphql';
import { IsString, IsEmail, IsOptional, IsUUID, IsInt, Min, IsEnum } from 'class-validator';
import { DeveloperStatus } from '../developer.entity';

@InputType()
export class UpdateDeveloperInput {
  @Field(() => ID)
  @IsUUID()
  id: string;

  @Field({ nullable: true })
  @IsString()
  @IsOptional()
  full_name?: string;

  @Field({ nullable: true })
  @IsEmail()
  @IsOptional()
  email?: string;

  @Field({ nullable: true })
  @IsString()
  @IsOptional()
  phone?: string;

  @Field(() => ID, { nullable: true })
  @IsUUID()
  @IsOptional()
  company_id?: string;
  
  @Field({ nullable: true })
  @IsString()
  @IsOptional()
  license_number?: string;

  @Field(() => Number, { nullable: true })
  @IsInt()
  @Min(0)
  @IsOptional()
  experience_years?: number;

  @Field({ nullable: true })
  @IsString()
  @IsOptional()
  specialization?: string;

  @Field(() => DeveloperStatus, { nullable: true })
  @IsEnum(DeveloperStatus)
  @IsOptional()
  status?: DeveloperStatus;
} 