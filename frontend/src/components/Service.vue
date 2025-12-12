<script setup lang="ts">
import {shallowRef} from "vue";
import type {ServiceDto} from "src/dto/service.ts";
import Controller from "./controller/Controller.vue";
import Database from "./database/Database.vue";

defineProps<{
  service?: ServiceDto,
}>()

const tabs = [
  {name: 'Database', component: Database},
  {name: 'Controllers', component: Controller},
];

const selectedTab = shallowRef(tabs[0])
</script>

<template>
  <header v-if="service">
    <ul>
      <li class="clickable" v-for="tab in tabs" :class="{ selected: selectedTab === tab }" @click="selectedTab = tab">
        {{ tab.name }}
      </li>
    </ul>
  </header>
  <main>
    <h1 v-if="!service">Select service</h1>
    <template v-else>
      <component :is="selectedTab?.component" :service="service"/>
    </template>
  </main>
</template>

<style scoped>
header {
  grid-column: 2;
  grid-row: 1;
  background-color: whitesmoke;

  ul {
    padding: 0.5rem;
    margin: 0;
    list-style-type: none;
    display: flex;

    li {
      padding: 0.5rem;

      &.selected {
        font-weight: bold;
      }
    }
  }
}

main {
  grid-column: 2;
  grid-row: 2;

  h1 {
    margin: 0.5rem;
  }
}
</style>
