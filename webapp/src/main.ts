import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import { useAuthStore } from './stores/auth'

// 样式导入
import './assets/styles/variables.css'
import './assets/styles/base.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)

// 路由守卫依赖 authStore，在挂载前初始化登录态
useAuthStore().init()

// 全局错误边界：捕获未处理的渲染/生命周期异常，避免白屏
app.config.errorHandler = (err, _instance, info) => {
  console.error('[全局错误]', info, err)
}

app.mount('#app')
