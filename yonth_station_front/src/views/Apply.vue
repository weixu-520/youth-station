<template>
  <div class="apply-page">
    <button class="back-btn" @click="$router.back()">← 返回</button>
    <div class="form-card">
      <h2>✍️ 提交入住申请</h2>
      <div v-if="error" class="err">{{ error }}</div>
      <form @submit.prevent="submit">
        <div class="form-group">
          <label>选择驿站 <span class="req">*</span></label>
          <select v-model="form.stationId" required>
            <option :value="0" disabled>请选择驿站</option>
            <option v-for="s in availableStations" :key="s.stationId" :value="s.stationId">{{ s.stationName }} ({{ s.district }})</option>
          </select>
        </div>
        <div class="row">
          <div class="form-group">
            <label>入住日期 <span class="req">*</span></label>
            <input v-model="form.checkinDate" type="date" :min="today" required />
          </div>
          <div class="form-group">
            <label>退房日期 <span class="req">*</span></label>
            <input v-model="form.checkoutDate" type="date" :min="minCheckout" required />
          </div>
        </div>
        <div class="form-group">
          <label>来访目的 <span class="req">*</span></label>
          <select v-model="form.visitPurpose" required>
            <option :value="0" disabled>请选择</option>
            <option :value="1">求职</option>
            <option :value="2">创业</option>
            <option :value="3">研学</option>
          </select>
        </div>
        <template v-if="form.visitPurpose === 1">
          <div class="form-group">
            <label>证明类型</label>
            <select v-model="form.interviewProofType">
              <option :value="1">面试邮件</option>
              <option :value="2">截图证明</option>
              <option :value="3">招聘公告</option>
              <option :value="4">函件</option>
            </select>
          </div>
          <div class="form-group">
            <label>证明内容 / 链接</label>
            <input v-model="form.proofContent" placeholder="请输入证明内容或上传链接" />
          </div>
        </template>
        <template v-if="form.visitPurpose === 2">
          <div class="form-group">
            <label>创业计划简介</label>
            <textarea v-model="form.businessPlan" rows="3" placeholder="请简要描述你的创业计划..."></textarea>
          </div>
        </template>
        <div class="form-group">
          <label>备注</label>
          <textarea v-model="form.remark" rows="2" placeholder="其他需要说明的信息（选填）"></textarea>
        </div>
        <button type="submit" class="submit-btn" :disabled="submitting">{{ submitting ? '提交中...' : '提交申请' }}</button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getAvailableStations } from '../api/station'
import { apply } from '../api/application'

const route = useRoute()
const router = useRouter()
const availableStations = ref([])
const submitting = ref(false)
const error = ref('')

const today = new Date().toISOString().split('T')[0]

const form = reactive({
  stationId: Number(route.params.stationId) || 0,
  checkinDate: '',
  checkoutDate: '',
  visitPurpose: 0,
  interviewProofType: 1,
  proofContent: '',
  businessPlan: '',
  remark: ''
})

const minCheckout = computed(() => {
  if (!form.checkinDate) return today
  const d = new Date(form.checkinDate); d.setDate(d.getDate() + 1)
  return d.toISOString().split('T')[0]
})

onMounted(async () => {
  try { const r = await getAvailableStations(); availableStations.value = r.data || [] } catch {}
})

async function submit() {
  error.value = ''
  if (!form.stationId) { error.value = '请选择驿站'; return }
  if (!form.checkinDate || !form.checkoutDate) { error.value = '请选择日期'; return }
  if (form.checkinDate >= form.checkoutDate) { error.value = '退房日期必须晚于入住日期'; return }
  const days = Math.ceil((new Date(form.checkoutDate) - new Date(form.checkinDate)) / 86400000)
  if (days > 7) { error.value = '最多连续入住7天'; return }
  if (!form.visitPurpose) { error.value = '请选择来访目的'; return }
  if (form.visitPurpose === 1 && !form.proofContent) { error.value = '求职需提供面试证明'; return }
  if (form.visitPurpose === 2 && !form.businessPlan) { error.value = '创业需提供创业计划'; return }
  submitting.value = true
  try {
    const payload = {
      stationId: form.stationId,
      checkinDate: form.checkinDate,
      checkoutDate: form.checkoutDate,
      visitPurpose: form.visitPurpose,
      remark: form.remark
    }
    if (form.visitPurpose === 1) payload.interviewInfo = { type: form.interviewProofType, content: form.proofContent, fileUrl: '' }
    if (form.visitPurpose === 2) payload.businessPlan = form.businessPlan
    await apply(payload)
    alert('申请提交成功！')
    router.push('/home/applications')
  } catch (e) {
    error.value = e.message
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.apply-page { max-width: 600px; margin: 0 auto }
.back-btn { background: none; border: none; color: #4a90d9; cursor: pointer; font-size: 14px; margin-bottom: 12px }
.form-card { background: #fff; border-radius: 12px; padding: 24px; box-shadow: 0 1px 4px rgba(0,0,0,.04) }
.form-card h2 { margin: 0 0 20px }
.err { background: #ffebee; color: #c62828; padding: 10px 14px; border-radius: 6px; margin-bottom: 16px; font-size: 13px }
.form-group { margin-bottom: 16px }
.form-group label { display: block; margin-bottom: 4px; font-size: 13px; font-weight: 600; color: #333 }
.req { color: #f44336 }
.row { display: grid; grid-template-columns: 1fr 1fr; gap: 12px }
input, select, textarea { width: 100%; padding: 10px 12px; border: 1px solid #ddd; border-radius: 6px; font-size: 14px; box-sizing: border-box }
textarea { resize: vertical }
.submit-btn { width: 100%; padding: 14px; background: #1a1a2e; color: #fff; border: none; border-radius: 8px; font-size: 16px; cursor: pointer }
.submit-btn:disabled { opacity: .6; cursor: not-allowed }
</style>
