<script setup>
import Aside from '@/views/layout/aside/AsideIndex.vue'
</script>

<template>
  <el-container class="layout-cont">
    <el-container>
      <el-aside class="main-cont my-aside" width="220px">
        <div
          class="min-h-[60px] text-center transition-all duration-300 flex items-center justify-center gap-2"
          style="background: #191a23"
        >
          <div class="inline-flex font-bold text-2xl" style="color: #ffffff">NetRecon</div>
        </div>
        <Aside class="aside" />
      </el-aside>

      <el-container>
        <el-header> </el-header>

        <el-main> </el-main>
      </el-container>
    </el-container>
  </el-container>
</template>

<style scoped>
.layout-container-demo .el-header {
  position: relative;
  background-color: var(--el-color-primary-light-7);
  color: var(--el-text-color-primary);
}

.layout-container-demo .el-aside {
  color: var(--el-text-color-primary);
  background: var(--el-color-primary-light-8);
}

.layout-container-demo .el-menu {
  border-right: none;
}

.layout-container-demo .el-main {
  padding: 0;
}

.layout-container-demo .toolbar {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  right: 20px;
}
</style>
