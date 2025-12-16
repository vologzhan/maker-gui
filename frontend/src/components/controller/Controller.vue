<script setup lang="ts">
import type {ServiceDto} from "src/dto/service.ts";
import {defineProps, shallowRef} from "vue";
import Method from "./Method.vue";
import Request from "./Request.vue";
import Response from "./Response.vue";

const {service} = defineProps<{
  service: ServiceDto
}>();

const tabs = [
  {name: 'Methods', component: Method},
  {name: 'Requests', component: Request},
  {name: 'Responses', component: Response},
] as const

const selectedTab = shallowRef<(typeof tabs)[number]>(tabs[0])
</script>

<template>
  <div class="controller-container">
    <nav>
      <ul>
        <li class="clickable"
            v-for="tab in tabs"
            :key="tab.name"
            :class="{ selected: selectedTab === tab }"
            @click="selectedTab = tab"
        >
          {{ tab.name }}
        </li>
      </ul>
    </nav>

    <main>
      <component :is="selectedTab.component" :service="service"/>
    </main>
  </div>
</template>

<style scoped>
.controller-container {
  grid-column: 2;
  grid-row: 2;

  display: grid;
  grid-template-rows: auto 1fr;

  nav {
    grid-row: 1;
    background-color: whitesmoke;

    ul {
      padding: 0;
      margin: 0;
      list-style-type: none;
      display: flex;

      li {
        padding: 0.5rem;
      }
    }
  }

  main {
    padding: 1rem;
    grid-row: 2;
  }
}
</style>
