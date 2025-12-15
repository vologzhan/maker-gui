<script setup lang="ts">
import {defineProps, onMounted, ref, watch} from 'vue';
import {v4 as uuid} from 'uuid';
import pluralize from 'pluralize';
import type {ServiceDto} from "src/dto/service.ts";
import type {EntityDto} from "src/dto/entity.ts";
import type {AttributeDto} from "src/dto/attribute.ts";
import {CreateEntity, DeleteEntity, GetEntityList, UpdateEntity} from "src/http/controller/entity.ts";
import type {EntityResponse} from "src/http/response/entity.ts";
import type {AttributeResponse} from "src/http/response/attribute.ts";
import Attribute from "./Attribute.vue";
import {addAttribute, saveAttribute} from "./attribute.ts";

const entities = ref<EntityDto[]>([])

const props = defineProps<{
  service: ServiceDto
}>();

watch(() => props.service, async (service: ServiceDto) => {
  entities.value = await getList(service.id)
});

onMounted(async () => {
  entities.value = await getList(props.service.id);
});

async function save(entity: EntityDto) {
  if (entity.nameDb === "" || entity.namePlural === "") {
    return
  }

  if (entity.id === "") {
    await create(entity)
  } else {
    await update(entity)
  }

  for (const fkEntity of entities.value) {
    for (const fkAttr of fkEntity.attributes) {
      if (fkAttr.fk !== entity) {
        continue
      }

      await saveAttribute(fkAttr, entities.value)
    }
  }
}

async function create(entity: EntityDto) {
  entity.id = uuid()

  await CreateEntity(props.service.id, {
    id: entity.id,
    name_db: entity.nameDb,
    name_plural: entity.namePlural,
  })

  for (const attr of entity.attributes) {
    await saveAttribute(attr, entities.value)
  }
}

async function update(entity: EntityDto) {
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

      await saveAttribute(fkAttr, entities.value)
    }
  }
}

async function getList(serviceId: string) {
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

function add() {
  const entity: EntityDto = {id: "", nameDb: "", namePlural: "", namePluralAuto: true, attributes: []}
  entities.value.push(entity)

  addAttribute(entity, entities.value, {
    nameDb: "id",
    typeDb: "serial",
    default: "",
    primaryKey: true,
    type: "serial",
  })
}

async function changeName(entity: EntityDto) {
  entity.nameDb = entity.nameDb.trim()
  if (entity.nameDb === "") {
    return
  }

  if (entity.namePluralAuto) {
    entity.namePlural = calcNamePlural(entity.nameDb)
  }

  return save(entity)
}

async function changeNamePlural(entity: EntityDto) {
  entity.namePlural = entity.namePlural.trim()
  if (entity.namePlural === "") {
    return
  }

  return save(entity)
}

async function changeNamePluralAuto(entity: EntityDto) {
  if (!entity.namePluralAuto) {
    return
  }
  entity.namePlural = calcNamePlural(entity.nameDb)

  return save(entity)
}

function calcNamePlural(name: string): string {
  const plural = pluralize(name)
  return plural === name ? `${name}_list` : plural
}
</script>

<template>
  <div class="entity-item" v-for="entity in entities" :key="entity.id">
    <button @click="deleteEntity(entity)">
      <i class="fa-solid fa-ban"></i>
    </button>

    <label>
      table
      <input type="text" placeholder="<table_name>" v-model="entity.nameDb" @change="changeName(entity)"/>
    </label>

    <label :class="{ 'disabled': entity.namePluralAuto }">
      plural
      <input type="checkbox" title="auto"
             v-model="entity.namePluralAuto"
             @change="changeNamePluralAuto(entity)"
      />
      <input type="text" placeholder="<plural_name>"
             v-model="entity.namePlural"
             :disabled="entity.namePluralAuto"
             @change="changeNamePlural(entity)"
      />
    </label>

    <Attribute :entities="entities" :entity="entity"/>
  </div>

  <button v-show="props.service.id" @click="add()">
    <i class="fa-solid fa-plus"></i> Table
  </button>
</template>

<style scoped>
.entity-item {
  margin-bottom: 1rem;
  padding: 0.3rem;
  background-color: whitesmoke;
}

input {
  margin-left: 0;
}
</style>