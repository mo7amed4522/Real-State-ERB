import { Resolver, Mutation, Args, Query, ID } from '@nestjs/graphql';
import { UseGuards } from '@nestjs/common';
import { GqlAuthGuard } from '../auth/guards/gql-auth.guard';
import { InteractionService } from './interaction.service';
import { Offer } from './offer.entity';
import { CreateOfferInput } from './dto/create-offer.input';
import { UpdateOfferInput } from './dto/update-offer.input';

@Resolver(() => Offer)
export class OfferResolver {
  constructor(private readonly interactionService: InteractionService) {}

  @Mutation(() => Offer)
  @UseGuards(GqlAuthGuard)
  createOffer(@Args('createOfferInput') createOfferInput: CreateOfferInput) {
    return this.interactionService.createOffer(createOfferInput);
  }

  @Query(() => Offer, { name: 'offer' })
  findOneOffer(@Args('id', { type: () => ID }) id: string) {
    return this.interactionService.findOneOffer(id);
  }

  @Query(() => [Offer], { name: 'offers' })
  findAllOffers() {
    return this.interactionService.findAllOffers();
  }

  @Mutation(() => Offer)
  @UseGuards(GqlAuthGuard)
  updateOffer(@Args('updateOfferInput') updateOfferInput: UpdateOfferInput) {
    return this.interactionService.updateOffer(updateOfferInput.id, updateOfferInput);
  }

  @Mutation(() => String)
  @UseGuards(GqlAuthGuard)
  removeOffer(@Args('id', { type: () => ID }) id: string) {
    this.interactionService.removeOffer(id);
    return 'Offer removed successfully';
  }
} 