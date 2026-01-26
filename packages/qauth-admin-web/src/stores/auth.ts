import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { httpClient } from '@/apis'
import type { AuthUser } from '@/types'

const ACCESS_TOKEN_KEY = 'access_token'
const REFRESH_TOKEN_KEY = 'refresh_token'
const USER_KEY = 'auth_user'
const ROLE_KEY = 'user_role'

export const useAuthStore = defineStore('auth', () => {
   // State
   const accessToken = ref<string | null>(localStorage.getItem(ACCESS_TOKEN_KEY))
   const refreshToken = ref<string | null>(localStorage.getItem(REFRESH_TOKEN_KEY))
   const userRole = ref<string | null>(localStorage.getItem(ROLE_KEY))
   const user = ref<AuthUser | null>(
      (() => {
         const stored = localStorage.getItem(USER_KEY)
         return stored ? JSON.parse(stored) : null
      })()
   )
   const isRefreshing = ref(false)
   const refreshPromise = ref<Promise<boolean> | null>(null)

   // Getters
   const isAuthenticated = computed(() => !!accessToken.value && !!user.value)
   const userName = computed(() => user.value?.name || user.value?.student_id || '')

   // Actions
   /**
    * 设置认证信息
    */
   function setAuth(
      tokens: { accessToken: string; refreshToken: string },
      userData: AuthUser,
      role: string
   ) {
      console.log(arguments)
      accessToken.value = tokens.accessToken
      refreshToken.value = tokens.refreshToken
      user.value = userData
      userRole.value = role

      localStorage.setItem(ACCESS_TOKEN_KEY, tokens.accessToken)
      localStorage.setItem(REFRESH_TOKEN_KEY, tokens.refreshToken)
      localStorage.setItem(USER_KEY, JSON.stringify(userData))
      localStorage.setItem(ROLE_KEY, role)
   }

   /**
    * 更新 tokens
    */
   function updateTokens(newAccessToken: string, newRefreshToken: string) {
      accessToken.value = newAccessToken
      refreshToken.value = newRefreshToken

      localStorage.setItem(ACCESS_TOKEN_KEY, newAccessToken)
      localStorage.setItem(REFRESH_TOKEN_KEY, newRefreshToken)
   }

   /**
    * 清除认证信息
    */
   function clearAuth() {
      accessToken.value = null
      refreshToken.value = null
      user.value = null
      userRole.value = null

      localStorage.removeItem(ACCESS_TOKEN_KEY)
      localStorage.removeItem(REFRESH_TOKEN_KEY)
      localStorage.removeItem(USER_KEY)
      localStorage.removeItem(ROLE_KEY)
   }

   /**
    * 刷新 token
    */
   async function refreshAccessToken(): Promise<boolean> {
      // 如果已经在刷新中，返回现有的 Promise
      if (isRefreshing.value && refreshPromise.value) {
         return refreshPromise.value
      }

      if (!refreshToken.value) {
         clearAuth()
         return false
      }

      isRefreshing.value = true

      refreshPromise.value = (async () => {
         try {
            const response = await httpClient.post('/_/v1/auth/refresh-token', {
               refresh_token: refreshToken.value,
            })

            const data = response.data?.data
            if (data?.access_token && data?.refresh_token) {
               updateTokens(data.access_token, data.refresh_token)
               return true
            }

            clearAuth()
            return false
         } catch {
            clearAuth()
            return false
         } finally {
            isRefreshing.value = false
            refreshPromise.value = null
         }
      })()

      return refreshPromise.value
   }

   /**
    * 登出
    */
   function logout() {
      clearAuth()
   }

   return {
      // State
      accessToken,
      refreshToken,
      user,
      isRefreshing,

      // Getters
      isAuthenticated,
      userRole,
      userName,

      // Actions
      setAuth,
      updateTokens,
      clearAuth,
      refreshAccessToken,
      logout,
   }
})
