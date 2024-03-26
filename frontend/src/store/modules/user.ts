import {defineStore} from 'pinia';
import {store} from "@/store";
import http from "@/utils/http";
import {login} from "@/api/login";


export const useUserStore = defineStore('user', {
    state: () => ({
        userData: null,
    }),
    actions: {
        async login(loginData: { username: string; password: string; answer: string; captchaId: string }) {
            try {
                const response = await login(loginData);
                this.userData = response.data;
            } catch (error) {
                console.error('登录失败:', error);
                throw error;
            }
        },
    }
})

export function useUserStoreHook() {
    return useUserStore(store);
}
