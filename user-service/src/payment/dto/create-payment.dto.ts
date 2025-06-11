import { ApiProperty } from "@nestjs/swagger";
import { IsNotEmpty } from "class-validator";

export class CreatePaymentDto {
    @ApiProperty({
        description:'the user ID for the payment',
        example: 1,
    })
    @IsNotEmpty()
    user_id:number;

    @ApiProperty({
        description: 'The plan ID for the payment',
        example: 1,
        type: Number,
    })
    @IsNotEmpty()
    plan_id:number;

}
