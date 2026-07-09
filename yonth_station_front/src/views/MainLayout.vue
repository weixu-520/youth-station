<template>
  <div class="layout">
    <aside class="sidebar" :class="{ collapsed: collapsed }">
      <div class="sidebar-brand" @click="$router.push('/home')">
        <span class="logo">🏠</span>
        <span v-show="!collapsed" class="title">云 驿</span>
      </div>
      <nav class="sidebar-nav">
        <router-link to="/home" class="nav-item"><span class="icon">📊</span><span v-show="!collapsed">首页</span></router-link>
        <router-link to="/home/stations" class="nav-item"><span class="icon">🏘️</span><span v-show="!collapsed">驿站列表</span></router-link>
        <router-link to="/home/applications" class="nav-item"><span class="icon">📋</span><span v-show="!collapsed">我的申请</span></router-link>
        <router-link to="/home/smart-chat" class="nav-item"><span class="icon">🤖</span><span v-show="!collapsed">智能客服</span></router-link>
        <router-link to="/home/chat" class="nav-item"><span class="icon">💬</span><span v-show="!collapsed">联系客服</span></router-link>
        <router-link to="/home/profile" class="nav-item"><span class="icon">👤</span><span v-show="!collapsed">个人中心</span></router-link>
        <template v-if="isAdmin">
          <div class="nav-divider" v-show="!collapsed"></div>
          <div class="nav-label" v-show="!collapsed">管理</div>
          <router-link to="/home/admin/dashboard" class="nav-item admin"><span class="icon">📈</span><span v-show="!collapsed">数据概览</span></router-link>
          <router-link to="/home/admin/applications" class="nav-item admin"><span class="icon">✅</span><span v-show="!collapsed">申请审核</span></router-link>
          <router-link to="/home/admin/stations" class="nav-item admin"><span class="icon">⚙️</span><span v-show="!collapsed">驿站管理</span></router-link>
          <router-link to="/home/admin/knowledge" class="nav-item admin"><span class="icon">📚</span><span v-show="!collapsed">知识库管理</span></router-link>
        </template>
      </nav>
      <div class="sidebar-footer">
        <button class="collapse-btn" @click="collapsed=!collapsed">{{ collapsed ? '▶' : '◀' }}</button>
        <button class="logout-btn" @click="logout"><span>🚪</span><span v-show="!collapsed">退出</span></button>
      </div>
    </aside>
    <main class="content">
      <header class="topbar">
        <h2>{{ $route.meta.title }}</h2>
        <span class="user-tag">{{ userName }}<span v-if="isAdmin" class="admin-badge">管理员</span></span>
      </header>
      <div class="page"><router-view /></div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const collapsed = ref(false)
const userName = ref(localStorage.getItem('userName') || '用户')
const isAdmin = computed(() => localStorage.getItem('isAdmin') === 'true')

function logout() {
  localStorage.clear()
  router.push('/login')
}
</script>

<style scoped>
.layout { display: flex; min-height: 100vh; background: #f0f2f5 }
.sidebar { width: 220px; background: #1a1a2e; color: #ccc; display: flex; flex-direction: column; transition: width .2s; padding: 16px 0 }
.sidebar.collapsed { width: 60px }
.sidebar-brand { display: flex; align-items: center; gap: 10px; padding: 8px 18px; cursor: pointer }
.logo { font-size: 28px }
.title { font-size: 20px; font-weight: 700; color: #fff }
.sidebar-nav { flex: 1; padding: 8px 0 }
.nav-item { display: flex; align-items: center; gap: 12px; padding: 12px 18px; color: #aaa; text-decoration: none; font-size: 14px; transition: .15s }
.nav-item:hover, .nav-item.router-link-active { background: #16213e; color: #64b5f6 }
.nav-item .icon { font-size: 18px; width: 22px; text-align: center }
.nav-item.admin { color: #81c784 }
.nav-divider { margin: 8px 18px; border-top: 1px solid #333 }
.nav-label { padding: 8px 18px 4px; font-size: 11px; text-transform: uppercase; color: #666; letter-spacing: 1px }
.sidebar-footer { padding: 8px; display: flex; justify-content: space-between }
.collapse-btn, .logout-btn { background: none; border: none; color: #888; cursor: pointer; font-size: 14px; padding: 6px 10px }
.logout-btn { display: flex; align-items: center; gap: 6px }
/* content */
.content { flex: 1; display: flex; flex-direction: column; overflow: hidden }
.topbar { display: flex; align-items: center; justify-content: space-between; padding: 16px 24px; background: #fff; box-shadow: 0 1px 4px rgba(0,0,0,.06) }
.topbar h2 { margin: 0; font-size: 18px }
.user-tag { display: flex; align-items: center; gap: 8px; font-size: 14px }
.admin-badge { background: #ff9800; color: #fff; padding: 2px 8px; border-radius: 10px; font-size: 11px }
.page { flex: 1; padding: 20px 24px; overflow-y: auto }
</style>
