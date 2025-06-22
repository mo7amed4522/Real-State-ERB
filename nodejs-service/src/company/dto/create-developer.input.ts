import { InputType, Field, ID } from '@nestjs/graphql';
import { IsString, IsEmail, IsNotEmpty, IsOptional, IsUUID, IsInt, Min } from 'class-validator';

@InputType()
export class CreateDeveloperInput {
  @Field()
  @IsString()
  @IsNotEmpty()
  full_name: string;

  @Field()
  @IsEmail()
  email: string;

  @Field()
  @IsString()
  @IsNotEmpty()
  phone: string;

  @Field(() => ID, { nullable: true })
  @IsUUID()
  @IsOptional()
  company_id?: string;

  @Field()
  @IsString()
  @IsNotEmpty()
  license_number: string;

  @Field(() => Number)
  @IsInt()
  @Min(0)
  experience_years: number;

  @Field()
  @IsString()
  @IsNotEmpty()
  specialization: string;
} 