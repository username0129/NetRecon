import http from '@/utils/http'

// 添加计划任务信息
export const AddCron = (data) => {
  return http({
    url: '/api/v1/cron/postaddtask',
    method: 'POST',
    data: data
  })
}
