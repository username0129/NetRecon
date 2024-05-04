import http from '@/utils/http'

// 获取资产信息
export const FetchAsset = (data) => {
  return http({
    url: '/api/v1/asset/postfetchasset',
    method: 'POST',
    data: data
  })
}

// 添加资产信息
export const AddAsset = (data) => {
  return http({
    url: '/api/v1/asset/postaddasset',
    method: 'POST',
    data: data
  })
}

// 删除资产信息
export const DeleteAsset = (data) => {
  return http({
    url: '/api/v1/asset/postdeleteasset',
    method: 'POST',
    data: data
  })
}

// 删除资产信息
export const UpdateAsset = (data) => {
  return http({
    url: '/api/v1/asset/postupdateasset',
    method: 'POST',
    data: data
  })
}
