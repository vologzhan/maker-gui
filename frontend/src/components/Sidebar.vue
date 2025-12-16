<script setup lang="ts">
import Service from "./service/Service.vue";
import type {ServiceDto} from "src/dto/service.ts";

const service = defineModel<ServiceDto | undefined>('service')
</script>

<template>
  <aside>
    <Service v-model:service="service"/>
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
