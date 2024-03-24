<script setup lang="ts">
import {useAppStore} from "@/store/modules/app";
import {storeToRefs} from "pinia";

const appStore = useAppStore();

const {locale} = storeToRefs(appStore);
</script>

<template>
  <el-config-provider :locale="locale">
    <router-view/>
  </el-config-provider>
</template>

<style scoped>

</style>
