import { Injectable, NotFoundException, OnModuleInit } from '@nestjs/common';
import { CreatePlanDto } from './dto/create-plan.dto';
import { UpdatePlanDto } from './dto/update-plan.dto';
import { InjectRepository } from '@nestjs/typeorm';
import { Plan } from './entities/plan.entity';
import { Repository } from 'typeorm';
import { defaultPlans } from 'src/shared/constants/contants';

@Injectable()
export class PlansService implements OnModuleInit {
  constructor(
    @InjectRepository(Plan) private planRepository: Repository<Plan>,
  ) {}

  async onModuleInit() {
    const count = await this.planRepository.count();
    if (count === 0) {
    for (const plan of defaultPlans) {
      const newPlan = this.planRepository.create(plan);
      await this.planRepository.save(newPlan); // sequential save
    }
    console.log('Default plans created in order');
  }
    console.log('PlansService initialized');
  }
  async create(createPlanDto: CreatePlanDto) {
    const plan = await this.planRepository.create(createPlanDto);
    return await this.planRepository.save(plan);
  }

  async findAll() {
    return await this.planRepository.find();
  }

  async findOne(id: number) {
    const plan = await this.planRepository.findOne({ where: { id } });
    if (!plan) {
      throw new NotFoundException(`Plan with id ${id} not found`);
    }
    return plan;
  }

  update(id: number, updatePlanDto: UpdatePlanDto) {
    return `This action updates a #${id} plan`;
  }

  async remove(id: number) {
    const plan = await this.planRepository.findOne({ where: { id } });
    if (!plan) {
      throw new NotFoundException(`Plan with id ${id} not found`);
    }
    await this.planRepository.remove(plan);
    return ;
  }
}
