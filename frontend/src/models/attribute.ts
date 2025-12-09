import type {Entity} from "./entity.ts";

export interface Attribute {
    entity: Entity
    id: string
    nameDb: string
    typeDb: string
    default: string
    nullable: boolean
    primaryKey: boolean
    type: string
    length: number
    fk: Entity|null
}
