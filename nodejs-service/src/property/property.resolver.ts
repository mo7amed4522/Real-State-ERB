import { Resolver, Mutation, Args, Subscription, Query, Int, ObjectType, Field, ID } from '@nestjs/graphql';
import { UseGuards } from '@nestjs/common';
import { GqlAuthGuard } from '../auth/guards/gql-auth.guard';
import { CurrentUser } from '../auth/decorators/current-user.decorator';
import { User } from '../user/user.entity';
import { PropertyService } from './property.service';
import { Property } from './property.entity';
import { PubSub } from 'graphql-subscriptions';

const pubSub = new PubSub();

@ObjectType()
class Comment {
  @Field(() => ID)
  id: string;

  @Field()
  content: string;

  @Field()
  userId: string;

  @Field()
  createdAt: string;
}

@Resolver()
export class PropertyResolver {
  constructor(private readonly propertyService: PropertyService) {}

  @Mutation(() => User)
  @UseGuards(GqlAuthGuard)
  async toggleFavoriteProperty(
    @CurrentUser() user: User,
    @Args('propertyId', { type: () => Int }) propertyId: number,
  ): Promise<User> {
    return this.propertyService.toggleFavorite(user.id, propertyId);
  }

  @Mutation(() => Property)
  async createProperty(@Args('input') input: any) {
    const property = await this.propertyService.createProperty(input);
    await (pubSub as any).publish('postCreated', { postCreated: property });
    return property;
  }

  @Subscription(() => Property, {
    resolve: (value) => value.postCreated,
  })
  postCreated() {
    return (pubSub as any).asyncIterator('postCreated');
  }

  @Mutation(() => Comment)
  @UseGuards(GqlAuthGuard)
  async createComment(
    @CurrentUser() user: User,
    @Args('propertyId', { type: () => Int }) propertyId: number,
    @Args('content') content: string,
    @Args('parentId', { type: () => Int, nullable: true }) parentId?: number,
  ): Promise<any> {
    const comment = await this.propertyService.createComment(user.id, propertyId, content, parentId);
    await (pubSub as any).publish('commentCreated', { commentCreated: comment });
    return comment;
  }

  @Subscription(() => Comment, {
    resolve: (value) => value.commentCreated,
  })
  commentCreated() {
    return (pubSub as any).asyncIterator('commentCreated');
  }

  @Mutation(() => User)
  @UseGuards(GqlAuthGuard)
  async toggleCommentLike(
    @CurrentUser() user: User,
    @Args('commentId', { type: () => Int }) commentId: number,
  ): Promise<User> {
    return this.propertyService.toggleCommentLike(user.id, commentId);
  }

  @Mutation(() => User)
  async likeProperty(@Args('propertyId', { type: () => Int }) propertyId: number, @Args('userId', { type: () => Int }) userId: number) {
    // TODO: Implement likeProperty in PropertyService
    if (!(this.propertyService as any).likeProperty) {
      (this.propertyService as any).likeProperty = async (propertyId: number, userId: number) => {
        // Dummy implementation
        return {} as User;
      };
    }
    const like = await (this.propertyService as any).likeProperty(propertyId, userId);
    await (pubSub as any).publish('likeCreated', { likeCreated: like });
    return like;
  }

  @Subscription(() => User, {
    resolve: (value) => value.likeCreated,
  })
  likeCreated() {
    return (pubSub as any).asyncIterator('likeCreated');
  }

  @Mutation(() => Comment)
  async replyToComment(@Args('input') input: any) {
    // TODO: Implement replyToComment in PropertyService
    if (!(this.propertyService as any).replyToComment) {
      (this.propertyService as any).replyToComment = async (input: any) => {
        // Dummy implementation
        return {} as any;
      };
    }
    const reply = await (this.propertyService as any).replyToComment(input);
    await (pubSub as any).publish('replyCreated', { replyCreated: reply });
    return reply;
  }

  @Subscription(() => Comment, {
    resolve: (value) => value.replyCreated,
  })
  replyCreated() {
    return (pubSub as any).asyncIterator('replyCreated');
  }
} 