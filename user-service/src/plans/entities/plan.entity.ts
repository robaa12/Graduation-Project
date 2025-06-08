import { Store } from "src/store/entities/store.entity";
import { User } from "src/user/entities/user.entity";
import { Column, CreateDateColumn, Entity, OneToMany, PrimaryGeneratedColumn, UpdateDateColumn } from "typeorm";

@Entity()
export class Plan {
    @PrimaryGeneratedColumn()
    id:number;

    @Column({ type: 'varchar', length: 100 })
    name: string;

    @Column({ type: 'text' })
    description: string;

    @Column({ type: 'decimal', precision: 10, scale: 2 })
    price: number;

    @Column({ type: 'boolean', default: false })
    isActive: boolean;

    @Column({type:'int'})
    num_of_stores: number;

    @CreateDateColumn({ type: 'timestamp' })
    createdAt: Date;

    @UpdateDateColumn({ type: 'timestamp' })
    updatedAt: Date;

    @OneToMany(()=>User , (user) => user.plan, {onDelete:'SET NULL'})
    users: User[];
}
