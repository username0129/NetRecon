import { createPinia } from 'pinia'

const store = createPinia()

// 注册 store
export function setupStore(app) {
  app.use(store)
}

export { store }
