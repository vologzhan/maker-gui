import axios from "axios"
import type {EntityListResponse} from "../response/entity.ts"
import type {SuccessResponse} from "../response/success.ts"
import type {EntityCreateRequest, EntityUpdateRequest} from "../request/entity.ts"

const URL = "http://localhost:1551/api/entity"

export async function GetEntityList(serviceId: string) {
    try {
        const url = `${URL}/service/${serviceId}`

        console.log("Get entity list:", url)
        const response = await axios.get<EntityListResponse>(url)
        console.log("Get entity list success:", response.data)

        return response.data
    } catch (error) {
        console.error("Get entity list error:", error)

        throw new Error("Failed to get entity list")
    }
}

export async function CreateEntity(serviceId: string, req: EntityCreateRequest) {
    try {
        const url = `${URL}/service/${serviceId}`

        console.log("Create entity:", url, req)
        const response = await axios.post<SuccessResponse>(url, req)
        console.log("Create entity success:", response.data)

        return response.data
    } catch (error) {
        console.error("Create entity error:", error)

        throw new Error("Failed to create entity")
    }
}

export async function UpdateEntity(id: string, req: EntityUpdateRequest) {
    try {
        const url = `${URL}/${id}`

        console.log("Update entity:", url, req)
        const response = await axios.put<SuccessResponse>(url, req)
        console.log("Update entity success:", response.data)

        return response.data
    } catch (error) {
        console.error("Update entity error:", error)

        throw new Error("Failed to update entity")
    }
}

export async function DeleteEntity(id: string) {
    try {
        const url = `${URL}/${id}`

        console.log("Delete entity:", url)
        const response = await axios.delete<SuccessResponse>(url)
        console.log("Delete entity success:", response.data)

        return response.data
    } catch (error) {
        console.error("Delete entity error:", error)

        throw new Error("Failed to delete entity")
    }
}
