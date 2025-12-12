import type {EntityDto} from "./entity.ts";

export interface AttributeDto {
    entity: EntityDto
    id: string
    nameDb: string
    typeDb: string
    default: string
    nullable: boolean
    primaryKey: boolean
    type: string
    length: number
    fk: EntityDto|null
}
