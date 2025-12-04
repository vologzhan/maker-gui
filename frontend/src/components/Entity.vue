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
  entities.value = await getEntityListHttp(serviceId)
});

function addEntity() {
  const entity = {
    id: "",
    name_db: "",
    attributes: [],
  }
  entities.value.push(entity)
  addAttribute(entity, {name_db: "uuid", type_db: "uuid", default: "uuid_generate_v4()", primary_key: true})
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
    fk_column: "",
    nullable: false,
    primary_key: false,
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

    switch (attribute.type_db) {
      case "string":
      case "text":
        return "''"
      case "int":
        return "0"
      case "bool":
        return "false"
      case "json":
        return "{}"
      case "uuid":
        return "uuid_generate_v4()"
      case "timestamp(0)":
        return "now()"
      case "serial":
      default:
        return ""
    }
  })()

  saveAttribute(entityId, attribute)
}

function changePrimaryKey(entityId: string, attribute: Attribute) {
  if (attribute.type_db === "uuid") {
    attribute.name_db = "uuid"
    attribute.default = "uuid_generate_v4()"
  } else {
    attribute.name_db = "id"
    attribute.default = ""
  }

  saveAttribute(entityId, attribute)
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
        <i class="fa-solid fa-ban"></i> Delete table
      </button>

      <div class="attribute-container" v-for="attribute in entity.attributes" :key="attribute.id">
        <input type="checkbox" hidden v-model="attribute.primary_key"/>
        <label class="sticky pk-label" v-show="attribute.primary_key">PK:</label>

        <button class="sticky" v-show="!attribute.primary_key" @click="deleteAttribute(entity, attribute)">
          <i class="fa-solid fa-ban"></i>
        </button>

        <input type="text" placeholder="column_name"
               :disabled="attribute.primary_key"
               v-model="attribute.name_db"
               @change="saveAttribute(entity.id, attribute)"
        />

        <select v-if="attribute.primary_key" v-model="attribute.type_db"
                @change="changePrimaryKey(entity.id, attribute)">
          <option disabled value="">Type</option>
          <option>uuid</option>
          <option>serial</option>
        </select>

        <select v-else v-model="attribute.type_db" @change="saveAttribute(entity.id, attribute)">
          <option disabled value="">Type</option>
          <option>string</option>
          <option>text</option>
          <option>int</option>
          <option>bool</option>
          <option>json</option>
          <option>uuid</option>
          <option value="timestamp(0)">datetime</option>
        </select>

        <label v-show="!attribute.primary_key">
          Null
          <input type="checkbox" class="sticky" v-model="attribute.nullable"
                 @change="saveAttribute(entity.id, attribute)"
          />
        </label>

        <input type="text" placeholder="default"
               :disabled="attribute.primary_key"
               v-model="attribute.default"
               @change="saveAttribute(entity.id, attribute)"
        />
        <button class="sticky"
                v-show="!attribute.primary_key"
                @click="setDefaultValueByType(entity.id, attribute)"
        >
          Default by type
        </button>

        <input type="text" placeholder="foreign_table"
               v-show="!attribute.primary_key"
               v-model="attribute.fk_table"
               @change="saveAttribute(entity.id, attribute)"
        />
        <input type="text" placeholder="foreign_column"
               v-show="!attribute.primary_key"
               v-model="attribute.fk_column"
               @change="saveAttribute(entity.id, attribute)"
        />
      </div>

      <div>
        <button class="sticky" @click="addAttribute(entity)">
          <i class="fa-solid fa-plus"></i> Column
        </button>
        <button
            :disabled="hasAttribute(entity, 'created_at')"
            @click="addAttribute(entity, {name_db: 'created_at', type_db: 'timestamp(0)', default: 'now()'})"
        >
          <i class="fa-solid fa-plus"></i> created_at
        </button>
        <button
            :disabled="hasAttribute(entity, 'updated_at')"
            @click="addAttribute(entity, {name_db: 'updated_at', type_db: 'timestamp(0)', default: 'null', nullable: true})"
        >
          <i class="fa-solid fa-plus"></i> updated_at
        </button>
        <button
            :disabled="hasAttribute(entity, 'deleted_at')"
            @click="addAttribute(entity, {name_db: 'deleted_at', type_db: 'timestamp(0)', default: 'null', nullable: true})"
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

.pk-label {
  margin-right: 0.5rem;
}
</style>