import { IsNotEmpty } from "class-validator";

export class CreateOrderPaymentDto {
    @IsNotEmpty()
    total_price: number;

    @IsNotEmpty()
    email: string;
    @IsNotEmpty()
    customer_name: string
    @IsNotEmpty()
    phone_number: string
    @IsNotEmpty()
    address: string
    @IsNotEmpty()
    payment_method: string
    @IsNotEmpty()
    note: string
    @IsNotEmpty()
    city: string
    @IsNotEmpty()
    governorate: string
    @IsNotEmpty()
    postal_code: string
    @IsNotEmpty()
    shipping_method: string
    @IsNotEmpty()
    order_items: OrderItem[]

    @IsNotEmpty()
    store_id: number;
    @IsNotEmpty()
    token:string;
}

export interface OrderItem {
    sku_id: number
    price: number
    quantity: number
}