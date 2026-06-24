<template>
  <div class="profile-page">
    <div class="card">
      <h2>👤 个人中心</h2>
      <div v-if="msg" :class="msgType">{{ msg }}</div>
      <form @submit.prevent="save" v-if="user">
        <div class="form-group">
          <label>用户名</label>
          <input v-model="form.userName" placeholder="用户名" />
        </div>
        <div class="row">
          <div class="form-group">
            <label>性别</label>
            <select v-model="form.gender">
              <option :value="0">未知</option>
              <option :value="1">男</option>
              <option :value="2">女</option>
            </select>
          </div>
          <div class="form-group">
            <label>出生日期</label>
            <input v-model="form.birthDate" type="date" />
          </div>
        </div>
        <div class="row">
          <div class="form-group">
            <label>学历</label>
            <select v-model="form.education">
              <option :value="0">未知</option>
              <option :value="1">大专</option>
              <option :value="2">本科</option>
              <option :value="3">硕士</option>
              <option :value="4">博士</option>
            </select>
          </div>
          <div class="form-group">
            <label>毕业年份</label>
            <input v-model="form.graduateYear" type="number" placeholder="如 2024" />
          </div>
        </div>
        <div class="form-group">
          <label>毕业院校</label>
          <input v-model="form.school" placeholder="毕业院校" />
        </div>
        <div class="form-group">
          <label>户籍城市</label>
          <input v-model="form.hukouCity" placeholder="如 南昌市" />
        </div>
        <div class="info-row">
          <span>📱 {{ user.phone || '未绑定' }}</span>
          <span>🆔 {{ user.idCard || '未认证' }}</span>
          <span :class="user.isAdmin ? 'admin' : ''">{{ user.isAdmin ? '👑 管理员' : '👤 普通用户' }}</span>
        </div>
        <button type="submit" class="save-btn" :disabled="saving">{{ saving ? '保存中...' : '保存修改' }}</button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getUserInfo, updateUserInfo } from '../api/user'

const user = ref(null)
const saving = ref(false)
const msg = ref('')
const msgType = ref('success')

const form = reactive({ userName: '', gender: 0, birthDate: '', education: 0, graduateYear: 0, school: '', hukouCity: '' })

onMounted(async () => {
  try { const r = await getUserInfo(); user.value = r.data; Object.assign(form, r.data) } catch {}
})

async function save() {
  msg.value = ''; saving.value = true
  try {
    const payload = {}
    for (const k of ['userName', 'birthDate', 'school', 'hukouCity']) if (form[k]) payload[k] = form[k]
    for (const k of ['gender', 'education', 'graduateYear']) if (form[k]) payload[k] = Number(form[k])
    await updateUserInfo(payload)
    msg.value = '保存成功'; msgType.value = 'success'
    localStorage.setItem('userName', form.userName)
    setTimeout(async () => { const r = await getUserInfo(); user.value = r.data }, 500)
  } catch (e) { msg.value = e.message; msgType.value = 'error' }
  finally { saving.value = false }
}
</script>

<style scoped>
.profile-page { max-width: 500px; margin: 0 auto }
.card { background: #fff; border-radius: 12px; padding: 24px; box-shadow: 0 1px 4px rgba(0,0,0,.04) }
.card h2 { margin: 0 0 16px }
.success { background: #e8f5e9; color: #2e7d32; padding: 10px; border-radius: 6px; margin-bottom: 12px; font-size: 13px }
.error { background: #ffebee; color: #c62828; padding: 10px; border-radius: 6px; margin-bottom: 12px; font-size: 13px }
.form-group { margin-bottom: 14px }
.form-group label { display: block; margin-bottom: 4px; font-size: 13px; font-weight: 600 }
.row { display: grid; grid-template-columns: 1fr 1fr; gap: 12px }
input, select { width: 100%; padding: 10px 12px; border: 1px solid #ddd; border-radius: 6px; font-size: 14px; box-sizing: border-box }
.info-row { display: flex; gap: 16px; margin: 16px 0; font-size: 13px; color: #666 }
.info-row .admin { color: #ff9800; font-weight: 600 }
.save-btn { width: 100%; padding: 12px; background: #1a1a2e; color: #fff; border: none; border-radius: 8px; font-size: 15px; cursor: pointer }
.save-btn:disabled { opacity: .6 }
</style>
