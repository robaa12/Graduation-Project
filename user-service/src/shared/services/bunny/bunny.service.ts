import { Injectable } from '@nestjs/common';
import axios from 'axios';
@Injectable()
export class BunnyService {
    private readonly storageZone = 'flocally';
    private readonly apiKey = 'f70fdcad-5903-4bbc-93383bcecd73-eb2a-4c4b';

    async uploadFile(file: Express.Multer.File): Promise<string> {
      const fileName = Date.now() + '-' + file.originalname;        
      const uploadUrl = `https://storage.bunnycdn.com/${this.storageZone}/${fileName}`;
      try {
          await axios.put(uploadUrl, file.buffer, {
              headers: {
              AccessKey: this.apiKey,
              'Content-Type': file.mimetype,
          },
      });
          return `https://${this.storageZone}.b-cdn.net/${fileName}`;
      } catch (error) {
              throw new Error(`File upload failed: ${error.message}`);
          }
    }

    async uploadVideo(file: Express.Multer.File): Promise<string> {
        try {          
          const uploadUrl = `https://storage.bunnycdn.com/${this.storageZone}/${file.originalname}`;
          const response = await axios.put(uploadUrl, file.buffer, {
            headers: {
              AccessKey: this.apiKey,
              'Content-Type': file.mimetype,
            },
          });
          console.log(response);
          
            return `https://${this.storageZone}.b-cdn.net/${file.originalname}`;
        } catch (error) {
          throw new Error(`Bunny.net upload failed: ${error.message}`);
        }
      }

    async uploadMultipleFiles(files: Express.Multer.File[]): Promise<string[]> {
        const uploadPromises = files.map(async (file) => {
            const uploadUrl = `https://storage.bunnycdn.com/${this.storageZone}/${file.originalname}`;
            try {
            await axios.put(uploadUrl, file.buffer, {
                    headers: {
                    AccessKey: this.apiKey,
                    'Content-Type': file.mimetype,
                },
            });

            return `https://${this.storageZone}.b-cdn.net/${file.originalname}`;
          } catch (error) {
            throw new Error(`File upload failed: ${error.message}`);
          }
        });
    
        return Promise.all(uploadPromises);
    }

    async deleteFile(fileName: string): Promise<void> {
        const deleteUrl = `https://storage.bunnycdn.com/${this.storageZone}/${fileName}`;
        return axios.delete(deleteUrl, {
            headers: {
                AccessKey: this.apiKey,
            },
        })
    }

}
