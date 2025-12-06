<script setup lang="ts">
import {defineProps, ref, watch} from 'vue';
import {v4 as uuid} from 'uuid';
import type {Service} from "../http/service.ts";
import {createEntityHttp, deleteEntityHttp, type Entity, getEntityListHttp, updateEntityHttp} from "../http/entity.ts";
import {type Attribute, createAttributeHttp, deleteAttributeHttp, updateAttributeHttp} from "../http/attribute.ts";

const entities = ref<Entity[]>([])

const props = defineProps<{
  service: Service
}>();

watch(() => props.service.id, async (serviceId: string) => {
  entities.value = await getEntityList(serviceId)
});

async function getEntityList(serviceId: string) {
  const entityList = await getEntityListHttp(serviceId)

  for (const entity of entityList) {
    for (const attribute of entity.attributes) {
      if (attribute.primary_key) {
        attribute.type = attribute.type_db === "uuid" ? "PK uuid" : "PK serial"
      } else if (attribute.fk_table) {
        attribute.type = attribute.fk_type
      } else if (attribute.type_db === "timestamp(0)") {
        attribute.type = "datetime"
      } else if (attribute.type_db === "text") {
        attribute.type = "string"
        attribute.length = 0
      } else if (attribute.type_db.startsWith("varchar")) {
        attribute.type = "string"
        attribute.length = Number(
            attribute.type_db.substring(attribute.type_db.indexOf('(') + 1, attribute.type_db.indexOf(')'))
        )
      } else {
        attribute.type = attribute.type_db
      }
    }
  }

  return entityList
}

function addEntity() {
  const entity = {
    id: "",
    name_db: "",
    attributes: [],
  }
  entities.value.push(entity)
  addAttribute(entity, {
    name_db: "uuid",
    type: "PK uuid",
    type_db: "uuid",
    default: "uuid_generate_v4()",
    primary_key: true,
  })
}

async function saveEntity(entity: Entity) {
  if (entity.name_db.trim() === "") {
    return
  }

  if (entity.id !== "") {
    await updateEntityHttp(entity)
    return
  }

  entity.id = uuid()
  await createEntityHttp(props.service.id, entity)

  for (const attr of entity.attributes) {
    await saveAttribute(entity.id, attr)
  }
}

async function deleteEntity(entity: Entity) {
  entities.value.splice(entities.value.indexOf(entity), 1);

  if (entity.id !== "") {
    await deleteEntityHttp(entity.id)
  }
}

function addAttribute(entity: Entity, options?: Partial<Attribute>) {
  const defaults: Attribute = {
    id: "",
    name_db: "",
    type_db: "",
    default: "",
    fk_table: "",
    fk_type: "",
    nullable: false,
    primary_key: false,
    type: "",
    length: 0,
  }

  const attr: Attribute = {
    ...defaults,
    ...options,
  }

  entity.attributes.push(attr)
  saveAttribute(entity.id, attr)
}

function hasAttribute(entity: Entity, nameDb: string) {
  return entity.attributes.some(attribute => attribute.name_db === nameDb)
}

async function saveAttribute(entityId: string, attribute: Attribute) {
  if (entityId === "" || attribute.name_db === "" || attribute.type_db === "") {
    return
  }

  if (attribute.id == "") {
    attribute.id = uuid()
    await createAttributeHttp(entityId, attribute)
  } else {
    await updateAttributeHttp(attribute)
  }
}

function editType(entityId: string, attribute: Attribute) {
  if (attribute.type === "PK uuid") {
    attribute.name_db = "uuid"
    attribute.type_db = "uuid"
    attribute.default = "uuid_generate_v4()"
  } else if (attribute.type === "PK serial") {
    attribute.name_db = "id"
    attribute.type_db = "serial"
    attribute.default = ""
  } else if (attribute.type === "datetime") {
    attribute.type_db = "timestamp(0)"
  } else if (attribute.type === "string") {
    attribute.length = 255
    calcTypeString(attribute)
  } else if (attribute.type === "one-to-one" || attribute.type === "many-to-one") {
    attribute.fk_type = attribute.type
  } else if (attribute.type === "PK one-to-one") {
    attribute.fk_type = "one-to-one"
  } else {
    attribute.type_db = attribute.type
  }

  saveAttribute(entityId, attribute)
}

function editLen(entityId: string, attribute: Attribute) {
  calcTypeString(attribute)
  saveAttribute(entityId, attribute)
}

function calcTypeString(attribute: Attribute) {
  if (attribute.length < 1) {
    attribute.type_db = "text"
  } else (
      attribute.type_db = `varchar(${attribute.length})`
  )
}

async function deleteAttribute(entity: Entity, attribute: Attribute) {
  entity.attributes.splice(entity.attributes.indexOf(attribute), 1)

  if (attribute.id !== "") {
    await deleteAttributeHttp(attribute.id)
  }
}

function setDefaultValueByType(entityId: string, attribute: Attribute) {
  attribute.default = (() => {
    if (attribute.nullable) {
      return "null"
    }

    switch (attribute.type) {
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

  saveAttribute(entityId, attribute)
}

function isForeignKey(attr: Attribute) {
  return attr.type === "one-to-one" || attr.type === "many-to-one" || attr.type === "PK one-to-one"
}

function editForeignKey(entityId: string, attr: Attribute) {
  for (const fkEntity of entities.value) {
    if (fkEntity.name_db !== attr.fk_table) {
      continue
    }

    for (const fkAttr of fkEntity.attributes) {
      if (!fkAttr.primary_key) {
        continue
      }

      attr.type_db = fkAttr.type_db === "uuid" ? "uuid" : "int"

      break
    }

    break
  }

  saveAttribute(entityId, attr)
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
             v-model="entity.name_db"
             @change="saveEntity(entity)"
      />

      <button @click="deleteEntity(entity)">
        <i class="fa-solid fa-ban"></i> Table
      </button>

      <div class="attribute-container" v-for="attribute in entity.attributes" :key="attribute.id">
        <input type="checkbox" hidden v-model="attribute.primary_key"/>

        <button :disabled="attribute.primary_key" @click="deleteAttribute(entity, attribute)">
          <i class="fa-solid fa-ban"></i>
        </button>

        <input type="text" placeholder="column_name"
               v-model="attribute.name_db"
               :disabled="attribute.primary_key"
               @change="saveAttribute(entity.id, attribute)"
        />

        <select v-model="attribute.type" @change="editType(entity.id, attribute)">
          <option disabled value="">Type</option>
          <option v-show="!attribute.primary_key">string</option>
          <option v-show="!attribute.primary_key">int</option>
          <option v-show="!attribute.primary_key">bool</option>
          <option v-show="!attribute.primary_key">json</option>
          <option v-show="!attribute.primary_key">datetime</option>
          <option>{{ attribute.primary_key ? "PK uuid" : "uuid" }}</option>
          <option v-show="attribute.primary_key">PK serial</option>
          <option>{{ attribute.primary_key ? "PK one-to-one" : "one-to-one" }}</option>
          <option v-show="!attribute.primary_key">many-to-one</option>
        </select>

        <label class="disabled"> db:
          <input class="sticky" type="text" :disabled="true" :value="attribute.type_db"/>
        </label>

        <label :class="{ 'disabled': attribute.primary_key }"> Null
          <input type="checkbox" class="sticky"
                 v-model="attribute.nullable"
                 :disabled="attribute.primary_key"
                 @change="saveAttribute(entity.id, attribute)"
          />
        </label>

        <label :class="{ 'disabled': attribute.primary_key }"> Default:
          <input type="text" placeholder="<none>" class="sticky"
                 v-model="attribute.default"
                 :disabled="attribute.primary_key"
                 @change="saveAttribute(entity.id, attribute)"
          />
        </label>
        <button class="sticky" :disabled="attribute.primary_key" @click="setDefaultValueByType(entity.id, attribute)">
          Default by type
        </button>

        <label v-show="attribute.type === 'string'" @change="editLen(entity.id, attribute)"> Len:
          <input type="number" placeholder="<text>" class="sticky" v-model="attribute.length"
                 @change="saveAttribute(entity.id, attribute)"/>
          < 1 for text
        </label>

        <label v-show="isForeignKey(attribute)" @change="editForeignKey(entity.id, attribute)"> FK table:
          <select class="sticky" v-model="attribute.fk_table">
            <option disabled value="">FK table</option>
            <option v-for="fkEntity in entities">{{ fkEntity.name_db }}</option>
          </select>
        </label>
      </div>

      <div>
        <button class="sticky" @click="addAttribute(entity)">
          <i class="fa-solid fa-plus"></i> Column
        </button>
        <button
            v-show="!hasAttribute(entity, 'created_at')"
            @click="addAttribute(entity, {name_db: 'created_at', type: 'datetime', type_db: 'timestamp(0)', default: 'now()'})"
        >
          <i class="fa-solid fa-plus"></i> created_at
        </button>
        <button
            v-show="!hasAttribute(entity, 'updated_at')"
            @click="addAttribute(entity, {name_db: 'updated_at', type: 'datetime', type_db: 'timestamp(0)', default: 'null', nullable: true})"
        >
          <i class="fa-solid fa-plus"></i> updated_at
        </button>
        <button
            v-show="!hasAttribute(entity, 'deleted_at')"
            @click="addAttribute(entity, {name_db: 'deleted_at', type: 'datetime', type_db: 'timestamp(0)', default: 'null', nullable: true})"
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