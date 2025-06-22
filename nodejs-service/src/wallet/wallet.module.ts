import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { WalletService } from './wallet.service';
import { WalletResolver } from './wallet.resolver';
import { Wallet, WalletTransaction, WalletLock } from './wallet.entity';
import { CommonModule } from '../common/common.module';

@Module({
  imports: [
    TypeOrmModule.forFeature([Wallet, WalletTransaction, WalletLock]),
    ConfigModule,
    CommonModule,
  ],
  providers: [WalletService, WalletResolver],
  exports: [WalletService],
})
export class WalletModule {} 