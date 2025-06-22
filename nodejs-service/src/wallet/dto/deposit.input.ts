import { InputType, Field } from '@nestjs/graphql';
import { IsNotEmpty, IsNumber, IsPositive, IsOptional, IsString, IsEnum } from 'class-validator';
import { PaymentMethod } from '../wallet.entity';

@InputType()
export class DepositInput {
  @Field()
  @IsNotEmpty()
  @IsNumber()
  @IsPositive()
  amount: number;

  @Field(() => PaymentMethod)
  @IsNotEmpty()
  @IsEnum(PaymentMethod)
  paymentMethod: PaymentMethod;

  @Field({ nullable: true })
  @IsOptional()
  @IsString()
  description?: string;

  @Field({ nullable: true })
  @IsOptional()
  @IsString()
  encryptedBankDetails?: string;

  @Field({ nullable: true })
  @IsOptional()
  @IsString()
  encryptedAccountName?: string;
} 