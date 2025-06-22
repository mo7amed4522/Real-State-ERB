import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { User } from './user.entity';
import * as bcrypt from 'bcrypt';

@Injectable()
export class UserService {
  constructor(
    @InjectRepository(User)
    private readonly userRepository: Repository<User>,
  ) {}

  async findById(id: string): Promise<User | undefined> {
    return this.userRepository.findOne({ where: { id } });
  }

  async findByEmail(email: string): Promise<User | undefined> {
    return this.userRepository.findOne({ where: { email } });
  }

  async findByResetToken(token: string): Promise<User | undefined> {
    return this.userRepository.findOne({ where: { resetPasswordToken: token } });
  }

  async findAll(): Promise<User[]> {
    return this.userRepository.find();
  }

  async create(details: Partial<User>): Promise<User> {
    if (details.password) {
      const salt = await bcrypt.genSalt();
      details.password = await bcrypt.hash(details.password, salt);
    }
    const newUser = this.userRepository.create(details);
    return this.userRepository.save(newUser);
  }

  async validatePassword(email: string, pass: string): Promise<User | null> {
    const user = await this.userRepository.findOne({ where: { email } });
    if (user && await bcrypt.compare(pass, user.password)) {
      const { password, ...result } = user;
      return result as User;
    }
    return null;
  }

  async update(id: string, details: Partial<User>): Promise<User> {
    await this.userRepository.update(id, details);
    return this.findById(id);
  }

  async delete(id: string): Promise<boolean> {
    const result = await this.userRepository.delete(id);
    return result.affected > 0;
  }

  async changePassword(userId: string, oldP: string, newP: string): Promise<boolean> {
    const user = await this.userRepository.findOne({ where: { id: userId }});
    if (!user) {
      throw new Error('User not found');
    }
    const isValid = await bcrypt.compare(oldP, user.password);
    if (!isValid) {
      throw new Error('Invalid old password');
    }
    const salt = await bcrypt.genSalt();
    user.password = await bcrypt.hash(newP, salt);
    await this.userRepository.save(user);
    return true;
  }
} 