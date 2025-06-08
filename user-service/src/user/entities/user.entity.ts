import { Plan } from "src/plans/entities/plan.entity";
import { Store } from "src/store/entities/store.entity";
import { Column, CreateDateColumn, DataSource, Entity, JoinColumn, ManyToOne, OneToMany, PrimaryGeneratedColumn, UpdateDateColumn } from "typeorm";

@Entity()
export class User {
    @PrimaryGeneratedColumn()
    id: number;

    @Column()
    firstName: string;

    @Column()
    lastName: string;

    @Column({
        default: false
    })
    isActive: boolean;

    @Column({
        unique: true
    })
    email: string;

    @Column()
    password: string;

    @Column({
        default: false
    })
    is_banned: boolean;

    @Column({
        nullable: true
    })
    phoneNumber: string;

    @Column({
        nullable: true
    })
    address: string;

    @Column({
        nullable: true
    })
    otp: string;

    @Column({
        nullable: true
    })
    otpExpiry: Date;

    @Column({nullable: true})
    plan_id:number;

    @CreateDateColumn()
    createAt: Date;

    @UpdateDateColumn()
    updateAt: Date;

    @OneToMany(() => Store, Store => Store.user)
    stores: Store[];

    @ManyToOne(()=>Plan , (plan)=> plan.users)
    @JoinColumn({ name: 'plan_id' })
    plan:Plan;
}
