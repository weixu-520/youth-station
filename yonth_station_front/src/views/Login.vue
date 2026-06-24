<template>
  <div class="auth-page">
    <div class="auth-container">
      <!-- 左侧品牌面板 -->
      <div class="brand-panel">
        <div class="brand-top">
          <div class="brand-logo">🏠</div>
          <h1 class="brand-title">云 驿</h1>
          <p class="brand-desc">青年人才驿站 · 让每一个梦想都有落脚的地方</p>
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

      <!-- 右侧登录表单 -->
      <div class="form-panel">
        <div class="form-header">
          <div class="welcome-icon">👋</div>
          <h2>欢迎回来</h2>
          <p>登录你的账号，查看申请进度或提交新的驿站预约</p>
        </div>

        <div v-if="errorMsg" class="error-msg">
          <span>⚠️</span> {{ errorMsg }}
        </div>

        <form @submit.prevent="handleLogin">
          <div class="form-group">
            <label>账号</label>
            <div class="input-wrapper">
              <span class="icon">👤</span>
              <input
                v-model="form.account"
                type="text"
                placeholder="用户名 / 手机号"
                autocomplete="username"
              />
            </div>
          </div>

          <div class="form-group">
            <label>密码</label>
            <div class="input-wrapper">
              <span class="icon">🔒</span>
              <input
                v-model="form.password"
                type="password"
                placeholder="输入密码"
                autocomplete="current-password"
              />
            </div>
          </div>

          <div class="form-extra">
            <label>
              <input type="checkbox" v-model="rememberMe" /> 记住登录
            </label>
            <a href="#">忘记密码？</a>
          </div>

          <button type="submit" class="btn btn-primary" :disabled="loading">
            {{ loading ? '验证中...' : '登 录' }}
          </button>
        </form>

        <div class="form-footer">
          还没有账号？<router-link to="/register">立即注册 →</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { login } from '../api/auth'

const router = useRouter()

const features = [
  { icon: '🏘️', title: '覆盖全区', desc: '南昌各大区域均有驿站可选' },
  { icon: '💰', title: '免费入住', desc: '应届生享受 3-7 天免费住宿' },
  { icon: '⚡', title: '快速审核', desc: '提交申请后 1-2 个工作日反馈' },
]

const form = reactive({
  account: '',
  password: ''
})

const rememberMe = ref(false)
const loading = ref(false)
const errorMsg = ref('')

async function handleLogin() {
  errorMsg.value = ''
  if (!form.account.trim()) { errorMsg.value = '请输入账号'; return }
  if (!form.password) { errorMsg.value = '请输入密码'; return }

  loading.value = true
  try {
    const res = await login(form.account.trim(), form.password)
    localStorage.setItem('token', res.data.token)
    localStorage.setItem('userId', res.data.userId)
    localStorage.setItem('userName', res.data.userName)
    localStorage.setItem('isAdmin', res.data.isAdmin || false)
    router.push('/home')
  } catch (e) {
    errorMsg.value = e.message || '登录失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>
