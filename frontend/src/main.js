import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './assets/tailwind.css'

import { setupStore } from './stores'
import { setupElIcons } from './plugins/icons'

import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

const app = createApp(App)

app.use(ElementPlus)

// 注册 pinia
setupStore(app)
// 注册 Element-plus 图标
setupElIcons(app)

app.use(router)

app.mount('#app')
