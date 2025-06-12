import { Column, Entity, ManyToOne, PrimaryGeneratedColumn } from "typeorm";
import { User } from "./user.entity";

@Entity()
export class UserGallery {
    @PrimaryGeneratedColumn()
    id: number;

    @Column({ type: 'varchar', length: 255 })
    imageUrl: string;

    @ManyToOne(()=> User, (user) => user.images, {onDelete: 'CASCADE'})
    user: User;
}