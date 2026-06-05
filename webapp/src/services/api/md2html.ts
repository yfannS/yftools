import { request } from './client'

export const md2htmlApi = {
  /** 服务端 Markdown 转换（降级/兼容场景） */
  convert(markdown: string, theme?: string) {
    return request<{ success: boolean; data: { html: string; formatted: string } }>({
      method: 'POST',
      path: '/api/tools/md2html/convert',
      body: { markdown, theme },
    })
  },

  /** 保存历史记录 */
  saveHistory(markdown: string, html: string, theme: string) {
    return request<{ success: boolean; id: number; message: string }>({
      method: 'POST',
      path: '/api/tools/md2html/history',
      body: { markdown, html, theme },
      auth: true,
    })
  },

  /** 获取历史记录 */
  getHistory(page = 1, pageSize = 20) {
    return request<{ success: boolean; data: any[]; total: number }>({
      method: 'GET',
      path: `/api/tools/md2html/history?page=${page}&pageSize=${pageSize}`,
      auth: true,
    }).then(res => {
      // 兼容不同的响应格式
      const data = (res as any).data || res
      const total = (res as any).total || (Array.isArray(data) ? data.length : 0)
      return { data: Array.isArray(data) ? data : [], total }
    })
  },

  /** 删除历史记录 */
  deleteHistory(id: number) {
    return request<{ success: boolean; message: string }>({
      method: 'DELETE',
      path: `/api/tools/md2html/history/${id}`,
      auth: true,
    })
  },

  /** 获取主题列表 */
  getThemes() {
    return request<{ success: boolean; data: { id: string; name: string; description: string }[] }>({
      method: 'GET',
      path: '/api/tools/md2html/themes',
    })
  },
}
