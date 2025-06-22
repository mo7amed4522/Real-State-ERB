import { Injectable, NotFoundException, BadRequestException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository, In } from 'typeorm';
import { Like } from './like.entity';
import { Company } from '../company/company.entity';
import { Developer } from '../company/developer.entity';
import { User } from '../user/user.entity';
import { Comment, CommentableType } from './comment.entity';
import { CreateCommentInput } from './dto/create-comment.input';
import { Reaction, Emoji } from './reaction.entity';
import { Building } from './building.entity';
import { CreateBuildingInput } from './dto/create-building.input';
import { UpdateBuildingInput } from './dto/update-building.input';
import { Offer } from './offer.entity';
import { CreateOfferInput } from './dto/create-offer.input';
import { UpdateOfferInput } from './dto/update-offer.input';
import { ToggleLikeInput, LikeableType } from './dto/toggle-like.input';
import { GetBuildingsArgs } from './dto/get-buildings.args';

@Injectable()
export class InteractionService {
  constructor(
    @InjectRepository(Like)
    private likeRepository: Repository<Like>,
    @InjectRepository(Company)
    private companyRepository: Repository<Company>,
    @InjectRepository(Developer)
    private developerRepository: Repository<Developer>,
    @InjectRepository(Comment)
    private commentRepository: Repository<Comment>,
    @InjectRepository(Reaction)
    private reactionRepository: Repository<Reaction>,
    @InjectRepository(Building)
    private buildingRepository: Repository<Building>,
    @InjectRepository(Offer)
    private offerRepository: Repository<Offer>,
  ) {}

  async toggleLike(userId: string, input: ToggleLikeInput): Promise<boolean> {
    const { entityId, entityType } = input;
    
    let existingLike: Like;
    let counterUpdateArgs: { id: string, repository: Repository<any>, field: 'total_likes' };

    switch (entityType) {
      case LikeableType.COMPANY:
        existingLike = await this.likeRepository.findOne({ where: { user_id: userId, company_id: entityId } });
        counterUpdateArgs = { id: entityId, repository: this.companyRepository, field: 'total_likes' };
        break;
      case LikeableType.DEVELOPER:
        existingLike = await this.likeRepository.findOne({ where: { user_id: userId, developer_id: entityId } });
        counterUpdateArgs = { id: entityId, repository: this.developerRepository, field: 'total_likes' };
        break;
      case LikeableType.BUILDING:
        existingLike = await this.likeRepository.findOne({ where: { user_id: userId, building_id: entityId } });
        counterUpdateArgs = { id: entityId, repository: this.buildingRepository, field: 'total_likes' };
        break;
      default:
        throw new BadRequestException('Invalid likeable type');
    }

    if (existingLike) {
      // Unlike
      await this.likeRepository.remove(existingLike);
      await counterUpdateArgs.repository.decrement({ id: counterUpdateArgs.id }, counterUpdateArgs.field, 1);
      return false; // unliked
    } else {
      // Like
      const newLike = this.likeRepository.create({
        user_id: userId,
        [`${entityType.toLowerCase()}_id`]: entityId,
      });
      await this.likeRepository.save(newLike);
      await counterUpdateArgs.repository.increment({ id: counterUpdateArgs.id }, counterUpdateArgs.field, 1);
      return true; // liked
    }
  }

  async createComment(userId: string, input: CreateCommentInput): Promise<Comment> {
    const newComment = this.commentRepository.create({
      user_id: userId,
      ...input,
    });
    const savedComment = await this.commentRepository.save(newComment);
    await this.updateCounter(input.target_id, input.target_type, 'total_comments', 1);
    return savedComment;
  }

  async deleteComment(commentId: string, userId: string): Promise<{ id: string; message: string }> {
    const comment = await this.commentRepository.findOne({
      where: { id: commentId, user_id: userId },
    });
    if (!comment) {
      throw new NotFoundException(`Comment with ID ${commentId} not found or user does not have permission to delete.`);
    }
    await this.commentRepository.remove(comment);
    await this.updateCounter(comment.target_id, comment.target_type, 'total_comments', -1);
    return { id: commentId, message: 'Comment removed successfully' };
  }

  async toggleReaction(
    userId: string,
    commentId: string,
    emoji: Emoji,
  ): Promise<{ reacted: boolean }> {
    const existingReaction = await this.reactionRepository.findOne({
      where: { user_id: userId, comment_id: commentId, emoji },
    });

    if (existingReaction) {
      // Remove reaction
      await this.reactionRepository.remove(existingReaction);
      await this.commentRepository.decrement({ id: commentId }, 'total_reactions', 1);
      return { reacted: false };
    } else {
      // Add reaction
      const newReaction = this.reactionRepository.create({
        user_id: userId,
        comment_id: commentId,
        emoji,
      });
      await this.reactionRepository.save(newReaction);
      await this.commentRepository.increment({ id: commentId }, 'total_reactions', 1);
      return { reacted: true };
    }
  }

  private async updateCounter(
    targetId: string,
    targetType: LikeableType | CommentableType,
    field: 'total_likes' | 'total_comments',
    change: 1 | -1,
  ) {
    if (targetType === 'company') {
      await this.companyRepository.increment({ id: targetId }, field, change);
    } else if (targetType === 'developer') {
      await this.developerRepository.increment({ id: targetId }, field, change);
    }
  }

  // Building Methods
  async createBuilding(input: CreateBuildingInput): Promise<Building> {
    const building = this.buildingRepository.create(input);
    const savedBuilding = await this.buildingRepository.save(building);
    
    // Update counters
    await this.companyRepository.increment({ id: input.company_id }, 'total_buildings', 1);
    await this.developerRepository.increment({ id: input.developer_id }, 'total_buildings', 1);
    
    return savedBuilding;
  }

  async findOneBuilding(id: string): Promise<Building> {
    const building = await this.buildingRepository.findOne({ where: { id } });
    if (!building) {
      throw new NotFoundException(`Building with ID ${id} not found`);
    }
    // Increment views
    await this.buildingRepository.increment({ id }, 'total_views', 1);
    
    // Return a fresh instance to include the updated view count
    return this.buildingRepository.findOne({ where: { id } });
  }

  async findAllBuildings(args: GetBuildingsArgs, userId?: string): Promise<Building[]> {
    const { city, region } = args;
    const where: any = {};
    if (city) where.city = city;
    if (region) where.region = region;
    
    const buildings = await this.buildingRepository.find({ where });

    if (userId) {
      const buildingIds = buildings.map(b => b.id);
      if (buildingIds.length === 0) return buildings;

      const likedBuildingIds = await this.likeRepository.find({
        where: { user_id: userId, building_id: In(buildingIds) },
        select: ['building_id'],
      }).then(likes => likes.map(l => l.building_id));

      buildings.forEach(building => {
        building.is_liked = likedBuildingIds.includes(building.id);
      });
    }

    return buildings;
  }

  async updateBuilding(id: string, input: UpdateBuildingInput): Promise<Building> {
    const building = await this.findOneBuilding(id);
    Object.assign(building, input);
    return this.buildingRepository.save(building);
  }

  async removeBuilding(id: string): Promise<{ id: string; message: string }> {
    const building = await this.findOneBuilding(id);
    await this.buildingRepository.remove(building);

    // Update counters
    await this.companyRepository.decrement({ id: building.company_id }, 'total_buildings', 1);
    await this.developerRepository.decrement({ id: building.developer_id }, 'total_buildings', 1);

    return { id, message: 'Building removed successfully' };
  }

  // Offer Methods
  async createOffer(input: CreateOfferInput): Promise<Offer> {
    const offer = this.offerRepository.create({
      ...input,
      offer_date: new Date(input.offer_date),
    });
    const savedOffer = await this.offerRepository.save(offer);

    // Update counters
    await this.companyRepository.increment({ id: input.company_id }, 'total_offers', 1);
    await this.developerRepository.increment({ id: input.developer_id }, 'total_offers', 1);

    return savedOffer;
  }

  async findOneOffer(id: string): Promise<Offer> {
    const offer = await this.offerRepository.findOne({ where: { id } });
    if (!offer) {
      throw new NotFoundException(`Offer with ID ${id} not found`);
    }
    return offer;
  }

  async findAllOffers(): Promise<Offer[]> {
    return this.offerRepository.find();
  }
  
  async updateOffer(id: string, input: UpdateOfferInput): Promise<Offer> {
    const offer = await this.findOneOffer(id);
    Object.assign(offer, {
      ...input,
      accepted_date: input.accepted_date ? new Date(input.accepted_date) : null,
    });
    return this.offerRepository.save(offer);
  }

  async removeOffer(id: string): Promise<{ id: string; message: string }> {
    const offer = await this.findOneOffer(id);
    await this.offerRepository.remove(offer);

    // Update counters
    await this.companyRepository.decrement({ id: offer.company_id }, 'total_offers', 1);
    await this.developerRepository.decrement({ id: offer.developer_id }, 'total_offers', 1);
    
    return { id, message: 'Offer removed successfully' };
  }

  async checkIfLiked(entityId: string, entityType: LikeableType, userId: string): Promise<boolean> {
    const like = await this.likeRepository.findOne({ 
      where: {
        user_id: userId, 
        [`${entityType.toLowerCase()}_id`]: entityId,
      }
    });
    return !!like;
  }
} 