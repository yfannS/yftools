<template>
  <div class="login-view">
    <!-- 动态背景光晕 -->
    <div class="ambient-bg">
      <div class="orb orb-1"></div>
      <div class="orb orb-2"></div>
      <div class="orb orb-3"></div>
    </div>

    <!-- 网格纹理 -->
    <div class="grid-texture"></div>

    <div class="login-card" :class="{ 'shake': shakeCard }">
      <!-- 顶部装饰线 -->
      <div class="card-top-accent"></div>

      <div class="card-header">
        <div class="tabs">
          <button
            class="tab"
            :class="{ active: mode === 'login' }"
            @click="switchMode('login')"
          >登录</button>
          <button
            class="tab"
            :class="{ active: mode === 'register' }"
            @click="switchMode('register')"
          >注册</button>
          <div class="tab-indicator" :class="mode"></div>
        </div>
        <router-link to="/md2html" class="back-btn" title="返回编辑器">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="19" y1="12" x2="5" y2="12"/><polyline points="12 19 5 12 12 5"/></svg>
        </router-link>
      </div>

      <form @submit.prevent="mode === 'login' ? onLogin() : onRegister()" class="login-form">
        <Transition name="form-fade" mode="out-in">
          <div :key="mode" class="form-fields">
            <div class="form-group" :class="{ focused: focusedField === 'username', filled: username }">
              <label for="login-username">用户名</label>
              <div class="input-wrap">
                <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                <input
                  id="login-username"
                  v-model="username"
                  type="text"
                  required
                  autocomplete="username"
                  @focus="focusedField = 'username'"
                  @blur="focusedField = ''"
                  placeholder="请输入用户名"
                />
                <div class="input-glow"></div>
              </div>
            </div>

            <div class="form-group" :class="{ focused: focusedField === 'password', filled: password }">
              <label for="login-password">密码</label>
              <div class="input-wrap">
                <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
                <input
                  id="login-password"
                  v-model="password"
                  type="password"
                  required
                  :autocomplete="mode === 'login' ? 'current-password' : 'new-password'"
                  @focus="focusedField = 'password'"
                  @blur="focusedField = ''"
                  placeholder="请输入密码"
                />
                <div class="input-glow"></div>
              </div>
            </div>
          </div>
        </Transition>

        <button
          type="submit"
          class="submit-btn"
          :disabled="loading"
          :class="{ loading: loading }"
        >
          <span class="btn-text">{{ submitText }}</span>
          <span v-if="loading" class="btn-loader">
            <span class="loader-dot"></span>
            <span class="loader-dot"></span>
            <span class="loader-dot"></span>
          </span>
        </button>

        <Transition name="message-slide">
          <p v-if="success" class="message success">{{ success }}</p>
        </Transition>
        <Transition name="message-slide">
          <p v-if="error" class="message error">{{ error }}</p>
        </Transition>
      </form>
    </div>

    <!-- 底部版权 -->
    <div class="login-footer">
      <span>Markdown Toolbox</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from '@/composables/useToast'
import { getApiErrorMessage } from '@/services/api/errorHandling'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const mode = ref<'login' | 'register'>('login')
const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')
const success = ref('')
const focusedField = ref('')
const shakeCard = ref(false)

const submitText = computed(() => {
  if (loading.value) return ''
  return mode.value === 'login' ? '登录' : '注册'
})

function switchMode(newMode: 'login' | 'register') {
  if (mode.value === newMode) return
  mode.value = newMode
  error.value = ''
  success.value = ''
}

function triggerShake() {
  shakeCard.value = true
  setTimeout(() => { shakeCard.value = false }, 500)
}

async function onLogin() {
  loading.value = true
  error.value = ''
  success.value = ''
  try {
    await authStore.login(username.value, password.value)
    const redirect = (router.currentRoute.value.query.redirect as string) || '/md2html'
    router.push(redirect)
  } catch (e) {
    error.value = getApiErrorMessage(e, '登录失败')
    showToast(error.value, 'err')
    triggerShake()
  } finally {
    loading.value = false
  }
}

async function onRegister() {
  loading.value = true
  error.value = ''
  success.value = ''
  try {
    await authStore.register(username.value, password.value)
    success.value = '注册成功，请登录'
    showToast(success.value, 'ok')
    mode.value = 'login'
  } catch (e) {
    error.value = getApiErrorMessage(e, '注册失败')
    showToast(error.value, 'err')
    triggerShake()
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
/* ===== 页面容器 ===== */
.login-view {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  padding: 24px;
  position: relative;
  overflow: hidden;
}

/* ===== 动态背景光晕 ===== */
.ambient-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
  z-index: 0;
}

.orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.4;
  animation: orbFloat 12s ease-in-out infinite;
}

.orb-1 {
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(100, 116, 139, 0.12) 0%, transparent 70%);
  top: -100px;
  left: -100px;
  animation-delay: 0s;
}

.orb-2 {
  width: 350px;
  height: 350px;
  background: radial-gradient(circle, rgba(100, 116, 139, 0.08) 0%, transparent 70%);
  bottom: -80px;
  right: -80px;
  animation-delay: -4s;
  animation-duration: 15s;
}

.orb-3 {
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, rgba(148, 163, 184, 0.06) 0%, transparent 70%);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: -8s;
  animation-duration: 18s;
}

@keyframes orbFloat {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  25% {
    transform: translate(30px, -20px) scale(1.05);
  }
  50% {
    transform: translate(-20px, 30px) scale(0.95);
  }
  75% {
    transform: translate(20px, 20px) scale(1.02);
  }
}

/* ===== 网格纹理 ===== */
.grid-texture {
  position: absolute;
  inset: 0;
  pointer-events: none;
  z-index: 0;
  opacity: 0.03;
  background-image:
    linear-gradient(var(--border-strong) 1px, transparent 1px),
    linear-gradient(90deg, var(--border-strong) 1px, transparent 1px);
  background-size: 40px 40px;
  mask-image: radial-gradient(ellipse at center, black 30%, transparent 70%);
  -webkit-mask-image: radial-gradient(ellipse at center, black 30%, transparent 70%);
}

/* ===== 卡片 ===== */
.login-card {
  position: relative;
  width: 100%;
  max-width: 380px;
  background: rgba(var(--surface-rgb, 255, 255, 255), 0.7);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 32px;
  z-index: 1;
  animation: cardAppear 0.6s cubic-bezier(0.16, 1, 0.3, 1) both;
  box-shadow:
    0 1px 2px rgba(0, 0, 0, 0.02),
    0 4px 8px rgba(0, 0, 0, 0.02),
    0 8px 16px rgba(0, 0, 0, 0.02),
    0 16px 32px rgba(0, 0, 0, 0.03);
}

:root[data-theme="dark"] .login-card {
  background: rgba(17, 17, 17, 0.7);
  box-shadow:
    0 1px 2px rgba(0, 0, 0, 0.1),
    0 4px 8px rgba(0, 0, 0, 0.1),
    0 8px 16px rgba(0, 0, 0, 0.1),
    0 16px 32px rgba(0, 0, 0, 0.15);
}

@keyframes cardAppear {
  from {
    opacity: 0;
    transform: translateY(24px) scale(0.97);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

/* 抖动动画 */
.login-card.shake {
  animation: cardShake 0.5s cubic-bezier(0.36, 0.07, 0.19, 0.97) both;
}

@keyframes cardShake {
  10%, 90% { transform: translateX(-1px); }
  20%, 80% { transform: translateX(2px); }
  30%, 50%, 70% { transform: translateX(-3px); }
  40%, 60% { transform: translateX(3px); }
}

/* 顶部装饰线 */
.card-top-accent {
  position: absolute;
  top: 0;
  left: 32px;
  right: 32px;
  height: 2px;
  background: linear-gradient(90deg, transparent, var(--accent), transparent);
  opacity: 0.5;
  border-radius: 1px;
}

/* ===== 卡片头部 ===== */
.card-header {
  display: flex;
  align-items: center;
  margin-bottom: 28px;
}

.tabs {
  display: flex;
  gap: 0;
  flex: 1;
  position: relative;
}

.tab {
  flex: 1;
  height: 40px;
  border: none;
  background: transparent;
  color: var(--text-tertiary);
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  position: relative;
  z-index: 1;
  transition: color 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  letter-spacing: -0.01em;
}

.tab:hover {
  color: var(--text-secondary);
}

.tab.active {
  color: var(--accent);
}

/* 滑动指示器 */
.tab-indicator {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 50%;
  height: 2px;
  background: var(--accent);
  border-radius: 1px;
  transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  box-shadow: 0 -1px 8px rgba(100, 116, 139, 0.25);
}

.tab-indicator.register {
  transform: translateX(100%);
}

.back-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: var(--radius-sm);
  color: var(--text-tertiary);
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  flex-shrink: 0;
  margin-left: 8px;
}

.back-btn:hover {
  color: var(--text);
  background: var(--surface-raised);
  transform: translateX(-2px);
}

.back-btn svg {
  transition: transform 0.2s;
}

.back-btn:hover svg {
  transform: translateX(-2px);
}

/* ===== 表单 ===== */
.login-form {
  display: flex;
  flex-direction: column;
}

.form-fields {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 6px;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  transition: color 0.2s;
}

.form-group.focused label {
  color: var(--accent);
}

.input-wrap {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 12px;
  color: var(--text-tertiary);
  transition: color 0.2s;
  z-index: 2;
  pointer-events: none;
}

.form-group.focused .input-icon {
  color: var(--accent);
}

.form-group input {
  width: 100%;
  height: 42px;
  padding: 0 12px 0 38px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 14px;
  background: var(--bg);
  color: var(--text);
  outline: none;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  position: relative;
  z-index: 1;
}

.form-group input::placeholder {
  color: var(--text-tertiary);
  opacity: 0.6;
}

.form-group input:focus {
  border-color: var(--accent);
  background: var(--surface);
}

/* 输入框发光效果 */
.input-glow {
  position: absolute;
  inset: -1px;
  border-radius: calc(var(--radius-sm) + 1px);
  opacity: 0;
  transition: opacity 0.3s;
  pointer-events: none;
  z-index: 0;
  box-shadow: 0 0 0 3px var(--accent-soft);
}

.form-group.focused .input-glow {
  opacity: 1;
}

/* ===== 提交按钮 ===== */
.submit-btn {
  position: relative;
  width: 100%;
  height: 42px;
  border: none;
  border-radius: var(--radius-sm);
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-hover) 100%);
  color: white;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  margin-top: 8px;
  overflow: hidden;
  transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  box-shadow: 0 2px 8px rgba(100, 116, 139, 0.15);
}

.submit-btn::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(255,255,255,0.1) 0%, transparent 50%);
  opacity: 0;
  transition: opacity 0.2s;
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 16px rgba(100, 116, 139, 0.2);
}

.submit-btn:hover:not(:disabled)::before {
  opacity: 1;
}

.submit-btn:active:not(:disabled) {
  transform: translateY(0) scale(0.98);
  box-shadow: 0 1px 4px rgba(100, 116, 139, 0.15);
  transition-duration: 0.1s;
}

.submit-btn:disabled {
  opacity: 0.8;
  cursor: not-allowed;
}

.btn-text {
  display: inline-block;
  transition: opacity 0.2s, transform 0.2s;
}

.submit-btn.loading .btn-text {
  opacity: 0;
  transform: scale(0.8);
}

/* 按钮加载动画 */
.btn-loader {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}

.submit-btn.loading .btn-loader {
  opacity: 1;
}

.loader-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: white;
  animation: loaderBounce 1.2s ease-in-out infinite;
}

.loader-dot:nth-child(1) { animation-delay: 0s; }
.loader-dot:nth-child(2) { animation-delay: 0.15s; }
.loader-dot:nth-child(3) { animation-delay: 0.3s; }

@keyframes loaderBounce {
  0%, 80%, 100% { transform: scale(0.6); opacity: 0.5; }
  40% { transform: scale(1); opacity: 1; }
}

/* ===== 消息提示 ===== */
.message {
  font-size: 13px;
  margin-top: 14px;
  padding: 10px 14px;
  border-radius: var(--radius-sm);
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 8px;
}

.message::before {
  content: '';
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.message.success {
  color: var(--ok);
  background: rgba(90, 143, 106, 0.06);
  border: 1px solid rgba(90, 143, 106, 0.12);
}

.message.success::before {
  background: var(--ok);
  box-shadow: 0 0 6px rgba(90, 143, 106, 0.3);
}

.message.error {
  color: var(--danger);
  background: rgba(180, 86, 86, 0.06);
  border: 1px solid rgba(180, 86, 86, 0.12);
}

.message.error::before {
  background: var(--danger);
  box-shadow: 0 0 6px rgba(180, 86, 86, 0.3);
}

/* ===== 底部 ===== */
.login-footer {
  position: relative;
  z-index: 1;
  margin-top: 24px;
  font-size: 12px;
  color: var(--text-tertiary);
  opacity: 0.6;
  letter-spacing: 0.05em;
  animation: footerFade 0.8s 0.3s cubic-bezier(0.16, 1, 0.3, 1) both;
}

@keyframes footerFade {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 0.6; transform: translateY(0); }
}

/* ===== Vue Transition 动画 ===== */

/* form-fade: 登录/注册切换 */
.form-fade-enter-active,
.form-fade-leave-active {
  transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.form-fade-enter-from {
  opacity: 0;
  transform: translateX(12px);
}

.form-fade-leave-to {
  opacity: 0;
  transform: translateX(-12px);
}

/* message-slide: 提示消息 */
.message-slide-enter-active,
.message-slide-leave-active {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.message-slide-enter-from,
.message-slide-leave-to {
  opacity: 0;
  transform: translateY(-8px) scale(0.97);
  max-height: 0;
  margin-top: 0;
  padding-top: 0;
  padding-bottom: 0;
}

/* 响应式 */
@media (max-width: 480px) {
  .login-card {
    padding: 24px 20px;
  }

  .card-top-accent {
    left: 20px;
    right: 20px;
  }
}
</style>
