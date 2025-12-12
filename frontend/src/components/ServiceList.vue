<script setup lang="ts">
import {onMounted, ref} from "vue";
import {v4 as uuid} from "uuid";
import type {ServiceDto} from "src/dto/service.ts";
import type {ServiceResponse} from "src/http/response/service.ts";
import {CreateService, GetServiceList} from "src/http/controller/service.ts";

const services = ref<ServiceDto[]>([])
const selectedService = defineModel<ServiceDto | undefined>('service')
const newServiceName = ref("")

onMounted(async () => {
  services.value = await getList()
})

async function getList() {
  const res = await GetServiceList()

  return res.items.map((service: ServiceResponse): ServiceDto => {
    return {id: service.id, name: service.name}
  })
}

async function create() {
  const name = newServiceName.value.trim()
  if (name === "") {
    return
  }

  const service: ServiceDto = {id: uuid(), name: name}
  await CreateService(service)

  services.value.push(service)
  selectedService.value = service
  newServiceName.value = ""
}
</script>

<template>
  <aside>
    <ul>
      <li
          v-for="service in services"
          :key="service.id"
          class="clickable"
          :class="{ selected: service.id === selectedService?.id }"
          @click="selectedService = service"
      >
        {{ service.name }}
      </li>
    </ul>

    <form @submit.prevent="create()">
      <input placeholder="<service-name>" v-model="newServiceName"/>
      <button type="submit">
        <i class="fa-solid fa-plus"></i> Service
      </button>
    </form>
  </aside>
</template>

<style scoped>
aside {
  grid-column: 1;
  grid-row: 1 / 3;
  padding: 0.3rem;
  overflow: hidden;
  background-color: whitesmoke;

  ul {
    padding: 0;
    margin: 0;
    list-style-type: none;
  }

  li.selected {
    font-weight: bold;
  }
}
</style>
