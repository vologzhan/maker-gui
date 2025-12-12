import type {AttributeDto} from "./attribute.ts";

export interface EntityDto {
    id: string
    nameDb: string
    namePlural: string
    namePluralAuto: boolean
    attributes: AttributeDto[]
}
