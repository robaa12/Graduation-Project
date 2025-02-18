import { PartialType } from '@nestjs/mapped-types';
import { CreateUserDto } from './create-user.dto';
import { ApiProperty } from '@nestjs/swagger';

export class UpdateUserDto extends PartialType(CreateUserDto) {
    @ApiProperty({
        description: 'The email of the User',
        example: 'jhon@doe.con'
    })
    email?: string;
    @ApiProperty({
        description: 'The first name of the User',
        example: 'John'
    })
    firstName?: string;
    @ApiProperty({
        description: 'The last name of the User',
        example: 'Doe'
    })
    lastName?: string;
    
    @ApiProperty({
        description: 'The phone number of the User',
        example: '1234567890'
    })
    phoneNumber?: string;

    @ApiProperty({
        description: 'The address of the User',
        example: '1234 Main Street'
    })
    address?: string;

    @ApiProperty({
        description: 'The store id of the User',
        example: 1
    })
    storeId?: number;
}
