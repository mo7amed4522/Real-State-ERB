import { Entity, PrimaryGeneratedColumn, Column, ManyToMany } from 'typeorm';
import { User } from '../user/user.entity';

@Entity('properties')
export class Property {
  @PrimaryGeneratedColumn()
  id: number;

  @Column('jsonb', { nullable: false, default: '{}' })
  title: Record<string, string>;

  @Column('jsonb', { nullable: false, default: '{}' })
  description: Record<string, string>;

  // Add other fields as needed for completeness
  // ...

  @ManyToMany(() => User, (user) => user.favoriteProperties)
  favoritedBy: User[];
} 