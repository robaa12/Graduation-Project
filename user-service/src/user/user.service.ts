import { CreatePaymentDto } from './../payment/dto/create-payment.dto';
import { User } from './entities/user.entity';
import {
  BadRequestException,
  Inject,
  Injectable,
  Logger,
  NotFoundException,
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
import { PlansService } from 'src/plans/plans.service';
import { UserPlanPayment } from './entities/user-plan-payment.entity';
import { Plan } from 'src/plans/entities/plan.entity';
@Injectable()
export class UserService {
  constructor(
    @InjectRepository(User) private userRepository: Repository<User>,
    @InjectRepository(UserPlanPayment) private userPlanPaymnetRepository: Repository<UserPlanPayment>,
    private MailerService: EmailService,
    private PlanService:PlansService,
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
      plan_expire_date: user.plan_expire_date,
      plan:user.plan,
      stores: user.stores.map((store) => ({
        id: store.id,
        name: store.store_name,
      })),
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

      const plan = await this.PlanService.findOne(createUserDto.plan_id);
      // Create user with all data at once
      let user = this.userRepository.create({
        ...createUserDto,
        password: hashedPassword,
        otp: otp,
        otpExpiry: otpExpiry,
        plan
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
      relations: ['stores' , 'plan'],
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
      relations: ['stores' , 'plan'],
    });
    if(!user) {
      throw new NotFoundException('User not found');
    }
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

  async createPayment(user:IUser , plan:Plan , response:any) {
    const userPlanPayment = this.userPlanPaymnetRepository.create({
      user_id: user.id,
      plan_id: plan.id,
      amount: plan.price,
      status: response.data.status,
      charge_id: response.data.id,
      currency: response.data.currency,
    });
    return await this.userPlanPaymnetRepository.save(userPlanPayment);

  }

  async updatePayment(tap_id:string , response:any){
    const payment = await this.userPlanPaymnetRepository.findOne({
      where: { charge_id: tap_id },
      relations: ['user', 'plan'],
    });
    const user = await this.findOne(payment.user_id);
    if (!payment) {
      throw new NotFoundException('Payment not found');
    }
    payment.status = response.status;
    user.plan_expire_date = new Date(Date.now() + 30 * 24 * 60 * 60 * 1000); // Extend plan by 30 days
    user.plan = payment.plan;
    await this.userRepository.save(user);
    await this.userPlanPaymnetRepository.save(payment);
    return payment;
  }
}
