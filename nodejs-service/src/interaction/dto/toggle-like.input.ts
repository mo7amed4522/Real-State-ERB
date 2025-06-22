import { InputType, Field, ID } from '@nestjs/graphql';
import { IsUUID, IsEnum } from 'class-validator';

export enum LikeableType {
  COMPANY = 'COMPANY',
  DEVELOPER = 'DEVELOPER',
  BUILDING = 'BUILDING',
}

@InputType()
export class ToggleLikeInput {
  @Field(() => ID)
  @IsUUID()
  entityId: string;

  @Field(() => LikeableType)
  @IsEnum(LikeableType)
  entityType: LikeableType;
} 