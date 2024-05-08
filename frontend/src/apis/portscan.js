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

// 删除端口扫描任务结果
export const DeletePortScanResult = (data) => {
  return http({
    url: '/api/v1/portscan/postdeleteresult',
    method: 'POST',
    data: data
  })
}

// 批量删除端口扫描任务结果
export const DeletePortScanResults = (data) => {
  return http({
    url: '/api/v1/portscan/postdeleteresults',
    method: 'POST',
    data: data
  })
}

// 删除端口扫描任务结果
export const ExportPortScanResult = (data) => {
  return http({
    url: '/api/v1/portscan/postexportdata',
    method: 'POST',
    data: data,
    responseType: 'blob'
  })
}
