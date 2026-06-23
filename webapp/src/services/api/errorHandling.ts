import { ApiError } from './client'

export const AUTH_ERROR_CODES = new Set([
  'AUTH_TOKEN_REQUIRED',
  'AUTH_TOKEN_MALFORMED',
  'AUTH_TOKEN_INVALID',
  'AUTH_SESSION_EXPIRED',
])

export function isApiError(error: unknown): error is ApiError {
  return error instanceof ApiError
}

export function getApiErrorCode(error: unknown): string | undefined {
  if (!isApiError(error)) return undefined
  return typeof error.code === 'string' ? error.code : undefined
}

export function getApiErrorData<T = Record<string, unknown>>(error: unknown): T | undefined {
  if (!isApiError(error)) return undefined
  return error.data as T | undefined
}

export function isAuthErrorCode(code?: string) {
  return !!code && AUTH_ERROR_CODES.has(code)
}

export function getApiErrorMessage(error: unknown, fallback = '操作失败，请稍后再试') {
  if (!isApiError(error)) return fallback

  return mapApiErrorMessage(error.code, error.message || fallback, error.data)
}

export function mapApiErrorMessage(code?: string, fallback?: string, data?: Record<string, unknown>) {
  switch (code) {
    case 'INVALID_PARAMS':
      return '请检查输入内容后重试'
    case 'USERNAME_EXISTS':
      return '用户名已存在，请换一个用户名'
    case 'INVALID_CREDENTIALS':
      return buildInvalidCredentialsMessage(data)
    case 'RATE_LIMITED':
      return buildRateLimitedMessage(data)
    case 'AUTH_TOKEN_REQUIRED':
      return '请先登录'
    case 'AUTH_TOKEN_MALFORMED':
    case 'AUTH_TOKEN_INVALID':
      return '登录状态无效，请重新登录'
    case 'AUTH_SESSION_EXPIRED':
      return '登录已失效，请重新登录'
    case 'USER_NOT_FOUND':
      return '用户不存在'
    case 'INVALID_RECORD_ID':
      return '记录 ID 无效'
    case 'HISTORY_NOT_FOUND':
      return '历史记录不存在或已删除'
    case 'HISTORY_LIST_FAILED':
      return '获取历史记录失败，请稍后再试'
    case 'HISTORY_SAVE_FAILED':
      return '保存历史记录失败'
    case 'HISTORY_DELETE_FAILED':
      return '删除历史记录失败'
    case 'HISTORY_RENAME_FAILED':
      return '修改标题失败'
    case 'JSON_FORMAT_FAILED':
      return 'JSON 格式化失败'
    case 'JSON_VALIDATE_FAILED':
      return 'JSON 校验失败'
    case 'JSON_HISTORY_FAILED':
      return 'JSON 历史记录操作失败'
    case 'MARKDOWN_CONVERT_FAILED':
      return 'Markdown 转换失败'
    case 'INTERNAL_ERROR':
      return '服务器开小差了，请稍后再试'
    default:
      return fallback || '操作失败，请稍后再试'
  }
}

function buildInvalidCredentialsMessage(data?: Record<string, unknown>) {
  const remaining = toNonNegativeInt(data?.remaining_attempts)
  if (remaining === undefined) {
    return '用户名或密码错误'
  }

  if (remaining <= 0) {
    const retryAfter = toNonNegativeInt(data?.retry_after_seconds)
    if (retryAfter && retryAfter > 0) {
      return `用户名或密码错误，请在 ${formatRetryAfter(retryAfter)}后重试`
    }
    return '用户名或密码错误，当前已无法继续尝试，请稍后再试'
  }

  return `用户名或密码错误，还可尝试 ${remaining} 次`
}

function buildRateLimitedMessage(data?: Record<string, unknown>) {
  const retryAfter = toNonNegativeInt(data?.retry_after_seconds)
  if (retryAfter && retryAfter > 0) {
    return `尝试次数过多，请在 ${formatRetryAfter(retryAfter)}后重试`
  }
  return '操作过于频繁，请稍后再试'
}

function toNonNegativeInt(value: unknown) {
  if (typeof value !== 'number' || !Number.isFinite(value)) return undefined
  return Math.max(0, Math.floor(value))
}

function formatRetryAfter(totalSeconds: number) {
  if (totalSeconds < 60) {
    return `${totalSeconds} 秒`
  }

  const minutes = Math.floor(totalSeconds / 60)
  const seconds = totalSeconds % 60
  if (seconds === 0) {
    return `${minutes} 分钟`
  }

  return `${minutes} 分 ${seconds} 秒`
}
