import { Resolver, Query, Mutation, Args, ResolveField, Parent } from '@nestjs/graphql';
import { UseGuards } from '@nestjs/common';
import { GqlAuthGuard } from '../auth/guards/gql-auth.guard';
import { CurrentUser } from '../auth/decorators/current-user.decorator';
import { User } from '../user/user.entity';
import { WalletService } from './wallet.service';
import { CreateWalletInput } from './dto/create-wallet.input';
import { DepositInput } from './dto/deposit.input';
import { WithdrawInput } from './dto/withdraw.input';
import { TransferInput } from './dto/transfer.input';
import { Wallet, WalletTransaction, TransactionType, TransactionStatus, PaymentMethod } from './wallet.entity';

@Resolver(() => Wallet)
@UseGuards(GqlAuthGuard)
export class WalletResolver {
  constructor(private walletService: WalletService) {}

  @Mutation(() => Wallet)
  async createWallet(
    @CurrentUser() user: User,
    @Args('input') input: CreateWalletInput
  ): Promise<Wallet> {
    return this.walletService.createWallet(user.id, input);
  }

  @Query(() => Wallet)
  async getWallet(
    @CurrentUser() user: User,
    @Args('walletId') walletId: string
  ): Promise<Wallet> {
    return this.walletService.getWallet(user.id, walletId);
  }

  @Query(() => Wallet)
  async getUserWallet(@CurrentUser() user: User): Promise<Wallet> {
    return this.walletService.getUserWallet(user.id);
  }

  @Query(() => [WalletTransaction])
  async getTransactionHistory(
    @CurrentUser() user: User,
    @Args('walletId') walletId: string,
    @Args('limit', { defaultValue: 50 }) limit: number,
    @Args('offset', { defaultValue: 0 }) offset: number
  ): Promise<WalletTransaction[]> {
    return this.walletService.getTransactionHistory(user.id, walletId, limit, offset);
  }

  @Query(() => WalletTransaction)
  async getTransaction(
    @CurrentUser() user: User,
    @Args('transactionId') transactionId: string
  ): Promise<WalletTransaction> {
    return this.walletService.getTransaction(transactionId, user.id);
  }

  @Query(() => Object)
  async getWalletBalance(
    @CurrentUser() user: User,
    @Args('walletId') walletId: string
  ): Promise<{ balance: number; frozenBalance: number }> {
    return this.walletService.getWalletBalance(user.id, walletId);
  }

  @Mutation(() => WalletTransaction)
  async deposit(
    @CurrentUser() user: User,
    @Args('walletId') walletId: string,
    @Args('input') input: DepositInput
  ): Promise<WalletTransaction> {
    return this.walletService.deposit(user.id, walletId, input);
  }

  @Mutation(() => WalletTransaction)
  async confirmDeposit(
    @CurrentUser() user: User,
    @Args('transactionId') transactionId: string,
    @Args('paymentIntentId') paymentIntentId: string
  ): Promise<WalletTransaction> {
    return this.walletService.confirmDeposit(transactionId, paymentIntentId);
  }

  @Mutation(() => WalletTransaction)
  async withdraw(
    @CurrentUser() user: User,
    @Args('walletId') walletId: string,
    @Args('input') input: WithdrawInput
  ): Promise<WalletTransaction> {
    return this.walletService.withdraw(user.id, walletId, input);
  }

  @Mutation(() => WalletTransaction)
  async transfer(
    @CurrentUser() user: User,
    @Args('fromWalletId') fromWalletId: string,
    @Args('input') input: TransferInput
  ): Promise<WalletTransaction> {
    return this.walletService.transfer(user.id, fromWalletId, input);
  }

  @Mutation(() => WalletTransaction)
  async cancelTransaction(
    @CurrentUser() user: User,
    @Args('transactionId') transactionId: string
  ): Promise<WalletTransaction> {
    return this.walletService.cancelTransaction(transactionId, user.id);
  }

  @Mutation(() => Wallet)
  async updateBankDetails(
    @CurrentUser() user: User,
    @Args('walletId') walletId: string,
    @Args('encryptedBankDetails') encryptedBankDetails: string
  ): Promise<Wallet> {
    return this.walletService.updateBankDetails(user.id, walletId, encryptedBankDetails);
  }

  @ResolveField(() => [WalletTransaction])
  async transactions(@Parent() wallet: Wallet): Promise<WalletTransaction[]> {
    return this.walletService.getTransactionHistory(wallet.userId, wallet.id, 10, 0);
  }
}

@Resolver(() => WalletTransaction)
export class WalletTransactionResolver {
  @Query(() => [WalletTransaction])
  async getAllTransactions(
    @Args('limit', { defaultValue: 50 }) limit: number,
    @Args('offset', { defaultValue: 0 }) offset: number
  ): Promise<WalletTransaction[]> {
    // This would typically be admin-only
    return [];
  }
} 