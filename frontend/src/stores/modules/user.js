import {defineStore} from 'pinia'
import {store} from '@/stores'
import {ref} from 'vue'

export const useUserStore = defineStore('user', {
    state: () => ({
        userInfo: ref(JSON.parse(localStorage.getItem('userInfo')) || {}),
        token: ref(localStorage.getItem('token') || '')
    }),
    actions: {
        setUserInfo(info) {
            this.userInfo = info
            localStorage.setItem('userInfo', JSON.stringify(info));
        },
        setToken(token) {
            this.token = token
            localStorage.setItem('token', token);
        },
        logout() {
            this.userInfo = {};
            this.token = '';
            localStorage.removeItem('userInfo');
            localStorage.removeItem('token');
        }
    },
    getters: {}
})

export function useUserStoreHook() {
    return useUserStore(store)
}
