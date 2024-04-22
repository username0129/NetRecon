<script setup>
import { computed, onMounted, reactive, ref, toRefs } from 'vue'

// 当前组件名称 MenuItem，用于单个菜单
defineOptions({
  name: 'MenuItem'
})

const props = defineProps({
  routeInfo: Object,
  myTheme: Object
})

const theme = reactive(props.myTheme)
const iconComponent = computed(() => props.routeInfo.meta.icon)

// 将CSS变量设置在组件挂载后
onMounted(() => {
  const { activeBackground, activeText, normalText, hoverBackground, hoverText } = toRefs(theme)
  document.documentElement.style.setProperty('--active-background', activeBackground.value)
  document.documentElement.style.setProperty('--active-text', activeText.value)
  document.documentElement.style.setProperty('--normal-text', normalText.value)
  document.documentElement.style.setProperty('--hover-background', hoverBackground.value)
  document.documentElement.style.setProperty('--hover-text', hoverText.value)
})

const activeBackground = ref(props.myTheme.activeBackground)
const activeText = ref(props.myTheme.activeText)
const normalText = ref(props.myTheme.normalText)
const hoverBackground = ref(props.myTheme.hoverBackground)
const hoverText = ref(props.myTheme.hoverText)
</script>

<template>
  <el-menu-item :index="routeInfo.name">
    <div class="my-menu-item">
      <el-icon v-if="iconComponent">
        <component :is="iconComponent" />
      </el-icon>
      <span class="my-menu-item-title">{{ routeInfo.meta.title }}</span>
    </div>
  </el-menu-item>
</template>

<style lang="scss" scoped>
.my-menu-item {
  width: 100%;
  flex: 1;
  height: 44px;
  display: flex;
  justify-content: flex-start;
  align-items: center;
  padding-left: 4px;

  .my-menu-item-title {
    flex: 1;
  }
}

.el-menu-item {
  color: var(--normal-text);

  &.is-active {
    .my-menu-item {
      background: var(--active-background) !important;
      border-radius: 4px;
      box-shadow: 0 0 2px 1px var(--active-background);

      i,
      span {
        color: var(--active-text);
      }
    }
  }

  &:hover {
    .my-menu-item {
      background: var(--hover-background);
      border-radius: 4px;
      box-shadow: 0 0 2px 1px var(--hover-background);
      color: var(--hover-text);
    }
  }
}
</style>
