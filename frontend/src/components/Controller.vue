<script setup lang="ts">
import type {ServiceDto} from "src/dto/service.ts";

const {service} = defineProps<{
  service: ServiceDto,
}>()
</script>

<template>
  Controller: {{service}}
</template>

<style scoped>
</style>
