import { Module } from '@nestjs/common';
import { GraphQLModule } from '@nestjs/graphql';
import { ApolloDriver, ApolloDriverConfig } from '@nestjs/apollo';
import { TypeOrmModule } from '@nestjs/typeorm';
import { join } from 'path';
import { UserModule } from './user/user.module';
import { AuthModule } from './auth/auth.module';
import { CommonModule } from './common/common.module';
import { PropertyModule } from './property/property.module';
import { CompanyModule } from './company/company.module';
import { FileModule } from './file/file.module';
import { InteractionModule } from './interaction/interaction.module';

@Module({
  imports: [
    GraphQLModule.forRoot<ApolloDriverConfig>({
      driver: ApolloDriver,
      autoSchemaFile: join(process.cwd(), 'src/schema.gql'),
      playground: true,
    }),
    TypeOrmModule.forRoot({
      type: 'postgres',
      url: process.env.DATABASE_URL,
      autoLoadEntities: true,
      synchronize: true, // Be cautious with this in production
    }),
    UserModule,
    AuthModule,
    CommonModule,
    PropertyModule,
    CompanyModule,
    FileModule,
    InteractionModule,
  ],
  controllers: [],
  providers: [],
  exports: [],
})
export class AppModule {} 