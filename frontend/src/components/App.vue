<script setup lang="ts">
import {ref} from "vue";
import Sidebar from "./Sidebar.vue";
import Main from "./AppMain.vue";
import type {ServiceDto} from "src/dto/service.ts";

const service = ref<ServiceDto | undefined>()
</script>

<template>
  <Sidebar v-model:service="service"/>
  <Main :service="service"/>
</template>

<style scoped>
</style>
