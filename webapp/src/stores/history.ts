import { defineStore } from 'pinia'
import { ref } from 'vue'
import { md2htmlApi, type HistoryListItem, type HistoryDetail } from '@/services/api/md2html'

export type { HistoryListItem, HistoryDetail }

export const useHistoryStore = defineStore('history', () => {
  const items = ref<HistoryListItem[]>([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(20)
  const loading = ref(false)
  const detailLoading = ref(false)

  /** 获取历史记录列表（轻量，不含 markdown/html） */
  async function fetchHistory(p = 1) {
    loading.value = true
    try {
      const res = await md2htmlApi.getHistory(p, pageSize.value)
      items.value = res.data || []
      total.value = res.total || 0
      page.value = res.page || p
    } catch (e) {
      console.error('Failed to fetch history:', e)
    } finally {
      loading.value = false
    }
  }

  /** 获取历史记录详情（含完整 markdown） */
  async function getDetail(id: number): Promise<HistoryDetail | null> {
    detailLoading.value = true
    try {
      const detail = await md2htmlApi.getHistoryDetail(id)
      return detail
    } catch (e) {
      console.error('Failed to fetch history detail:', e)
      return null
    } finally {
      detailLoading.value = false
    }
  }

  /** 删除历史记录 */
  async function deleteItem(id: number) {
    try {
      await md2htmlApi.deleteHistory(id)
      await fetchHistory(page.value)
    } catch (e) {
      console.error('Failed to delete history item:', e)
    }
  }

  return { items, total, page, pageSize, loading, detailLoading, fetchHistory, getDetail, deleteItem }
})
