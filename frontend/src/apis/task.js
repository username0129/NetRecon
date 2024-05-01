import http from '@/utils/http'

// 获取验证码
export const FetchTasks = (data) => {
  return http({
    url: '/api/v1/task/postfetchtasks',
    method: 'POST',
    data: data
  })
}
