import { ApiProperty } from "@nestjs/swagger";
import { IsArray, IsNotEmpty, IsNumber } from "class-validator";

export class AddGalleryImagesDto {
    @ApiProperty({
        description: 'Store ID',
        type: Number,
        example: 1
    })
    @IsNotEmpty()
    @IsNumber()
    storeId: number;

    @ApiProperty({
        description: 'Array of image URLs',
        type: [String],
        example: ['https://example.com/image1.jpg', 'https://example.com/image2.jpg']
    })
    @IsNotEmpty()
    @IsArray()
    imageUrls: string[];
}
