import http from '@/utils/http'

// 检查初始化
export const checkInit = () => {
  return http({
    url: '/api/v1/init/getinit',
    method: 'GET'
  })
}

// 进行初始化
export const init = (initData) => {
  return http({
    url: '/api/v1/init/postinit',
    method: 'POST',
    data: initData
  })
}
