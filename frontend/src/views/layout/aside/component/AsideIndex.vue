<script setup>
import MenuItem from './MenuItem.vue'
import SubMenu from './SubMenu.vue'
import { computed } from 'vue'

// 当前组件名称 AsideComponent，根据是否含有子路由显示对应的菜单
defineOptions({
  name: 'AsideComponent'
})

const props = defineProps({
  routeInfo: Object, // 路由信息和子菜单数据
  myTheme: Object // 包含了主题相关的样式设置
})

// 判断路由是否含有子路由
const hasChildren = computed(() => props.routeInfo.children?.length > 0)

// 排除那些被标记为 hidden 的子路由。
const visibleChildren = computed(
  () => props.routeInfo.children?.filter((item) => !item.hidden) || []
)

// 根据是否含有子路由来判断是使用 SubmenuItem（有子路由） 还是 MenuItem（无子路由）
const menuComponent = computed(() => {
  return hasChildren.value && visibleChildren.value.length ? SubMenu : MenuItem
})
</script>

<template>
  <component :is="menuComponent" :my-theme="myTheme" :route-info="routeInfo">
    <template v-if="hasChildren">
      <AsideComponent
        v-for="item in visibleChildren"
        :key="item.name"
        :route-info="item"
        :my-theme="myTheme"
      />
    </template>
  </component>
</template>
