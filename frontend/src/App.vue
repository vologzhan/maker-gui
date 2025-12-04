<script setup lang="ts">
import {ref} from 'vue'
import Service from "./components/Service.vue"
import Entity from "./components/Entity.vue"
import {type Service as ServiceDto} from './http/service'

const currentService = ref<ServiceDto>({
  id: "",
  name: "",
})
</script>

<template>
  <div class="app flex-container">
    <Service @service-selected="currentService = $event"></Service>
    <Entity :service="currentService"></Entity>
  </div>
</template>

<style scoped>
</style>
