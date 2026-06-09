import { ref, watch, type Ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useEditorStore } from '@/stores/editor'
import { md2htmlApi } from '@/services/api/md2html'

type ToastFn = (msg: string, type?: string) => void

/**
 * 自动保存 composable
 * 
 * 逻辑：
 * 1. 登录态下，markdown 内容变化后触发保存
 * 2. 3s 防抖，避免频繁保存
 * 3. 仅在 dirty（用户修改过）时才保存
 * 4. 保存只传 markdown + theme，不传 html
 * 5. 后端自动提取标题和计算字符数
 * 6. 后端 UPSERT 机制保证 5 分钟内同标题不重复
 */
export function useAutoSave(toast: ToastFn, renderReady: Ref<boolean>) {
  const authStore = useAuthStore()
  const editorStore = useEditorStore()

  const saveTimer = ref<ReturnType<typeof setTimeout> | null>(null)
  const isSaving = ref(false)
  const lastSavedContent = ref('')
  const lastSaveTime = ref(0)

  const SAVE_DEBOUNCE = 3000 // 3s 防抖
  const MIN_INTERVAL = 5000  // 最小保存间隔 5s

  function saveToServer() {
    if (!authStore.isLoggedIn) return
    if (!editorStore.content.trim()) return
    if (!editorStore.dirty) return
    if (editorStore.content === lastSavedContent.value) return

    const now = Date.now()
    if (now - lastSaveTime.value < MIN_INTERVAL) return

    isSaving.value = true
    editorStore.markSaving()
    lastSavedContent.value = editorStore.content
    lastSaveTime.value = now

    md2htmlApi.saveHistory(editorStore.content, 'default')
      .then(() => {
        editorStore.markSaved()
      })
      .catch((err) => {
        console.error('Auto-save failed:', err)
        // 保存失败重置状态以便下次重试
        lastSavedContent.value = ''
        editorStore.markSaveError()
      })
      .finally(() => {
        isSaving.value = false
      })
  }

  function scheduleSave() {
    if (!authStore.isLoggedIn) return
    if (!renderReady.value) return

    if (saveTimer.value) {
      clearTimeout(saveTimer.value)
    }
    saveTimer.value = setTimeout(() => {
      saveToServer()
    }, SAVE_DEBOUNCE)
  }

  function retrySave() {
    if (saveTimer.value) {
      clearTimeout(saveTimer.value)
      saveTimer.value = null
    }
    lastSavedContent.value = ''
    lastSaveTime.value = 0
    saveToServer()
  }

  // 监听内容变化
  watch(
    () => editorStore.content,
    () => {
      if (authStore.isLoggedIn && editorStore.content.trim()) {
        scheduleSave()
      }
    }
  )

  // 监听登录状态变化
  watch(
    () => authStore.isLoggedIn,
    (loggedIn) => {
      if (loggedIn && editorStore.content.trim()) {
        // 刚登录时立即触发一次保存
        setTimeout(() => saveToServer(), 1000)
      }
    }
  )

  function dispose() {
    if (saveTimer.value) {
      clearTimeout(saveTimer.value)
      saveTimer.value = null
    }
    // 退出时如果有未保存的修改，尝试保存一次
    if (authStore.isLoggedIn && editorStore.dirty && editorStore.content.trim()) {
      md2htmlApi.saveHistory(editorStore.content, 'default').then(() => {
        editorStore.markSaved()
      }).catch(() => {
        editorStore.markSaveError()
      })
    }
  }

  return { isSaving, scheduleSave, retrySave, dispose }
}
