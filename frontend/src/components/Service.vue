<script setup lang="ts">
import {onMounted, ref} from "vue"
import {v4 as uuid} from 'uuid'
import {CreateService, GetServiceList} from "../http/controller/service.ts";
import type {Service} from "../dto/service.ts";
import type {ServiceResponse} from "../http/response/service.ts";

const services = ref<Service[]>([])
const newServiceName = ref("")

const emit = defineEmits<{
  (e: 'service-selected', service: Service): void
}>()

onMounted(async () => {
  services.value = await getServiceList()
})

async function getServiceList() {
  const res = await GetServiceList()

  return res.items.map((service: ServiceResponse): Service => {
    return {id: service.id, name: service.name}
  })
}

async function createService() {
  const name = newServiceName.value.trim()
  if (name === "") {
    return
  }

  const service: Service = {id: uuid(), name: name}
  await CreateService(service)
  services.value.push(service)
  newServiceName.value = ""
}
</script>

<template>
  <aside>
    <div
        class="clickable"
        v-for="service in services"
        :key="service.id"
        @click="emit('service-selected', service)"
    >
      {{ service.name }}
    </div>

    <div>
      <input placeholder="service-name" v-model="newServiceName"/>
      <button @click="createService()">
        <i class="fa-solid fa-plus"></i> Service
      </button>
    </div>
  </aside>
</template>

<style scoped>
aside {
  float: left;
  width: 15%;
  background-color: whitesmoke;
  padding: 0.5rem;
  height: 100vh;
  position: sticky;
  top: 0;
}
</style>
