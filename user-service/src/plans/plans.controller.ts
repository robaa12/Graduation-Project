import { Plan } from 'src/plans/entities/plan.entity';
import { Controller, Get, Post, Body, Patch, Param, Delete } from '@nestjs/common';
import { PlansService } from './plans.service';
import { CreatePlanDto } from './dto/create-plan.dto';
import { UpdatePlanDto } from './dto/update-plan.dto';

@Controller('plans')
export class PlansController {
  constructor(private readonly plansService: PlansService) {}

  @Post()
  async create(@Body() createPlanDto: CreatePlanDto) {
    const plan = await this.plansService.create(createPlanDto);
    return {
      message: 'Plan created successfully',
      data: plan,
    };
  }

  @Get()
  async findAll() {
    const plans = await this.plansService.findAll();
    return {
      message: 'Plans retrieved successfully',
      data: plans,
    };
  }

  @Get(':id')
  async findOne(@Param('id') id: number) {
    const plan = await this.plansService.findOne(id);
    return {
      message: `Plan with id ${id} retrieved successfully`,
      data: plan,
    };
  }

  @Delete(':id')
  async remove(@Param('id') id: number) {
    await this.plansService.remove(id);
    return {
      message: `Plan with id ${id} removed successfully`,
      data:null
    };
  }

}
