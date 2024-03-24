import App from './App.vue'
import router from "@/router";
import {setupStore} from "@/store";
import './tailwind.css'

const app = createApp(App);

setupStore(app)

app.use(router).mount("#app")

export default app
