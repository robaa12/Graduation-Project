import { BunnyService } from './shared/services/bunny/bunny.service';
import { BadRequestException, Injectable } from '@nestjs/common';

@Injectable()
export class AppService {
  constructor(
    private BunnyService: BunnyService,
  ) {}
  getHello(): string {
    return 'Hello World!';
  }

  async uploadImage(file: Express.Multer.File): Promise<string> {
    try {
      const result = await this.BunnyService.uploadFile(file);
      return result;
    } catch (error) {
      throw new BadRequestException(`Failed to upload image: ${error.message}`);
    }
  }
}
