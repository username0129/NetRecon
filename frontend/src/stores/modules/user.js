import { defineStore } from 'pinia'
import { fetchUserInfo } from '@/apis/user.js'

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
    },
    async fetchUserInfo() {
      try {
        const response = await fetchUserInfo()
        if (response.code == 200) {
          this.setUserInfo(response.data)
        }
        return response
      } catch (error) {
        console.error('获取用户信息失败:', error)
      }
    }
  }
})
