import { defineStore } from 'pinia'
import { FetchCurrentUserInfo } from '@/apis/user.js'
import router from '@/router/index.js'

export const useUserStore = defineStore('user', {
  state: () => ({
    userInfo: {},
    token: localStorage.getItem('token') || ''
  }),
  actions: {
    setUserInfo(info) {
      this.userInfo = info
    },
    setToken(token) {
      this.token = token
      localStorage.setItem('token', token)
    },
    logout() {
      this.userInfo = {}
      this.token = ''
      localStorage.removeItem('token')
      router.push({ name: 'Login', replace: true })
      window.location.reload()
    },
    async FetchCurrentUserInfo() {
      try {
        const response = await FetchCurrentUserInfo()
        if (response.code === 200) {
          this.setUserInfo(response.data)
        }
        return response
      } catch (error) {
        console.error('获取用户信息失败:', error)
      }
    }
  },
  getters: {}
})
