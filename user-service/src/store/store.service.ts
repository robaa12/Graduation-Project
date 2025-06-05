import {  CreateStoreThemeDto } from './dto/create-store-theme.dto';
import { Injectable, NotFoundException } from '@nestjs/common';
import { CreateStoreDto } from './dto/create-store.dto';
import { UpdateStoreDto } from './dto/update-store.dto';
import { InjectRepository } from '@nestjs/typeorm';
import { EmailService } from 'src/shared/services/email/email.service';
import { Repository } from 'typeorm';
import { Store } from './entities/store.entity';
import { UserService } from 'src/user/user.service';
import { InjectModel } from '@nestjs/mongoose';
import { StoreThemeSchema } from './entities/store-theme.entity';
import { Model } from 'mongoose';

@Injectable()
export class StoreService {
  constructor(
    @InjectRepository(Store)  private storeRepository: Repository<Store>,
    private MailerService:EmailService,
    private readonly UserService: UserService,
    @InjectModel('StoreTheme') private storeThemeModel: Model<StoreThemeSchema>
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
    const store =  await this.storeRepository.findOne(
      {where:{id} }
    );
    if(!store){
      throw new NotFoundException('Store not found');
    }
    return store;
  }

  async createStoreTheme(CreateStoreThemeDto:CreateStoreThemeDto){
    const store = await this.storeRepository.findOne({where:{id:CreateStoreThemeDto.storeId}});
    if(!store){
      throw new NotFoundException('Store not found');
    }
    const storeTheme = await this.storeThemeModel.create(CreateStoreThemeDto)
    return storeTheme;
  }

  async findStoreThemes(storeId:number){
    return await this.storeThemeModel.find({storeId});
  }

  async findStoreThemeById(id:string){
    return await this.storeThemeModel.findOne({_id:id});
  }

  async updateStoreTheme(id:string,updateStoreThemeDto:CreateStoreThemeDto){
    return await this.storeThemeModel.findOneAndUpdate({_id:id},updateStoreThemeDto);
  }
  async removeStoreTheme(id:string){
    return await this.storeThemeModel.findOneAndDelete({_id:id});
  }
}
