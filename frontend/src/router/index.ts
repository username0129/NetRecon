import {createRouter, createWebHashHistory, RouteRecordRaw} from "vue-router";
import {close, start} from '@/utils/nporgress'

// 静态路由
const myRoutes: RouteRecordRaw[] = [
    {
        path: '/',
        redirect: '/login'
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/view/login/index.vue')
    },
]

// 创建路由
const router = createRouter({
    history: createWebHashHistory(),
    routes: myRoutes,
    scrollBehavior: () => ({left: 0, top: 0}),    // 切换路由时保证页面显示在最上端
});

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