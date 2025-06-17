import { Entity, PrimaryColumn, PrimaryGeneratedColumn } from "typeorm";
@Entity()
export class StoreOrderPayment {
    @PrimaryGeneratedColumn()
    id:number;

    @PrimaryColumn()
    store_id: number;

    @PrimaryColumn()
    order_id: number;

    @PrimaryColumn()
    charge_id: string;
}