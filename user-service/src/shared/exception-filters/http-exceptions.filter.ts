import {
    ArgumentsHost,
    BadRequestException,
    Catch,
    ExceptionFilter,
    HttpException,
    Logger,
  } from '@nestjs/common';
  
  @Catch(HttpException)
  export class HttpExceptionFilter implements ExceptionFilter {
    logger = new Logger('HttpExceptionFilter');
    catch(exception: HttpException, host: ArgumentsHost) {
      const ctx = host.switchToHttp();
      const response = ctx.getResponse();
      const request = ctx.getRequest();
      const statusCode = exception.getStatus() || 500;
  
      const errorResponse = {
        status: false,
        code: statusCode || 500,
        message: exception.message,
      };
  
      if (exception instanceof BadRequestException) {
        errorResponse['message'] =
          exception.getResponse()['message'] || 'Bad Request';
      }
  
      response.status(statusCode).json(errorResponse);
    }
  }