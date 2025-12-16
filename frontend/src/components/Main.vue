<script setup lang="ts">
import {shallowRef} from "vue";
import type {ServiceDto} from "src/dto/service.ts";
import Controller from "./controller/Controller.vue";
import Entity from "./entity/Entity.vue";

const {service} = defineProps<{
  service?: ServiceDto,
}>()

const tabs = [
  {name: 'Database', component: Entity},
  {name: 'Controllers', component: Controller},
] as const

const selectedTab = shallowRef<(typeof tabs)[number]>(tabs[0])
</script>

<template>
  <header>
    <h1 v-if="!service">Select service</h1>
    <ul v-else>
      <li class="clickable"
          v-for="tab in tabs"
          :key="tab.name"
          :class="{ selected: selectedTab === tab }"
          @click="selectedTab = tab"
      >
        {{ tab.name }}
      </li>
    </ul>
  </header>

  <template v-if="service">
    <component :is="selectedTab.component" :service="service"/>
  </template>
</template>

<style scoped>
header {
  grid-column: 2;
  grid-row: 1;
  background-color: whitesmoke;

  h1 {
    margin: 0.5rem;
  }

  ul {
    padding: 0.5rem;
    margin: 0;
    list-style-type: none;
    display: flex;

    li {
      padding: 0.5rem;
    }
  }
}
</style>
