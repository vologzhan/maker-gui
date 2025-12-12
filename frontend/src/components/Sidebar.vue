<script setup lang="ts">
import {onMounted, ref} from "vue";
import ServiceList from "./service/ServiceList.vue";
import ServiceCreate from "./service/ServiceCreate.vue";
import {GetServiceList} from "src/http/controller/service.ts";
import type {ServiceResponse} from "src/http/response/service.ts";
import type {ServiceDto} from "src/dto/service.ts";

const services = ref<ServiceDto[]>([])
const selectedService = defineModel<ServiceDto | undefined>('service')

onMounted(async () => {
  services.value = await getList()
})

function onCreated(service: ServiceDto) {
  services.value.push(service);
  selectedService.value = service;
}

async function getList() {
  const res = await GetServiceList()
  return res.items.map((service: ServiceResponse): ServiceDto => {
    return {id: service.id, name: service.name}
  })
}
</script>

<template>
  <aside>
    <ServiceList :services="services" v-model:service="selectedService"/>
    <ServiceCreate @created="onCreated"/>
  </aside>
</template>

<style scoped>
aside {
  grid-column: 1;
  grid-row: 1 / 3;
  padding: 0.3rem;
  overflow: hidden;
  background-color: whitesmoke;
}
</style>
