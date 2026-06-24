import api from './index'

export function getStationList(data) {
  return api.post('/station/list', data)
}
export function getStationDetail(stationId) {
  return api.get(`/station/detail/${stationId}`)
}
export function getAvailableStations() {
  return api.get('/station/available')
}
// 点赞
export function like(stationId) {
  return api.post('/station/like', { stationId })
}
export function unlike(stationId) {
  return api.post('/station/unlike', { stationId })
}
export function getLikeCount(stationId) {
  return api.get(`/station/like/count/${stationId}`)
}
// 评论
export function addComment(data) {
  return api.post('/station/comment', data)
}
export function getComments(params) {
  return api.get('/station/comment/list', { params })
}
