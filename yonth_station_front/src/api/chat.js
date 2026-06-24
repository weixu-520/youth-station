import api from './index'

export function getHistory() {
  return api.get('/chat/history')
}
export function getChatUsers() {
  return api.get('/chat/users')
}
