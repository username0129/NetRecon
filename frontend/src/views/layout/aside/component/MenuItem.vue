<script setup>
import { computed, onMounted, reactive, toRefs } from 'vue'
import { ElIcon } from 'element-plus'

// 定义组件名称为MenuItem
defineOptions({
  name: 'MenuItem'
})

const props = defineProps({
  routeInfo: Object,
  myTheme: Object
})

const theme = reactive(props.myTheme)
const iconComponent = computed(() => props.routeInfo.meta.icon)

onMounted(() => {
  const { activeBackground, activeText, normalText, hoverBackground, hoverText } = toRefs(theme)
  document.documentElement.style.setProperty('--active-background', activeBackground.value)
  document.documentElement.style.setProperty('--active-text', activeText.value)
  document.documentElement.style.setProperty('--normal-text', normalText.value)
  document.documentElement.style.setProperty('--hover-background', hoverBackground.value)
  document.documentElement.style.setProperty('--hover-text', hoverText.value)
})
</script>

<template>
  <el-menu-item :index="routeInfo.name">
    <div class="my-menu-item">
      <el-icon v-if="iconComponent">
        <component :is="iconComponent" />
      </el-icon>
      <span>{{ routeInfo.meta.title }}</span>
    </div>
  </el-menu-item>
</template>

<style lang="scss" scoped>
.my-menu-item {
  display: flex;
  align-items: center;
  height: 45px;
  padding-left: 4px;
  width: 100%;
}

.el-menu-item {
  color: var(--normal-text);

  &.is-active {
    background: var(--active-background) !important;
    border-radius: 4px;
    box-shadow: 0 0 2px 1px var(--active-background);

    .my-menu-item {
      color: var(--active-text);
    }
  }

  &:hover {
    background: var(--hover-background);
    box-shadow: 0 0 2px 1px var(--hover-background);
    color: var(--hover-text);
  }
}
</style>
