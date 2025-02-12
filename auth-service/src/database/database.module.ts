import { DataSource } from 'typeorm';
import { Global, Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ConfigService } from '@nestjs/config';

@Global() 
@Module({
    imports:[TypeOrmModule.forRootAsync({
        useFactory: (configService:ConfigService) => ({
            type: 'postgres',
            host: configService.getOrThrow('DATABASE_HOST'),
            port: configService.getOrThrow('DATABASE_PORT'),
            username: configService.getOrThrow('DATABASE_USER'),
            password: configService.getOrThrow('DATABASE_PASSWORD'),
            database: configService.getOrThrow('DATABASE_NAME'),
            autoLoadEntities: true,
            synchronize: true,
        }),
        inject: [ConfigService],
    })],
    providers: [],
    exports: [],
})
export class DatabaseModule {}
