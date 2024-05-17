import http from '@/utils/http'

// 获取任务数量
export const FetchTaskCount = () => {
  return http({
    url: '/api/v1/echarts/postfetchtaskcount',
    method: 'POST'
  })
}
// 获取资产信息
export const FetchDomainCount = () => {
  return http({
    url: '/api/v1/echarts/postfetchdomaincount',
    method: 'POST'
  })
}

// 获取资产信息
export const FetchPortCount = () => {
  return http({
    url: '/api/v1/echarts/postfetchportcount',
    method: 'POST'
  })
}
