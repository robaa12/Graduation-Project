import { ApiProperty } from "@nestjs/swagger";
import { isNotEmpty, IsNotEmpty } from "class-validator";

export class CreateStoreThemeDto {
  @ApiProperty({
    type:Number,
    example:1
  })
  @IsNotEmpty()
  storeId:number;

  @ApiProperty({
    type:Object,
  })
  @IsNotEmpty()
  theme:any;

  @ApiProperty({
    type:Boolean,
    default:false
  })
  @IsNotEmpty()
  isActive:boolean = true;
}