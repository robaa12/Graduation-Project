import { Store } from "src/store/entities/store.entity";

export interface IUser{
    id: number;
    firstName: string;
    lastName: string;
    stores:{id:number,name:string}[]
    isActive: boolean;
    email: string;
    is_banned: boolean;
    phoneNumber: string;
    address: string;
    createAt:Date;
    updateAt:Date;
}