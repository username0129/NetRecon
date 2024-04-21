import { useRouteStore } from '@/stores/modules/route.js'
import { useUserStore } from '@/stores/modules/user.js'
import { close, remove, start } from '@/utils/nprogress.js'
import router from '@/router/index.js'

export async function ensureDynamicRoutes() {
  const routeStore = useRouteStore()
  const userStore = useUserStore()

  if (routeStore.flag) {
    return // 如果已经加载，则无需重新加载
  }

  if (userStore.token) {
    try {
      await userStore.fetchUserInfo()
      await routeStore.setRoutes() // 确保动态路由加载逻辑正确
      routeStore.routes.forEach((route) => {
        if (!router.hasRoute(route.name)) {
          router.addRoute(route) // 只添加尚未添加的路由
        }
      })
      routeStore.flag = true
    } catch (error) {
      console.error('Error loading dynamic routes:', error)
    }
  }
}

// 路由守卫
router.beforeEach(async (to) => {
  start()
  document.title = `${to.meta.title} | NetRecone`

  const userStore = useUserStore()
  const token = userStore.token
  const whiteList = ['Login', 'Init']

  // 用户未登录且访问路径不在白名单中，跳转到登陆页面
  if (!token && !whiteList.includes(to.name)) {
    close()
    return { name: 'Login', query: { redirect: to.fullPath } }
  }

  // 用户已登录，且访问 / 路径时，跳转到仪表盘页面
  if (token && to.path === '/') {
    close()
    return { name: 'Dashboard', replace: true }
  }

  // 访问路径不存在时，跳转到 404 页面
  if (!router.hasRoute(to.name)) {
    close()
    return { name: '404' }
  }

  return true
})

router.afterEach(() => {
  close()
})

router.onError(() => {
  remove()
})
