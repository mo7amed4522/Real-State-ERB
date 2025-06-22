import { InputType, Field, ID } from '@nestjs/graphql';
import { IsString, IsNotEmpty, IsUUID, IsOptional, IsEnum } from 'class-validator';
import { CommentableType } from '../comment.entity';

@InputType()
export class CreateCommentInput {
  @Field(() => ID, { nullable: true })
  @IsUUID()
  @IsOptional()
  parent_id?: string;

  @Field(() => CommentableType)
  @IsEnum(CommentableType)
  target_type: CommentableType;

  @Field(() => ID)
  @IsUUID()
  target_id: string;

  @Field()
  @IsString()
  @IsNotEmpty()
  content: string;
} 