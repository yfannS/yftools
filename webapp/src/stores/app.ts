import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useAppStore = defineStore('app', () => {
  const theme = ref<'light' | 'dark'>('light')
  const sidebarOpen = ref(false)
  const currentTool = ref('md2html')

  function setTheme(t: 'light' | 'dark') {
    theme.value = t
    document.documentElement.setAttribute('data-theme', t)
    localStorage.setItem('m2h_webapp_theme', t)
  }

  function toggleTheme() {
    setTheme(theme.value === 'light' ? 'dark' : 'light')
  }

  function initTheme() {
    const saved = localStorage.getItem('m2h_webapp_theme')
    if (saved === 'dark' || saved === 'light') {
      setTheme(saved)
    } else if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
      setTheme('dark')
    } else {
      setTheme('light')
    }
  }

  // 监听系统主题变化
  watch(
    () => window.matchMedia('(prefers-color-scheme: dark)'),
    () => {
      if (!localStorage.getItem('m2h_webapp_theme')) {
        initTheme()
      }
    }
  )

  return { theme, sidebarOpen, currentTool, setTheme, toggleTheme, initTheme }
})
