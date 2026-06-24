<template>
  <div>
    <h2>📈 数据概览</h2>
    <div class="dash-grid">
      <div class="dash-card"><div class="num">{{ data.totalApplications }}</div><div class="lbl">全部申请</div></div>
      <div class="dash-card"><div class="num">{{ data.todayApplications }}</div><div class="lbl">今日申请</div></div>
      <div class="dash-card"><div class="num warn">{{ data.pendingApplications }}</div><div class="lbl">待审核</div></div>
      <div class="dash-card"><div class="num">{{ data.checkedInApplications }}</div><div class="lbl">已入住</div></div>
      <div class="dash-card"><div class="num">{{ data.totalStations }}</div><div class="lbl">总驿站</div></div>
      <div class="dash-card"><div class="num green">{{ data.activeStations }}</div><div class="lbl">运营中</div></div>
    </div>
    <div class="quick-links">
      <button @click="$router.push('/home/admin/applications')">📋 审核申请</button>
      <button @click="$router.push('/home/admin/stations')">⚙️ 管理驿站</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getDashboard } from '../api/admin'

const data = ref({})
onMounted(async () => {
  try { const r = await getDashboard(); data.value = r.data } catch {}
})
</script>

<style scoped>
.dash-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; margin-bottom: 24px }
.dash-card { background: #fff; padding: 24px; border-radius: 10px; text-align: center; box-shadow: 0 1px 4px rgba(0,0,0,.04) }
.num { font-size: 32px; font-weight: 700; color: #1a1a2e }
.num.warn { color: #ff9800 }
.num.green { color: #4caf50 }
.lbl { font-size: 13px; color: #888; margin-top: 6px }
.quick-links { display: flex; gap: 12px }
.quick-links button { padding: 12px 24px; border: none; border-radius: 8px; cursor: pointer; font-size: 14px; background: #1a1a2e; color: #fff }
</style>
