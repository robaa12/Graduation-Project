import { Category } from "src/category/entities/category.entity";
import { Plan } from "src/plans/entities/plan.entity";
import { User } from "src/user/entities/user.entity";
import { Column, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn } from "typeorm";

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
    store_currency: string;

    @ManyToOne(()=> User, user=>user.stores , {onDelete: 'CASCADE'})
    user: User;

    @ManyToOne(()=> Category, category=>category.stores , {onDelete: 'CASCADE'})
    @JoinColumn({ name: 'category_id' })
    category: Category;
}
