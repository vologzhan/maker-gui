import axios from 'axios'
import type {ServiceListResponse} from "../response/service.ts"
import type {ServiceCreateRequest} from "../request/service.ts"
import type {SuccessResponse} from "../response/success.ts"

const URL = "http://localhost:1551/api/service"

export async function GetServiceList() {
    try {
        console.log("Get service list:", URL)
        const response = await axios.get<ServiceListResponse>(URL)
        console.log("Get service list success:", response.data)

        return response.data
    } catch (error) {
        console.error("Get service list error:", error)

        throw new Error("Failed to get service list")
    }
}

export async function CreateService(req: ServiceCreateRequest) {
    try {
        console.log("Create service:", URL, req)
        const response = await axios.post<SuccessResponse>(URL, req)
        console.log("Create service success:", response.data)

        return response.data
    } catch (error) {
        console.error("Create service error:", error)

        throw new Error("Failed to create service")
    }
}
