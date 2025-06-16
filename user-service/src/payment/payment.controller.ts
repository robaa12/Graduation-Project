import { Controller, Get, Post, Body, Patch, Param, Delete, Query } from '@nestjs/common';
import { PaymentService } from './payment.service';
import { CreatePaymentDto } from './dto/create-payment.dto';
import { UpdatePaymentDto } from './dto/update-payment.dto';
import { ApiProperty } from '@nestjs/swagger';
import { CreateOrderPaymentDto } from './dto/create-order-payment.dto';

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
  @Post('order')
  async createOrderPayment(@Body() createOrderPaymentDto: CreateOrderPaymentDto) {
    const data =  await this.paymentService.createOrderPayment(createOrderPaymentDto);
    return {
      message: 'Payment created successfully',
      data: {...data},
    }
  }

  @Get('order/:id')
  async getOrderPayment(@Param('id') id: string) {
    const orderPayment = await this.paymentService.retrieveOrderPayment(id);
    return {
      message: 'Order payment retrieved successfully',
      data: {...orderPayment},
    };
  }

  @Get('callback')
  async paymentCallback(@Query('tap_id') tapChargeId: string) {
    const chargeDetails = await this.paymentService.retrieveCharge(tapChargeId);
    return {
      message: 'Payment status retrieved',
      status: chargeDetails.charge.status,
      data: {...chargeDetails},
    };
  }
}
