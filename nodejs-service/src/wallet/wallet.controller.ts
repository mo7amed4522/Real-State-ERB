import { Controller, Post, Headers, RawBodyRequest, Req, HttpCode, HttpStatus } from '@nestjs/common';
import { Request } from 'express';
import Stripe from 'stripe';
import { WalletService } from './wallet.service';
import { ConfigService } from '@nestjs/config';

@Controller('wallet')
export class WalletController {
  private stripe: Stripe;

  constructor(
    private walletService: WalletService,
    private configService: ConfigService,
  ) {
    const stripeSecretKey = this.configService.get<string>('STRIPE_SECRET_KEY');
    this.stripe = new Stripe(stripeSecretKey, {
      apiVersion: '2022-11-15',
    });
  }

  @Post('webhook')
  @HttpCode(HttpStatus.OK)
  async handleWebhook(
    @Req() request: RawBodyRequest<Request>,
    @Headers('stripe-signature') signature: string,
  ): Promise<void> {
    const webhookSecret = this.configService.get<string>('STRIPE_WEBHOOK_SECRET');
    
    if (!webhookSecret) {
      throw new Error('STRIPE_WEBHOOK_SECRET is required');
    }

    let event: Stripe.Event;

    try {
      event = this.stripe.webhooks.constructEvent(
        request.rawBody,
        signature,
        webhookSecret,
      );
    } catch (err) {
      throw new Error(`Webhook signature verification failed: ${err.message}`);
    }

    await this.walletService.handleStripeWebhook(event);
  }
} 