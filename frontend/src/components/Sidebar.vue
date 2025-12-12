<script setup lang="ts">
import {ref} from "vue";
import ServiceList from "./service/ServiceList.vue";
import ServiceCreate from "./service/ServiceCreate.vue";
import type {ServiceDto} from "src/dto/service.ts";

const services = ref<ServiceDto[]>([])
const service = defineModel<ServiceDto | undefined>('service')
</script>

<template>
  <aside>
    <ServiceList v-model:services="services" v-model:service="service"/>
    <ServiceCreate v-model:services="services" v-model:service="service"/>
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
