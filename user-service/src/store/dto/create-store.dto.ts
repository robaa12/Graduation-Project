import { ApiProperty } from "@nestjs/swagger";
import { IsNotEmpty, IsNumber } from "class-validator";

export class CreateStoreDto {
    @ApiProperty({type:String,description:'store name'})
    @IsNotEmpty()
    store_name:string

    @ApiProperty({type:String,description:'store description'})
    @IsNotEmpty()
    description:string

    @ApiProperty({type:String,description:'store phone number'})
    @IsNotEmpty()
    business_phone:string

    @ApiProperty({type:Number,description:'category id'})
    @IsNotEmpty()
    @IsNumber()
    category_id:number

    @ApiProperty({type:String,description:'store currency'})
    @IsNotEmpty()
    store_currency:string

    @ApiProperty({type:Number,description:'user id'})
    @IsNotEmpty()
    user_id:number
}
