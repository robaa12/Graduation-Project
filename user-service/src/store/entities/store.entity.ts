import { Category } from "src/category/entities/category.entity";
import { Plan } from "src/plans/entities/plan.entity";
import { User } from "src/user/entities/user.entity";
import { Column, CreateDateColumn, Entity, JoinColumn, ManyToOne, OneToMany, PrimaryGeneratedColumn, UpdateDateColumn } from "typeorm";
import { StoreGallery } from "./user-gallery.entity";

@Entity()
export class Store {
    @PrimaryGeneratedColumn()
    id:number;

    @Column()
    store_name: string;

    @Column({nullable:true})
    href: string;

    @Column({unique:true})
    slug:string;

    @Column()
    description: string;

    @Column()
    business_phone:string;

    @Column()
    category_id:number;

    @Column()
    store_currency: string;

    @CreateDateColumn({ type: 'timestamp' })
    created_at: Date;

    @UpdateDateColumn({ type: 'timestamp' })
    updated_at: Date;

    @ManyToOne(()=> User, user=>user.stores , {onDelete: 'CASCADE'})
    user: User;

    @OneToMany(() => StoreGallery, gallery => gallery.store, { onDelete: 'CASCADE' })
    images: StoreGallery[];

    @ManyToOne(()=> Category, category=>category.stores , {onDelete: 'CASCADE'})
    @JoinColumn({ name: 'category_id' })
    category: Category;
}
