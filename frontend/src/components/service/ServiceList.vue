<script setup lang="ts">
import {onMounted} from "vue";
import type {ServiceDto} from "src/dto/service.ts";
import type {ServiceResponse} from "src/http/response/service.ts";
import {GetServiceList} from "src/http/controller/service.ts";

const services = defineModel<ServiceDto[]>('services', {default: []});
const selectedService = defineModel<ServiceDto | undefined>('service');

onMounted(async () => {
  services.value = await getList();
});

async function getList() {
  const res = await GetServiceList();
  return res.items.map((service: ServiceResponse): ServiceDto => {
    return {id: service.id, name: service.name};
  });
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
