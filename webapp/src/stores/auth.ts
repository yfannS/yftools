import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { AUTH_EXPIRED_EVENT } from '@/services/api/client'
import { authApi } from '@/services/api/auth'

const TOKEN_KEY = 'm2h_webapp_token'
const USERNAME_KEY = 'm2h_webapp_username'
const EXPIRES_AT_KEY = 'm2h_webapp_expires_at'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(null)
  const username = ref<string | null>(null)
  const expiresAt = ref<string | null>(null)
  let expireTimer: ReturnType<typeof setTimeout> | null = null
  const isLoggedIn = computed(() => !!token.value)

  function init() {
    syncFromStorage()
    if (isExpired(expiresAt.value)) {
      logout(false)
      return
    }
    scheduleExpiryTimer()
  }

  async function login(user: string, password: string) {
    const data = await authApi.login(user, password)
    token.value = data.token
    username.value = data.username
    expiresAt.value = data.expires_at
    localStorage.setItem(TOKEN_KEY, data.token)
    localStorage.setItem(USERNAME_KEY, data.username)
    localStorage.setItem(EXPIRES_AT_KEY, data.expires_at)
    scheduleExpiryTimer()
  }

  async function register(user: string, password: string) {
    await authApi.register(user, password)
  }

  function logout(notifyExpired = false) {
    clearExpiryTimer()
    token.value = null
    username.value = null
    expiresAt.value = null
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USERNAME_KEY)
    localStorage.removeItem(EXPIRES_AT_KEY)

    if (notifyExpired) {
      window.dispatchEvent(new CustomEvent(AUTH_EXPIRED_EVENT, {
        detail: {
          code: 'AUTH_SESSION_EXPIRED',
          message: '登录已失效，请重新登录',
        },
      }))
    }
  }

  function syncFromStorage() {
    const saved = localStorage.getItem(TOKEN_KEY)
    const savedUser = localStorage.getItem(USERNAME_KEY)
    const savedExpiresAt = localStorage.getItem(EXPIRES_AT_KEY)
    token.value = saved
    username.value = savedUser
    expiresAt.value = savedExpiresAt
  }

  function scheduleExpiryTimer() {
    clearExpiryTimer()
    if (!expiresAt.value) return

    const expireAt = new Date(expiresAt.value).getTime()
    if (Number.isNaN(expireAt)) return

    const delay = expireAt - Date.now()
    if (delay <= 0) {
      logout(true)
      return
    }

    expireTimer = setTimeout(() => {
      logout(true)
    }, delay)
  }

  function clearExpiryTimer() {
    if (!expireTimer) return
    clearTimeout(expireTimer)
    expireTimer = null
  }

  function isExpired(value: string | null) {
    if (!value) return false
    const expireAt = new Date(value).getTime()
    if (Number.isNaN(expireAt)) return false
    return expireAt <= Date.now()
  }

  /** 实时校验登录态：token 存在且未过期，过期则清理并返回 false */
  function checkAuth(): boolean {
    if (!token.value) return false
    if (isExpired(expiresAt.value)) {
      logout(false)
      return false
    }
    return true
  }

  return { token, username, expiresAt, isLoggedIn, init, login, register, logout, syncFromStorage, checkAuth }
})
