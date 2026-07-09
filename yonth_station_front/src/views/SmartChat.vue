<template>
  <div class="smart-chat">
    <h2>🤖 智能客服</h2>
    <p class="subtitle">基于知识库的 AI 助手，随时为你解答驿站相关问题</p>

    <div class="chat-box" ref="chatBox">
      <div v-if="messages.length === 0" class="welcome">
        <div class="welcome-icon">🤖</div>
        <h3>你好！我是云驿智能客服</h3>
        <p>可以问我关于驿站入住、申请流程、政策等问题</p>
        <div class="quick-questions">
          <span v-for="q in quickQuestions" :key="q" @click="ask(q)">{{ q }}</span>
        </div>
      </div>

      <div v-for="(m, i) in messages" :key="i" :class="['bubble', m.role]">
        <div class="avatar">{{ m.role === 'user' ? '👤' : '🤖' }}</div>
        <div class="content" v-text="m.content"></div>
      </div>

      <div v-if="loading" class="bubble ai">
        <div class="avatar">🤖</div>
        <div class="content typing"><span></span><span></span><span></span></div>
      </div>
    </div>

    <div class="input-row">
      <input v-model="text" @keyup.enter="ask(text)" placeholder="输入你的问题..." :disabled="loading" />
      <button @click="ask(text)" :disabled="loading || !text.trim()">发送</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'

const messages = ref([])
const text = ref('')
const loading = ref(false)
const sessionId = ref('')
const chatBox = ref(null)

const quickQuestions = [
  '如何申请入住？',
  '最多可以住几天？',
  '需要准备什么材料？',
  '押金怎么退还？'
]

function scrollDown() {
  nextTick(() => { const el = chatBox.value; if (el) el.scrollTop = el.scrollHeight })
}

async function ask(question) {
  if (!question || !question.trim()) return
  if (loading.value) return

  messages.value.push({ role: 'user', content: question.trim() })
  text.value = ''
  loading.value = true
  scrollDown()

  try {
    const token = localStorage.getItem('token')
    const res = await fetch('http://127.0.0.1:8888/api/v1/chat/ask/stream', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
      body: JSON.stringify({ question: question.trim(), sessionId: sessionId.value })
    })
    const reader = res.body.getReader()
    const decoder = new TextDecoder()
    let buffer = ''
    let aiMsg = null
    while (true) {
      const { done, value } = await reader.read()
      if (value) {
        buffer += decoder.decode(value, { stream: true })
      }
      const lines = buffer.split('\n')
      buffer = lines.pop() || ''
      for (const line of lines) {
        if (!line.startsWith('data: ')) continue
        try {
          const d = JSON.parse(line.slice(6))
          if (d.type === 'session') sessionId.value = d.sessionId
          else if (d.type === 'text') {
            if (!aiMsg) { aiMsg = { role: 'ai', content: '' }; messages.value.push(aiMsg) }
            aiMsg.content += d.content; scrollDown()
          }
          else if (d.type === 'error') {
            if (!aiMsg) { aiMsg = { role: 'ai', content: '' }; messages.value.push(aiMsg) }
            aiMsg.content = d.msg; break
          }
        } catch {}
      }
      if (done) {
        // 处理 buffer 中剩余的数据（最后可能没有 \n 结尾）
        if (buffer.startsWith('data: ')) {
          try {
            const d = JSON.parse(buffer.slice(6))
            if (d.type === 'text' && aiMsg) { aiMsg.content += d.content }
          } catch {}
        }
        break
      }
    }
    if (!aiMsg) { aiMsg = { role: 'ai', content: '抱歉，我暂时无法回答这个问题。' }; messages.value.push(aiMsg) }
  } catch {
    if (!aiMsg) { aiMsg = { role: 'ai', content: '' }; messages.value.push(aiMsg) }
    aiMsg.content = '网络错误，请稍后重试。'
  } finally {
    loading.value = false
    scrollDown()
  }
}

onMounted(() => {
  sessionId.value = localStorage.getItem('aiSessionId') || ''
})
</script>

<style scoped>
.smart-chat { max-width: 700px; margin: 0 auto }
.smart-chat h2 { margin-bottom: 4px }
.subtitle { color: #888; font-size: 13px; margin: 0 0 16px }
.chat-box { background: #fff; border-radius: 12px; padding: 20px; height: 450px; overflow-y: auto; margin-bottom: 12px; box-shadow: 0 1px 4px rgba(0,0,0,.04) }
.welcome { text-align: center; padding: 40px 0 }
.welcome-icon { font-size: 48px; margin-bottom: 12px }
.welcome h3 { margin: 0 0 8px }
.welcome p { color: #888; font-size: 14px; margin-bottom: 16px }
.quick-questions { display: flex; flex-wrap: wrap; gap: 8px; justify-content: center }
.quick-questions span { padding: 6px 14px; background: #e3f2fd; color: #1976d2; border-radius: 16px; font-size: 13px; cursor: pointer; transition: .15s }
.quick-questions span:hover { background: #bbdefb }
.bubble { display: flex; gap: 10px; margin-bottom: 16px }
.bubble.user { flex-direction: row-reverse }
.bubble .avatar { width: 32px; height: 32px; border-radius: 50%; background: #f0f0f0; display: flex; align-items: center; justify-content: center; font-size: 16px; flex-shrink: 0 }
.bubble.user .avatar { background: #1a1a2e }
.bubble.ai .content { background: #f5f5f5; color: #333 }
.bubble.user .content { background: #1a1a2e; color: #fff }
.bubble .content { padding: 10px 16px; border-radius: 14px; max-width: 80%; font-size: 14px; line-height: 1.5; white-space: pre-wrap; word-break: break-word }
.bubble.user .content { border-bottom-right-radius: 4px }
.bubble.ai .content { border-bottom-left-radius: 4px }
.typing { display: flex; gap: 4px; align-items: center; padding: 14px 20px !important }
.typing span { width: 6px; height: 6px; border-radius: 50%; background: #bbb; animation: blink 1.4s infinite }
.typing span:nth-child(2) { animation-delay: .2s }
.typing span:nth-child(3) { animation-delay: .4s }
@keyframes blink { 0%,60%,100% { opacity: .3 } 30% { opacity: 1 } }
.input-row { display: flex; gap: 8px }
.input-row input { flex: 1; padding: 12px; border: 1px solid #ddd; border-radius: 8px; font-size: 14px }
.input-row button { padding: 12px 24px; background: #1a1a2e; color: #fff; border: none; border-radius: 8px; cursor: pointer; font-size: 14px }
.input-row button:disabled { opacity: .5 }
</style>
