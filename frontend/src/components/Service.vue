<script setup lang="ts">
import {onMounted, ref} from "vue";
import {v4 as uuid} from "uuid";
import {CreateService, GetServiceList} from "src/http/controller/service.ts";
import type {ServiceDto} from "src/dto/service.ts";
import type {ServiceResponse} from "src/http/response/service.ts";

const services = ref<ServiceDto[]>([])
const selectedService = defineModel<ServiceDto | undefined>('service');
const newName = ref("")

onMounted(async () => {
  services.value = await getList();
});

async function getList() {
  const res = await GetServiceList();
  return res.items.map((service: ServiceResponse): ServiceDto => {
    return {id: service.id, name: service.name};
  });
}

async function create() {
  const name = newName.value.trim()
  if (name === "") {
    return
  }

  const service: ServiceDto = {id: uuid(), name: name}
  await CreateService(service)

  services.value.push(service)
  selectedService.value = service
  newName.value = ""
}
</script>

<template>
  <ul>
    <li class="clickable"
        v-for="service in services"
        :key="service.id"
        :class="{ selected: service.id === selectedService?.id }"
        @click="selectedService = service"
    >
      {{ service.name }}
    </li>
  </ul>

  <form @submit.prevent="create()">
    <input placeholder="<service-name>" v-model="newName"/>
    <button type="submit">
      <i class="fa-solid fa-plus"></i> Service
    </button>
  </form>
</template>

<style scoped>
ul {
  padding: 0;
  margin: 0;
  list-style-type: none;
}

.selected {
  font-weight: bold;
}
</style>
