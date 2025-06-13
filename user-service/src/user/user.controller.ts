import {
  Controller,
  Get,
  Post,
  Body,
  Patch,
  Param,
  Delete,
  HttpCode,
} from '@nestjs/common';
import { UserService } from './user.service';
import { CreateUserDto } from './dto/create-user.dto';
import { UpdateUserDto } from './dto/update-user.dto';
import { ApiOperation, ApiTags } from '@nestjs/swagger';
@ApiTags('User')
@Controller('user')
export class UserController {
  constructor(private readonly userService: UserService) {}

  @ApiOperation({ summary: 'Create User' })
  @Post('')
  async create(@Body() createUserDto: CreateUserDto) {
    let user = await this.userService.create(createUserDto);
    console.log(user);
    return {
      id: user.id, // Match Go `ID` field
      email: user.email, // Match Go `Email`
      first_name: user.firstName, // Match Go `FirstName`
      last_name: user.lastName, // Match Go `LastName`
      stores: user.stores, // Match Go `StoresID`
    };
  }

  @ApiOperation({
    summary: 'Login User',
    requestBody: {
      content: {
        'application/json': {
          schema: { type: 'object' },
          example: {
            email: 'jhon@doe.com',
            password: 'Password',
          },
        },
      },
    },
  })
  @Post('login')
  @HttpCode(200)
  async Login(@Body() body: any) {
    let user = await this.userService.login(body);
    return {
      message: 'User logged in successfully',
      data: user,
    };
  }

  @ApiOperation({
    summary: 'Verify Email',
    requestBody: {
      content: {
        'application/json': {
          schema: { type: 'object' },
          example: {
            email: 'jhon@doe.com',
            otp: '123456',
          },
        },
      },
    },
  })
  @Post('verify-email')
  async verifyEmail(@Body() body: any) {
    let user = await this.userService.verifyEmail(body.email, body.otp);
    return {
      message: 'Email verified successfully',
      data: user,
    };
  }

  @Post('gallery')
  async addPhotoToGallery(@Body() body: { userId: number; imageUrl: string },) {
    const image = await this.userService.addPhotoToGallery(body.userId, body.imageUrl);
    return {
      message: 'Image added to gallery successfully',
      data: image,
    }
  }


  @ApiOperation({ summary: 'Get All Users' })
  @Get('')
  async findAll() {
    const users = await this.userService.findAll();
    return {
      message: 'All Users fetched successfully',
      data: users,
    };
  }

  @Get('gallery/:userId')
  async getGallery(@Param('userId') userId: number) {
    const gallery = await this.userService.getUserGallery(userId);
    return {
      message: 'Gallery fetched successfully',
      data: gallery,
    };
  }

  @ApiOperation({ summary: 'Get User By Id' })
  @Get(':id')
  async findOne(@Param('id') id: string) {
    let user = await this.userService.findOne(+id);
    return {
      message: 'User fetched successfully',
      data: user,
    };
  }

  @ApiOperation({ summary: 'Update User' })
  @Patch(':id')
  async update(@Param('id') id: string, @Body() updateUserDto: UpdateUserDto) {
    let user = await this.userService.update(+id, updateUserDto);
    return {
      message: 'User updated successfully',
      data: user,
    };
  }

  @Delete('gallery/:photoId')
  async deletePhotoFromGallery(@Param('photoId') photoId: number) {
    const result = await this.userService.deletePhotoFromGallery( photoId);
    return {
      message: 'Photo deleted from gallery successfully',
      data: result,
    };

  }

  @ApiOperation({ summary: 'Delete User' })
  @Delete(':id')
  async remove(@Param('id') id: string) {
    return this.userService.remove(+id);
  }
}
