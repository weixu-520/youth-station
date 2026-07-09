<template>
  <div class="knowledge-page">
    <h2>📚 知识库管理</h2>
    <p class="subtitle">上传文档到知识库，智能客服将基于这些内容回答问题</p>

    <div v-if="msg" :class="msgType">{{ msg }}</div>

    <div class="form-card">
      <div class="form-group">
        <label>文档标题 <span class="req">*</span></label>
        <input v-model="form.title" placeholder="例如：南昌市青年驿站入住须知" />
      </div>
      <div class="form-group">
        <label>分类</label>
        <select v-model="form.category">
          <option value="">未分类</option>
          <option value="入住政策">入住政策</option>
          <option value="申请流程">申请流程</option>
          <option value="押金说明">押金说明</option>
          <option value="常见问题">常见问题</option>
          <option value="驿站介绍">驿站介绍</option>
        </select>
      </div>
      <div class="form-group">
        <label>文档内容 <span class="req">*</span></label>
        <textarea v-model="form.content" rows="10" placeholder="输入知识库文档内容..."></textarea>
      </div>
      <button class="submit-btn" @click="upload" :disabled="uploading || !form.title.trim() || !form.content.trim()">
        {{ uploading ? '上传中...' : '上传到知识库' }}
      </button>
    </div>

    <div class="tips">
      <h4>💡 提示</h4>
      <ul>
        <li>文档内容会被向量化后存入知识库，智能客服将基于相似度检索相关文档</li>
        <li>建议每篇文档聚焦一个主题，内容清晰简洁</li>
        <li>相似度阈值 {{ threshold }}，只有匹配度高于此值的文档才会被引用</li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { uploadKnowledge } from '../api/chat'

const form = reactive({ title: '', content: '', category: '' })
const uploading = ref(false)
const msg = ref('')
const msgType = ref('success')
const threshold = ref(0.7)

async function upload() {
  if (!form.title.trim() || !form.content.trim()) return
  uploading.value = true
  msg.value = ''
  try {
    await uploadKnowledge({ ...form })
    msg.value = '上传成功！文档已加入知识库'
    msgType.value = 'success'
    form.title = ''
    form.content = ''
    form.category = ''
  } catch (e) {
    msg.value = e.message || '上传失败'
    msgType.value = 'error'
  } finally {
    uploading.value = false
  }
}
</script>

<style scoped>
.knowledge-page { max-width: 650px; margin: 0 auto }
.knowledge-page h2 { margin-bottom: 4px }
.subtitle { color: #888; font-size: 13px; margin: 0 0 16px }
.success { background: #e8f5e9; color: #2e7d32; padding: 10px 14px; border-radius: 6px; margin-bottom: 12px; font-size: 13px }
.error { background: #ffebee; color: #c62828; padding: 10px 14px; border-radius: 6px; margin-bottom: 12px; font-size: 13px }
.form-card { background: #fff; border-radius: 12px; padding: 24px; box-shadow: 0 1px 4px rgba(0,0,0,.04); margin-bottom: 16px }
.form-group { margin-bottom: 16px }
.form-group label { display: block; margin-bottom: 4px; font-size: 13px; font-weight: 600 }
.req { color: #f44336 }
input, select, textarea { width: 100%; padding: 10px 12px; border: 1px solid #ddd; border-radius: 6px; font-size: 14px; box-sizing: border-box }
textarea { resize: vertical; font-family: inherit }
.submit-btn { width: 100%; padding: 12px; background: #1a1a2e; color: #fff; border: none; border-radius: 8px; font-size: 15px; cursor: pointer }
.submit-btn:disabled { opacity: .5 }
.tips { background: #fffbe6; border: 1px solid #ffd54f; border-radius: 8px; padding: 16px }
.tips h4 { margin: 0 0 8px }
.tips ul { margin: 0; padding-left: 18px; font-size: 13px; color: #666 }
.tips li { margin-bottom: 4px }
</style>
