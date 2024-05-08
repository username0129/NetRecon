import http from '@/utils/http'

// 获取端口扫描任务执行结果
export const FetchOperationResult = (data) => {
  return http({
    url: '/api/v1/operation/postfetchresult',
    method: 'POST',
    data: data
  })
}

// 删除端口扫描任务结果
export const DeleteOperationResults = (data) => {
  return http({
    url: '/api/v1/operation/postdeleteresults',
    method: 'POST',
    data: data
  })
}
