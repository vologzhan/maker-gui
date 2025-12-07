<script setup lang="ts">
import {defineProps, ref, watch} from 'vue';
import {v4 as uuid} from 'uuid';
import type {Service} from "../dto/service.ts";
import type {Entity} from "../dto/entity.ts";
import type {Attribute} from "../dto/attribute.ts";
import {CreateEntity, DeleteEntity, GetEntityList, UpdateEntity} from "../http/controller/entity.ts";
import {CreateAttribute, DeleteAttribute, UpdateAttribute} from "../http/controller/attribute.ts";
import type {EntityResponse} from "../http/response/entity.ts";
import type {AttributeResponse} from "../http/response/attribute.ts";

const entities = ref<Entity[]>([])

const props = defineProps<{
  service: Service
}>();

watch(() => props.service.id, async (serviceId: string) => {
  entities.value = await getEntityList(serviceId)
});

async function saveEntity(entity: Entity) {
  if (entity.nameDb.trim() === "") {
    return
  }

  if (entity.id === "") {
    await createEntity(entity)
  } else {
    await updateEntity(entity)
  }

  for (const fkEntity of entities.value) {
    for (const fkAttr of fkEntity.attributes) {
      if (fkAttr.fk !== entity) {
        continue
      }

      await saveAttribute(fkAttr)
    }
  }
}

async function createEntity(entity: Entity) {
  entity.id = uuid()

  await CreateEntity(props.service.id, {
    id: entity.id,
    name_db: entity.nameDb,
  })

  for (const attr of entity.attributes) {
    await saveAttribute(attr)
  }
}

async function updateEntity(entity: Entity) {
  await UpdateEntity(entity.id, {
    name_db: entity.nameDb,
  })
}

async function deleteEntity(entity: Entity) {
  entities.value.splice(entities.value.indexOf(entity), 1);

  if (entity.id !== "") {
    await DeleteEntity(entity.id)
  }

  for (const fkEntity of entities.value) {
    for (const fkAttr of fkEntity.attributes) {
      if (fkAttr.fk !== entity) {
        continue
      }

      fkAttr.fk = null
      fkAttr.type = fkAttr.typeDb

      await saveAttribute(fkAttr)
    }
  }
}

async function getEntityList(serviceId: string) {
  const res = await GetEntityList(serviceId)

  const relations = new Map<string, Attribute[]>

  const entities = res.items.map((entityRes: EntityResponse) => {
    const entity: Entity = {
      id: entityRes.id,
      nameDb: entityRes.name_db,
      attributes: [],
    }

    entity.attributes = entityRes.attributes.map((attrRes: AttributeResponse): Attribute => {
      const attr = convertAttributeResponseToDto(entity, attrRes)

      if (attrRes.fk_table === "") {
        return attr
      }

      let relAttributes = relations.get(attrRes.fk_table)
      if (!relAttributes) {
        relAttributes = []
        relations.set(attrRes.fk_table, relAttributes)
      }

      relAttributes.push(attr)

      return attr
    })

    return entity
  })

  entities.forEach((entity) => {
    relations.get(entity.nameDb)?.forEach((attr) => {
      attr.fk = entity
    })
  })

  return entities
}

function convertAttributeResponseToDto(entity: Entity, res: AttributeResponse) {
  const attr: Attribute = {
    entity: entity,
    id: res.id,
    nameDb: res.name_db,
    typeDb: res.type_db,
    default: res.default,
    nullable: res.nullable,
    primaryKey: res.primary_key,
    type: "",
    length: 0,
    fk: null,
  }

  if (res.fk_type !== "") {
    attr.type = res.fk_type
  } else if (attr.typeDb === "timestamp(0)") {
    attr.type = "datetime"
  } else if (attr.typeDb === "text") {
    attr.type = "string"
  } else if (attr.typeDb.startsWith("varchar")) {
    attr.type = "string"
    attr.length = Number(
        attr.typeDb.substring(attr.typeDb.indexOf('(') + 1, attr.typeDb.indexOf(')'))
    )
  } else {
    attr.type = attr.typeDb
  }

  return attr
}

async function saveAttribute(attr: Attribute) {
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

  for (const fkEntity of entities.value) {
    for (const fkAttr of fkEntity.attributes) {
      if (fkAttr.fk !== attr.entity) {
        continue
      }

      if (calcFkTypeByDbType(fkAttr.typeDb) !== newFkType) {
        fkAttr.typeDb = newFkType

        await saveAttribute(fkAttr)
      }
    }
  }
}

async function createAttribute(attr: Attribute) {
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

async function updateAttribute(attr: Attribute) {
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

async function deleteAttribute(attr: Attribute) {
  attr.entity.attributes.splice(attr.entity.attributes.indexOf(attr), 1)

  if (attr.id !== "") {
    await DeleteAttribute(attr.id)
  }
}

function addEntity() {
  const entity: Entity = {id: "", nameDb: "", attributes: []}
  entities.value.push(entity)

  addAttribute(entity, {
    nameDb: "id",
    typeDb: "serial",
    default: "",
    primaryKey: true,
    type: "serial",
  })
}

async function addAttribute(entity: Entity, options?: Partial<Attribute>) {
  const defaults: Attribute = {
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

  const attr: Attribute = {
    ...defaults,
    ...options,
  }

  entity.attributes.push(attr)

  return  saveAttribute(attr)
}

function hasAttribute(entity: Entity, nameDb: string) {
  return entity.attributes.some(attr => attr.nameDb === nameDb)
}

async function editType(attr: Attribute) {
  attr.length = 0
  deleteDefaultNullIfNotNullable(attr)

  if (!isForeignKey(attr)) {
    attr.fk = null
  }

  if (attr.primaryKey) {
    return await editPrimaryKey(attr)
  }

  if (isForeignKey(attr)) {
    if (attr.fk === null) {
      attr.typeDb = ""
    }
  } else if (attr.type === "datetime") {
    attr.typeDb = "timestamp(0)"
  } else if (attr.type === "string") {
    attr.length = 255
    attr.typeDb = calcTypeDbString(attr)
  } else {
    attr.typeDb = attr.type
  }

  return saveAttribute(attr)
}

async function editPrimaryKey(attr: Attribute) {
  if (attr.type === "one-to-one") {
    attr.nameDb = ""
    attr.typeDb = ""

    return
  }

  if (attr.type === "serial") {
    attr.nameDb = "id"
    attr.typeDb = attr.type
  } else if (attr.type === "uuid") {
    attr.nameDb = "uuid"
    attr.typeDb = attr.type
    attr.default = "uuid_generate_v4()"
  }

  await saveAttribute(attr)
}

function editLen(attr: Attribute) {
  attr.typeDb = calcTypeDbString(attr)
  saveAttribute(attr)
}

function calcTypeDbString(attr: Attribute) {
  if (attr.length < 1) {
    return "text"
  }
  return `varchar(${attr.length})`
}

async function editNullable(attr: Attribute) {
  deleteDefaultNullIfNotNullable(attr)

  return saveAttribute(attr)
}

function deleteDefaultNullIfNotNullable(attr: Attribute) {
  if (attr.nullable && attr.default === "null") {
    return
  }
  attr.default = ""
}

function setDefaultValueByType(attr: Attribute) {
  attr.default = (() => {
    if (attr.nullable) {
      return "null"
    }

    switch (attr.type) {
      case "string":
        return "''"
      case "int":
        return "0"
      case "bool":
        return "false"
      case "json":
        return "{}"
      case "datetime":
        return "now()"
      case "uuid":
        return "uuid_generate_v4()"
      default:
        return ""
    }
  })()

  saveAttribute(attr)
}

function isForeignKey(attr: Attribute) {
  return attr.type === "one-to-one" || attr.type === "many-to-one"
}

function calcFkTypeByDbType(typeDb: string) {
  if (typeDb === "uuid") {
    return "uuid"
  }
  return "int"
}

function editForeignKey(attr: Attribute) {
  if (attr.fk === null) {
    return
  }

  for (const fkAttr of attr.fk.attributes) {
    if (!fkAttr.primaryKey) {
      continue
    }

    attr.typeDb = calcFkTypeByDbType(fkAttr.typeDb)
  }

  saveAttribute(attr)
}
</script>

<template>
  <main>
    <h2 v-if="!props.service.id">Select service</h2>
    <h2 v-else>Service: {{ props.service.name }}</h2>

    <button v-show="props.service.id" class="fixed-button" @click="addEntity()">
      <i class="fa-solid fa-plus"></i> Table
    </button>

    <div class="entity-container" v-for="entity in entities" :key="entity.id">
      Table:
      <input type="text" class="sticky" placeholder="table_name"
             v-model="entity.nameDb"
             @change="saveEntity(entity)"
      />

      <button @click="deleteEntity(entity)">
        <i class="fa-solid fa-ban"></i> Table
      </button>

      <div class="attribute-container" v-for="attr in entity.attributes" :key="attr.id">
        <input type="checkbox" hidden v-model="attr.primaryKey"/>

        <button :disabled="attr.primaryKey" @click="deleteAttribute(attr)">
          <i class="fa-solid fa-ban"></i>
        </button>

        <input type="text" placeholder="column_name"
               v-model="attr.nameDb"
               :disabled="attr.primaryKey && !isForeignKey(attr)"
               @change="saveAttribute(attr)"
        />

        <select v-model="attr.type" @change="editType(attr)">
          <option disabled value="">Type</option>
          <option v-show="!attr.primaryKey">string</option>
          <option v-show="!attr.primaryKey">int</option>
          <option v-show="!attr.primaryKey">bool</option>
          <option v-show="!attr.primaryKey">json</option>
          <option v-show="!attr.primaryKey">datetime</option>
          <option value="serial" v-show="attr.primaryKey">PK serial</option>
          <option value="uuid">{{ attr.primaryKey ? "PK uuid" : "uuid" }}</option>
          <option value="one-to-one">{{ attr.primaryKey ? "PK one-to-one" : "one-to-one" }}</option>
          <option v-show="!attr.primaryKey">many-to-one</option>
        </select>

        <label class="disabled"> db:
          <input class="sticky" type="text" :disabled="true" :value="attr.typeDb"/>
        </label>

        <label :class="{ 'disabled': attr.primaryKey }"> Null
          <input type="checkbox" class="sticky"
                 v-model="attr.nullable"
                 :disabled="attr.primaryKey"
                 @change="editNullable(attr)"
          />
        </label>

        <label :class="{ 'disabled': attr.primaryKey }"> Default:
          <input type="text" placeholder="<none>" class="sticky"
                 v-model="attr.default"
                 :disabled="attr.primaryKey"
                 @change="saveAttribute(attr)"
          />
        </label>
        <button class="sticky" :disabled="attr.primaryKey" @click="setDefaultValueByType(attr)">
          Default by type
        </button>

        <label v-show="attr.type === 'string'"> Len:
          <input type="number" placeholder="<text>" class="sticky" v-model="attr.length" @change="editLen(attr)"/>
          < 1 for text
        </label>

        <label v-show="isForeignKey(attr)"> FK table:
          <select class="sticky" v-model="attr.fk" @change="editForeignKey(attr)">
            <option disabled value="">FK table</option>
            <option v-for="fkEntity in entities" :value="fkEntity">{{ fkEntity.nameDb }}</option>
          </select>
        </label>
      </div>

      <div>
        <button class="sticky" @click="addAttribute(entity)">
          <i class="fa-solid fa-plus"></i> Column
        </button>
        <button
            v-show="!hasAttribute(entity, 'created_at')"
            @click="addAttribute(entity, {nameDb: 'created_at', typeDb: 'timestamp(0)', default: 'now()', type: 'datetime'})"
        >
          <i class="fa-solid fa-plus"></i> created_at
        </button>
        <button
            v-show="!hasAttribute(entity, 'updated_at')"
            @click="addAttribute(entity, {nameDb: 'updated_at', typeDb: 'timestamp(0)', default: 'null', nullable: true, type: 'datetime'})"
        >
          <i class="fa-solid fa-plus"></i> updated_at
        </button>
        <button
            v-show="!hasAttribute(entity, 'deleted_at')"
            @click="addAttribute(entity, {nameDb: 'deleted_at', typeDb: 'timestamp(0)', default: 'null', nullable: true, type: 'datetime'})"
        >
          <i class="fa-solid fa-plus"></i> deleted_at
        </button>
      </div>
    </div>
  </main>
</template>

<style scoped>
main {
  margin: 1rem;
  width: 100%;
  padding-right: 3rem;
}

.entity-container {
  padding: 0.5rem;
  margin-right: 2rem;
  margin-bottom: 2rem;
  background-color: whitesmoke;
  width: 100%;

  input, label, select, button {
    margin-left: 0.5rem;
  }

  .sticky {
    margin-left: 0;
  }
}

.attribute-container {
  margin: 0.5rem;
}

.fixed-button {
  position: fixed;
  top: 2rem; /* Отступ от верхнего края */
  right: 2rem; /* Отступ от правого края */
}

.disabled {
  color: gray;
}
</style>