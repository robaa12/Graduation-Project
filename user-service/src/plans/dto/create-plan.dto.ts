import { ApiProperty } from '@nestjs/swagger';
import { IsString, IsNumber, IsBoolean, IsOptional, IsNotEmpty, Min, MaxLength } from 'class-validator';

export class CreatePlanDto {

    @ApiProperty({
        description: 'The name of the plan',
        example: 'Basic Plan',
        maxLength: 100
    })
    @IsString()
    @IsNotEmpty()
    @MaxLength(100)
    name: string;

    @ApiProperty({
        description: 'A brief description of the plan',
        example: 'This plan includes basic features for small businesses.',
        maxLength: 500
    })
    @IsString()
    description: string;

    @ApiProperty({
        description: 'The price of the plan in the specified currency',
        example: 29.99,
        minimum: 0
    })
    @IsNumber()
    @Min(0)
    price: number;

    @ApiProperty({
        description: 'Indicates whether the plan is currently active',
        example: true,
        default: false
    })
    @IsBoolean()
    @IsOptional()
    isActive?: boolean;

    @ApiProperty({
        description: 'The number of stores allowed under this plan',
        example: 5,
        minimum: 1
    })
    @IsNumber()
    @Min(1)
    num_of_stores: number;
}
