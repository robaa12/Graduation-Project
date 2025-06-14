import { Controller, Get, Post, UploadedFile, UseInterceptors } from '@nestjs/common';
import { AppService } from './app.service';
import { EmailService } from './shared/services/email/email.service';
import { FileInterceptor } from '@nestjs/platform-express';
import { ApiProperty } from '@nestjs/swagger';

@Controller()
export class AppController {
  constructor(private readonly appService: AppService , private EmailService:EmailService) {}

  @Get()
  getHello(): string {
    return this.appService.getHello();
  }

  @ApiProperty({
    description: 'Upload an image file',
    type: 'string',
    format: 'binary',
    required: true,
  })
  @Post('upload/file')
  @UseInterceptors(FileInterceptor('file'))
  async uploadImage(@UploadedFile()file: Express.Multer.File) {
    const url = await this.appService.uploadImage(file);
    return {
      message: 'File uploaded successfully',
      data: url,
    }
  }

}
