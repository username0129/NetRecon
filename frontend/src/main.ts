import App from './App.vue'
import router from "@/router";
import './tailwind.css'

import {setupStore} from "@/store";
import {setupElIcons} from "@/plugins/icons";
import {createApp} from "vue";
import {setupI18n} from "@/plugins/i18n";


const app = createApp(App);

// 注册 pinia
setupStore(app)
// 注册 Element-plus 图标
setupElIcons(app);
// 国际化
setupI18n(app)

app.use(router).mount("#app")

export default app
