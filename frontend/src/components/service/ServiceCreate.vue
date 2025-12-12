<script setup lang="ts">
import {ref} from "vue";
import {v4 as uuid} from "uuid";
import type {ServiceDto} from "src/dto/service.ts";
import {CreateService} from "src/http/controller/service.ts";

const emit = defineEmits<{
  created: [service: ServiceDto];
}>();

const newName = ref("")

async function create() {
  const name = newName.value.trim()
  if (name === "") {
    return
  }

  const service: ServiceDto = {id: uuid(), name: name}
  await CreateService(service)

  emit('created', service)
  newName.value = ""
}
</script>

<template>
  <form @submit.prevent="create()">
    <input placeholder="<service-name>" v-model="newName"/>
    <button type="submit">
      <i class="fa-solid fa-plus"></i> Service
    </button>
  </form>
</template>

<style scoped>
</style>
