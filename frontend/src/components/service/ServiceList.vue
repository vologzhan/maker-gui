<script setup lang="ts">
import type {ServiceDto} from "src/dto/service.ts";

defineProps<{
  services: ServiceDto[],
}>()

const selectedService = defineModel<ServiceDto | undefined>('service')
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
