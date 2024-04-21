import http from '@/utils/http'

// 获取用户信息
export const fetchUserInfo = () => {
  return http({
    url: '/api/v1/user/getuserinfo',
    method: 'GET'
  })
}
