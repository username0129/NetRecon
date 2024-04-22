<script setup>
import Aside from '@/views/layout/aside/AsideIndex.vue'
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { fmtTitle } from '@/utils/fmtRouterTitle.js'

defineOptions({
  name: 'LayoutIndex'
})

const route = useRoute()

const matched = computed(() => route.matched)
</script>

<template>
  <el-container class="layout-cont">
    <el-aside class="main-cont my-aside" width="220px">
      <div class="logo-container">
        <div class="app-title">NetRecon</div>
      </div>
      <Aside />
    </el-aside>

    <!-- 分块滑动功能 -->
    <el-main class="main-cont main-right">
      <transition :duration="{ enter: 800, leave: 100 }" mode="out-in" name="el-fade-in-linear">
        <div class="main-header">
          <el-header>
            <el-row>
              <el-col>
                <el-breadcrumb class="breadcrumb">
                  <el-breadcrumb-item
                    v-for="item in matched.slice(1, matched.length)"
                    :key="item.path"
                  >
                    {{ fmtTitle(item.meta.title, route) }}
                  </el-breadcrumb-item>
                </el-breadcrumb>
              </el-col>
            </el-row>
          </el-header>
        </div>
      </transition>
    </el-main>
  </el-container>
</template>

<style lang="scss" scoped>
.logo-container {
  min-height: 60px;
  text-align: center;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #191a23;
}

.app-title {
  font-size: 24px;
  font-weight: bold;
  color: #ffffff;
}

.main-header {
  width: calc(100% - 220px);
  position: fixed;
  top: 0;
  z-index: 50;
  box-sizing: border-box;
}

.button {
  font-size: 12px;
  color: #666;
  background: rgb(250, 250, 250);
  width: 25px;
  padding: 4px 8px;
  border: 1px solid #eaeaea;
  margin-right: 4px;
  border-radius: 4px;
}

:deep .el-overlay {
  background-color: hsla(0, 0%, 100%, 0.9);
}
</style>
