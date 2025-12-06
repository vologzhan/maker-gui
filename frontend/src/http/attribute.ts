import axios from "axios";

const URL = "http://localhost:1551/api/attribute";

export interface Attribute {
    id: string
    name_db: string
    type_db: string
    default: string
    fk_table: string
    fk_type: string
    nullable: boolean
    primary_key: boolean

    type: string
    length: number
}

export async function createAttributeHttp(entityId: string, attribute: Attribute) {
    try {
        const url = `${URL}/entity/${entityId}`

        console.log("Create attribute:", url, attribute);
        const response = await axios.post(url, attribute);
        console.log("Create attribute success:", response.data);
    } catch (error) {
        console.error("Create attribute error:", error);

        throw new Error("Failed to create attribute");
    }
}

export async function updateAttributeHttp(attribute: Attribute) {
    try {
        const url = `${URL}/${attribute.id}`

        console.log("Update attribute:", url, attribute);
        const response = await axios.put(url, attribute);
        console.log("Update attribute success:", response.data);
    } catch (error) {
        console.error("Update attribute error:", error);

        throw new Error("Failed to update attribute");
    }
}

export async function deleteAttributeHttp(id: string) {
    try {
        const url = `${URL}/${id}`

        console.log("Delete attribute:", url);
        const response = await axios.delete(url);
        console.log("Delete attribute success:", response.data);
    } catch (error) {
        console.error("Delete attribute error:", error);

        throw new Error("Failed to delete attribute");
    }
}
