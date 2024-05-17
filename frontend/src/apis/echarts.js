import http from '@/utils/http'

// 获取资产信息
export const FetchTaskCount = () => {
  return http({
    url: '/api/v1/echarts/postfetchtaskcount',
    method: 'POST',
  })
}
