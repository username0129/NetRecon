import { createRouter, createWebHistory } from 'vue-router'
import { close, start } from '@/utils/nprogress.js'

// 静态路由
const constRoutes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/LoginView.vue')
  },
  {
    path: '/init',
    name: 'Init',
    component: () => import('@/views/InitView.vue')
  }
]

// 创建路由
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: constRoutes,
  scrollBehavior: () => ({ left: 0, top: 0 }) // 切换路由时保证页面显示在最上端
})

// 在路由跳转开始时显示 NProgress
router.beforeEach((to, from, next) => {
  start()
  next()
})

// 在路由跳转结束时隐藏NProgress
router.afterEach(() => {
  close()
})

export default router
