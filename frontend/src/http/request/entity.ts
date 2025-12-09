export interface EntityCreateRequest {
    id: string
    name_db: string
    name_plural: string
}

export interface EntityUpdateRequest {
    name_db: string
    name_plural: string
}
