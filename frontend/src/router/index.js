import { createRouter, createWebHistory } from 'vue-router'

// 静态路由
const constRoutes = [
  {
    path: '/login',
    name: 'Login',
    meta: { title: '登录页面' },
    component: () => import('@/views/login/IndexView.vue')
  },
  {
    path: '/init',
    name: 'Init',
    meta: { title: '初始化' },
    component: () => import('@/views/init/IndexView.vue')
  },
  {
    path: '/404',
    name: '404',
    meta: { title: '404' },
    component: () => import('@/views/init/IndexView.vue')
  }
]

// 创建路由
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: constRoutes,
  scrollBehavior: () => ({ left: 0, top: 0 }) // 切换路由时保证页面显示在最上端
})

export default router
