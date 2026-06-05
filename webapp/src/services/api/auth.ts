import { request } from './client'

export const authApi = {
  register(username: string, password: string) {
    return request<{ success: boolean; message: string }>({
      method: 'POST',
      path: '/api/auth/register',
      body: { username, password },
    })
  },

  login(username: string, password: string) {
    return request<{ success: boolean; data: { token: string; username: string } }>({
      method: 'POST',
      path: '/api/auth/login',
      body: { username, password },
    }).then(res => {
      // 兼容不同的响应格式
      if ((res as any).token) return { token: (res as any).token, username: (res as any).username }
      if ((res as any).data) return (res as any).data
      return res as any
    })
  },

  getProfile() {
    return request<{ success: boolean; data: { id: number; username: string; created_at: string } }>({
      method: 'GET',
      path: '/api/auth/profile',
      auth: true,
    })
  },
}
