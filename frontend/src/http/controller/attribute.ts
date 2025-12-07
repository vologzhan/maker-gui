import axios from "axios"
import type {AttributeCreateRequest, AttributeUpdateRequest} from "../request/attribute.ts"
import type {SuccessResponse} from "../response/success.ts"

const URL = "http://localhost:1551/api/attribute"

export async function CreateAttribute(entityId: string, req: AttributeCreateRequest) {
    try {
        const url = `${URL}/entity/${entityId}`

        console.log("Create attribute:", url, req)
        const response = await axios.post<SuccessResponse>(url, req)
        console.log("Create attribute success:", response.data)

        return response.data
    } catch (error) {
        console.error("Create attribute error:", error)

        throw new Error("Failed to create attribute")
    }
}

export async function UpdateAttribute(id: string, req: AttributeUpdateRequest) {
    try {
        const url = `${URL}/${id}`

        console.log("Update attribute:", url, req)
        const response = await axios.put<SuccessResponse>(url, req)
        console.log("Update attribute success:", response.data)

        return response.data
    } catch (error) {
        console.error("Update attribute error:", error)

        throw new Error("Failed to update attribute")
    }
}

export async function DeleteAttribute(id: string) {
    try {
        const url = `${URL}/${id}`

        console.log("Delete attribute:", url)
        const response = await axios.delete<SuccessResponse>(url)
        console.log("Delete attribute success:", response.data)

        return response.data
    } catch (error) {
        console.error("Delete attribute error:", error)

        throw new Error("Failed to delete attribute")
    }
}
