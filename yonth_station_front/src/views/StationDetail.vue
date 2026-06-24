<template>
  <div>
    <button class="back-btn" @click="$router.back()">← 返回</button>
    <div v-if="station" class="detail-card">
      <div class="detail-header">
        <div class="detail-img">{{ station.stationName[0] }}</div>
        <div class="detail-title">
          <h1>{{ station.stationName }}</h1>
          <div class="rating">⭐ {{ station.averageRating }} ({{ station.totalReviews }} 评价)</div>
        </div>
      </div>

      <!-- 点赞 -->
      <div class="like-row">
        <button :class="['like-btn', liked ? 'liked' : '']" @click="toggleLike">
          {{ liked ? '❤️' : '🤍' }} <span>{{ likeCount }}</span>
        </button>
        <button class="contact-btn" @click="$router.push('/home/chat')">💬 联系客服</button>
      </div>

      <div class="detail-grid">
        <div class="info-block">
          <h3>📍 位置信息</h3>
          <p>区域：{{ station.district }}</p>
          <p>地址：{{ station.address }}</p>
          <p>地铁：{{ station.nearbyMetro || '暂无' }}</p>
        </div>
        <div class="info-block">
          <h3>📞 联系方式</h3>
          <p>电话：{{ station.contactPhone }}</p>
          <p>时间：{{ station.businessHours }}</p>
        </div>
        <div class="info-block">
          <h3>🏠 房间情况</h3>
          <p>总房间：{{ station.totalRooms }} 间</p>
          <p>可预约：{{ station.availableRooms }} 间</p>
          <p>本周剩余名额：{{ station.remainingQuota }} / {{ station.weeklyQuota }}</p>
        </div>
        <div class="info-block">
          <h3>🛋️ 配套设施</h3>
          <div class="amenities-list">
            <span v-for="a in amenities" :key="a" class="tag">{{ a }}</span>
          </div>
        </div>
      </div>
      <div class="info-block full">
        <h3>📝 驿站介绍</h3>
        <p>{{ station.description }}</p>
      </div>
      <button class="apply-btn" @click="$router.push(`/home/apply/${station.stationId}`)">✍️ 立即申请</button>

      <!-- 评论区 -->
      <div class="comment-section">
        <h3>💬 评论 ({{ commentTotal }})</h3>
        <div class="comment-form">
          <textarea v-model="commentText" placeholder="写下你的评论..." rows="3"></textarea>
          <button @click="submitComment" :disabled="!commentText.trim() || commenting">{{ commenting ? '发送中...' : '发表评论' }}</button>
        </div>
        <div v-if="comments.length" class="comment-list">
          <div class="comment-item" v-for="c in comments" :key="c.id">
            <div class="comment-avatar">{{ c.userName?.[0] || '?' }}</div>
            <div class="comment-body">
              <div class="comment-meta"><strong>{{ c.userName }}</strong> · {{ formatTime(c.createdAt) }}</div>
              <p>{{ c.content }}</p>
            </div>
          </div>
        </div>
        <div v-if="commentTotal > comments.length" class="load-more">
          <button @click="loadComments(commentPage+1)">加载更多</button>
        </div>
      </div>
    </div>
    <div v-else class="loading">加载中...</div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getStationDetail, like, unlike, getLikeCount, addComment, getComments } from '../api/station'

const route = useRoute()
const station = ref(null)
const liked = ref(false)
const likeCount = ref(0)
const comments = ref([])
const commentTotal = ref(0)
const commentPage = ref(1)
const commentText = ref('')
const commenting = ref(false)

const amenities = computed(() => {
  try { return JSON.parse(station.value?.amenities || '[]') } catch { return [] }
})

function formatTime(ts) {
  const d = new Date(ts * 1000)
  return d.toLocaleDateString('zh-CN') + ' ' + d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

async function toggleLike() {
  try {
    if (liked.value) {
      const r = await unlike(station.value.stationId)
      liked.value = false
      likeCount.value = r.data?.count ?? likeCount.value
    } else {
      const r = await like(station.value.stationId)
      liked.value = true
      likeCount.value = r.data?.count ?? likeCount.value
    }
  } catch (e) {
    // 后端返回 "已点赞"/"未点赞" 时 code=400，前端静默更新状态
    if (e.message === '已点赞') { liked.value = true } else if (e.message === '未点赞') { liked.value = false }
  }
}

async function loadComments(page = 1) {
  try {
    const r = await getComments({ stationId: route.params.id, page, pageSize: 10 })
    if (page === 1) comments.value = r.data?.list || []
    else comments.value.push(...(r.data?.list || []))
    commentTotal.value = r.data?.total || 0
    commentPage.value = page
  } catch {}
}

async function submitComment() {
  if (!commentText.value.trim()) return
  commenting.value = true
  try {
    await addComment({ stationId: station.value.stationId, content: commentText.value.trim() })
    commentText.value = ''
    loadComments(1)
  } catch (e) { alert(e.message) }
  finally { commenting.value = false }
}

onMounted(async () => {
  try {
    const [detailRes, countRes] = await Promise.all([
      getStationDetail(route.params.id),
      getLikeCount(route.params.id).catch(() => ({ data: { count: 0 } }))
    ])
    station.value = detailRes.data
    likeCount.value = countRes.data?.count || 0
    loadComments(1)
  } catch {}
})
</script>

<style scoped>
.back-btn { background: none; border: none; color: #4a90d9; cursor: pointer; font-size: 14px; margin-bottom: 16px }
.detail-card { background: #fff; border-radius: 12px; padding: 24px; box-shadow: 0 1px 4px rgba(0,0,0,.04) }
.detail-header { display: flex; gap: 20px; margin-bottom: 20px }
.detail-img { width: 80px; height: 80px; border-radius: 12px; background: linear-gradient(135deg, #1a1a2e, #4a90d9); display: flex; align-items: center; justify-content: center; font-size: 36px; color: #fff; font-weight: 700; flex-shrink: 0 }
.detail-title h1 { margin: 0 0 6px; font-size: 22px }
.rating { color: #ff9800; font-size: 14px }
.like-row { display: flex; gap: 12px; margin-bottom: 20px; align-items: center }
.like-btn { padding: 8px 20px; border: 1.5px solid #ddd; border-radius: 20px; background: #fff; cursor: pointer; font-size: 15px; transition: .2s }
.like-btn.liked { border-color: #e91e63; background: #fce4ec; color: #e91e63 }
.like-btn:hover { transform: scale(1.05) }
.contact-btn { padding: 8px 20px; border: 1.5px solid #4a90d9; border-radius: 20px; background: #e3f2fd; color: #1976d2; cursor: pointer; font-size: 14px }
.detail-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin-bottom: 16px }
.info-block { background: #f8f9fa; border-radius: 8px; padding: 16px }
.info-block.full { grid-column: 1 / -1 }
.info-block h3 { margin: 0 0 8px; font-size: 14px; color: #333 }
.info-block p { margin: 4px 0; font-size: 13px; color: #666 }
.amenities-list { display: flex; flex-wrap: wrap; gap: 6px }
.tag { padding: 3px 10px; background: #e3f2fd; color: #1976d2; border-radius: 12px; font-size: 12px }
.apply-btn { margin: 20px 0; width: 100%; padding: 14px; background: #1a1a2e; color: #fff; border: none; border-radius: 8px; font-size: 16px; cursor: pointer; transition: .15s }
.apply-btn:hover { background: #16213e }
.comment-section { margin-top: 24px; border-top: 1px solid #eee; padding-top: 20px }
.comment-section h3 { margin: 0 0 14px }
.comment-form { display: flex; gap: 10px; margin-bottom: 16px }
.comment-form textarea { flex: 1; padding: 10px; border: 1px solid #ddd; border-radius: 8px; resize: vertical; font-size: 13px }
.comment-form button { padding: 8px 20px; background: #1a1a2e; color: #fff; border: none; border-radius: 8px; cursor: pointer; white-space: nowrap }
.comment-form button:disabled { opacity: .5 }
.comment-item { display: flex; gap: 12px; padding: 12px 0; border-bottom: 1px solid #f5f5f5 }
.comment-avatar { width: 36px; height: 36px; border-radius: 50%; background: #e3f2fd; color: #1976d2; display: flex; align-items: center; justify-content: center; font-weight: 700; flex-shrink: 0 }
.comment-meta { font-size: 12px; color: #888; margin-bottom: 4px }
.comment-body p { margin: 0; font-size: 14px; line-height: 1.5 }
.load-more { text-align: center; padding: 12px }
.load-more button { background: none; border: none; color: #4a90d9; cursor: pointer }
.loading { text-align: center; padding: 40px; color: #888 }
</style>
