<template>
  <div>
    <div class="welcome-card">
      <h1>👋 你好，{{ userName }}</h1>
      <p>欢迎使用南昌市青年驿站服务平台</p>
    </div>
    <div class="stats-row">
      <div class="stat-card" @click="$router.push('/home/stations')">
        <div class="stat-icon">🏘️</div>
        <div class="stat-num">{{ stationCount }}</div>
        <div class="stat-label">运营驿站</div>
      </div>
      <div class="stat-card" @click="$router.push('/home/applications')">
        <div class="stat-icon">📋</div>
        <div class="stat-num">{{ myAppCount }}</div>
        <div class="stat-label">我的申请</div>
      </div>
      <div class="stat-card" @click="$router.push('/home/apply')">
        <div class="stat-icon">✍️</div>
        <div class="stat-num">+</div>
        <div class="stat-label">提交申请</div>
      </div>
    </div>
    <div class="quick-section" v-if="recentApps.length">
      <h3>📌 最近申请</h3>
      <div class="app-card" v-for="a in recentApps" :key="a.applicationId" @click="$router.push(`/home/application/${a.applicationId}`)">
        <div class="app-info">
          <strong>{{ a.stationName }}</strong>
          <span>{{ a.checkinDate }} ~ {{ a.checkoutDate }}</span>
        </div>
        <span class="status-tag" :class="statusClass(a.status)">{{ a.statusDesc }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getAvailableStations } from '../api/station'
import { getMyApplications } from '../api/application'

const userName = ref(localStorage.getItem('userName') || '用户')
const stationCount = ref(0)
const myAppCount = ref(0)
const recentApps = ref([])

const statusMap = { 0: 'pending', 1: 'passed', 2: 'rejected', 3: 'cancelled', 4: 'checked-in', 5: 'checked-out' }
function statusClass(s) { return statusMap[s] || '' }

onMounted(async () => {
  try { const r = await getAvailableStations(); stationCount.value = r.data?.length || 0 } catch {}
  try {
    const r = await getMyApplications({ page: 1, pageSize: 5 })
    recentApps.value = r.data?.list || []
    myAppCount.value = r.data?.total || 0
  } catch {}
})
</script>

<style scoped>
.welcome-card { background: linear-gradient(135deg, #1a1a2e, #16213e); color: #fff; padding: 32px; border-radius: 12px; margin-bottom: 20px }
.welcome-card h1 { margin: 0 0 8px; font-size: 24px }
.welcome-card p { margin: 0; opacity: .7 }
.stats-row { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; margin-bottom: 24px }
.stat-card { background: #fff; padding: 24px; border-radius: 10px; text-align: center; cursor: pointer; transition: .15s; box-shadow: 0 1px 4px rgba(0,0,0,.04) }
.stat-card:hover { transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,0,0,.1) }
.stat-icon { font-size: 32px; margin-bottom: 8px }
.stat-num { font-size: 28px; font-weight: 700; color: #1a1a2e }
.stat-label { font-size: 13px; color: #888; margin-top: 4px }
.quick-section { background: #fff; border-radius: 10px; padding: 20px; box-shadow: 0 1px 4px rgba(0,0,0,.04) }
.quick-section h3 { margin: 0 0 12px }
.app-card { display: flex; justify-content: space-between; align-items: center; padding: 12px 0; border-bottom: 1px solid #f0f0f0; cursor: pointer }
.app-card:last-child { border: none }
.app-info { display: flex; flex-direction: column; gap: 2px }
.app-info span { font-size: 12px; color: #888 }
.status-tag { padding: 4px 10px; border-radius: 12px; font-size: 12px }
.status-tag.pending { background: #fff3e0; color: #ff9800 }
.status-tag.passed { background: #e8f5e9; color: #4caf50 }
.status-tag.rejected { background: #ffebee; color: #f44336 }
.status-tag.checked-in { background: #e3f2fd; color: #2196f3 }
.status-tag.checked-out { background: #f3e5f5; color: #9c27b0 }
.status-tag.cancelled { background: #f5f5f5; color: #999 }
</style>
