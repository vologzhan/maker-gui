<script setup lang="ts">
import {ref} from "vue";
import {v4 as uuid} from "uuid";
import type {ServiceDto} from "src/dto/service.ts";
import {CreateService} from "src/http/controller/service.ts";

const services = defineModel<ServiceDto[]>('services', {required: true});
const service = defineModel<ServiceDto | undefined>('service');

const newName = ref("")

async function create() {
  const name = newName.value.trim()
  if (name === "") {
    return
  }

  const serv: ServiceDto = {id: uuid(), name: name}
  await CreateService(serv)

  services.value.push(serv)
  service.value = serv

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
