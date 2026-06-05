import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/services/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(null)
  const username = ref<string | null>(null)
  const isLoggedIn = computed(() => !!token.value)

  function init() {
    const saved = localStorage.getItem('m2h_webapp_token')
    const savedUser = localStorage.getItem('m2h_webapp_username')
    if (saved) {
      token.value = saved
      username.value = savedUser
    }
  }

  async function login(user: string, password: string) {
    const res = await authApi.login(user, password)
    token.value = res.token
    username.value = res.username
    localStorage.setItem('m2h_webapp_token', res.token)
    localStorage.setItem('m2h_webapp_username', res.username)
  }

  async function register(user: string, password: string) {
    await authApi.register(user, password)
  }

  function logout() {
    token.value = null
    username.value = null
    localStorage.removeItem('m2h_webapp_token')
    localStorage.removeItem('m2h_webapp_username')
  }

  return { token, username, isLoggedIn, init, login, register, logout }
})
