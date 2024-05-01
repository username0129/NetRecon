import axios from 'axios'

// 创建Axios实例
const http = axios.create({
  baseURL: 'http://103.228.64.175:8081', // API 基础地址
  timeout: 10000 // 请求超时时间
})

// 请求拦截器
http.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    // 请求错误处理
    return Promise.reject(error)
  }
)

// 响应拦截器
http.interceptors.response.use(
  (response) => {
    if (response.status === 200) {
      return response.data
    } else {
      return Promise.reject(response.data)
    }
  },
  (error) => {
    console.error('There was an error!', error)
    return Promise.reject(error)
  }
)

export default http
