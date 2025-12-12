<script setup lang="ts">
import {onMounted} from "vue";
import type {ServiceDto} from "src/dto/service.ts";
import type {ServiceResponse} from "src/http/response/service.ts";
import {GetServiceList} from "src/http/controller/service.ts";

const services = defineModel<ServiceDto[]>('services');
const service = defineModel<ServiceDto | undefined>('service');

onMounted(async () => {
  services.value = await getList();
});

async function getList() {
  const res = await GetServiceList();
  return res.items.map((serv: ServiceResponse): ServiceDto => {
    return {id: serv.id, name: serv.name};
  });
}
</script>

<template>
  <ul>
    <li class="clickable"
        v-for="serv in services"
        :key="serv.id"
        :class="{ selected: serv.id === service?.id }"
        @click="service = serv"
    >
      {{ serv.name }}
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
