import http from '@/utils/http'

// 提交端口扫描任务
export const SubmitPortScanTask = (data) => {
  return http({
    url: '/api/v1/portscan/postportscan',
    method: 'POST',
    data: data
  })
}


// 获取端口扫描任务执行结果
export const FetchPortScanResult = (data) => {
  return http({
    url: '/api/v1/portscan/postfetchresult',
    method: 'POST',
    data: data
  })
}
