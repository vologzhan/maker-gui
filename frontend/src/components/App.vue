<script setup lang="ts">
import {ref} from "vue";
import ServiceList from "./ServiceList.vue";
import Service from "./Service.vue";
import type {ServiceDto} from "src/dto/service.ts";

const service = ref<ServiceDto | undefined>()
</script>

<template>
  <ServiceList v-model:service="service"/>
  <Service :service="service"/>
</template>

<style scoped>
</style>
