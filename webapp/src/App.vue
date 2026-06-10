<template>
  <router-view />
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from '@/composables/useToast'
import { AUTH_EXPIRED_EVENT } from '@/services/api/client'
import { mapApiErrorMessage } from '@/services/api/errorHandling'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()
authStore.init()

function handleAuthExpired(event: Event) {
  const customEvent = event as CustomEvent<{ code?: string; message?: string }>
  const code = customEvent.detail?.code
  const nextMessage = mapApiErrorMessage(code, customEvent.detail?.message || '登录已失效，请重新登录')

  authStore.logout()
  showToast(nextMessage, 'err')

  if (router.currentRoute.value.name === 'login') return

  router.push({
    name: 'login',
    query: { redirect: router.currentRoute.value.fullPath },
  })
}

onMounted(() => {
  window.addEventListener(AUTH_EXPIRED_EVENT, handleAuthExpired as EventListener)
})

onBeforeUnmount(() => {
  window.removeEventListener(AUTH_EXPIRED_EVENT, handleAuthExpired as EventListener)
})
</script>
