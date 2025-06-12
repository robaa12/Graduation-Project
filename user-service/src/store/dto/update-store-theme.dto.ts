import { PartialType } from "@nestjs/swagger";
import { CreateStoreThemeDto } from "./create-store-theme.dto";

export class UpdateStoreThemeDto extends PartialType(CreateStoreThemeDto) {
    theme?:any;
    isActive?:boolean;
}