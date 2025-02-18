import { Controller, Get } from '@nestjs/common';
import { AppService } from './app.service';
import { EmailService } from './shared/services/email/email.service';

@Controller()
export class AppController {
  constructor(private readonly appService: AppService , private EmailService:EmailService) {}

  @Get()
  getHello(): string {
    return this.appService.getHello();
  }

}
