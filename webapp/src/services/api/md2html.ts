import { request } from './client'

// ==================== 历史记录相关类型 ====================

/** 历史列表项（轻量，不含 markdown/html 大字段） */
export interface HistoryListItem {
  id: number
  title: string
  char_count: number
  theme: string
  created_at: string
  updated_at: string
}

/** 历史详情（含完整 markdown） */
export interface HistoryDetail {
  id: number
  title: string
  markdown: string
  html?: string
  char_count: number
  theme: string
  created_at: string
  updated_at: string
}

/** 分页响应 */
export interface HistoryPageData {
  data: HistoryListItem[]
  total: number
  page: number
  size: number
}

/** 保存响应 */
export interface SaveHistoryData {
  id: number
}

// ==================== API 方法 ====================

export const md2htmlApi = {
  /** 服务端 Markdown 转换（降级/兼容场景） */
  convert(markdown: string, theme?: string) {
    return request<{ html: string; formatted: string }>({
      method: 'POST',
      path: '/api/tools/md2html/convert',
      body: { markdown, theme },
    })
  },

  /**
   * 保存历史记录
   * html 为可选参数（前端本地渲染时不传）
   */
  saveHistory(markdown: string, theme?: string, html?: string) {
    return request<SaveHistoryData>({
      method: 'POST',
      path: '/api/tools/md2html/history',
      body: { markdown, theme, html },
      auth: true,
    })
  },

  /**
   * 获取历史记录列表（轻量，不含 markdown/html）
   * 返回 title + char_count + 时间，用于列表展示
   */
  getHistory(page = 1, pageSize = 20) {
    return request<HistoryPageData>({
      method: 'GET',
      path: `/api/tools/md2html/history?page=${page}&pageSize=${pageSize}`,
      auth: true,
    })
  },

  /**
   * 获取历史记录详情（含完整 markdown）
   * 点击列表项时调用，获取完整内容后加载到编辑器
   */
  getHistoryDetail(id: number) {
    return request<HistoryDetail>({
      method: 'GET',
      path: `/api/tools/md2html/history/${id}`,
      auth: true,
    })
  },

  /** 删除历史记录 */
  deleteHistory(id: number) {
    return request<{ message: string }>({
      method: 'DELETE',
      path: `/api/tools/md2html/history/${id}`,
      auth: true,
    })
  },

  /** 获取主题列表 */
  getThemes() {
    return request<{ id: string; name: string; description: string }[]>({
      method: 'GET',
      path: '/api/tools/md2html/themes',
    })
  },
}
