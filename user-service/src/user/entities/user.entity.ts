import { Plan } from "src/plans/entities/plan.entity";
import { Store } from "src/store/entities/store.entity";
import { Column, CreateDateColumn, DataSource, Entity, JoinColumn, ManyToOne, OneToMany, PrimaryGeneratedColumn, UpdateDateColumn } from "typeorm";
import { UserPlanPayment } from "./user-plan-payment.entity";
import { UserGallery } from "./user-gallery.entity";

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

    @Column({type:'date' , nullable: true})
    plan_expire_date:Date;

    @CreateDateColumn()
    createAt: Date;

    @UpdateDateColumn()
    updateAt: Date;

    @OneToMany(() => Store, Store => Store.user)
    stores: Store[];

    @OneToMany(() => UserGallery, gallery => gallery.user, { onDelete: 'CASCADE' })
    images: UserGallery[];

    @ManyToOne(()=>Plan , (plan)=> plan.users)
    @JoinColumn({ name: 'plan_id' })
    plan:Plan;

    @OneToMany(()=>UserPlanPayment , (userPlanPayment)=> userPlanPayment.user, { onDelete: 'CASCADE' })
    payments: UserPlanPayment[];
}
