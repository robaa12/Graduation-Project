import {
  Controller,
  Get,
  Post,
  Body,
  Patch,
  Param,
  Delete,
} from '@nestjs/common';
import { CategoryService } from './category.service';
import { CreateCategoryDto } from './dto/create-category.dto';
import { UpdateCategoryDto } from './dto/update-category.dto';
import { ApiOperation } from '@nestjs/swagger';

@Controller('category')
export class CategoryController {
  constructor(private readonly categoryService: CategoryService) {}

  @ApiOperation({ summary: 'Create a new category' })
  @Post()
  async create(@Body() createCategoryDto: CreateCategoryDto) {
    const category =
      await this.categoryService.createCategory(createCategoryDto);
    return {
      message: 'Category created successfully',
      data: { category },
    };
  }

  @ApiOperation({ summary: 'Get all categories' })
  @Get()
  async findAll() {
    const categories = await this.categoryService.findAll();
    return {
      message: 'All categories fetched successfully',
      data: { categories },
    };
  }
}
