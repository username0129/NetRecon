import { createApp } from 'vue'
import App from '@/App.vue'
import router from '@/router'
import { setupStore } from './stores'
import { setupElIcons } from './plugins/icons'
import { ensureDynamicRoutes } from '@/permission.js'
import ElementPlus from 'element-plus'

import 'element-plus/dist/index.css'
import '@/permission'
import '@/assets/tailwind.css'
import '@/style/element_visiable.scss'

const app = createApp(App)

app.use(ElementPlus)

// 注册 pinia
setupStore(app)
// 注册 Element-plus 图标
setupElIcons(app)

await ensureDynamicRoutes()
app.use(router)

app.mount('#app')
