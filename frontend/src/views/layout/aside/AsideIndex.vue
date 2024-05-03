<script setup>
import AsideComponent from '@/views/layout/aside/component/AsideIndex.vue'
import { reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useRouteStore } from '@/stores/modules/route.js'

// 组件名称：AsideIndex（侧边栏）
defineOptions({
  name: 'AsideIndex'
})

const route = useRoute()
const router = useRouter()
const routeStore = useRouteStore()

const myTheme = reactive({
  background: '#191a23',
  activeBackground: 'var(--el-color-primary)',
  activeText: '#fff',
  normalText: '#fff',
  hoverBackground: 'rgba(64, 158, 255, 0.08)',
  hoverText: '#fff'
})

// 当前活动路由
const activeMenu = ref(route.name)

watch(route, (newRoute) => {
  if (newRoute.name !== activeMenu.value) {
    activeMenu.value = newRoute.name
  }
})

const selectMenuItem = (index) => {
  if (index !== route.name) {
    router.push({ name: index, replace: true })
  }
}
</script>

<template>
  <div class="aside-container">
    <el-scrollbar>
      <transition :duration="{ enter: 800, leave: 100 }" mode="out-in" name="el-fade-in-linear">
        <el-menu
          class="el-menu-vertical"
          :default-active="activeMenu"
          :background-color="myTheme.background"
          :active-text-color="myTheme.activeText"
          unique-opened
          @select="selectMenuItem"
        >
          <template v-for="item in routeStore.routes[0].children">
            <aside-component
              v-if="!item.meta.hidden"
              :key="item.name"
              :route-info="item"
              :my-theme="myTheme"
            />
          </template>
        </el-menu>
      </transition>
    </el-scrollbar>
  </div>
</template>

<style lang="scss" scoped>
.aside-container {
  background: #191a23;
  height: calc(100vh - 60px);
}

.el-sub-menu__title:hover,
.el-menu-item:hover {
  background: transparent;
}

.el-scrollbar {
  .el-scrollbar__view {
    height: 100%;
  }
}
</style>
