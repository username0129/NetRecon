import http from '@/utils/http'

// 获取验证码
export const captcha = () => {
    return http({
            url: '/api/v1/captcha/getcaptcha',
            method: 'GET'
        }
    )
}

export const login = (loginData: { username: string; password: string; answer: string; captchaId: string }) => {
    return http({
        url: '/api/v1/auth/postlogin',
        method: 'POST',
        data: loginData
    })
}