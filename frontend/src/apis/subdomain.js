import http from '@/utils/http'

// 提交子域名扫描任务
export const SubmitSubdomainTask = (data) => {
  return http({
    url: '/api/v1/subdomain/postbrutesubdomains',
    method: 'POST',
    data: data
  })
}

// 获取子域名扫描任务执行结果
export const FetchSubdomainResult = (data) => {
  return http({
    url: '/api/v1/subdomain/postfetchresult',
    method: 'POST',
    data: data
  })
}

// 删除端口扫描任务结果
export const DeleteSubdomainResult = (data) => {
  return http({
    url: '/api/v1/subdomain/postdeleteresult',
    method: 'POST',
    data: data
  })
}

// 批量删除端口扫描任务结果
export const DeleteSubdomainResults = (data) => {
  return http({
    url: '/api/v1/subdomain/postdeleteresults',
    method: 'POST',
    data: data
  })
}

// 删除端口扫描任务结果
export const ExportSubdomainResult = (data) => {
  return http({
    url: '/api/v1/subdomain/postexportdata',
    method: 'POST',
    data: data,
    responseType: 'blob'
  })
}
