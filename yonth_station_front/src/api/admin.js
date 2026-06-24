import api from './index'

export function getDashboard() {
  return api.get('/admin/dashboard')
}
export function getApplicationList(data) {
  return api.post('/admin/application/list', data)
}
export function auditApplication(data) {
  return api.post('/admin/application/audit', data)
}
export function getStationList(data) {
  return api.post('/admin/station/list', data)
}
export function updateStation(stationId, data) {
  return api.put(`/admin/station/${stationId}`, data)
}
