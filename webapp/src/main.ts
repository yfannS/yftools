import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'

// 样式导入
import './assets/styles/variables.css'
import './assets/styles/base.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
