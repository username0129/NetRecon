import http from '@/utils/http'

// 获取验证码
export const getRoutes = () => {
  return http({
    url: '/api/v1/route/getroute',
    method: 'GET'
  })
}
