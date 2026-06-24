<template>
  <div>
    <h2>📋 申请审核</h2>
    <div class="toolbar">
      <select v-model="statusFilter">
        <option :value="null">全部</option>
        <option :value="0">待审核</option>
        <option :value="1">已通过</option>
        <option :value="2">已拒绝</option>
      </select>
    </div>
    <div class="table-card">
      <table>
        <thead><tr><th>ID</th><th>申请人</th><th>驿站</th><th>日期</th><th>目的</th><th>状态</th><th>操作</th></tr></thead>
        <tbody>
          <tr v-for="a in apps" :key="a.applicationId">
            <td>{{ a.applicationId }}</td>
            <td>{{ a.userName }}</td>
            <td>{{ a.stationName }}</td>
            <td>{{ a.checkinDate }}~{{ a.checkoutDate }}</td>
            <td>{{ ['','求职','创业','研学'][a.visitPurpose] }}</td>
            <td><span class="stag" :class="sc(a.status)">{{ a.statusDesc }}</span></td>
            <td>
              <template v-if="a.status===0">
                <button class="act pass" @click="audit(a.applicationId, 1)">通过</button>
                <button class="act reject" @click="openReject(a.applicationId)">拒绝</button>
              </template>
              <span v-else class="done">-</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <!-- 拒绝弹窗 -->
    <div v-if="rejecting" class="modal-mask" @click.self="rejecting=null">
      <div class="modal">
        <h3>拒绝申请 #{{ rejecting }}</h3>
        <textarea v-model="rejectReason" rows="3" placeholder="请输入拒绝原因..."></textarea>
        <div class="modal-btns">
          <button @click="rejecting=null">取消</button>
          <button class="danger" @click="doReject">确认拒绝</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { getApplicationList, auditApplication } from '../api/admin'

const apps = ref([])
const statusFilter = ref(null)
const rejecting = ref(null)
const rejectReason = ref('')
const sm = { 0: 'pending', 1: 'passed', 2: 'rejected', 3: 'cancelled', 4: 'checked-in', 5: 'checked-out' }
function sc(s) { return sm[s] || '' }

async function fetch() {
  try {
    const params = { page: 1, pageSize: 50 }
    if (statusFilter.value !== null) params.status = statusFilter.value
    const r = await getApplicationList(params)
    apps.value = r.data?.list || []
  } catch {}
}
function openReject(id) { rejecting.value = id; rejectReason.value = '' }
async function doReject() {
  if (!rejectReason.value.trim()) return alert('请输入拒绝原因')
  try { await auditApplication({ applicationId: rejecting.value, result: 2, rejectReason: rejectReason.value }); rejecting.value = null; fetch() } catch (e) { alert(e.message) }
}
async function audit(id, result) {
  try { await auditApplication({ applicationId: id, result }); fetch() } catch (e) { alert(e.message) }
}
watch(statusFilter, fetch)
onMounted(fetch)
</script>

<style scoped>
.toolbar { margin-bottom: 16px }
.toolbar select { padding: 8px 12px; border: 1px solid #ddd; border-radius: 6px }
.table-card { background: #fff; border-radius: 10px; overflow: hidden; box-shadow: 0 1px 4px rgba(0,0,0,.04) }
table { width: 100%; border-collapse: collapse }
th, td { padding: 10px 14px; font-size: 13px; text-align: left; border-bottom: 1px solid #f0f0f0 }
th { background: #f8f9fa }
.stag { padding: 2px 8px; border-radius: 8px; font-size: 11px }
.stag.pending { background: #fff3e0; color: #ff9800 }
.stag.passed { background: #e8f5e9; color: #4caf50 }
.stag.rejected { background: #ffebee; color: #f44336 }
.act { padding: 3px 10px; border: none; border-radius: 4px; cursor: pointer; font-size: 11px; margin-right: 4px }
.act.pass { background: #e8f5e9; color: #2e7d32 }
.act.reject { background: #ffebee; color: #c62828 }
.done { color: #ccc }
.modal-mask { position: fixed; inset: 0; background: rgba(0,0,0,.3); display: flex; align-items: center; justify-content: center; z-index: 100 }
.modal { background: #fff; border-radius: 12px; padding: 24px; width: 400px }
.modal h3 { margin: 0 0 12px }
.modal textarea { width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 6px; resize: vertical; box-sizing: border-box }
.modal-btns { display: flex; gap: 8px; justify-content: flex-end; margin-top: 12px }
.modal-btns button { padding: 8px 16px; border: 1px solid #ddd; border-radius: 6px; cursor: pointer; background: #fff }
.modal-btns .danger { background: #f44336; color: #fff; border-color: #f44336 }
</style>
