import { Injectable, Logger, BadRequestException, ForbiddenException, InternalServerErrorException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository, DataSource } from 'typeorm';
import { ConfigService } from '@nestjs/config';
import Stripe from 'stripe';
import Redlock from 'redlock';
import Redis from 'ioredis';
import { Wallet, WalletTransaction, TransactionType, TransactionStatus, PaymentMethod } from './wallet.entity';
import { CreateWalletInput } from './dto/create-wallet.input';
import { DepositInput } from './dto/deposit.input';
import { WithdrawInput } from './dto/withdraw.input';
import { TransferInput } from './dto/transfer.input';
import { EncryptionService } from '../common/encryption.service';

@Injectable()
export class WalletService {
  private readonly logger = new Logger(WalletService.name);
  private stripe: Stripe;
  private redlock: Redlock;
  private redis: Redis;

  constructor(
    @InjectRepository(Wallet)
    private walletRepository: Repository<Wallet>,
    @InjectRepository(WalletTransaction)
    private transactionRepository: Repository<WalletTransaction>,
    private dataSource: DataSource,
    private configService: ConfigService,
    private encryptionService: EncryptionService,
  ) {
    this.initializeStripe();
    this.initializeRedis();
  }

  private async initializeStripe() {
    const stripeSecretKey = this.configService.get<string>('STRIPE_SECRET_KEY');
    if (!stripeSecretKey) {
      throw new Error('STRIPE_SECRET_KEY is required');
    }
    this.stripe = new Stripe(stripeSecretKey, {
      apiVersion: '2022-11-15',
    });
  }

  private async initializeRedis() {
    const redisUrl = this.configService.get<string>('REDIS_URL') || 'redis://localhost:6379';
    this.redis = new Redis(redisUrl);
    
    this.redlock = new Redlock([this.redis], {
      driftFactor: 0.01,
      retryCount: 10,
      retryDelay: 200,
      retryJitter: 200,
      automaticExtensionThreshold: 500,
    });

    this.redlock.on('error', (error) => {
      this.logger.error('RedLock error:', error);
    });
  }

  async createWallet(userId: string, input: CreateWalletInput): Promise<Wallet> {
    const existingWallet = await this.walletRepository.findOne({
      where: { userId, isActive: true }
    });

    if (existingWallet) {
      throw new BadRequestException('User already has an active wallet');
    }

    let stripeCustomerId = input.stripeCustomerId;
    if (!stripeCustomerId) {
      const customer = await this.stripe.customers.create({
        metadata: { userId }
      });
      stripeCustomerId = customer.id;
    }

    const wallet = this.walletRepository.create({
      userId,
      balance: 0,
      frozenBalance: 0,
      stripeCustomerId,
      encryptedBankDetails: input.encryptedBankDetails ? 
        this.encryptionService.encrypt(input.encryptedBankDetails) : null,
      isActive: true
    });

    return this.walletRepository.save(wallet);
  }

  async getWallet(userId: string, walletId: string): Promise<Wallet> {
    const wallet = await this.walletRepository.findOne({
      where: { id: walletId, userId, isActive: true },
      relations: ['transactions']
    });

    if (!wallet) {
      throw new BadRequestException('Wallet not found');
    }

    return wallet;
  }

  async getUserWallet(userId: string): Promise<Wallet> {
    const wallet = await this.walletRepository.findOne({
      where: { userId, isActive: true },
      relations: ['transactions']
    });

    if (!wallet) {
      throw new BadRequestException('Wallet not found');
    }

    return wallet;
  }

  async deposit(userId: string, walletId: string, input: DepositInput): Promise<WalletTransaction> {
    const lockKey = `wallet:${walletId}:lock`;
    let lock;

    try {
      // Acquire distributed lock
      lock = await this.redlock.acquire([lockKey], 30000); // 30 seconds lock

      const wallet = await this.getWallet(userId, walletId);
      
      // Create payment intent with Stripe
      const paymentIntent = await this.stripe.paymentIntents.create({
        amount: Math.round(input.amount * 100), // Convert to cents
        currency: 'usd',
        customer: wallet.stripeCustomerId,
        metadata: {
          userId,
          walletId,
          transactionType: TransactionType.DEPOSIT,
          description: input.description || 'Wallet deposit'
        }
      });

      // Create transaction record
      const transaction = this.transactionRepository.create({
        walletId,
        userId,
        transactionType: TransactionType.DEPOSIT,
        status: TransactionStatus.PENDING,
        amount: input.amount,
        stripePaymentIntentId: paymentIntent.id,
        paymentMethod: input.paymentMethod,
        description: input.description,
        encryptedBankId: input.encryptedBankDetails ? 
          this.encryptionService.encrypt(input.encryptedBankDetails) : null,
        encryptedAccountName: input.encryptedAccountName ? 
          this.encryptionService.encrypt(input.encryptedAccountName) : null,
        metadata: {
          stripePaymentIntentId: paymentIntent.id,
          clientSecret: paymentIntent.client_secret
        }
      });

      const savedTransaction = await this.transactionRepository.save(transaction);

      this.logger.log(`Deposit initiated: ${savedTransaction.id} for wallet ${walletId}`);

      return savedTransaction;

    } catch (error) {
      this.logger.error('Deposit error:', error);
      throw new InternalServerErrorException('Failed to process deposit');
    } finally {
      if (lock) {
        await lock.release();
      }
    }
  }

  async confirmDeposit(transactionId: string, paymentIntentId: string): Promise<WalletTransaction> {
    const lockKey = `transaction:${transactionId}:lock`;
    let lock;

    try {
      lock = await this.redlock.acquire([lockKey], 30000);

      const transaction = await this.transactionRepository.findOne({
        where: { id: transactionId, stripePaymentIntentId: paymentIntentId }
      });

      if (!transaction) {
        throw new BadRequestException('Transaction not found');
      }

      if (transaction.status !== TransactionStatus.PENDING) {
        throw new BadRequestException('Transaction already processed');
      }

      // Verify payment intent with Stripe
      const paymentIntent = await this.stripe.paymentIntents.retrieve(paymentIntentId);
      
      if (paymentIntent.status !== 'succeeded') {
        throw new BadRequestException('Payment not completed');
      }

      // Use database transaction for atomicity
      const queryRunner = this.dataSource.createQueryRunner();
      await queryRunner.connect();
      await queryRunner.startTransaction();

      try {
        // Update transaction status
        transaction.status = TransactionStatus.COMPLETED;
        await queryRunner.manager.save(WalletTransaction, transaction);

        // Update wallet balance
        const wallet = await queryRunner.manager.findOne(Wallet, {
          where: { id: transaction.walletId }
        });

        wallet.balance += transaction.amount;
        await queryRunner.manager.save(Wallet, wallet);

        await queryRunner.commitTransaction();

        this.logger.log(`Deposit confirmed: ${transactionId} for wallet ${transaction.walletId}`);

        return transaction;

      } catch (error) {
        await queryRunner.rollbackTransaction();
        throw error;
      } finally {
        await queryRunner.release();
      }

    } catch (error) {
      this.logger.error('Confirm deposit error:', error);
      throw new InternalServerErrorException('Failed to confirm deposit');
    } finally {
      if (lock) {
        await lock.release();
      }
    }
  }

  async withdraw(userId: string, walletId: string, input: WithdrawInput): Promise<WalletTransaction> {
    const lockKey = `wallet:${walletId}:lock`;
    let lock;

    try {
      lock = await this.redlock.acquire([lockKey], 30000);

      const wallet = await this.getWallet(userId, walletId);

      if (wallet.balance < input.amount) {
        throw new BadRequestException('Insufficient balance');
      }

      // Create transfer with Stripe
      const transfer = await this.stripe.transfers.create({
        amount: Math.round(input.amount * 100),
        currency: 'usd',
        destination: wallet.stripeCustomerId,
        metadata: {
          userId,
          walletId,
          transactionType: TransactionType.WITHDRAWAL,
          description: input.description || 'Wallet withdrawal'
        }
      });

      // Use database transaction
      const queryRunner = this.dataSource.createQueryRunner();
      await queryRunner.connect();
      await queryRunner.startTransaction();

      try {
        // Create transaction record
        const transaction = this.transactionRepository.create({
          walletId,
          userId,
          transactionType: TransactionType.WITHDRAWAL,
          status: TransactionStatus.PROCESSING,
          amount: input.amount,
          stripeTransferId: transfer.id,
          paymentMethod: input.paymentMethod,
          description: input.description,
          encryptedBankId: input.encryptedBankDetails ? 
            this.encryptionService.encrypt(input.encryptedBankDetails) : null,
          encryptedAccountName: input.encryptedAccountName ? 
            this.encryptionService.encrypt(input.encryptedAccountName) : null,
          metadata: { stripeTransferId: transfer.id }
        });

        const savedTransaction = await queryRunner.manager.save(WalletTransaction, transaction);

        // Update wallet balance
        wallet.balance -= input.amount;
        await queryRunner.manager.save(Wallet, wallet);

        await queryRunner.commitTransaction();

        this.logger.log(`Withdrawal initiated: ${savedTransaction.id} for wallet ${walletId}`);

        return savedTransaction;

      } catch (error) {
        await queryRunner.rollbackTransaction();
        throw error;
      } finally {
        await queryRunner.release();
      }

    } catch (error) {
      this.logger.error('Withdrawal error:', error);
      throw new InternalServerErrorException('Failed to process withdrawal');
    } finally {
      if (lock) {
        await lock.release();
      }
    }
  }

  async transfer(userId: string, fromWalletId: string, input: TransferInput): Promise<WalletTransaction> {
    const lockKey = `wallet:${fromWalletId}:lock`;
    let lock;

    try {
      lock = await this.redlock.acquire([lockKey], 30000);

      const fromWallet = await this.getWallet(userId, fromWalletId);
      const toWallet = await this.walletRepository.findOne({
        where: { id: input.receiverWalletId, isActive: true }
      });

      if (!toWallet) {
        throw new BadRequestException('Receiver wallet not found');
      }

      if (fromWallet.balance < input.amount) {
        throw new BadRequestException('Insufficient balance');
      }

      if (fromWalletId === input.receiverWalletId) {
        throw new BadRequestException('Cannot transfer to same wallet');
      }

      // Use database transaction for atomic transfer
      const queryRunner = this.dataSource.createQueryRunner();
      await queryRunner.connect();
      await queryRunner.startTransaction();

      try {
        // Create transaction record
        const transaction = this.transactionRepository.create({
          walletId: fromWalletId,
          userId,
          transactionType: TransactionType.TRANSFER,
          status: TransactionStatus.COMPLETED,
          amount: input.amount,
          description: input.description,
          encryptedReceiverDetails: input.encryptedReceiverDetails ? 
            this.encryptionService.encrypt(input.encryptedReceiverDetails) : null,
          metadata: {
            receiverWalletId: input.receiverWalletId,
            receiverUserId: toWallet.userId
          }
        });

        const savedTransaction = await queryRunner.manager.save(WalletTransaction, transaction);

        // Update sender wallet balance
        fromWallet.balance -= input.amount;
        await queryRunner.manager.save(Wallet, fromWallet);

        // Update receiver wallet balance
        toWallet.balance += input.amount;
        await queryRunner.manager.save(Wallet, toWallet);

        await queryRunner.commitTransaction();

        this.logger.log(`Transfer completed: ${savedTransaction.id} from ${fromWalletId} to ${input.receiverWalletId}`);

        return savedTransaction;

      } catch (error) {
        await queryRunner.rollbackTransaction();
        throw error;
      } finally {
        await queryRunner.release();
      }

    } catch (error) {
      this.logger.error('Transfer error:', error);
      throw new InternalServerErrorException('Failed to process transfer');
    } finally {
      if (lock) {
        await lock.release();
      }
    }
  }

  async getTransactionHistory(userId: string, walletId: string, limit = 50, offset = 0): Promise<WalletTransaction[]> {
    const wallet = await this.getWallet(userId, walletId);

    return this.transactionRepository.find({
      where: { walletId },
      order: { createdAt: 'DESC' },
      take: limit,
      skip: offset
    });
  }

  async getTransaction(transactionId: string, userId: string): Promise<WalletTransaction> {
    const transaction = await this.transactionRepository.findOne({
      where: { id: transactionId, userId }
    });

    if (!transaction) {
      throw new BadRequestException('Transaction not found');
    }

    return transaction;
  }

  async cancelTransaction(transactionId: string, userId: string): Promise<WalletTransaction> {
    const lockKey = `transaction:${transactionId}:lock`;
    let lock;

    try {
      lock = await this.redlock.acquire([lockKey], 30000);

      const transaction = await this.transactionRepository.findOne({
        where: { id: transactionId, userId }
      });

      if (!transaction) {
        throw new BadRequestException('Transaction not found');
      }

      if (transaction.status !== TransactionStatus.PENDING) {
        throw new BadRequestException('Transaction cannot be cancelled');
      }

      transaction.status = TransactionStatus.CANCELLED;
      return this.transactionRepository.save(transaction);

    } catch (error) {
      this.logger.error('Cancel transaction error:', error);
      throw new InternalServerErrorException('Failed to cancel transaction');
    } finally {
      if (lock) {
        await lock.release();
      }
    }
  }

  async updateBankDetails(userId: string, walletId: string, encryptedBankDetails: string): Promise<Wallet> {
    const wallet = await this.getWallet(userId, walletId);
    
    wallet.encryptedBankDetails = this.encryptionService.encrypt(encryptedBankDetails);
    return this.walletRepository.save(wallet);
  }

  async getWalletBalance(userId: string, walletId: string): Promise<{ balance: number; frozenBalance: number }> {
    const wallet = await this.getWallet(userId, walletId);
    return {
      balance: wallet.balance,
      frozenBalance: wallet.frozenBalance
    };
  }

  // Webhook handler for Stripe events
  async handleStripeWebhook(event: Stripe.Event): Promise<void> {
    switch (event.type) {
      case 'payment_intent.succeeded':
        await this.handlePaymentIntentSucceeded(event.data.object as Stripe.PaymentIntent);
        break;
      case 'payment_intent.payment_failed':
        await this.handlePaymentIntentFailed(event.data.object as Stripe.PaymentIntent);
        break;
      case 'transfer.created':
        await this.handleTransferCreated(event.data.object as Stripe.Transfer);
        break;
      case 'transfer.failed':
        await this.handleTransferFailed(event.data.object as Stripe.Transfer);
        break;
    }
  }

  private async handlePaymentIntentSucceeded(paymentIntent: Stripe.PaymentIntent): Promise<void> {
    const transaction = await this.transactionRepository.findOne({
      where: { stripePaymentIntentId: paymentIntent.id }
    });

    if (transaction) {
      await this.confirmDeposit(transaction.id, paymentIntent.id);
    }
  }

  private async handlePaymentIntentFailed(paymentIntent: Stripe.PaymentIntent): Promise<void> {
    const transaction = await this.transactionRepository.findOne({
      where: { stripePaymentIntentId: paymentIntent.id }
    });

    if (transaction) {
      transaction.status = TransactionStatus.FAILED;
      transaction.failureReason = 'Payment failed';
      await this.transactionRepository.save(transaction);
    }
  }

  private async handleTransferCreated(transfer: Stripe.Transfer): Promise<void> {
    const transaction = await this.transactionRepository.findOne({
      where: { stripeTransferId: transfer.id }
    });

    if (transaction) {
      transaction.status = TransactionStatus.COMPLETED;
      await this.transactionRepository.save(transaction);
    }
  }

  private async handleTransferFailed(transfer: Stripe.Transfer): Promise<void> {
    const transaction = await this.transactionRepository.findOne({
      where: { stripeTransferId: transfer.id }
    });

    if (transaction) {
      transaction.status = TransactionStatus.FAILED;
      transaction.failureReason = 'Transfer failed';
      await this.transactionRepository.save(transaction);
    }
  }
} 