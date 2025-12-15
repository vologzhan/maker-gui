import type {AttributeDto} from "src/dto/attribute.ts";
import {v4 as uuid} from "uuid";
import {CreateAttribute, UpdateAttribute} from "src/http/controller/attribute.ts";
import type {EntityDto} from "src/dto/entity.ts";

export async function saveAttribute(attr: AttributeDto, entities: EntityDto[]) {
    if (attr.entity.id === "" || attr.nameDb === "" || attr.typeDb === "") {
        return
    }

    if (attr.id === "") {
        await createAttribute(attr)
    } else {
        await updateAttribute(attr)
    }

    if (!attr.primaryKey) {
        return
    }

    const newFkType = calcFkTypeByDbType(attr.typeDb)

    for (const fkEntity of entities) {
        for (const fkAttr of fkEntity.attributes) {
            if (fkAttr.fk !== attr.entity) {
                continue
            }

            if (calcFkTypeByDbType(fkAttr.typeDb) !== newFkType) {
                fkAttr.typeDb = newFkType

                await saveAttribute(fkAttr, entities)
            }
        }
    }
}

async function createAttribute(attr: AttributeDto) {
    attr.id = uuid()

    await CreateAttribute(attr.entity.id, {
        id: attr.id,
        name_db: attr.nameDb,
        type_db: attr.typeDb,
        default: attr.default,
        fk_table: attr.fk ? attr.fk.nameDb : "",
        fk_type: attr.fk ? attr.type : "",
        nullable: attr.nullable,
        primary_key: attr.primaryKey,
    })
}

async function updateAttribute(attr: AttributeDto) {
    await UpdateAttribute(attr.id, {
            name_db: attr.nameDb,
            type_db: attr.typeDb,
            default: attr.default,
            fk_table: attr.fk ? attr.fk.nameDb : "",
            fk_type: attr.fk ? attr.type : "",
            nullable: attr.nullable,
            primary_key: attr.primaryKey,
        }
    )
}

export async function addAttribute(entity: EntityDto, entities: EntityDto[], options?: Partial<AttributeDto>) {
    const defaults: AttributeDto = {
        entity: entity,
        id: "",
        nameDb: "",
        typeDb: "",
        default: "",
        nullable: false,
        primaryKey: false,
        type: "",
        length: 0,
        fk: null,
    }

    const attr: AttributeDto = {
        ...defaults,
        ...options,
    }

    entity.attributes.push(attr)

    return saveAttribute(attr, entities)
}

export function calcFkTypeByDbType(typeDb: string) {
    if (typeDb === "uuid") {
        return "uuid"
    }
    return "int"
}
