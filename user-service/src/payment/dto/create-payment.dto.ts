import { ApiProperty } from "@nestjs/swagger";

export class CreatePaymentDto {
    @ApiProperty({
        description:'the user ID for the payment',
        example: 1,
    })
    user_id:number;

    @ApiProperty({
        description: 'The plan ID for the payment',
        example: 1,
        type: Number,
    })
    plan_id:number;

}
