/**
 * API 模块入口
 * 统一的 HTTP 客户端，支持 Token 自动刷新
 */
import { Env } from '@/config/env'
import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios'

// Token 存储键
const ACCESS_TOKEN_KEY = 'access_token'
const REFRESH_TOKEN_KEY = 'refresh_token'

// 刷新 token 状态
let isRefreshing = false
let refreshSubscribers: Array<(token: string) => void> = []

// 添加等待刷新完成的请求
function subscribeTokenRefresh(callback: (token: string) => void) {
   refreshSubscribers.push(callback)
}

// 通知所有等待的请求
function onTokenRefreshed(token: string) {
   refreshSubscribers.forEach(callback => callback(token))
   refreshSubscribers = []
}

// 通知刷新失败
function onRefreshFailed() {
   refreshSubscribers = []
}

// 创建 HTTP 客户端
export const httpClient = axios.create({
   baseURL: Env.API_BASE_URL || 'http://localhost:8080',
   timeout: 30000,
   headers: {
      'Content-Type': 'application/json',
   },
})

// 请求拦截器 - 添加 Token
httpClient.interceptors.request.use(
   config => {
      const token = localStorage.getItem(ACCESS_TOKEN_KEY)
      if (token) {
         config.headers.Authorization = `Bearer ${token}`
      }
      return config
   },
   error => Promise.reject(error)
)

// 响应拦截器 - 处理 401 错误和自动刷新 Token
httpClient.interceptors.response.use(
   response => response,
   async (error: AxiosError) => {
      const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean }

      // 如果是 401 错误且不是刷新 token 的请求
      if (
         error.response?.status === 401 &&
         originalRequest &&
         !originalRequest._retry &&
         !originalRequest.url?.includes('/refresh-token')
      ) {
         // 标记为已重试
         originalRequest._retry = true

         // 如果正在刷新 token，将请求加入等待队列
         if (isRefreshing) {
            return new Promise(resolve => {
               subscribeTokenRefresh((token: string) => {
                  originalRequest.headers.Authorization = `Bearer ${token}`
                  resolve(httpClient(originalRequest))
               })
            })
         }

         // 开始刷新 token
         isRefreshing = true

         const refreshTokenValue = localStorage.getItem(REFRESH_TOKEN_KEY)

         if (!refreshTokenValue) {
            // 没有 refresh token，跳转登录
            isRefreshing = false
            onRefreshFailed()
            clearAuthAndRedirect()
            return Promise.reject(error)
         }

         try {
            // 调用刷新 token 接口
            const response = await axios.post(
               `${Env.API_BASE_URL || 'http://localhost:8080'}/_/v1/auth/refresh-token`,
               { refresh_token: refreshTokenValue },
               { headers: { 'Content-Type': 'application/json' } }
            )

            const data = response.data?.data
            if (data?.access_token && data?.refresh_token) {
               // 更新本地存储
               localStorage.setItem(ACCESS_TOKEN_KEY, data.access_token)
               localStorage.setItem(REFRESH_TOKEN_KEY, data.refresh_token)

               // 通知所有等待的请求
               onTokenRefreshed(data.access_token)

               // 重试原始请求
               originalRequest.headers.Authorization = `Bearer ${data.access_token}`
               return httpClient(originalRequest)
            } else {
               throw new Error('Invalid refresh response')
            }
         } catch (refreshError) {
            // 刷新失败，清除认证信息并跳转登录
            onRefreshFailed()
            clearAuthAndRedirect()
            return Promise.reject(refreshError)
         } finally {
            isRefreshing = false
         }
      }

      return Promise.reject(error)
   }
)

// 清除认证信息并跳转登录页
function clearAuthAndRedirect() {
   localStorage.removeItem(ACCESS_TOKEN_KEY)
   localStorage.removeItem(REFRESH_TOKEN_KEY)
   localStorage.removeItem('auth_user')

   // 如果不在登录相关页面，跳转到登录页
   if (!window.location.pathname.startsWith('/auth')) {
      window.location.href = '/auth/login'
   }
}

// 模拟延迟函数（用于兼容旧代码）
export const delay = (ms: number = 100) => new Promise(resolve => setTimeout(resolve, ms))

// 模拟 API 响应（用于兼容旧代码或开发测试）
export async function mockResponse<T>(data: T, delayMs: number = 100): Promise<T> {
   await delay(delayMs)
   return data
}

// 导出所有 API 模块
export * from './auth'
export * from './dashboard'
export * from './users'
export * from './roles'
export * from './oauth'
export * from './audit'
export * from './notifications'
export * from './organizations'
export * from './profile'
export * from './settings'
export * from './app-group'
