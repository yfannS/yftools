import { request } from './client'

export interface LoginData {
  token: string
  username: string
  expires_at: string
}

export interface ProfileData {
  id: number
  username: string
  created_at: string
}

export const authApi = {
  register(username: string, password: string) {
    return request<{ message: string }>({
      method: 'POST',
      path: '/api/auth/register',
      body: { username, password },
    })
  },

  login(username: string, password: string) {
    return request<LoginData>({
      method: 'POST',
      path: '/api/auth/login',
      auth: false,
      body: { username, password },
    })
  },

  getProfile() {
    return request<ProfileData>({
      method: 'GET',
      path: '/api/auth/profile',
      auth: true,
    })
  },
}
