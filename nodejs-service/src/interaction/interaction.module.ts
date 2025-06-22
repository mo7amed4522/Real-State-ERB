import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Building } from './building.entity';
import { Offer } from './offer.entity';
import { Like } from './like.entity';
import { Comment } from './comment.entity';
import { Reaction } from './reaction.entity';
import { Company } from '../company/company.entity';
import { Developer } from '../company/developer.entity';
import { InteractionService } from './interaction.service';
import { LikeResolver } from './like.resolver';
import { CommentResolver } from './comment.resolver';
import { ReactionResolver } from './reaction.resolver';
import { BuildingResolver } from './building.resolver';
import { OfferResolver } from './offer.resolver';
// Import services and resolvers later

@Module({
  imports: [
    TypeOrmModule.forFeature([
      Building,
      Offer,
      Like,
      Comment,
      Reaction,
      Company,
      Developer,
    ]),
  ],
  providers: [
    InteractionService,
    LikeResolver,
    CommentResolver,
    ReactionResolver,
    BuildingResolver,
    OfferResolver,
  ],
  // Add providers and exports later
})
export class InteractionModule {} 