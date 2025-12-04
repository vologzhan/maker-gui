import axios from "axios";
import {type Attribute} from "./attribute.ts";

const URL = "http://localhost:1551/api/entity";

interface EntityList {
    items: Entity[];
}

export interface Entity {
    id: string
    name_db: string
    attributes: Attribute[]
}

export async function getEntityListHttp(serviceId: string) {
    try {
        const url = `${URL}/service/${serviceId}`

        console.log("Get entity list:", url);
        const response = await axios.get<EntityList>(url);
        console.log("Get entity list success:", response.data);

        return response.data.items;
    } catch (error) {
        console.error("Get entity list error:", error);

        throw new Error("Failed to get entity list");
    }
}

export async function createEntityHttp(serviceId: string, entity: Entity) {
    try {
        const url = `${URL}/service/${serviceId}`

        console.log("Create entity:", url, entity);
        const response = await axios.post(url, entity);
        console.log("Create entity success:", response.data);
    } catch (error) {
        console.error("Create entity error:", error);

        throw new Error("Failed to create entity");
    }
}

export async function updateEntityHttp(entity: Entity) {
    try {
        const url = `${URL}/${entity.id}`

        console.log("Update entity:", url, entity);
        const response = await axios.put(url, entity);
        console.log("Update entity success:", response.data);
    } catch (error) {
        console.error("Update entity error:", error);

        throw new Error("Failed to update entity");
    }
}

export async function deleteEntityHttp(id: string) {
    try {
        const url = `${URL}/${id}`

        console.log("Delete entity:", url);
        const response = await axios.delete(url);
        console.log("Delete entity success:", response.data);
    } catch (error) {
        console.error("Delete entity error:", error);

        throw new Error("Failed to delete entity");
    }
}
