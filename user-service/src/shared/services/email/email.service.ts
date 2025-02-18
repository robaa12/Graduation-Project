import { MailerService } from '@nestjs-modules/mailer';
import { Injectable } from '@nestjs/common';

@Injectable()
export class EmailService {
    constructor(private readonly mailService: MailerService) {}

    async sendVerficationMail(otp:string , email:string) {
        const message = `Welcome to Motager. Your OTP is ${otp} Please Use this OTP to verify your email address.`;
        await this.mailService.sendMail({
        from: 'Motager <no-reply-elms450@zohomail.com>',
        to: email,
        subject: `Motager Email Verification`,
        text: message,
        });
    }
}
