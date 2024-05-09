import http from '@/utils/http'

// 获取文件列表
export const FetchFiles = (data) => {
  return http({
    url: '/api/v1/upload/postfetchfiles',
    method: 'POST',
    data: data
  })
}

// 上传文件
export const UploadFile = (data) => {
  return http({
    url: '/api/v1/upload/postuploadfile',
    method: 'POST',
    data: data
  })
}

// 删除文件
export const DeleteFile = (data) => {
  return http({
    url: '/api/v1/upload/postdeletefile',
    method: 'POST',
    data: data
  })
}
