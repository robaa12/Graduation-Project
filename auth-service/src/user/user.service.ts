import { User } from './entities/user.entity';
import {
  BadRequestException,
  Inject,
  Injectable,
  Logger,
  UnauthorizedException,
} from '@nestjs/common';
import { CreateUserDto } from './dto/create-user.dto';
import { UpdateUserDto } from './dto/update-user.dto';
import { MoreThan, QueryFailedError, Repository } from 'typeorm';
import { InjectRepository } from '@nestjs/typeorm';
import * as bcrypt from 'bcryptjs';
import { JwtService } from '@nestjs/jwt';
import { EmailService } from 'src/shared/services/email/email.service';
import { IUser } from 'src/shared/interfaces/uesr.interface';
import { DuplicatedValueException } from 'src/shared/exception-filters/duplicate-value-exception.filter';
@Injectable()
export class UserService {
  constructor(
    @InjectRepository(User) private userRepository: Repository<User>,
    private MailerService: EmailService,
    private jwtService: JwtService,
  ) {}

  private castToUser(user: User): IUser {
    return {
      id: user.id,
      firstName: user.firstName,
      lastName: user.lastName,
      isActive: user.isActive,
      email: user.email,
      is_banned: user.is_banned,
      phoneNumber: user.phoneNumber,
      stores_id: user.stores.map((store) => store.id),
      address: user.address,
      createAt: user.createAt,
      updateAt: user.updateAt,
    };
  }
  async create(createUserDto: CreateUserDto): Promise<IUser> {
    try {
      // Hash the password
      const hashedPassword = await bcrypt.hash(createUserDto.password, 10);

      // Generate OTP
      const otp = Math.floor(100000 + Math.random() * 900000).toString();
      const otpExpiry = new Date(Date.now() + 10 * 60 * 1000);

      // Create user with all data at once
      let user = this.userRepository.create({
        ...createUserDto,
        password: hashedPassword,
        otp: otp,
        otpExpiry: otpExpiry,
      });

      // Save user
      user = await this.userRepository.save(user);

      // Send verification email asynchronously without waiting
      this.MailerService.sendVerficationMail(otp, user.email).catch((error) => {
        Logger.error('Failed to send verification email', error);
      });

      // Transform and return user immediately
      return this.castToUser({ ...user, stores: [] });
    } catch (e) {
      if (e instanceof QueryFailedError) {
        throw new DuplicatedValueException('Email is already registered');
      }
      throw e;
    }
  }

  async verifyEmail(email: string, otp: string): Promise<IUser> {
    const user = await this.userRepository.findOne({
      where: { email, otp, otpExpiry: MoreThan(new Date(Date.now())) },
    });
    if (!user) {
      throw new BadRequestException('Invalid OTP');
    }
    user.isActive = true;
    user.otp = null;
    user.otpExpiry = null;
    return this.castToUser(await this.userRepository.save(user));
  }
  async login(data: any) {
    let user = await this.userRepository.findOne({
      where: { email: data.email },
      relations: ['stores'],
    });
    if (!user) {
      throw new UnauthorizedException('User not found');
    }
    let isPasswordMatch = await bcrypt.compare(data.password, user.password);
    if (!isPasswordMatch) {
      throw new BadRequestException('Invalid Credentials');
    }
    return this.castToUser(user);
  }
  async findAll() {
    return await this.userRepository.find();
  }

  async findOne(id: number) {
    const user = await this.userRepository.findOne({
      where: { id },
      relations: ['stores'],
    });
    console.log(user);
    return this.castToUser(user);
  }

  async update(id: number, updateUserDto: UpdateUserDto): Promise<IUser> {
    let user = await this.userRepository.update(id, updateUserDto);
    let updatedUser = await this.userRepository.findOneBy({ id });
    return this.castToUser(updatedUser);
  }

  async remove(id: number) {
    return await this.userRepository.delete(id);
  }
}
