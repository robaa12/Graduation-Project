import { Controller, Get, Post, Body, Patch, Param, Delete } from '@nestjs/common';
import { StoreService } from './store.service';
import { CreateStoreDto } from './dto/create-store.dto';
import { UpdateStoreDto } from './dto/update-store.dto';
import { CreateStoreThemeDto } from './dto/create-store-theme.dto';
import { ApiOperation } from '@nestjs/swagger';

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
  @Get('user/:userId')
  @ApiOperation({summary:'Find Stores by User ID'})
  async findStoreByUserId(@Param('userId') userId: number) {
      let store = await this.storeService.findAllStoresByUserId(userId);
      return {
        message: 'Store fetched successfully',
        data : store
      }
  }

    @Post('theme')
    @ApiOperation({summary:'Create Store Theme'})
    async createStoreTheme(@Body() CreateStoreThemeDto:CreateStoreThemeDto){ {
        let storeTheme = await this.storeService.createStoreTheme(CreateStoreThemeDto);
        return {
          message: 'Store theme created successfully',
          data : storeTheme
        }
    }
  }

  @Get('theme/:storeId')
  @ApiOperation({summary:'Find Store Themes'})
  async findStoreThemes(@Param('storeId') storeId:string){
    let storeThemes = await this.storeService.findStoreThemes(+storeId);
    return {
      message: 'Store themes fetched successfully',
      data : storeThemes
    }
  }

  @Patch('theme/:id')
  @ApiOperation({summary:'Update Store Theme'})
  async update(@Param('id') id: string, @Body() CreateStoreThemeDto: CreateStoreThemeDto) {
    const theme = await this.storeService.updateStoreTheme(id, CreateStoreThemeDto);
    return {
      message: 'Store theme updated successfully',
      data : theme
    }
  }

  @Delete(':id')
  @ApiOperation({summary:'Delete Store'})
  async deleteStore(@Param('id') id: number) {
    const store = await this.storeService.deleteStore(id);
    return {
      message: 'Store deleted successfully',
    }
  }

  @Delete('theme/:id')
  @ApiOperation({summary:'Delete Store Theme'})
  async remove(@Param('id') id: string) {
    const theme = await this.storeService.removeStoreTheme(id);
    return {
      message: 'Store theme deleted successfully',
    }
  }
}
