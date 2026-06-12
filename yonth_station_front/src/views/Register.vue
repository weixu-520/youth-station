<template>
  <div class="auth-page">
    <div class="auth-container">
      <!-- 左侧品牌面板 -->
      <div class="brand-panel">
        <div class="brand-top">
          <div class="brand-logo">🌟</div>
          <h1 class="brand-title">加入云驿</h1>
          <p class="brand-desc">创建你的专属账号，开启青年驿站之旅</p>
        </div>

        <div class="brand-features">
          <div class="brand-feature" v-for="f in features" :key="f.icon">
            <span class="feat-icon">{{ f.icon }}</span>
            <div class="feat-text">
              <h4>{{ f.title }}</h4>
              <span>{{ f.desc }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧注册表单 -->
      <div class="form-panel">
        <div class="form-header">
          <div class="welcome-icon">✨</div>
          <h2>创建账号</h2>
          <p>填写基本信息，30 秒完成注册</p>
        </div>

        <div v-if="errorMsg" class="error-msg">
          <span>⚠️</span> {{ errorMsg }}
        </div>
        <div v-if="successMsg" class="success-msg">
          <span>✅</span> {{ successMsg }}
        </div>

        <form @submit.prevent="handleRegister">
          <div class="form-group">
            <label>用户名 <span class="required">*</span></label>
            <div class="input-wrapper">
              <span class="icon">👤</span>
              <input
                v-model="form.userName"
                type="text"
                placeholder="给自己取一个名字"
                autocomplete="username"
              />
            </div>
          </div>

          <div class="form-group">
            <label>手机号 <span class="optional">（选填，建议填写以接收通知）</span></label>
            <div class="input-wrapper">
              <span class="icon">📱</span>
              <input
                v-model="form.phone"
                type="tel"
                placeholder="11 位手机号码"
                autocomplete="tel"
              />
            </div>
          </div>

          <div class="form-group">
            <label>密码 <span class="required">*</span></label>
            <div class="input-wrapper">
              <span class="icon">🔒</span>
              <input
                v-model="form.password"
                type="password"
                placeholder="至少 6 位，建议使用字母+数字组合"
                autocomplete="new-password"
              />
            </div>
          </div>

          <div class="form-group">
            <label>确认密码 <span class="required">*</span></label>
            <div class="input-wrapper">
              <span class="icon">🔒</span>
              <input
                v-model="confirmPassword"
                type="password"
                placeholder="请再次输入密码"
                autocomplete="new-password"
              />
            </div>
          </div>

          <button type="submit" class="btn btn-primary" :disabled="loading">
            {{ loading ? '注册中...' : '注 册' }}
          </button>
        </form>

        <div class="form-footer">
          已有账号？<router-link to="/login">立即登录 →</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { register } from '../api/auth'

const router = useRouter()

const features = [
  { icon: '📝', title: '简单注册', desc: '仅需用户名和密码即可完成' },
  { icon: '🔐', title: '安全可靠', desc: '密码加密存储，保障账号安全' },
  { icon: '🎯', title: '即刻申请', desc: '注册后即可浏览驿站并提交申请' },
]

const form = reactive({
  userName: '',
  phone: '',
  password: ''
})

const confirmPassword = ref('')
const loading = ref(false)
const errorMsg = ref('')
const successMsg = ref('')

async function handleRegister() {
  errorMsg.value = ''
  successMsg.value = ''

  if (!form.userName.trim()) { errorMsg.value = '请输入用户名'; return }
  if (form.userName.trim().length < 2) { errorMsg.value = '用户名至少 2 个字符'; return }
  if (!form.password) { errorMsg.value = '请输入密码'; return }
  if (form.password.length < 6) { errorMsg.value = '密码至少 6 位'; return }
  if (form.password !== confirmPassword.value) { errorMsg.value = '两次输入的密码不一致'; return }
  if (form.phone && !/^1[3-9]\d{9}$/.test(form.phone)) { errorMsg.value = '手机号格式不正确'; return }

  loading.value = true
  try {
    const payload = { userName: form.userName.trim(), password: form.password }
    if (form.phone) payload.phone = form.phone

    await register(payload)
    successMsg.value = '注册成功！即将跳转登录页...'
    setTimeout(() => router.push('/login'), 1500)
  } catch (e) {
    errorMsg.value = e.message || '注册失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>
