import { Store } from "src/store/entities/store.entity";
import { Column, Entity, OneToMany, PrimaryGeneratedColumn } from "typeorm";

@Entity()
export class Category {
    @PrimaryGeneratedColumn()
    id:number;

    @Column({nullable:false})
    name:string

    @Column({nullable:true})
    description:string

    @OneToMany(()=>Store , store=>store.category)
    stores:Store[]
}
