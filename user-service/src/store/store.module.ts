import { forwardRef, Module } from '@nestjs/common';
import { StoreService } from './store.service';
import { StoreController } from './store.controller';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Store } from './entities/store.entity';
import { EmailService } from 'src/shared/services/email/email.service';
import { UserModule } from 'src/user/user.module';
import { MongooseModule } from '@nestjs/mongoose';
import { CategoryModule } from 'src/category/category.module';
import { PlansModule } from 'src/plans/plans.module';
import { StoreTheme } from './entities/store-theme.entity';

@Module({
  imports: [TypeOrmModule.forFeature([Store]) , forwardRef(()=>CategoryModule), forwardRef(()=>PlansModule), UserModule , MongooseModule.forFeature([{ name: 'StoreTheme', schema: StoreTheme }])],
  controllers: [StoreController],
  providers: [StoreService , EmailService ],
})
export class StoreModule {}
