<template>
  <div>
    <button class="back-btn" @click="$router.back()">← 返回</button>
    <div v-if="app" class="detail-card">
      <h2>📋 申请详情 #{{ app.applicationId }}</h2>
      <div class="stag-header"><span class="stag" :class="sc(app.status)">{{ app.statusDesc }}</span></div>
      <div class="grid">
        <div class="block"><h4>申请人</h4><p>{{ app.userName }}</p></div>
        <div class="block"><h4>驿站</h4><p>{{ app.stationName }}</p></div>
        <div class="block"><h4>入住日期</h4><p>{{ app.checkinDate }}</p></div>
        <div class="block"><h4>退房日期</h4><p>{{ app.checkoutDate }}</p></div>
        <div class="block"><h4>来访目的</h4><p>{{ ['','求职','创业','研学'][app.visitPurpose] || '-' }}</p></div>
        <div class="block"><h4>押金</h4><p>{{ app.depositAmount ? (app.depositAmount/100).toFixed(2)+'元' : '未缴纳' }}</p></div>
        <div class="block" v-if="app.rejectReason"><h4>拒绝原因</h4><p style="color:#f44336">{{ app.rejectReason }}</p></div>
        <div class="block" v-if="app.auditBy"><h4>审核人</h4><p>{{ app.auditBy }}</p></div>
      </div>
      <div class="block"><h4>联系方式</h4><p>📱 {{ app.userPhone || '-' }} | 🆔 {{ app.userIdCard || '-' }}</p></div>
      <div class="block"><h4>毕业院校 / 户籍</h4><p>{{ app.userSchool || '-' }} / {{ app.userHukou || '-' }}</p></div>
    </div>
    <div v-else class="loading">加载中...</div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getApplicationDetail } from '../api/application'

const route = useRoute()
const app = ref(null)
const sm = { 0: 'pending', 1: 'passed', 2: 'rejected', 3: 'cancelled', 4: 'checked-in', 5: 'checked-out' }
function sc(s) { return sm[s] || '' }

onMounted(async () => {
  try { const r = await getApplicationDetail(route.params.id); app.value = r.data } catch {}
})
</script>

<style scoped>
.back-btn { background: none; border: none; color: #4a90d9; cursor: pointer; font-size: 14px; margin-bottom: 12px }
.detail-card { background: #fff; border-radius: 12px; padding: 24px; box-shadow: 0 1px 4px rgba(0,0,0,.04) }
.detail-card h2 { margin: 0 0 16px }
.stag-header { margin-bottom: 20px }
.stag { padding: 4px 14px; border-radius: 12px; font-size: 13px }
.stag.pending { background: #fff3e0; color: #ff9800 }
.stag.passed { background: #e8f5e9; color: #4caf50 }
.stag.rejected { background: #ffebee; color: #f44336 }
.stag.cancelled { background: #f5f5f5; color: #999 }
.stag.checked-in { background: #e3f2fd; color: #2196f3 }
.stag.checked-out { background: #f3e5f5; color: #9c27b0 }
.grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; margin-bottom: 12px }
.block { background: #f8f9fa; border-radius: 8px; padding: 12px 16px }
.block h4 { margin: 0 0 4px; font-size: 12px; color: #888; text-transform: uppercase }
.block p { margin: 0; font-size: 14px }
.loading { text-align: center; padding: 40px; color: #888 }
</style>
