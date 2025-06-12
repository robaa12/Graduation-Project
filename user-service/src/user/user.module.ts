import { forwardRef, Module } from '@nestjs/common';
import { UserService } from './user.service';
import { UserController } from './user.controller';
import { User } from './entities/user.entity';
import { TypeOrmModule } from '@nestjs/typeorm';
import { EmailService } from 'src/shared/services/email/email.service';
import { PlansModule } from 'src/plans/plans.module';
import { UserPlanPayment } from './entities/user-plan-payment.entity';
import { UserGallery } from './entities/user-gallery.entity';

@Module({
  imports: [TypeOrmModule.forFeature([User , UserPlanPayment , UserGallery]) , forwardRef(()=> PlansModule)],
  controllers: [UserController],
  providers: [UserService , EmailService],
  exports:[UserService]
})
export class UserModule {}
