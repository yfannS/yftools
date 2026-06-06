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
    const data = await authApi.login(user, password)
    token.value = data.token
    username.value = data.username
    localStorage.setItem('m2h_webapp_token', data.token)
    localStorage.setItem('m2h_webapp_username', data.username)
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
