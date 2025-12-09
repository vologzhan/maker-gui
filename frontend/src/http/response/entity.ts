import type {AttributeResponse} from "./attribute.ts"

export interface EntityListResponse {
    items: EntityResponse[]
}

export interface EntityResponse {
    id: string
    name_db: string
    name_plural: string
    attributes: AttributeResponse[]
}
