import { Plan } from "src/plans/entities/plan.entity";
import { Column, CreateDateColumn, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn, UpdateDateColumn } from "typeorm";
import { User } from "./user.entity";

@Entity()
export class UserPlanPayment {
    @PrimaryGeneratedColumn()
    id: number;
    
    @Column()
    user_id: number;
    
    @Column()
    plan_id: number;
    
    @Column({nullable: true})
    paymentDate: Date;
    
    @Column({ type: 'decimal', precision: 10, scale: 2 })
    amount: number;
    
    @Column()
    status: string;
    
    @Column()
    charge_id: string;
    
    @Column()
    currency:string;
    
    @CreateDateColumn()
    createdAt: Date;
    
    @UpdateDateColumn()
    updatedAt: Date;

    @ManyToOne(()=> User, user => user.payments, { onDelete: 'CASCADE' })
    @JoinColumn({ name: 'user_id' })
    user: User;

    @ManyToOne(() => Plan, plan => plan.payments, { onDelete: 'CASCADE' })
    @JoinColumn({ name: 'plan_id' })
    plan: Plan;
}