import { Injectable, NotFoundException } from '@nestjs/common';
import { CreateStoreDto } from './dto/create-store.dto';
import { UpdateStoreDto } from './dto/update-store.dto';
import { InjectRepository } from '@nestjs/typeorm';
import { EmailService } from 'src/shared/services/email/email.service';
import { Repository } from 'typeorm';
import { Store } from './entities/store.entity';
import { UserService } from 'src/user/user.service';

@Injectable()
export class StoreService {
  constructor(
    @InjectRepository(Store)  private storeRepository: Repository<Store>,
    private MailerService:EmailService,
    private readonly UserService: UserService
  ) {}

  async createStore (createStoreDto: CreateStoreDto):Promise<Store> {
    const user = await this.UserService.findOne(createStoreDto.user_id);
    if(!user){
      throw new NotFoundException('User not found');
    }
    const store = this.storeRepository.create({
      ...createStoreDto,
      user
    })
    return await this.storeRepository.save(store);
  }

  async findAll():Promise<Store[]> {
    return await this.storeRepository.find();
  }

  async findOne(id: number):Promise<Store> {
    return await this.storeRepository.findOne(
      {where:{id} }
    );
  }
}
