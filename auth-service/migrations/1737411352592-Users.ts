import { Logger } from "@nestjs/common";
import { MigrationInterface, QueryRunner } from "typeorm";

export class Users1737411352592 implements MigrationInterface {
    private readonly logger = new Logger(Users1737411352592.name);
    public async up(queryRunner: QueryRunner): Promise<void> {
        
        await queryRunner.query(`CREATE TABLE "users" (
            "id" SERIAL NOT NULL,
            "name" character varying NOT NULL,
            "email" character varying NOT NULL,
            "password" character varying NOT NULL,
            "is_active" boolean NOT NULL DEFAULT true,
            "is_banned" boolean NOT NULL DEFAULT false,
            "phone_number" character varying,
            "address" character varying,
            "store_id" integer,
            "created_at" TIMESTAMP NOT NULL DEFAULT now(),
            "updated_at" TIMESTAMP NOT NULL DEFAULT now(),
            PRIMARY KEY ("id"))`);
        this.logger.log('Users table created');
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`DROP TABLE "users"`);
        this.logger.log('Users table dropped');
    }

}
