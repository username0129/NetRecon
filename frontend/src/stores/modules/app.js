import { defineStore } from 'pinia'
import defaultSettings from '@/settings'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import en from 'element-plus/es/locale/lang/en'

export const useAppStore = defineStore('app', {
  state: () => ({
    language: defaultSettings.language
  }),
  actions: {
    changeLanguage(val) {
      this.language = val
    }
  },
  getters: {
    locale() {
      return this.language === 'en' ? en : zhCn
    }
  }
})
