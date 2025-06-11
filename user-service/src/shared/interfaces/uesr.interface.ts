import { Plan } from "src/plans/entities/plan.entity";
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
    plan_expire_date:Date;
    plan:Plan
    createAt:Date;
    updateAt:Date;
}