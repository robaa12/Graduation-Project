import { PlansService } from 'src/plans/plans.service';
import { UserService } from 'src/user/user.service';
import { BadRequestException, HttpException, Injectable } from '@nestjs/common';
import { CreatePaymentDto } from './dto/create-payment.dto';
import { UpdatePaymentDto } from './dto/update-payment.dto';
import axios from 'axios';
import * as dotenv from 'dotenv';
dotenv.config();
@Injectable()
export class PaymentService {
  constructor(
    private UserService: UserService,
    private PlansService: PlansService,
  ) {}
 async createCharge(createPaymentDto:CreatePaymentDto) {
  const plan = await this.PlansService.findOne(createPaymentDto.plan_id);
  const user = await this.UserService.findOne(createPaymentDto.user_id);
    try {
      const response = await axios.post(
        `${process.env.PAYMENT_URL}/charges`,
        {
          amount:plan.price,
          currency:'EGP',
          threeDSecure: true,
          save_card: false,
          description: 'Payment for Order',
          statement_descriptor: 'MyStore',
          metadata: { order_id: '1234' },
          reference: {
            transaction: 'txn_0001',
            order: 'ord_0001',
          },
          receipt: {
            email: true,
            sms: false,
          },
          customer: {
            first_name: 'John',
            last_name: 'Doe',
            email: user.email,
            phone: {
              country_code: '20',
              number: '50000000',
            },
          },
          source: {
            id: 'src_all', 
          },
          post: {
            "url": "http://localhost:3000/payment/callback", 
          },  
          redirect: {
            url: 'https://motager-v2.vercel.app/ar',
          },
        },
        {
          headers: {
            Authorization: `Bearer ${process.env.PAYMENT_SECRET}`,
            'Content-Type': 'application/json',
          },
        },
      );
      const userPlanPayment = await this.UserService.createPayment(user , plan , response)
      return response.data;
    } catch (error) {
      console.log(error);
      throw new BadRequestException(error.response?.data || 'Tap payment error', error.response?.status || 500);
    }
  }

  async retrieveCharge(chargeId: string) {
    try {
      const response = await axios.get(`${process.env.PAYMENT_URL}/charges/${chargeId}`, {
        headers: {
          Authorization: `Bearer ${process.env.PAYMENT_SECRET}`,
        },
      });
      
      return response.data;
    } catch (error) {
      throw new HttpException(error.response?.data || 'Error retrieving charge', error.response?.status || 500);
    }
  }
}

