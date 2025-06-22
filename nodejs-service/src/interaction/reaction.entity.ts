import { Entity, PrimaryGeneratedColumn, Column, CreateDateColumn, ManyToOne, JoinColumn } from 'typeorm';
import { User } from '../user/user.entity';
import { Comment } from './comment.entity';

export enum Emoji {
  THUMBS_UP = '👍',
  HEART = '❤️',
  LAUGH = '😂',
  WOW = '😮',
  SAD = '😢',
  ANGRY = '😡',
}

@Entity('reactions')
export class Reaction {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ type: 'uuid' })
  comment_id: string;

  @ManyToOne(() => Comment)
  @JoinColumn({ name: 'comment_id' })
  comment: Comment;

  @Column({ type: 'uuid' })
  user_id: string;
  
  @ManyToOne(() => User)
  @JoinColumn({ name: 'user_id' })
  user: User;

  @Column({ type: 'enum', enum: Emoji })
  emoji: Emoji;

  @CreateDateColumn()
  created_at: Date;
} 