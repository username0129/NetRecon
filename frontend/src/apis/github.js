import http from '@/utils/http'

export function Commits(page) {
  return http({
    url: 'https://api.github.com/repos/flipped-aurora/gin-vue-admin/commits?page=' + page,
    method: 'get'
  })
}

export function Members() {
  return http({
    url: 'https://api.github.com/orgs/FLIPPED-AURORA/members',
    method: 'get'
  })
}
