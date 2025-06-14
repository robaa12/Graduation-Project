import { Column, Entity, ManyToOne, PrimaryGeneratedColumn } from "typeorm";
import { Store } from "./store.entity";

@Entity()
export class StoreGallery {
    @PrimaryGeneratedColumn()
    id: number;

    @Column({ type: 'varchar', length: 255 })
    imageUrl: string;

    @ManyToOne(()=> Store, (store) => store.images, {onDelete: 'CASCADE'})
    store: Store;
}