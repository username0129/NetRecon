import { useAppStore } from '@/stores/modules/app'
import { createI18n } from 'vue-i18n'
import enLocale from './lang/en.js'
import zhCnLocale from './lang/zh_cn.js'
import { storeToRefs } from 'pinia'

const appStore = useAppStore()

const { language } = storeToRefs(appStore)

const messages = {
  'zh-cn': {
    ...zhCnLocale
  },
  en: {
    ...enLocale
  }
}

const i18n = createI18n({
  locale: language,
  legacy: false,
  messages: messages,
  globalInjection: true
})

export default i18n
