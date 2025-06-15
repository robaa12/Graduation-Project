import {
  Controller,
  Get,
  Post,
  Body,
  Patch,
  Param,
  Delete,
} from '@nestjs/common';
import { StoreService } from './store.service';
import { CreateStoreDto } from './dto/create-store.dto';
import { UpdateStoreDto } from './dto/update-store.dto';
import { CreateStoreThemeDto } from './dto/create-store-theme.dto';
import { ApiOperation } from '@nestjs/swagger';
import { UpdateStoreThemeDto } from './dto/update-store-theme.dto';
import { AddGalleryImagesDto } from './dto/add-gallery-images.dto';

@Controller('store')
export class StoreController {
  constructor(private readonly storeService: StoreService) {}

  @Post('')
  async create(@Body() createStoreDto: CreateStoreDto) {
    let store = await this.storeService.createStore(createStoreDto);
    return {
      message: 'Store created successfully',
      data: store,
    };
  }

  @Get('')
  async findAll() {
    const stores = await this.storeService.findAll();
    return {
      message: 'All Stores fetched successfully',
      data: stores,
    };
  }

  @Get('slug/:slug')
  @ApiOperation({summary:'Find Store by Slug'})
  async findStoreBySlug(@Param('slug') slug: string) {
      let store = await this.storeService.findStoreBySlug(slug);
      return {
        message: 'Store fetched successfully',
        data : store
      }
  }
  
  @Get(':id')
  async findOne(@Param('id') id: string) {
    let store = await this.storeService.findOne(+id);
    return {
      message: 'Store fetched successfully',
      data: store,
    };
  }
  @Get('user/:userId')
  @ApiOperation({ summary: 'Find Stores by User ID' })
  async findStoreByUserId(@Param('userId') userId: number) {
    let store = await this.storeService.findAllStoresByUserId(userId);
    return {
      message: 'Store fetched successfully',
      data: store,
    };
  }

  @Post('theme')
  @ApiOperation({ summary: 'Create Store Theme' })
  async createStoreTheme(@Body() CreateStoreThemeDto: CreateStoreThemeDto) {
    {
      let storeTheme =
        await this.storeService.createStoreTheme(CreateStoreThemeDto);
      return {
        message: 'Store theme created successfully',
        data: storeTheme,
      };
    }
  }

  @Post('gallery')
  async addPhotoToGallery(@Body() body: { storeId: number; imageUrl: string }) {
    const image = await this.storeService.addPhotoToGallery(
      body.storeId,
      body.imageUrl,
    );
    return {
      message: 'Image added to gallery successfully',
      data: image,
    };
  }

  @Post('gallery/bulk')
  @ApiOperation({ summary: 'Add multiple images to gallery' })
  async addPhotosToGallery(@Body() addGalleryImagesDto: AddGalleryImagesDto) {
    const images = await this.storeService.addPhotosToGallery(
      addGalleryImagesDto.storeId,
      addGalleryImagesDto.imageUrls,
    );
    return {
      message: 'Images added to gallery successfully',
      data: images,
    };
  }

  @Get('gallery/:storeId')
  async getGallery(@Param('storeId') storeId: number) {
    const gallery = await this.storeService.getStoreGallery(storeId);
    return {
      message: 'Gallery fetched successfully',
      data: gallery,
    };
  }

  @Get('theme/:storeId')
  @ApiOperation({ summary: 'Find Store Themes' })
  async findStoreThemes(@Param('storeId') storeId: string) {
    let storeThemes = await this.storeService.findStoreThemes(+storeId);
    return {
      message: 'Store themes fetched successfully',
      data: storeThemes,
    };
  }
  @Get('theme/:storeId/active')
  @ApiOperation({ summary: 'Find Active Store Theme' })
  async findActiveStoreTheme(@Param('storeId') storeId: string) {
    let storeTheme = await this.storeService.fincStoreActiveTheme(+storeId);
    return {
      message: 'Active Store theme fetched successfully',
      data: storeTheme,
    };
  }
  @Get('theme/slug/:slug/active')
  @ApiOperation({ summary: 'Find Active Store Theme' })
  async findActiveStoreThemeBySlug(@Param('slug') slug: string) {
    let storeTheme =
      await this.storeService.findStoreActiveThemeByStoreSlug(slug);
    return {
      message: 'Active Store theme fetched successfully',
      data: storeTheme,
    };
  }

  @Delete(':id')
  @ApiOperation({ summary: 'Delete Store' })
  async deleteStore(@Param('id') id: number) {
    const store = await this.storeService.deleteStore(id);
    return {
      message: 'Store deleted successfully',
    };
  }

  @Delete('gallery/:photoId')
  async deletePhotoFromGallery(@Param('photoId') photoId: number) {
    const result = await this.storeService.deletePhotoFromGallery(photoId);
    return {
      message: 'Photo deleted from gallery successfully',
      data: result,
    };
  }

  @Delete('theme/:id')
  @ApiOperation({ summary: 'Delete Store Theme' })
  async remove(@Param('id') id: string) {
    const theme = await this.storeService.removeStoreTheme(id);
    return {
      message: 'Store theme deleted successfully',
    };
  }
}
