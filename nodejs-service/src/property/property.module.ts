import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Property } from './property.entity';
import { PropertyService } from './property.service';
import { PropertyResolver } from './property.resolver';
import { User } from 'src/user/user.entity';
import { ConfigModule } from '@nestjs/config';
import { Comment } from './comment.entity';

@Module({
  imports: [TypeOrmModule.forFeature([Property, User, Comment]), ConfigModule],
  providers: [PropertyService, PropertyResolver],
  exports: [TypeOrmModule, PropertyService],
})
export class PropertyModule {} 