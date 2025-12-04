import axios from 'axios'

const URL = "http://localhost:1551/api/service";

interface ServiceList {
    items: Service[];
}

export interface Service {
    id: string
    name: string
}

export async function getServiceListHttp() {
    try {
        console.log("Get service list:", URL)
        const response = await axios.get<ServiceList>(URL);
        console.log("Get service list success:", response.data);

        return response.data.items;
    } catch (error) {
        console.error("Get service list error:", error);

        throw new Error("Failed to get service list");
    }
}

export async function createServiceHttp(service: Service) {
    try {
        console.log("Create service:", URL, service)
        const response = await axios.post(URL, service);
        console.log("Create service success:", response.data);
    } catch (error) {
        console.error("Create service error:", error);

        throw new Error("Failed to create service");
    }
}
