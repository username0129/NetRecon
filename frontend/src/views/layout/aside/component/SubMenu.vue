<script setup>
import { computed, onMounted, reactive, toRefs } from 'vue'

// 有子项的复合菜单项：SubMenu
defineOptions({
  name: 'SubMenu'
})

const props = defineProps({
  routeInfo: Object,
  myTheme: Object
})

const theme = reactive(props.myTheme)
const iconComponent = computed(() => props.routeInfo.meta.icon)

// 设置CSS变量
onMounted(() => {
  const { activeBackground, activeText, normalText } = toRefs(theme)
  document.documentElement.style.setProperty('--activeBackground', activeBackground.value)
  document.documentElement.style.setProperty('--activeText', activeText.value)
  document.documentElement.style.setProperty('--normalText', normalText.value)
})
</script>

<template>
  <el-sub-menu :index="routeInfo.name">
    <template #title>
      <div class="my-subMenu">
        <el-icon v-if="iconComponent">
          <component :is="iconComponent" />
        </el-icon>
        <span>{{ routeInfo.meta.title }}</span>
      </div>
    </template>
    <slot />
  </el-sub-menu>
</template>

<style lang="scss" scoped>
.el-sub-menu {
  :deep(.el-sub-menu__title) {
    padding: 6px;
    color: var(--normalText);
  }
}

.is-active:not(.is-opened) {
  :deep(.el-sub-menu__title .my-subMenu) {
    flex: 1;
    height: 100%;
    line-height: 44px;
    background: var(--activeBackground) !important;
    border-radius: 4px;
    box-shadow: 0 0 2px 1px var(--activeBackground) !important;
    color: var(--activeText);
  }
}

.my-subMenu {
  padding-left: 4px;
}
</style>
