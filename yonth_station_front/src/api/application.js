import api from './index'

export function apply(data) {
  return api.post('/application/apply', data)
}
export function getMyApplications(data) {
  return api.post('/application/my-list', data)
}
export function getApplicationDetail(id) {
  return api.get(`/application/detail/${id}`)
}
export function cancelApplication(id) {
  return api.post(`/application/cancel/${id}`)
}
export function paymentNotify(data) {
  return api.post('/application/payment/notify', data)
}
export function checkin(data) {
  return api.post('/application/checkin', data)
}
export function checkout(data) {
  return api.post('/application/checkout', data)
}
export function refundDeposit(id) {
  return api.post(`/application/refund/${id}`)
}
