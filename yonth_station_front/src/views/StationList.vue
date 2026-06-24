<template>
  <div>
    <div class="toolbar">
      <select v-model="district" class="filter">
        <option value="">全部区域</option>
        <option v-for="d in districts" :key="d" :value="d">{{ d }}</option>
      </select>
      <input v-model="search" placeholder="搜索驿站名称..." class="search-input" />
    </div>
    <div class="grid">
      <div class="card" v-for="s in filteredStations" :key="s.stationId" @click="$router.push(`/home/station/${s.stationId}`)">
        <div class="card-img">{{ s.stationName[0] }}</div>
        <div class="card-body">
          <h3>{{ s.stationName }}</h3>
          <p class="addr">📍 {{ s.district }} · {{ s.address }}</p>
          <p class="meta">🚇 {{ s.nearbyMetro || '暂无' }}</p>
          <div class="card-footer">
            <span :class="['badge', s.status === 1 ? 'on' : 'off']">{{ s.status === 1 ? '运营中' : '已关闭' }}</span>
            <span class="rooms">🏠 剩余 {{ s.availableRooms }} 间</span>
          </div>
        </div>
      </div>
    </div>
    <div class="pagination">
      <button :disabled="page===1" @click="page--; fetchData()">上一页</button>
      <span>第 {{ page }} / {{ totalPages }} 页</span>
      <button :disabled="page>=totalPages" @click="page++; fetchData()">下一页</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getStationList } from '../api/station'

const stations = ref([])
const district = ref('')
const search = ref('')
const page = ref(1)
const pageSize = ref(9)
const total = ref(0)
const districts = ['东湖区', '红谷滩区', '青山湖区', '西湖区', '高新区', '青云谱区', '新建区', '南昌县']

const totalPages = computed(() => Math.ceil(total.value / pageSize.value) || 1)

const filteredStations = computed(() => {
  let list = stations.value
  if (district.value) list = list.filter(s => s.district === district.value)
  if (search.value) list = list.filter(s => s.stationName.includes(search.value))
  return list
})

async function fetchData() {
  try {
    const res = await getStationList({ page: page.value, pageSize: pageSize.value, status: 1 })
    stations.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch {}
}
onMounted(fetchData)
</script>

<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px }
.filter, .search-input { padding: 8px 12px; border: 1px solid #ddd; border-radius: 6px; font-size: 14px }
.search-input { flex: 1; max-width: 300px }
.grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px }
.card { background: #fff; border-radius: 10px; overflow: hidden; cursor: pointer; transition: .15s; box-shadow: 0 1px 4px rgba(0,0,0,.04) }
.card:hover { transform: translateY(-2px); box-shadow: 0 4px 14px rgba(0,0,0,.1) }
.card-img { height: 100px; background: linear-gradient(135deg, #1a1a2e, #4a90d9); display: flex; align-items: center; justify-content: center; font-size: 40px; color: #fff; font-weight: 700 }
.card-body { padding: 16px }
.card-body h3 { margin: 0 0 6px; font-size: 16px }
.addr, .meta { font-size: 12px; color: #888; margin: 2px 0 }
.card-footer { display: flex; justify-content: space-between; align-items: center; margin-top: 10px }
.badge { padding: 3px 8px; border-radius: 8px; font-size: 11px }
.badge.on { background: #e8f5e9; color: #4caf50 }
.badge.off { background: #ffebee; color: #f44336 }
.rooms { font-size: 12px; color: #666 }
.pagination { display: flex; justify-content: center; align-items: center; gap: 16px; margin-top: 20px }
.pagination button { padding: 6px 16px; border: 1px solid #ddd; background: #fff; border-radius: 6px; cursor: pointer }
.pagination button:disabled { opacity: .4; cursor: not-allowed }
</style>
