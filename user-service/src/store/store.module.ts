import { Module } from '@nestjs/common';
import { StoreService } from './store.service';
import { StoreController } from './store.controller';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Store } from './entities/store.entity';
import { EmailService } from 'src/shared/services/email/email.service';
import { UserModule } from 'src/user/user.module';

@Module({
  imports: [TypeOrmModule.forFeature([Store]) , UserModule],
  controllers: [StoreController],
  providers: [StoreService , EmailService],
})
export class StoreModule {}
