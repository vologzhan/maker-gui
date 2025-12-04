<script setup lang="ts">
import {onMounted, ref} from "vue"
import {createServiceHttp, getServiceListHttp, type Service} from "../http/service"
import {v4 as uuid} from 'uuid'

const emit = defineEmits<{
  (e: 'service-selected', service: Service): void
}>()

const services = ref<Service[]>([])
const newServiceName = ref("")

onMounted(async () => {
  services.value = await getServiceListHttp()
})

async function createService() {
  const name = newServiceName.value.trim()
  if (name === "") {
    return
  }

  const service: Service = {
    id: uuid(),
    name: name,
  }

  await createServiceHttp(service)

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
