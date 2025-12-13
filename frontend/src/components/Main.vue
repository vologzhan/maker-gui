<script setup lang="ts">
import {shallowRef} from "vue";
import type {ServiceDto} from "src/dto/service.ts";
import Controller from "./Controller.vue";
import Entity from "./Entity.vue";

const {service} = defineProps<{
  service?: ServiceDto,
}>()

const tabs = [
  {name: 'Database', component: Entity},
  {name: 'Controllers', component: Controller},
] as const

const tab = shallowRef<(typeof tabs)[number]>(tabs[0])
</script>

<template>
  <main v-if="!service">
    <h1>Select service</h1>
  </main>

  <template v-else>
    <header>
      <ul>
        <li class="clickable"
            v-for="t in tabs"
            :key="t.name"
            :class="{ selected: tab === t }"
            @click="tab = t"
        >
          {{ t.name }}
        </li>
      </ul>
    </header>
    <main>
      <component :is="tab.component" :service="service"/>
    </main>
  </template>
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
