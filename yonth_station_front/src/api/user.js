import api from './index'

export function getUserInfo() {
  return api.get('/user/info')
}
export function updateUserInfo(data) {
  return api.put('/user/info', data)
}
