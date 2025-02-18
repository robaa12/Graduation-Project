import { ApiProperty } from "@nestjs/swagger";
import { IsNotEmpty } from "class-validator";

export class CreateStoreThemeDto {
        @ApiProperty({
            description:'Name of the theme',
            example:'Theme 1'
        })
        @IsNotEmpty()
        name: string;
        
        @ApiProperty({
            description:'Image of the theme',
            example:'https://www.google.com'
        })
        @IsNotEmpty()
        img:string;

        @ApiProperty({
            description:'Local path of the theme',
            example:'https://www.google.com'
        })
        @IsNotEmpty()
        localPath:string;

        @ApiProperty({
            description:'Pages of the theme',
            example:[{
                name:'Home',
                section:['header','footer'],
                body:[{
                    type:'carousel',
                    name:'Carousel',
                    data:[{
                        title:'Title',
                        subtitle:'Subtitle',
                        imageUrl:'https://www.google.com'
                    }]
                }]
            }]
        })
        @IsNotEmpty()
        pages:Array<{
            _id:false,
            name:string,
            section:string[],
            body:Array<{
                _id:false,
                type:string,
                name:string,
                data:Array<{
                    title:string,
                    subtitle:string,
                    imageUrl:string
                }>
            }>
        }>

        @ApiProperty({
            description:'Id of the store',
            example:1
        })
        @IsNotEmpty()
        storeId:number;
}