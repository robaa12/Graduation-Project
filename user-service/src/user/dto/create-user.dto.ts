import { ApiProperty } from "@nestjs/swagger";
import { IsNotEmpty } from "class-validator";

export class CreateUserDto {
    @ApiProperty({
        description: 'The first name of the User',
        example: 'John',
        required: true,
    })
    @IsNotEmpty()
    firstName: string;

    @ApiProperty({
        description: 'The last name of the User',
        example: 'Doe',
        required: true,
    })
    @IsNotEmpty()
    lastName: string;

    @ApiProperty({
        description: 'The password of the User',
        example: 'password',
        required: true,
    })
    @IsNotEmpty()
    password:string;

    @ApiProperty({
        description: 'The email of the User',
        example: 'jhon@doe.com',
        required: true,
    })
    @IsNotEmpty()
    email: string;

    
    plan_id:number = 1;
}


