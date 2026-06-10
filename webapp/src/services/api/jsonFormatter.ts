import { request } from './client'

// ==================== JSON 格式化相关类型 ====================

export interface JsonFormatData {
  output: string
  size: number
  minified: boolean
}

export interface JsonValidateData {
  valid: boolean
  error?: string
  keys?: number
  depth?: number
}

export interface JsonHistoryListItem {
  id: number
  title: string
  char_count: number
  created_at: string
  updated_at: string
}

export interface JsonHistoryDetail {
  id: number
  title: string
  input: string
  output: string
  char_count: number
  created_at: string
  updated_at: string
}

export interface JsonHistoryPageData {
  data: JsonHistoryListItem[]
  total: number
  page: number
  size: number
}

// ==================== API 方法 ====================

export const jsonFormatterApi = {
  /** 格式化/压缩 JSON */
  format(input: string, indent = 2, minify = false) {
    return request<JsonFormatData>({
      method: 'POST',
      path: '/api/tools/json-formatter/format',
      body: { input, indent, minify },
    })
  },

  /** 校验 JSON */
  validate(input: string) {
    return request<JsonValidateData>({
      method: 'POST',
      path: '/api/tools/json-formatter/validate',
      body: { input },
    })
  },

  /** 获取历史记录列表 */
  getHistory(page = 1, pageSize = 20) {
    return request<JsonHistoryPageData>({
      method: 'GET',
      path: `/api/tools/json-formatter/history?page=${page}&pageSize=${pageSize}`,
      auth: true,
    })
  },

  /** 获取历史记录详情 */
  getHistoryDetail(id: number) {
    return request<JsonHistoryDetail>({
      method: 'GET',
      path: `/api/tools/json-formatter/history/${id}`,
      auth: true,
    })
  },

  /** 保存历史记录 */
  saveHistory(input: string, output: string) {
    return request<{ id: number }>({
      method: 'POST',
      path: '/api/tools/json-formatter/history',
      body: { input, output },
      auth: true,
    })
  },

  /** 删除历史记录 */
  deleteHistory(id: number) {
    return request<{ message: string }>({
      method: 'DELETE',
      path: `/api/tools/json-formatter/history/${id}`,
      auth: true,
    })
  },

  /** 修改历史记录标题 */
  renameHistory(id: number, title: string) {
    return request<{ message: string }>({
      method: 'PATCH',
      path: `/api/tools/json-formatter/history/${id}`,
      body: { title },
      auth: true,
    })
  },
}
