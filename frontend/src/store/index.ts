import type {App} from "vue";
import {createPinia} from "pinia";

const store = createPinia();

// 注册 store
export function setupStore(app: App<Element>) {
    app.use(store);
}

export {store}