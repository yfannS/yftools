const BASE_URL = import.meta.env.VITE_API_BASE || ''
const AUTH_EXPIRED_EVENT = 'app:auth-expired'

/** 默认请求超时（毫秒） */
const DEFAULT_TIMEOUT_MS = 15_000

interface RequestOptions {
  method: string
  path: string
  body?: unknown
  auth?: boolean
  /** 不自动解包 data 字段，返回完整响应 */
  raw?: boolean
  /** 自定义超时毫秒数，0 表示不超时 */
  timeout?: number
  /** 外部取消信号，路由切换时可用于取消在途请求 */
  signal?: AbortSignal
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
 *
 * 超时：默认 15s，可通过 opts.timeout 覆盖（0 = 不超时）。
 * 取消：传入 opts.signal 可在路由切换时 abort 请求。
 */
export async function request<T>(opts: RequestOptions): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }

  const token = localStorage.getItem('m2h_webapp_token')
  if (token && opts.auth !== false) {
    headers['Authorization'] = `Bearer ${token}`
  }

  // 组合超时 signal 与外部 signal
  const timeoutMs = opts.timeout === undefined ? DEFAULT_TIMEOUT_MS : opts.timeout
  let combinedSignal: AbortSignal | undefined

  if (timeoutMs > 0 && opts.signal) {
    // 两路 signal 取最先触发的
    const timeoutSignal = AbortSignal.timeout(timeoutMs)
    combinedSignal = AbortSignal.any
      ? AbortSignal.any([timeoutSignal, opts.signal])
      : opts.signal // 降级：不支持 any 时优先外部 signal
  } else if (timeoutMs > 0) {
    combinedSignal = AbortSignal.timeout(timeoutMs)
  } else {
    combinedSignal = opts.signal
  }

  const res = await fetch(`${BASE_URL}${opts.path}`, {
    method: opts.method,
    headers,
    body: opts.body ? JSON.stringify(opts.body) : undefined,
    signal: combinedSignal,
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
