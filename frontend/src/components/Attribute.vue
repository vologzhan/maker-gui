<script setup lang="ts">
import {defineProps} from 'vue';
import type {EntityDto} from "src/dto/entity.ts";
import type {AttributeDto} from "src/dto/attribute.ts";
import {DeleteAttribute} from "src/http/controller/attribute.ts";
import {addAttribute, calcFkTypeByDbType, saveAttribute} from "./attribute.ts";

const props = defineProps<{
  entities: EntityDto[]
  entity: EntityDto
}>();

async function deleteAttribute(attr: AttributeDto) {
  if (attr.id !== "") {
    await DeleteAttribute(attr.id)
  }

  attr.entity.attributes.splice(attr.entity.attributes.indexOf(attr), 1)
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

  return saveAttribute(attr, props.entities)
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

  await saveAttribute(attr, props.entities)
}

function editLen(attr: AttributeDto) {
  attr.typeDb = calcTypeDbString(attr)
  saveAttribute(attr, props.entities)
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

  return saveAttribute(attr, props.entities)
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

  saveAttribute(attr, props.entities)
}

function isForeignKey(attr: AttributeDto) {
  return attr.type === "one-to-one" || attr.type === "many-to-one"
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

  saveAttribute(attr, props.entities)
}
</script>

<template>
  <div class="attribute-item" v-for="attr in entity.attributes" :key="attr.id">
    <input type="checkbox" hidden v-model="attr.primaryKey"/>

    <button :disabled="attr.primaryKey" @click="deleteAttribute(attr)">
      <i class="fa-solid fa-ban"></i>
    </button>

    <input type="text" placeholder="column_name"
           v-model="attr.nameDb"
           :disabled="attr.primaryKey && !isForeignKey(attr)"
           @change="saveAttribute(attr, entities)"
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

    <label class="disabled">
      db
      <input type="text" :disabled="true" :value="attr.typeDb"/>
    </label>

    <label :class="{ 'disabled': attr.primaryKey }">
      null
      <input type="checkbox" v-model="attr.nullable" :disabled="attr.primaryKey" @change="editNullable(attr)"/>
    </label>

    <label :class="{ 'disabled': attr.primaryKey }">
      default
      <input type="text" placeholder="<none>" v-model="attr.default" :disabled="attr.primaryKey"
             @change="saveAttribute(attr, entities)"/>
    </label>

    <button :disabled="attr.primaryKey" @click="setDefaultValueByType(attr)">
      Default by type
    </button>

    <label v-show="attr.type === 'string'">
      len
      <input type="number" placeholder="<text>" v-model="attr.length" @change="editLen(attr)"/>
      < 1 for text
    </label>

    <label v-show="isForeignKey(attr)">
      FK table
      <select v-model="attr.fk" @change="editForeignKey(attr)">
        <option disabled value="">FK table</option>
        <option v-for="fkEntity in entities" :value="fkEntity">{{ fkEntity.nameDb }}</option>
      </select>
    </label>
  </div>

  <button @click="addAttribute(entity, entities)">
    <i class="fa-solid fa-plus"></i> Column
  </button>

  <button
      v-show="!hasAttribute(entity, 'created_at')"
      @click="addAttribute(entity, entities, {nameDb: 'created_at', typeDb: 'timestamp(0)', default: 'now()', type: 'datetime'})"
  >
    <i class="fa-solid fa-plus"></i> created_at
  </button>

  <button
      v-show="!hasAttribute(entity, 'updated_at')"
      @click="addAttribute(entity, entities, {nameDb: 'updated_at', typeDb: 'timestamp(0)', default: 'null', nullable: true, type: 'datetime'})"
  >
    <i class="fa-solid fa-plus"></i> updated_at
  </button>

  <button
      v-show="!hasAttribute(entity, 'deleted_at')"
      @click="addAttribute(entity, entities, {nameDb: 'deleted_at', typeDb: 'timestamp(0)', default: 'null', nullable: true, type: 'datetime'})"
  >
    <i class="fa-solid fa-plus"></i> deleted_at
  </button>
</template>

<style scoped>
.attribute-item {
  margin: 0.5rem;
}

input {
  margin-left: 0;
}
</style>