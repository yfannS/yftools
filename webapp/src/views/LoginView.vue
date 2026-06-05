<template>
  <div class="login-view">
    <div class="login-card">
      <div class="tabs">
        <button class="tab" :class="{ active: mode === 'login' }" @click="mode = 'login'; error = ''">登录</button>
        <button class="tab" :class="{ active: mode === 'register' }" @click="mode = 'register'; error = ''">注册</button>
      </div>
      <form @submit.prevent="mode === 'login' ? onLogin() : onRegister()">
        <div class="form-group">
          <label>用户名</label>
          <input v-model="username" type="text" required autocomplete="username" />
        </div>
        <div class="form-group">
          <label>密码</label>
          <input v-model="password" type="password" required autocomplete="current-password" />
        </div>
        <button type="submit" class="submit-btn" :disabled="loading">
          {{ loading ? (mode === 'login' ? '登录中...' : '注册中...') : (mode === 'login' ? '登录' : '注册') }}
        </button>
        <p v-if="error" class="error">{{ error }}</p>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const mode = ref<'login' | 'register'>('login')
const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

async function onLogin() {
  loading.value = true
  error.value = ''
  try {
    await authStore.login(username.value, password.value)
    const redirect = (router.currentRoute.value.query.redirect as string) || '/md2html'
    router.push(redirect)
  } catch (e: any) {
    error.value = e.message || '登录失败'
  } finally {
    loading.value = false
  }
}

async function onRegister() {
  loading.value = true
  error.value = ''
  try {
    await authStore.register(username.value, password.value)
    await authStore.login(username.value, password.value)
    router.push('/md2html')
  } catch (e: any) {
    error.value = e.message || '注册失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-view {
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1;
  padding: 24px;
}

.login-card {
  width: 100%;
  max-width: 360px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 32px;
}

.tabs {
  display: flex;
  gap: 0;
  margin-bottom: 24px;
  border-bottom: 1px solid var(--border);
}

.tab {
  flex: 1;
  height: 36px;
  border: none;
  background: transparent;
  color: var(--text-tertiary);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  margin-bottom: -1px;
  transition: all 0.12s;
}

.tab:hover { color: var(--text-secondary); }

.tab.active {
  color: var(--accent);
  border-bottom-color: var(--accent);
}

.form-group { margin-bottom: 16px; }

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 6px;
  color: var(--text-secondary);
}

.form-group input {
  width: 100%;
  height: 36px;
  padding: 0 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 14px;
  background: var(--bg);
  color: var(--text);
  outline: none;
}

.form-group input:focus { border-color: var(--accent); }

.submit-btn {
  width: 100%;
  height: 36px;
  border: none;
  border-radius: var(--radius-sm);
  background: var(--accent);
  color: white;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.12s;
}

.submit-btn:hover { background: var(--accent-hover); }
.submit-btn:disabled { opacity: 0.6; cursor: not-allowed; }

.error { color: var(--danger); font-size: 13px; margin-top: 12px; }
</style>
