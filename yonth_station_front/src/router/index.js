import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', redirect: '/home' },
  { path: '/login', name: 'Login', component: () => import('../views/Login.vue'), meta: { title: '登录' } },
  { path: '/register', name: 'Register', component: () => import('../views/Register.vue'), meta: { title: '注册' } },
  {
    path: '/home',
    component: () => import('../views/MainLayout.vue'),
    children: [
      { path: '', name: 'Home', component: () => import('../views/Home.vue'), meta: { title: '首页' } },
      { path: 'stations', name: 'StationList', component: () => import('../views/StationList.vue'), meta: { title: '驿站列表' } },
      { path: 'station/:id', name: 'StationDetail', component: () => import('../views/StationDetail.vue'), meta: { title: '驿站详情' } },
      { path: 'apply/:stationId?', name: 'Apply', component: () => import('../views/Apply.vue'), meta: { title: '提交申请' } },
      { path: 'applications', name: 'MyApplications', component: () => import('../views/MyApplications.vue'), meta: { title: '我的申请' } },
      { path: 'application/:id', name: 'ApplicationDetail', component: () => import('../views/ApplicationDetail.vue'), meta: { title: '申请详情' } },
      { path: 'chat', name: 'Chat', component: () => import('../views/Chat.vue'), meta: { title: '联系客服' } },
      { path: 'smart-chat', name: 'SmartChat', component: () => import('../views/SmartChat.vue'), meta: { title: '智能客服' } },
      { path: 'profile', name: 'UserProfile', component: () => import('../views/UserProfile.vue'), meta: { title: '个人中心' } },
      { path: 'admin/dashboard', name: 'AdminDashboard', component: () => import('../views/AdminDashboard.vue'), meta: { title: '管理后台' } },
      { path: 'admin/applications', name: 'AdminApplications', component: () => import('../views/AdminApplications.vue'), meta: { title: '申请管理' } },
      { path: 'admin/stations', name: 'AdminStations', component: () => import('../views/AdminStations.vue'), meta: { title: '驿站管理' } },
      { path: 'admin/knowledge', name: 'KnowledgeUpload', component: () => import('../views/KnowledgeUpload.vue'), meta: { title: '知识库管理' } },
    ]
  },
  { path: '/:pathMatch(.*)*', redirect: '/home' }
]

const router = createRouter({ history: createWebHistory(), routes })

router.beforeEach((to, from, next) => {
  document.title = to.meta.title ? `${to.meta.title} — 云驿` : '云驿'
  const token = localStorage.getItem('token')
  const publicPages = ['/login', '/register']
  if (!token && !publicPages.includes(to.path)) return next('/login')
  if (token && publicPages.includes(to.path)) return next('/home')
  next()
})

export default router
