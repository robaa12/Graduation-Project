import { UserModule } from './user/user.module';
import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { DatabaseModule } from './database/database.module';
import { ConfigModule } from '@nestjs/config';
import { JwtModule } from '@nestjs/jwt';
import { MailerModule } from '@nestjs-modules/mailer';
import { EmailService } from './shared/services/email/email.service';
import { StoreModule } from './store/store.module';
import { CategoryModule } from './category/category.module';
import { MongooseModule } from '@nestjs/mongoose';
let mongoUrl = 'mongodb://admin:adminpassword@mongo-db:27017/users?authSource=admin';
@Module({
  imports: [
  ConfigModule.forRoot({ envFilePath: '.env',isGlobal:true}), DatabaseModule , UserModule ,
    JwtModule.register({
    global: true,
    secret: process.env.JWT_SECRET,
    signOptions: { expiresIn: '365d' },
  }),
  MailerModule.forRoot({
    transport: {
      host: process.env.EMAIL_HOST,
      port: 465,
      secure: true,
      auth: {
        user: process.env.EMAIL_USERNAME,
        pass: process.env.EMAIL_PASSWORD,
      },
      tls: {
        rejectUnauthorized: false, // Add this if you face SSL certificate issues
      },
    },
  }),
  MongooseModule.forRoot(mongoUrl, { dbName: 'themes'} ),
  StoreModule,
  CategoryModule,
],
  controllers: [AppController],
  providers: [AppService , EmailService],
})
export class AppModule {}
