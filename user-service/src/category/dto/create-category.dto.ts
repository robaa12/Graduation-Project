import { ApiProperty } from "@nestjs/swagger";
import { IsNotEmpty } from "class-validator";

export class CreateCategoryDto {
    @ApiProperty({
        description: 'The name of the category',
        example: 'Electronics',
        maxLength: 100
    })
    @IsNotEmpty()
    name: string;
    
    @ApiProperty({
        description: 'A brief description of the category',
        example: 'Devices and gadgets related to electronics',
        maxLength: 500
    })
    @IsNotEmpty()
    description: string;
}
