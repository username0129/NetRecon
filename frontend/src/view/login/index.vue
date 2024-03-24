<script setup lang="ts">

import {useCounterStore} from "@/store/counter";

const store = useCounterStore();

const {count} = storeToRefs(store);
const {increment} = store;
</script>

<template>
  <div>
    <p>Count: {{ count }}</p>
    <button @click="increment">Increment</button>
  </div>
</template>

<style scoped>

</style>
