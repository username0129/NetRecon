import http from '@/utils/http'

// 获取资产信息
export const FofaSearch = (data) => {
  return http({
    url: '/api/v1/fofa/postfofasearch',
    method: 'POST',
    data: data
  })
}

// 获取资产信息
export const FofaExportData = (data) => {
  return http({
    url: '/api/v1/fofa/postexportdata',
    method: 'POST',
    data: data,
    responseType: 'blob'
  })
}
