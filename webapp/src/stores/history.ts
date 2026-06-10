import { defineStore } from 'pinia'
import { ref } from 'vue'
import { showToast } from '@/composables/useToast'
import { md2htmlApi, type HistoryListItem, type HistoryDetail } from '@/services/api/md2html'
import { getApiErrorCode, getApiErrorMessage } from '@/services/api/errorHandling'

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
      items.value = []
      total.value = 0
      const code = getApiErrorCode(e)
      if (code !== 'AUTH_SESSION_EXPIRED' && code !== 'AUTH_TOKEN_INVALID') {
        showToast(getApiErrorMessage(e, '获取历史记录失败'), 'err')
      }
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
      const code = getApiErrorCode(e)
      if (code !== 'AUTH_SESSION_EXPIRED' && code !== 'AUTH_TOKEN_INVALID') {
        showToast(getApiErrorMessage(e, '加载历史记录失败'), 'err')
      }
      return null
    } finally {
      detailLoading.value = false
    }
  }

  /** 删除历史记录 */
  async function deleteItem(id: number) {
    try {
      await md2htmlApi.deleteHistory(id)
      showToast('历史记录已删除', 'ok')
      await fetchHistory(page.value)
    } catch (e) {
      console.error('Failed to delete history item:', e)
      const code = getApiErrorCode(e)
      if (code !== 'AUTH_SESSION_EXPIRED' && code !== 'AUTH_TOKEN_INVALID') {
        showToast(getApiErrorMessage(e, '删除历史记录失败'), 'err')
      }
    }
  }

  /** 修改历史记录标题 */
  async function renameItem(id: number, title: string) {
    try {
      await md2htmlApi.renameHistory(id, title)
      // 本地更新列表中的标题，避免重新请求
      const item = items.value.find(i => i.id === id)
      if (item) item.title = title
      showToast('标题已修改', 'ok')
    } catch (e) {
      console.error('Failed to rename history item:', e)
      showToast(getApiErrorMessage(e, '修改标题失败'), 'err')
      throw e
    }
  }

  return { items, total, page, pageSize, loading, detailLoading, fetchHistory, getDetail, deleteItem, renameItem }
})
