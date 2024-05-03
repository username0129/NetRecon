import http from '@/utils/http'

// 获取任务列表
export const FetchTasks = (data) => {
  return http({
    url: '/api/v1/task/postfetchtasks',
    method: 'POST',
    data: data
  })
}

// 删除任务
export const DeleteTask = (data) => {
  return http({
    url: '/api/v1/task/postdeletetask',
    method: 'POST',
    data: data
  })
}

// 取消任务
export const CancelTask = (data) => {
  return http({
    url: '/api/v1/task/postcanceltask',
    method: 'POST',
    data: data
  })
}
