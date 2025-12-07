export interface AttributeCreateRequest {
    id: string
    name_db: string
    type_db: string
    default: string
    fk_table: string
    fk_type: string
    nullable: boolean
    primary_key: boolean
}

export interface AttributeUpdateRequest {
    name_db: string
    type_db: string
    default: string
    fk_table: string
    fk_type: string
    nullable: boolean
    primary_key: boolean
}