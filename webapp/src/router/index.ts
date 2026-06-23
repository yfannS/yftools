import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

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
    path: '/json-formatter',
    name: 'json-formatter',
    component: () => import('@/views/JsonFormatterView.vue'),
    meta: { title: 'JSON 格式化' }
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
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/views/NotFoundView.vue'),
    meta: { title: '页面不存在' }
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

  // 需要认证的页面：复用 authStore.checkAuth 实时校验登录态（含过期清理）
  if (to.meta.requiresAuth && !useAuthStore().checkAuth()) {
    next({ name: 'login', query: { redirect: to.fullPath } })
    return
  }

  next()
})

export default router
