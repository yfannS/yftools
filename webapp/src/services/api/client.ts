const BASE_URL = import.meta.env.VITE_API_BASE || ''

interface RequestOptions {
  method: string
  path: string
  body?: unknown
  auth?: boolean
  /** 不自动解包 data 字段，返回完整响应 */
  raw?: boolean
}

export class ApiError extends Error {
  status: number
  body: any

  constructor(status: number, body: any) {
    super(body?.message || body?.error || `API error ${status}`)
    this.status = status
    this.body = body
  }
}

/**
 * 统一请求函数
 * 
 * 后端响应格式：{ success: boolean, data?: T, message?: string }
 * 默认自动解包 data 字段，直接返回 T 类型数据。
 * 设置 raw: true 时返回完整响应体。
 */
export async function request<T>(opts: RequestOptions): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }

  const token = localStorage.getItem('m2h_webapp_token')
  if (token && opts.auth !== false) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(`${BASE_URL}${opts.path}`, {
    method: opts.method,
    headers,
    body: opts.body ? JSON.stringify(opts.body) : undefined,
  })

  if (res.status === 401) {
    // Token 失效：清除登录态，不中断编辑
    localStorage.removeItem('m2h_webapp_token')
    localStorage.removeItem('m2h_webapp_username')
  }

  if (!res.ok) {
    const err = await res.json().catch(() => ({ message: res.statusText }))
    throw new ApiError(res.status, err)
  }

  if (res.status === 204) return undefined as T

  const json = await res.json()

  // raw 模式：返回完整响应体
  if (opts.raw) return json as T

  // 自动解包：后端 { success, data, message } → 提取 data
  if (json && typeof json === 'object' && 'success' in json && 'data' in json) {
    return json.data as T
  }

  // 兼容：无 data 字段时返回整个响应
  return json as T
}
