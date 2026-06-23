const BASE_URL = import.meta.env.VITE_API_BASE || ''
const AUTH_EXPIRED_EVENT = 'app:auth-expired'

interface RequestOptions {
  method: string
  path: string
  body?: unknown
  auth?: boolean
  /** 不自动解包 data 字段，返回完整响应 */
  raw?: boolean
}

export interface ApiErrorBody {
  success?: boolean
  message?: string
  error?: string
  code?: string
  request_id?: string
  data?: Record<string, unknown>
  [key: string]: unknown
}

export class ApiError extends Error {
  status: number
  body: ApiErrorBody
  code?: string
  requestId?: string
  data?: Record<string, unknown>

  constructor(status: number, body: ApiErrorBody) {
    super(body?.message || body?.error || `API error ${status}`)
    this.status = status
    this.body = body
    this.code = typeof body?.code === 'string' ? body.code : undefined
    this.requestId = typeof body?.request_id === 'string' ? body.request_id : undefined
    this.data = body?.data
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

  if (!res.ok) {
    const err = await res.json().catch(() => ({ message: res.statusText }))
    const apiError = new ApiError(res.status, err)

    if (res.status === 401 && isAuthFailure(apiError.code)) {
      clearLocalAuth()
      window.dispatchEvent(new CustomEvent(AUTH_EXPIRED_EVENT, {
        detail: {
          status: apiError.status,
          code: apiError.code,
          message: apiError.message,
          path: opts.path,
        },
      }))
    }

    throw apiError
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

function clearLocalAuth() {
  localStorage.removeItem('m2h_webapp_token')
  localStorage.removeItem('m2h_webapp_username')
  localStorage.removeItem('m2h_webapp_expires_at')
}

function isAuthFailure(code?: string) {
  return code === 'AUTH_TOKEN_REQUIRED'
    || code === 'AUTH_TOKEN_MALFORMED'
    || code === 'AUTH_TOKEN_INVALID'
    || code === 'AUTH_SESSION_EXPIRED'
}

export { AUTH_EXPIRED_EVENT }
