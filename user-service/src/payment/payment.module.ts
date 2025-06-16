import { forwardRef, Module } from '@nestjs/common';
import { PaymentService } from './payment.service';
import { PaymentController } from './payment.controller';
import { UserModule } from 'src/user/user.module';
import { PlansModule } from 'src/plans/plans.module';
import { TypeOrmModule } from '@nestjs/typeorm';
import { StoreOrderPayment } from './entities/store-order-payment.entity';
import { StoreModule } from 'src/store/store.module';

@Module({
  imports: [TypeOrmModule.forFeature([StoreOrderPayment]),forwardRef(()=>UserModule) , forwardRef(()=>PlansModule) , forwardRef(()=>StoreModule), ],
  controllers: [PaymentController],
  providers: [PaymentService],
})
export class PaymentModule {}
