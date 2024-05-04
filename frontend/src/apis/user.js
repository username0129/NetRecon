import http from '@/utils/http'

// 获取用户信息
export const FetchCurrentUserInfo = () => {
  return http({
    url: '/api/v1/user/getuserinfo',
    method: 'GET'
  })
}

// 添加用户信息
export const AddUserInfo = (data) => {
  return http({
    url: '/api/v1/user/postadduserinfo',
    method: 'POST',
    data: data
  })
}

export const UpdateUserInfo = (data) => {
  return http({
    url: '/api/v1/user/postupdateuserinfo',
    method: 'POST',
    data: data
  })
}

export const FetchUserInfo = (data) => {
  return http({
    url: '/api/v1/user/postfetchusers',
    method: 'POST',
    data: data
  })
}

export const DeleteUserInfo = (data) => {
  return http({
    url: '/api/v1/user/postdeleteuserinfo',
    method: 'POST',
    data: data
  })
}

export const ResetPassword = (data) => {
  return http({
    url: '/api/v1/user/postresetpassword',
    method: 'POST',
    data: data
  })
}
