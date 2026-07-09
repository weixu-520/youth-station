<template>
  <div class="chat-page">
    <!-- 管理员视角：左侧用户列表 + 右侧聊天 -->
    <template v-if="isAdmin">
      <div class="admin-layout">
        <div class="user-sidebar">
          <h3>💬 咨询列表</h3>
          <div v-if="chatUsers.length === 0" class="empty">暂无用户咨询</div>
          <div
            v-for="u in chatUsers" :key="u.userId"
            :class="['user-item', activeUserId === u.userId ? 'active' : '']"
            @click="selectUser(u)"
          >{{ u.userName }}</div>
        </div>
        <div class="chat-main">
          <div v-if="activeUserId === 0" class="no-select">← 请选择要回复的用户</div>
          <template v-else>
            <div class="chat-header">回复：{{ activeUserName }}</div>
            <div class="chat-box" ref="chatBox">
              <div v-for="(m, i) in filteredMessages" :key="i" :class="['msg', m.fromUserId === userId ? 'mine' : 'other']">
                <div class="msg-meta">{{ m.fromUserId === userId ? '我' : activeUserName }} · {{ formatTime(m.createdAt) }}</div>
                <div class="msg-content">{{ m.content }}</div>
              </div>
            </div>
            <div class="input-row">
              <input v-model="text" @keyup.enter="send" placeholder="输入回复..." :disabled="!connected" />
              <button @click="send" :disabled="!connected || !text.trim()">发送</button>
            </div>
          </template>
        </div>
      </div>
    </template>

    <!-- 普通用户视角：直接聊客服 -->
    <template v-else>
      <h2>💬 联系客服</h2>
      <div v-if="!connected" class="status">连接中...</div>
      <div class="chat-box" ref="chatBox">
        <div v-for="(m, i) in messages" :key="i" :class="['msg', m.fromUserId === userId ? 'mine' : 'other']">
          <div class="msg-meta">{{ m.fromUserId === userId ? '我' : '客服' }} · {{ formatTime(m.createdAt) }}</div>
          <div class="msg-content">{{ m.content }}</div>
        </div>
        <div v-if="messages.length === 0" class="empty">暂无消息，在下方输入开始对话</div>
      </div>
      <div class="input-row">
        <input v-model="text" @keyup.enter="send" placeholder="输入消息..." :disabled="!connected" />
        <button @click="send" :disabled="!connected || !text.trim()">发送</button>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { getHistory, getChatUsers } from '../api/chat'

const userId = ref(Number(localStorage.getItem('userId')) || 0)
const isAdmin = ref(localStorage.getItem('isAdmin') === 'true')
const token = localStorage.getItem('token') || ''
const connected = ref(false)
const messages = ref([])
const text = ref('')
const chatBox = ref(null)
const chatUsers = ref([])
const activeUserId = ref(0)
const activeUserName = ref('')
const seenMsgIds = new Set()
let ws = null

// 管理员：只显示与当前选中用户的对话
const filteredMessages = computed(() => {
  if (!isAdmin.value || activeUserId.value === 0) return []
  return messages.value.filter(m =>
    m.fromUserId === activeUserId.value || m.toUserId === activeUserId.value
  )
})

function formatTime(ts) {
  if (!ts) return ''
  const d = new Date(ts * 1000)
  return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

function scrollDown() {
  nextTick(() => { const el = chatBox.value; if (el) el.scrollTop = el.scrollHeight })
}

function addMessage(msg) {
  const key = `${msg.fromUserId}|${msg.content}|${msg.createdAt}`
  if (seenMsgIds.has(key)) return
  seenMsgIds.add(key)
  messages.value.push(msg)
  scrollDown()
}

function connect() {
  const proto = location.protocol === 'https:' ? 'wss' : 'ws'
  ws = new WebSocket(`${proto}://${location.host}/api/v1/ws?token=${token}`)
  ws.onopen = () => { connected.value = true }
  ws.onclose = () => { connected.value = false }
  ws.onmessage = (e) => {
    try { addMessage(JSON.parse(e.data)) } catch {}
  }
}

function send() {
  if (!text.value.trim() || !ws || !connected.value) return
  const now = Math.floor(Date.now() / 1000)
  // 管理员发给指定用户，普通用户发给管理员(toUserId=0)
  const msg = {
    toUserId: isAdmin.value ? activeUserId.value : 0,
    targetType: isAdmin.value ? 2 : 1,
    content: text.value.trim(),
    createdAt: now
  }
  ws.send(JSON.stringify(msg))
  addMessage({ ...msg, fromUserId: userId.value })
  text.value = ''
}

// 管理员：加载聊天用户列表
async function loadChatUsers() {
  if (!isAdmin.value) return
  try { const r = await getChatUsers(); chatUsers.value = r.data || [] } catch {}
}

function selectUser(u) {
  activeUserId.value = u.userId
  activeUserName.value = u.userName
  scrollDown()
}

async function loadHistory() {
  try {
    const r = await getHistory()
    if (r.data?.length) r.data.forEach(m => addMessage(m))
  } catch {}
}

onMounted(async () => {
  await loadHistory()
  if (isAdmin.value) await loadChatUsers()
  connect()
})
onUnmounted(() => ws?.close())
</script>

<style scoped>
.chat-page { max-width: 800px; margin: 0 auto; height: calc(100vh - 120px); display: flex; flex-direction: column }
.chat-page h2 { margin-bottom: 16px; flex-shrink: 0 }
.status { color: #ff9800; margin-bottom: 8px; flex-shrink: 0 }
/* admin layout */
.admin-layout { display: flex; gap: 0; background: #fff; border-radius: 10px; overflow: hidden; box-shadow: 0 1px 4px rgba(0,0,0,.04); flex: 1; min-height: 0 }
.user-sidebar { width: 180px; border-right: 1px solid #eee; padding: 16px; overflow-y: auto; flex-shrink: 0 }
.user-sidebar h3 { font-size: 14px; margin: 0 0 12px }
.user-item { padding: 10px 12px; cursor: pointer; border-radius: 6px; font-size: 13px; margin-bottom: 4px }
.user-item:hover { background: #f0f2f5 }
.user-item.active { background: #1a1a2e; color: #fff }
.chat-main { flex: 1; display: flex; flex-direction: column; min-width: 0 }
.chat-header { padding: 14px 16px; border-bottom: 1px solid #eee; font-weight: 600; font-size: 14px }
.no-select { flex: 1; display: flex; align-items: center; justify-content: center; color: #ccc }
/* chat */
.chat-box { flex: 1; padding: 16px; overflow-y: auto }
.msg { margin-bottom: 12px; max-width: 75% }
.msg.mine { margin-left: auto }
.msg-meta { font-size: 11px; color: #999; margin-bottom: 2px }
.msg-content { padding: 8px 14px; border-radius: 14px; font-size: 14px; line-height: 1.4; word-break: break-word }
.msg.mine .msg-content { background: #1a1a2e; color: #fff; border-bottom-right-radius: 4px }
.msg.other .msg-content { background: #f0f0f0; color: #333; border-bottom-left-radius: 4px }
.input-row { display: flex; gap: 8px; padding: 12px; border-top: 1px solid #eee; flex-shrink: 0 }
.input-row input { flex: 1; padding: 10px; border: 1px solid #ddd; border-radius: 8px; font-size: 14px }
.input-row button { padding: 10px 20px; background: #1a1a2e; color: #fff; border: none; border-radius: 8px; cursor: pointer }
.input-row button:disabled { opacity: .5 }
.empty { text-align: center; padding: 40px 0; color: #ccc }
</style>
