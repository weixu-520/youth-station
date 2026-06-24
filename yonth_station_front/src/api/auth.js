import api from './index'

// 登录
export function login(account, password) {
  return api.post('/auth/login', { account, password })
}
// 注册
export function register(data) {
  return api.post('/auth/register', data)
}
