import { defineStore } from 'pinia'
import { store } from '@/stores'
import { ref } from 'vue'
import cookie from 'js-cookie'

export const useUserStore = defineStore('user', {
  state: () => ({
    userInfo: ref({
      uuid: '',
      nickName: '',
      authority: {}
    }),
    token: ref(cookie.get('x-token') || '')
  }),
  actions: {
    setUserInfo() {},
    setToken(token) {
      this.token = token
    }
  },
  getters: {}
})

export function useUserStoreHook() {
  return useUserStore(store)
}
