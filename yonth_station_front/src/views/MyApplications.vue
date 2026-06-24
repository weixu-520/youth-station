<template>
  <div>
    <div class="toolbar">
      <select v-model="statusFilter" class="filter">
        <option :value="-1">全部状态</option>
        <option :value="0">待审核</option>
        <option :value="1">已通过</option>
        <option :value="2">已拒绝</option>
        <option :value="3">已取消</option>
        <option :value="4">已入住</option>
        <option :value="5">已退房</option>
      </select>
      <button class="new-btn" @click="$router.push('/home/apply')">✍️ 新建申请</button>
    </div>
    <div v-if="apps.length === 0" class="empty">暂无申请记录</div>
    <div class="table-card" v-else>
      <table>
        <thead><tr><th>驿站</th><th>日期</th><th>状态</th><th>押金</th><th>操作</th></tr></thead>
        <tbody>
          <tr v-for="a in apps" :key="a.applicationId">
            <td>{{ a.stationName }}</td>
            <td>{{ a.checkinDate }} ~ {{ a.checkoutDate }}</td>
            <td><span class="stag" :class="sc(a.status)">{{ a.statusDesc }}</span></td>
            <td>{{ a.depositAmount ? (a.depositAmount/100).toFixed(2)+'元' : '-' }}</td>
            <td>
              <button class="act-btn" @click="$router.push(`/home/application/${a.applicationId}`)">详情</button>
              <button v-if="a.status===0" class="act-btn danger" @click="cancel(a.applicationId)">取消</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { getMyApplications, cancelApplication } from '../api/application'

const apps = ref([])
const statusFilter = ref(-1)
const sm = { 0: 'pending', 1: 'passed', 2: 'rejected', 3: 'cancelled', 4: 'checked-in', 5: 'checked-out' }
function sc(s) { return sm[s] || '' }

async function fetch() {
  try {
    const params = { page: 1, pageSize: 50 }
    if (statusFilter.value !== -1) params.status = statusFilter.value
    const r = await getMyApplications(params)
    apps.value = r.data?.list || []
  } catch {}
}
async function cancel(id) {
  if (!confirm('确认取消该申请？')) return
  try { await cancelApplication(id); fetch() } catch (e) { alert(e.message) }
}
watch(statusFilter, fetch)
onMounted(fetch)
</script>

<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; justify-content: space-between }
.filter { padding: 8px 12px; border: 1px solid #ddd; border-radius: 6px }
.new-btn { padding: 8px 20px; background: #1a1a2e; color: #fff; border: none; border-radius: 6px; cursor: pointer }
.table-card { background: #fff; border-radius: 10px; overflow: hidden; box-shadow: 0 1px 4px rgba(0,0,0,.04) }
table { width: 100%; border-collapse: collapse }
th, td { padding: 12px 16px; text-align: left; font-size: 14px; border-bottom: 1px solid #f0f0f0 }
th { background: #f8f9fa; font-weight: 600 }
.stag { padding: 3px 10px; border-radius: 10px; font-size: 12px }
.stag.pending { background: #fff3e0; color: #ff9800 }
.stag.passed { background: #e8f5e9; color: #4caf50 }
.stag.rejected { background: #ffebee; color: #f44336 }
.stag.cancelled { background: #f5f5f5; color: #999 }
.stag.checked-in { background: #e3f2fd; color: #2196f3 }
.stag.checked-out { background: #f3e5f5; color: #9c27b0 }
.act-btn { padding: 4px 12px; border: 1px solid #ddd; background: #fff; border-radius: 4px; cursor: pointer; font-size: 12px; margin-right: 4px }
.act-btn.danger { border-color: #ffcdd2; color: #f44336 }
.empty { text-align: center; padding: 60px; color: #999 }
</style>
