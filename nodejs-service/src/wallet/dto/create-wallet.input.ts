import { InputType, Field } from '@nestjs/graphql';
import { IsOptional, IsString, IsEnum } from 'class-validator';
import { PaymentMethod } from '../wallet.entity';

@InputType()
export class CreateWalletInput {
  @Field({ nullable: true })
  @IsOptional()
  @IsString()
  stripeCustomerId?: string;

  @Field({ nullable: true })
  @IsOptional()
  @IsString()
  encryptedBankDetails?: string;

  @Field(() => PaymentMethod, { nullable: true })
  @IsOptional()
  @IsEnum(PaymentMethod)
  paymentMethod?: PaymentMethod;
} 