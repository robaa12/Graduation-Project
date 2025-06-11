import { Controller, Get, Post, Body, Patch, Param, Delete, Query } from '@nestjs/common';
import { PaymentService } from './payment.service';
import { CreatePaymentDto } from './dto/create-payment.dto';
import { UpdatePaymentDto } from './dto/update-payment.dto';
import { ApiProperty } from '@nestjs/swagger';

@Controller('payment')
export class PaymentController {
  constructor(private readonly paymentService: PaymentService) {}

  @ApiProperty({
    description: 'Create a new payment charge',
  })
  @Post('')
  async createPayment(@Body() createPaymentDto: CreatePaymentDto) {
    return await this.paymentService.createCharge(createPaymentDto);
  }

  @Get('callback')
  async paymentCallback(@Query('tap_id') tapChargeId: string) {
    const chargeDetails = await this.paymentService.retrieveCharge(tapChargeId);
    return {
      message: 'Payment status retrieved',
      status: chargeDetails.status,
      charge: chargeDetails,
    };
  }
}
