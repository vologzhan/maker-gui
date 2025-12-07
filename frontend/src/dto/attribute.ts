import type {Entity} from "./entity.ts";

export interface Attribute {
    entity: Entity
    id: string
    nameDb: string
    typeDb: string
    default: string
    fkTable: string
    fkType: string
    nullable: boolean
    primaryKey: boolean
    type: string // secondary
    length: number // secondary
}
