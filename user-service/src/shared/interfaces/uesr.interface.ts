import { Store } from "src/store/entities/store.entity";

export interface IUser{
    id: number;
    firstName: string;
    lastName: string;
    stores_id:number[]
    isActive: boolean;
    email: string;
    is_banned: boolean;
    phoneNumber: string;
    address: string;
    createAt:Date;
    updateAt:Date;
}