import { InputType, Field, ID, Float } from '@nestjs/graphql';
import { IsString, IsOptional, IsUUID, IsNumber, Min, IsEnum, IsDateString, IsLatitude, IsLongitude } from 'class-validator';
import { BuildingStatus } from '../building.entity';

@InputType()
export class UpdateBuildingInput {
  @Field(() => ID)
  @IsUUID()
  id: string;

  @Field({ nullable: true })
  @IsString()
  @IsOptional()
  title?: string;

  @Field({ nullable: true })
  @IsString()
  @IsOptional()
  description?: string;

  @Field({ nullable: true })
  @IsString()
  @IsOptional()
  address?: string;

  @Field({ nullable: true })
  @IsString()
  @IsOptional()
  city?: string;

  @Field({ nullable: true })
  @IsString()
  @IsOptional()
  region?: string;

  @Field(() => Float, { nullable: true })
  @IsLatitude()
  @IsOptional()
  latitude?: number;

  @Field(() => Float, { nullable: true })
  @IsLongitude()
  @IsOptional()
  longitude?: number;

  @Field({ nullable: true })
  @IsNumber()
  @Min(0)
  @IsOptional()
  price?: number;

  @Field(() => BuildingStatus, { nullable: true })
  @IsEnum(BuildingStatus)
  @IsOptional()
  status?: BuildingStatus;

  @Field({ nullable: true })
  @IsDateString()
  @IsOptional()
  sold_at?: string;
} 