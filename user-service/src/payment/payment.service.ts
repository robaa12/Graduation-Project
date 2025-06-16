import { PlansService } from 'src/plans/plans.service';
import { UserService } from 'src/user/user.service';
import { BadRequestException, HttpException, Injectable } from '@nestjs/common';
import { CreatePaymentDto } from './dto/create-payment.dto';
import { UpdatePaymentDto } from './dto/update-payment.dto';
import axios from 'axios';
import * as dotenv from 'dotenv';
import { CreateOrderPaymentDto } from './dto/create-order-payment.dto';
import { InjectRepository } from '@nestjs/typeorm';
import { StoreOrderPayment } from './entities/store-order-payment.entity';
import { Repository } from 'typeorm';
import { StoreService } from 'src/store/store.service';
dotenv.config();
@Injectable()
export class PaymentService {
  constructor(
    private UserService: UserService,
    private PlansService: PlansService,
    @InjectRepository(StoreOrderPayment) private storeOrderPaymentRepository: Repository<StoreOrderPayment>,
    private StoreService:StoreService
  ) {}
 async createCharge(createPaymentDto:CreatePaymentDto) {  
  const plan = await this.PlansService.findOne(createPaymentDto.plan_id);
  const user = await this.UserService.findOne(createPaymentDto.user_id);
  console.log("Plan",plan);
  console.log("User",user);
  
  
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
              number: '1143101501',
            },
          },
          source: {
            id: 'src_all', 
          },
          post: {
            "url": "http://localhost:3000/payment/callback", 
          },  
          redirect: {
            url: 'http://localhost:3001/en/payment/success',
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
      throw new BadRequestException(error.response?.data || 'Tap payment error', error.response?.status || 400);
    }
  }

  async retrieveCharge(chargeId: string) {
    try {
      const response = await axios.get(`${process.env.PAYMENT_URL}/charges/${chargeId}`, {
        headers: {
          Authorization: `Bearer ${process.env.PAYMENT_SECRET}`,
        },
      });
      const userPlanPayment = await this.UserService.updatePayment(chargeId , response.data);
      return {
        charge: response.data,
        userPlanPayment: userPlanPayment,
      };
    } catch (error) {
      throw new HttpException(error.response?.data || 'Error retrieving charge', error.response?.status || 500);
    }
  }

  async createOrderPayment(createOrderPayment:CreateOrderPaymentDto){    
    try{
        const order = await axios.post(
        `${process.env.ORDER_SERVICE_URL}/stores/${createOrderPayment.store_id}/orders`,{
          ...createOrderPayment
        });
        
        const payment = await axios.post(
        `${process.env.PAYMENT_URL}/charges`,
        {
          amount: createOrderPayment.total_price,
          currency: 'EGP',
          threeDSecure: true,
          save_card: false,
          description: 'Payment for Order',
          statement_descriptor: 'MyStore',
          metadata: { order_id: order.data.id },
          reference: {
            transaction: 'txn_0001',
            order: order.data.id,
          },
          receipt: {
            email: true,
            sms: false,
          },
          customer: {
            first_name: createOrderPayment.customer_name.split(' ')[0],
            last_name: createOrderPayment.customer_name.split(' ')[1] || '',
            email: createOrderPayment.email,
            phone: {
              country_code: '20',
              number: createOrderPayment.phone_number,
            },
          },
          source: {
            id: 'src_all', 
          },
          post: {
            "url": "http://localhost:3000/payment/callback", 
          },  
          redirect: {
            url: 'http://localhost:3001/shop/order/success',
          },
        },
        {
          headers: {
            Authorization: `Bearer ${process.env.PAYMENT_SECRET}`,
            'Content-Type': 'application/json',
          },
        },
        );
        const storeOrderPayment = this.storeOrderPaymentRepository.create({
          store_id: createOrderPayment.store_id,
          order_id: order.data.order_id,
          charge_id: payment.data.id,
        });
        await this.storeOrderPaymentRepository.save(storeOrderPayment);
        return {
          order: order.data,
          payment: payment.data,
        }

    }catch(error){
      console.log(error);
      
      throw new BadRequestException(error.response?.data || 'Error creating order payment', error.response?.status || 400);
    }
   

  }

  async retrieveOrderPayment(chargeId: string) {
    try {
      const response = await axios.get(
        `${process.env.PAYMENT_URL}/charges/${chargeId}`,
        {
          headers: {
            Authorization: `Bearer ${process.env.PAYMENT_SECRET}`,
          },
        }
      );
      console.log(response.data);
      
      const orderPayment = await this.storeOrderPaymentRepository.findOne({
        where: { charge_id: chargeId },
      });
      const store = await this.StoreService.findOne(orderPayment.store_id);
      return {
        charge: response.data,
        store
      };
    } catch (error:any) {
      console.log(error);
      throw new BadRequestException(error.response?.data || 'Error retrieving order payment', error.response?.status || '');
    }
  }
}

