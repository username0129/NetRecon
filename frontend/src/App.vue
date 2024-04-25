<script setup>
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'

defineOptions({
  name: 'App'
})
</script>

<template>
  <div id="app">
    <router-view v-slot="{ Component }">
      <el-config-provider :locale="zhCn">
        <component :is="Component" />
      </el-config-provider>
    </router-view>
  </div>
</template>

<style lang="scss">
#app {
  background: #eee;
  height: 100vh;
  overflow: hidden;
  font-weight: 400;
}

.el-button {
  font-weight: 400;
}
</style>
