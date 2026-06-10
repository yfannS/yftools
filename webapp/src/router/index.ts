import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/HomeView.vue'),
    meta: { title: '工具箱' }
  },
  {
    path: '/md2html',
    name: 'md2html',
    component: () => import('@/views/Md2HtmlView.vue'),
    meta: { title: 'Markdown → HTML' }
  },
  {
    path: '/history',
    name: 'history',
    component: () => import('@/views/HistoryView.vue'),
    meta: { title: '历史记录', requiresAuth: true }
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/LoginView.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/settings',
    name: 'settings',
    component: () => import('@/views/SettingsView.vue'),
    meta: { title: '设置' }
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  // 更新页面标题
  document.title = `${to.meta.title || '工具箱'} - Toolbox`

  // 需要认证但未登录
  if (to.meta.requiresAuth) {
    const token = localStorage.getItem('m2h_webapp_token')
    const expiresAt = localStorage.getItem('m2h_webapp_expires_at')
    const expired = !!expiresAt && !Number.isNaN(new Date(expiresAt).getTime()) && new Date(expiresAt).getTime() <= Date.now()
    if (!token || expired) {
      if (expired) {
        localStorage.removeItem('m2h_webapp_token')
        localStorage.removeItem('m2h_webapp_username')
        localStorage.removeItem('m2h_webapp_expires_at')
      }
      next({ name: 'login', query: { redirect: to.fullPath } })
      return
    }
  }

  next()
})

export default router
