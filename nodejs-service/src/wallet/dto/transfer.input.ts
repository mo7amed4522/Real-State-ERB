import { InputType, Field } from '@nestjs/graphql';
import { IsNotEmpty, IsNumber, IsPositive, IsOptional, IsString, IsUUID } from 'class-validator';

@InputType()
export class TransferInput {
  @Field()
  @IsNotEmpty()
  @IsNumber()
  @IsPositive()
  amount: number;

  @Field()
  @IsNotEmpty()
  @IsUUID()
  receiverWalletId: string;

  @Field({ nullable: true })
  @IsOptional()
  @IsString()
  description?: string;

  @Field({ nullable: true })
  @IsOptional()
  @IsString()
  encryptedReceiverDetails?: string;
} 