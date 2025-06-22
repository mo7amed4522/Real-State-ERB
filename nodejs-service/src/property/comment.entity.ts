import { Entity, PrimaryGeneratedColumn, ManyToMany } from 'typeorm';
import { User } from '../user/user.entity';

@Entity('comments')
export class Comment {
  @PrimaryGeneratedColumn()
  id: number;

  @ManyToMany(() => User, (user) => user.likedComments)
  likedBy: User[];
} 