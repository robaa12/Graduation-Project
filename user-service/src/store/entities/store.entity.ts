import { Category } from "src/category/entities/category.entity";
import { User } from "src/user/entities/user.entity";
import { Column, Entity, ManyToOne, PrimaryGeneratedColumn } from "typeorm";

@Entity()
export class Store {
    @PrimaryGeneratedColumn()
    id:number;

    @Column()
    store_name: string;

    @Column({nullable:true})
    href: string;

    @Column({nullable:true})
    slug:string;

    @Column()
    description: string;

    @Column()
    business_phone:string;

    @Column()
    category_id:number;

    @Column()
    plan_id:number;

    @Column()
    store_currency: string;

    @ManyToOne(()=> User, user=>user.stores)
    user: User;

    @ManyToOne(()=> Category, category=>category.stores)
    category: Category;
}
