import api from './index'

export function getHistory() {
  return api.get('/chat/history')
}
export function getChatUsers() {
  return api.get('/chat/users')
}
// 智能客服
export function askAI(question, sessionId) {
  return api.post('/chat/ask', { question, sessionId })
}
export function uploadKnowledge(data) {
  return api.post('/chat/knowledge/upload', data)
}
