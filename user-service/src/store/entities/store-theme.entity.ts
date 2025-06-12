import { Prop, Schema, SchemaFactory } from '@nestjs/mongoose';
import { Document, Types } from 'mongoose';
@Schema()
export class StoreThemeSchema {
  @Prop()
  storeId:number;

  @Prop({
    type:Object
  })
  theme:any;

  @Prop({
    type:Boolean,
    default:false
  })
  isActive:boolean;
}
export const StoreTheme = SchemaFactory.createForClass(StoreThemeSchema);
