import { Entity, PrimaryGeneratedColumn, Column, CreateDateColumn, ManyToOne, JoinColumn, OneToMany } from 'typeorm';
import { User } from '../user/user.entity';

export enum CommentableType {
  COMPANY = 'company',
  DEVELOPER = 'developer',
}

@Entity('social_comments') // Renamed to avoid conflicts if 'comments' table exists
export class Comment {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ type: 'uuid', nullable: true })
  parent_id: string;

  @ManyToOne(() => Comment, comment => comment.replies, { nullable: true })
  @JoinColumn({ name: 'parent_id' })
  parent: Comment;
  
  @OneToMany(() => Comment, comment => comment.parent)
  replies: Comment[];

  @Column({ type: 'enum', enum: CommentableType })
  target_type: CommentableType;
  
  @Column({ type: 'uuid' })
  target_id: string;

  @Column({ type: 'uuid' })
  user_id: string;

  @ManyToOne(() => User)
  @JoinColumn({ name: 'user_id' })
  user: User;

  @Column('text')
  content: string;

  @Column({ type: 'int', default: 0 })
  total_reactions: number;

  @CreateDateColumn()
  created_at: Date;
} 