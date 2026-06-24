<template>
  <div>
    <h2>⚙️ 驿站管理</h2>
    <div class="table-card">
      <table>
        <thead><tr><th>ID</th><th>名称</th><th>区域</th><th>房间</th><th>状态</th><th>操作</th></tr></thead>
        <tbody>
          <tr v-for="s in stations" :key="s.stationId">
            <td>{{ s.stationId }}</td>
            <td>{{ s.stationName }}</td>
            <td>{{ s.district }}</td>
            <td>{{ s.availableRooms }}/{{ s.totalRooms }}</td>
            <td><span :class="s.status===1?'on':'off'">{{ s.status===1?'运营':'关闭' }}</span></td>
            <td><button class="edit-btn" @click="openEdit(s)">编辑</button></td>
          </tr>
        </tbody>
      </table>
    </div>
    <!-- 编辑弹窗 -->
    <div v-if="editing" class="modal-mask" @click.self="editing=null">
      <div class="modal">
        <h3>编辑驿站 — {{ editing.stationName }}</h3>
        <div class="form-grid">
          <div><label>名称</label><input v-model="editForm.stationName" /></div>
          <div><label>区域</label><input v-model="editForm.district" /></div>
          <div><label>地址</label><input v-model="editForm.address" /></div>
          <div><label>电话</label><input v-model="editForm.contactPhone" /></div>
          <div><label>营业时间</label><input v-model="editForm.businessHours" /></div>
          <div><label>总房间</label><input v-model="editForm.totalRooms" type="number" /></div>
          <div><label>可预约</label><input v-model="editForm.availableRooms" type="number" /></div>
          <div><label>状态</label><select v-model="editForm.status"><option :value="0">关闭</option><option :value="1">运营</option></select></div>
          <div><label>每周配额</label><input v-model="editForm.weeklyQuota" type="number" /></div>
          <div><label>剩余配额</label><input v-model="editForm.remainingQuota" type="number" /></div>
          <div><label>描述</label><textarea v-model="editForm.description" rows="2"></textarea></div>
          <div><label>配套设施(JSON)</label><input v-model="editForm.amenities" /></div>
          <div><label>附近地铁</label><input v-model="editForm.nearbyMetro" /></div>
          <div><label>图片URL</label><input v-model="editForm.imageUrl" /></div>
        </div>
        <div class="modal-btns">
          <button @click="editing=null">取消</button>
          <button class="primary" @click="saveEdit" :disabled="saving">{{ saving?'保存中...':'保存' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getStationList, updateStation } from '../api/admin'

const stations = ref([])
const editing = ref(null)
const saving = ref(false)
const editForm = ref({})

async function fetch() {
  try { const r = await getStationList({ page: 1, pageSize: 50 }); stations.value = r.data?.list || [] } catch {}
}
function openEdit(s) { editing.value = s; editForm.value = { ...s } }
async function saveEdit() {
  saving.value = true
  try {
    await updateStation(editing.value.stationId, { ...editForm.value, stationId: editing.value.stationId, latitude: editing.value.latitude||0, longitude: editing.value.longitude||0 })
    editing.value = null; fetch()
  } catch (e) { alert(e.message) } finally { saving.value = false }
}
onMounted(fetch)
</script>

<style scoped>
.table-card { background: #fff; border-radius: 10px; overflow: hidden; box-shadow: 0 1px 4px rgba(0,0,0,.04) }
table { width: 100%; border-collapse: collapse }
th, td { padding: 10px 14px; font-size: 13px; text-align: left; border-bottom: 1px solid #f0f0f0 }
th { background: #f8f9fa }
.on { color: #4caf50; font-weight: 600 }
.off { color: #f44336 }
.edit-btn { padding: 4px 12px; border: 1px solid #4a90d9; color: #4a90d9; background: #fff; border-radius: 4px; cursor: pointer; font-size: 12px }
.modal-mask { position: fixed; inset: 0; background: rgba(0,0,0,.3); display: flex; align-items: center; justify-content: center; z-index: 100 }
.modal { background: #fff; border-radius: 12px; padding: 24px; width: 600px; max-height: 80vh; overflow-y: auto }
.modal h3 { margin: 0 0 16px }
.form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px }
.form-grid div { display: flex; flex-direction: column }
.form-grid label { font-size: 12px; color: #666; margin-bottom: 2px }
.form-grid input, .form-grid select, .form-grid textarea { padding: 8px 10px; border: 1px solid #ddd; border-radius: 4px; font-size: 13px }
.modal-btns { display: flex; gap: 8px; justify-content: flex-end; margin-top: 16px }
.modal-btns button { padding: 8px 20px; border: 1px solid #ddd; border-radius: 6px; cursor: pointer; background: #fff }
.modal-btns .primary { background: #1a1a2e; color: #fff; border-color: #1a1a2e }
</style>
