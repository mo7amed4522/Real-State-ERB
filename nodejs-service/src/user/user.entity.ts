import { Entity, PrimaryGeneratedColumn, Column, CreateDateColumn, UpdateDateColumn, ManyToMany, JoinTable } from 'typeorm';
import { IsEmail, IsEnum, IsOptional, IsUrl, Length } from 'class-validator';
import { ObjectType, Field, ID, registerEnumType } from '@nestjs/graphql';
import { Property } from '../property/property.entity';
import { Comment } from '../property/comment.entity';

export enum Role {
  USER = 'user',
  ADMIN = 'admin',
}

registerEnumType(Role, {
  name: 'Role',
  description: 'User roles',
});

@ObjectType()
@Entity()
export class User {
  @Field(() => ID)
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Field()
  @Column({ unique: true })
  email: string;

  @Field({ nullable: true })
  @Column({ nullable: true })
  firstName?: string;

  @Field({ nullable: true })
  @Column({ nullable: true })
  lastName?: string;

  @Field(() => String, { nullable: true })
  @Column({ nullable: true })
  @IsOptional()
  @IsUrl()
  photoUrl?: string;

  @Field({ nullable: true })
  @Column({ nullable: true })
  phoneNumber?: string;

  // This field is for the database only. Not exposed via GraphQL.
  @Column({ nullable: true })
  password?: string;

  @Field({ nullable: true })
  @Column({ nullable: true })
  gender?: string;

  @Field({ nullable: true })
  @Column({ nullable: true })
  googleId?: string;

  // We can add other OAuth provider IDs here as well
  // e.g., @Column({ nullable: true }) facebookId?: string;

  @Field(() => [Role])
  @Column({
    type: 'simple-array',
    enum: Role,
    default: [Role.USER],
  })
  roles: Role[];

  @Column({ nullable: true })
  resetPasswordToken?: string;

  @Column({ type: 'timestamp', nullable: true })
  resetPasswordExpires?: Date;

  @ManyToMany(() => Property, { cascade: true })
  @JoinTable({
    name: 'user_favorites_property',
    joinColumn: { name: 'userId', referencedColumnName: 'id' },
    inverseJoinColumn: { name: 'propertyId', referencedColumnName: 'id' },
  })
  favoriteProperties: Property[];

  @ManyToMany(() => Comment, { cascade: true })
  @JoinTable({
    name: 'user_likes_comment',
    joinColumn: { name: 'userId', referencedColumnName: 'id' },
    inverseJoinColumn: { name: 'commentId', referencedColumnName: 'id' },
  })
  likedComments: Comment[];

  @Field(() => Date)
  @CreateDateColumn()
  createdAt: Date;
} 