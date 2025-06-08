import { forwardRef, Module } from '@nestjs/common';
import { UserService } from './user.service';
import { UserController } from './user.controller';
import { User } from './entities/user.entity';
import { TypeOrmModule } from '@nestjs/typeorm';
import { EmailService } from 'src/shared/services/email/email.service';
import { PlansModule } from 'src/plans/plans.module';

@Module({
  imports: [TypeOrmModule.forFeature([User]) , forwardRef(()=> PlansModule)],
  controllers: [UserController],
  providers: [UserService , EmailService],
  exports:[UserService]
})
export class UserModule {}
