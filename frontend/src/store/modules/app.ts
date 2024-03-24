import {defineStore} from 'pinia';
import defaultSettings from "@/settings";
import zhCn from "element-plus/es/locale/lang/zh-cn";
import en from "element-plus/es/locale/lang/en";
import {store} from "@/store";


export const useAppStore = defineStore('app', {
    state: () => ({
        language: useStorage('language', defaultSettings.language),
    }),
    actions: {
        changeLanguage(val: string) {
            this.language = val;
        },
    },
    getters: {
        locale: (state) => {
            return state.language === 'en' ? en : zhCn;
        }
    },
})

export function useAppStoreHook() {
    return useAppStore(store);
}
