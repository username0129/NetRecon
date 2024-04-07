import http from '@/utils/http'

// 获取验证码
export const captcha = () => {
  return http({
    url: '/api/v1/captcha/getcaptcha',
    method: 'GET'
  })
}

// 用户登录
export const login = (loginData) => {
  return http({
    url: '/api/v1/auth/postlogin',
    method: 'POST',
    data: loginData
  })
}
