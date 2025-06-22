import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { User } from '../user/user.entity';
import { Property } from './property.entity';
import { Comment } from './comment.entity';
import axios from 'axios';
import { ConfigService } from '@nestjs/config';
import { translateText } from '../common/translation.util';

@Injectable()
export class PropertyService {
  private goServiceUrl: string;

  constructor(
    @InjectRepository(User)
    private usersRepository: Repository<User>,
    @InjectRepository(Property)
    private propertiesRepository: Repository<Property>,
    @InjectRepository(Comment)
    private commentsRepository: Repository<Comment>,
    private configService: ConfigService,
  ) {
    this.goServiceUrl = 'http://go-service:8080/graphql';
  }

  async toggleFavorite(userId: string, propertyId: number): Promise<User> {
    const user = await this.usersRepository.findOne({
      where: { id: userId },
      relations: ['favoriteProperties'],
    });

    let property = await this.propertiesRepository.findOne({ where: { id: propertyId } });
    if (!property) {
      // If property is not in nodejs-service DB, it might exist in go-service
      // We create a placeholder here to establish the relationship
      property = this.propertiesRepository.create({ id: propertyId });
      await this.propertiesRepository.save(property);
    }
    
    const isFavorited = user.favoriteProperties.some((p) => p.id === property.id);

    let goMutation: string;

    if (isFavorited) {
      user.favoriteProperties = user.favoriteProperties.filter((p) => p.id !== property.id);
      goMutation = `
        mutation {
          decrementFavoriteCount(id: "${propertyId}") {
            id
            favoritesCount
          }
        }
      `;
    } else {
      user.favoriteProperties.push(property);
      goMutation = `
        mutation {
          incrementFavoriteCount(id: "${propertyId}") {
            id
            favoritesCount
          }
        }
      `;
    }

    try {
      await axios.post(this.goServiceUrl, { query: goMutation });
    } catch (error) {
      // Handle error, maybe log it or throw a specific exception
      console.error('Error calling go-service:', error.message);
      // Potentially revert the user favorite action if the call fails
      throw new Error('Failed to update favorite count.');
    }

    return this.usersRepository.save(user);
  }

  async createComment(userId: string, propertyId: number, content: string, parentId?: number): Promise<any> {
    const parentIdPart = parentId ? `, parentId: \"${parentId}\"` : '';
    const goMutation = `
      mutation {
        createComment(propertyId: \"${propertyId}\", content: \"${content}\", userId: \"${userId}\"${parentIdPart}) {
          id
          content
          userId
          parentId
          createdAt
        }
      }
    `;
    try {
      const response = await axios.post(this.goServiceUrl, { query: goMutation });
      return response.data.data.createComment;
    } catch (error) {
      console.error('Error calling go-service for createComment:', error.message);
      throw new Error('Failed to create comment.');
    }
  }

  async toggleCommentLike(userId: string, commentId: number): Promise<User> {
    const user = await this.usersRepository.findOne({
      where: { id: userId },
      relations: ['likedComments'],
    });

    let comment = await this.commentsRepository.findOne({ where: { id: commentId } });
    if (!comment) {
      comment = this.commentsRepository.create({ id: commentId });
      await this.commentsRepository.save(comment);
    }

    const isLiked = user.likedComments.some((c) => c.id === comment.id);
    let goMutation: string;

    if (isLiked) {
      user.likedComments = user.likedComments.filter((c) => c.id !== comment.id);
      goMutation = `
        mutation {
          decrementCommentLike(commentId: "${commentId}") {
            id
            likesCount
          }
        }
      `;
    } else {
      user.likedComments.push(comment);
      goMutation = `
        mutation {
          incrementCommentLike(commentId: "${commentId}") {
            id
            likesCount
          }
        }
      `;
    }

    try {
      await axios.post(this.goServiceUrl, { query: goMutation });
    } catch (error) {
      console.error('Error calling go-service for toggleCommentLike:', error.message);
      throw new Error('Failed to update comment like count.');
    }

    return this.usersRepository.save(user);
  }

  async createProperty(input: any): Promise<Property> {
    const langs = ['en', 'ar', 'fr', 'de', 'hi', 'ru', 'fil'];
    const titleTranslations = await translateText(input.title, 'en', langs);
    const descriptionTranslations = await translateText(input.description, 'en', langs);
    const property = this.propertiesRepository.create({
      ...input,
      title: titleTranslations,
      description: descriptionTranslations,
    });
    const saved = await this.propertiesRepository.save(property);
    return Array.isArray(saved) ? saved[0] : saved;
  }
} 