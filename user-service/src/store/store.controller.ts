import { Controller, Get, Post, Body, Patch, Param, Delete } from '@nestjs/common';
import { StoreService } from './store.service';
import { CreateStoreDto } from './dto/create-store.dto';
import { UpdateStoreDto } from './dto/update-store.dto';

@Controller('store')
export class StoreController {
  constructor(private readonly storeService: StoreService) {}

  @Post('')
  async create(@Body() createStoreDto: CreateStoreDto) {
      let store = await this.storeService.createStore(createStoreDto);
      return {
        message: 'Store created successfully',
        data : store
      }
  }

  @Get('')
  async findAll() {
    const stores = await this.storeService.findAll();
    return {
      message: 'All Stores fetched successfully',
      data : stores
    }
  }

  @Get(':id')
  async findOne(@Param('id') id: string) {
      let store = await  this.storeService.findOne(+id);
      return {
        message: 'Store fetched successfully',
        data : store
      }
    }

}
