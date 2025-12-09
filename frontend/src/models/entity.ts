import type {Attribute} from "./attribute.ts";

export interface Entity {
    id: string
    nameDb: string
    namePlural: string
    namePluralAuto: boolean
    attributes: Attribute[]
}
