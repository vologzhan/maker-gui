<script setup lang="ts">
import {defineProps, onMounted, ref, watch} from 'vue';
import {v4 as uuid} from 'uuid';
import pluralize from 'pluralize';
import type {ServiceDto} from "src/dto/service.ts";
import type {EntityDto} from "src/dto/entity.ts";
import type {AttributeDto} from "src/dto/attribute.ts";
import {CreateEntity, DeleteEntity, GetEntityList, UpdateEntity} from "src/http/controller/entity.ts";
import {CreateAttribute, DeleteAttribute, UpdateAttribute} from "src/http/controller/attribute.ts";
import type {EntityResponse} from "src/http/response/entity.ts";
import type {AttributeResponse} from "src/http/response/attribute.ts";

const entities = ref<EntityDto[]>([])

const props = defineProps<{
  service: ServiceDto
}>();

watch(() => props.service, async (service: ServiceDto) => {
  entities.value = await getEntityList(service.id)
});

onMounted(async () => {
  entities.value = await getEntityList(props.service.id);
});

async function saveEntity(entity: EntityDto) {
  if (entity.nameDb === "" || entity.namePlural === "") {
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

async function createEntity(entity: EntityDto) {
  entity.id = uuid()

  await CreateEntity(props.service.id, {
    id: entity.id,
    name_db: entity.nameDb,
    name_plural: entity.namePlural,
  })

  for (const attr of entity.attributes) {
    await saveAttribute(attr)
  }
}

async function updateEntity(entity: EntityDto) {
  await UpdateEntity(entity.id, {
    name_db: entity.nameDb,
    name_plural: entity.namePlural,
  })
}

async function deleteEntity(entity: EntityDto) {
  if (entity.id !== "") {
    await DeleteEntity(entity.id)
  }

  entities.value.splice(entities.value.indexOf(entity), 1);

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

  const relations = new Map<string, AttributeDto[]>

  const entities = res.items.map((entityRes: EntityResponse) => {
    const entity: EntityDto = {
      id: entityRes.id,
      nameDb: entityRes.name_db,
      namePlural: entityRes.name_plural,
      namePluralAuto: entityRes.name_plural === calcNamePlural(entityRes.name_db),
      attributes: [],
    }

    entity.attributes = entityRes.attributes.map((attrRes: AttributeResponse): AttributeDto => {
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

function convertAttributeResponseToDto(entity: EntityDto, res: AttributeResponse) {
  const attr: AttributeDto = {
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

async function saveAttribute(attr: AttributeDto) {
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

async function deleteAttribute(attr: AttributeDto) {
  if (attr.id !== "") {
    await DeleteAttribute(attr.id)
  }

  attr.entity.attributes.splice(attr.entity.attributes.indexOf(attr), 1)
}

function addEntity() {
  const entity: EntityDto = {id: "", nameDb: "", namePlural: "", namePluralAuto: true, attributes: []}
  entities.value.push(entity)

  addAttribute(entity, {
    nameDb: "id",
    typeDb: "serial",
    default: "",
    primaryKey: true,
    type: "serial",
  })
}

async function changeEntityName(entity: EntityDto) {
  entity.nameDb = entity.nameDb.trim()
  if (entity.nameDb === "") {
    return
  }

  if (entity.namePluralAuto) {
    entity.namePlural = calcNamePlural(entity.nameDb)
  }

  return saveEntity(entity)
}

async function changeEntityNamePlural(entity: EntityDto) {
  entity.namePlural = entity.namePlural.trim()
  if (entity.namePlural === "") {
    return
  }

  return saveEntity(entity)
}

async function changeEntityNamePluralAuto(entity: EntityDto) {
  if (!entity.namePluralAuto) {
    return
  }
  entity.namePlural = calcNamePlural(entity.nameDb)

  return saveEntity(entity)
}

function calcNamePlural(name: string): string {
  const plural = pluralize(name)
  return plural === name ? `${name}_list` : plural
}

async function addAttribute(entity: EntityDto, options?: Partial<AttributeDto>) {
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

  return saveAttribute(attr)
}

function hasAttribute(entity: EntityDto, nameDb: string) {
  return entity.attributes.some(attr => attr.nameDb === nameDb)
}

async function editType(attr: AttributeDto) {
  attr.length = 0
  if (!attr.nullable || attr.default !== "null") {
    attr.default = ""
  }

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

async function editPrimaryKey(attr: AttributeDto) {
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

function editLen(attr: AttributeDto) {
  attr.typeDb = calcTypeDbString(attr)
  saveAttribute(attr)
}

function calcTypeDbString(attr: AttributeDto) {
  if (attr.length < 1) {
    return "text"
  }
  return `varchar(${attr.length})`
}

async function editNullable(attr: AttributeDto) {
  if (!attr.nullable && attr.default === "null") {
    attr.default = ""
  }

  return saveAttribute(attr)
}

function setDefaultValueByType(attr: AttributeDto) {
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

function isForeignKey(attr: AttributeDto) {
  return attr.type === "one-to-one" || attr.type === "many-to-one"
}

function calcFkTypeByDbType(typeDb: string) {
  if (typeDb === "uuid") {
    return "uuid"
  }
  return "int"
}

function editForeignKey(attr: AttributeDto) {
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
    <div class="entity-container" v-for="entity in entities" :key="entity.id">
      <button @click="deleteEntity(entity)">
        <i class="fa-solid fa-ban"></i>
      </button>

      <label> table
        <input type="text" class="sticky" placeholder="<table_name>"
               v-model="entity.nameDb"
               @change="changeEntityName(entity)"
        />
      </label>

      <label :class="{ 'disabled': entity.namePluralAuto }"> plural
        <input type="checkbox" class="sticky" title="auto"
               v-model="entity.namePluralAuto"
               @change="changeEntityNamePluralAuto(entity)"
        />
        <input type="text" class="sticky" placeholder="<plural_name>"
               v-model="entity.namePlural"
               :disabled="entity.namePluralAuto"
               @change="changeEntityNamePlural(entity)"
        />
      </label>

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

        <label class="disabled"> db
          <input class="sticky" type="text" :disabled="true" :value="attr.typeDb"/>
        </label>

        <label :class="{ 'disabled': attr.primaryKey }"> null
          <input type="checkbox" class="sticky"
                 v-model="attr.nullable"
                 :disabled="attr.primaryKey"
                 @change="editNullable(attr)"
          />
        </label>

        <label :class="{ 'disabled': attr.primaryKey }"> default
          <input type="text" placeholder="<none>" class="sticky"
                 v-model="attr.default"
                 :disabled="attr.primaryKey"
                 @change="saveAttribute(attr)"
          />
        </label>
        <button class="sticky" :disabled="attr.primaryKey" @click="setDefaultValueByType(attr)">
          Default by type
        </button>

        <label v-show="attr.type === 'string'"> len
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

    <button v-show="props.service.id" @click="addEntity()">
      <i class="fa-solid fa-plus"></i> Table
    </button>
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

.disabled {
  color: gray;
}
</style>