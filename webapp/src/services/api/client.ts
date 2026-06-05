const BASE_URL = import.meta.env.VITE_API_BASE || ''

interface RequestOptions {
  method: string
  path: string
  body?: unknown
  auth?: boolean
}

export class ApiError extends Error {
  status: number
  body: any

  constructor(status: number, body: any) {
    super(body?.error || body?.message || `API error ${status}`)
    this.status = status
    this.body = body
  }
}

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
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new ApiError(res.status, err)
  }

  if (res.status === 204) return undefined as T
  return res.json()
}
