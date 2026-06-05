import { defineStore } from 'pinia'
import { ref } from 'vue'
import { md2htmlApi } from '@/services/api/md2html'

export interface HistoryItem {
  id: number
  user_id: number
  markdown: string
  html: string
  theme: string
  created_at: string
}

export const useHistoryStore = defineStore('history', () => {
  const items = ref<HistoryItem[]>([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(20)
  const loading = ref(false)

  async function fetchHistory(p = 1) {
    loading.value = true
    try {
      const res = await md2htmlApi.getHistory(p, pageSize.value)
      items.value = res.data
      total.value = res.total
      page.value = p
    } catch (e) {
      console.error('Failed to fetch history:', e)
    } finally {
      loading.value = false
    }
  }

  async function deleteItem(id: number) {
    try {
      await md2htmlApi.deleteHistory(id)
      await fetchHistory(page.value)
    } catch (e) {
      console.error('Failed to delete history item:', e)
    }
  }

  return { items, total, page, pageSize, loading, fetchHistory, deleteItem }
})
