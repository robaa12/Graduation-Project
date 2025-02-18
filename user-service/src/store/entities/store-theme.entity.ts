import { Prop, Schema, SchemaFactory } from "@nestjs/mongoose";

@Schema()
export class StoreThemeSchema {
    @Prop()
    name: string;

    @Prop()
    img:string;

    @Prop()
    localPath:string;

    @Prop()
    isActive:boolean;
    
    @Prop()
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

    @Prop()
    storeId:number;
}

const StoreTheme = SchemaFactory.createForClass(StoreThemeSchema);
export default StoreTheme;