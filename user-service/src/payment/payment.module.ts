import { forwardRef, Module } from '@nestjs/common';
import { PaymentService } from './payment.service';
import { PaymentController } from './payment.controller';
import { UserModule } from 'src/user/user.module';
import { PlansModule } from 'src/plans/plans.module';

@Module({
  imports: [forwardRef(()=>UserModule) , forwardRef(()=>PlansModule)],
  controllers: [PaymentController],
  providers: [PaymentService],
})
export class PaymentModule {}
