import { DataSource } from 'typeorm';
import { Global, Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ConfigService } from '@nestjs/config';

@Global()
@Module({
  imports: [
    TypeOrmModule.forRootAsync({
      useFactory: (configService: ConfigService) => ({
        type: 'postgres',
        host: process.env.DSN ? new URL(process.env.DSN).hostname : 'localhost',
        port: process.env.DSN ? parseInt(new URL(process.env.DSN).port) : 5432,
        username: process.env.DSN
          ? new URL(process.env.DSN).username
          : 'postgres',
        password: process.env.DSN
          ? new URL(process.env.DSN).password
          : 'password',
        database: process.env.DSN
          ? new URL(process.env.DSN).pathname.replace('/', '')
          : 'user_db',
        autoLoadEntities: true,
        synchronize: process.env.APP_ENV !== 'production',
      }),
      inject: [ConfigService],
    }),
  ],
  providers: [],
  exports: [],
})
export class DatabaseModule {}
