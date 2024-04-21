import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getRoutes } from '@/apis/route'
import { RouteHandle } from '@/utils/route.js'

export const useRouteStore = defineStore('route', {
  state: () => ({
    routes: ref([]),
    flag: false,
    routeMap: ref({}) // 包含路由到组件的映射关系
  }),
  actions: {
    // 格式化接收到的路由数据
    formatRoutes(routes, parent = null) {
      routes.forEach((route) => {
        if (route.name) {
          route.parent = parent
          this.routeMap[route.name] = route
          if (route.children?.length) {
            this.formatRoutes(route.children, route)
          }
        }
      })
    },

    // 获取并设置动态路由
    async setRoutes() {
      const baseRoute = [
        {
          path: '/layout',
          name: 'layout',
          component: () => import('@/views/layout/IndexView.vue'),
          meta: {
            title: '底层layout'
          },
          children: []
        }
      ]
      const response = await getRoutes()
      const routes = response.data

      this.formatRoutes(routes)
      baseRoute[0].children = routes

      RouteHandle(baseRoute)
      this.routes = [...baseRoute] // 确保响应性更新
      this.flag = true
    }
  },
  getters: {
    // 可以添加一些基于状态的计算属性
  }
})
