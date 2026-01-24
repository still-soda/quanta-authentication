/**
 * 认证相关 API
 */

import { httpClient } from './index'
import type {
   LoginRequest,
   LoginResponse,
   ForgotPasswordRequest,
   RegisterRequest,
   ApplicationResponse,
   RefreshTokenRequest,
   RefreshTokenResponse,
} from '@/types'

// API 基础路径
const AUTH_BASE = '/_/v1/auth'

/**
 * 用户登录
 */
export async function login(data: LoginRequest): Promise<LoginResponse> {
   const response = await httpClient.post(`${AUTH_BASE}/login`, data)
   return response.data?.data
}

/**
 * 刷新 Token
 */
export async function refreshToken(data: RefreshTokenRequest): Promise<RefreshTokenResponse> {
   const response = await httpClient.post(`${AUTH_BASE}/refresh-token`, data)
   return response.data?.data
}

/**
 * 用户注册（向服务端提交注册）
 */
export async function register(data: RegisterRequest): Promise<ApplicationResponse> {
   try {
      const response = await httpClient.post(`${AUTH_BASE}/register`, data)
      return {
         success: true,
         message: '注册成功',
         applicationId: response.data?.data?.id,
      }
   } catch (error: unknown) {
      // 处理错误响应
      const err = error as { response?: { data?: { message?: string } } }
      throw new Error(err.response?.data?.message || '注册失败，请稍后重试')
   }
}

/**
 * 提交忘记密码申请
 * 注意：目前服务端可能没有此接口，返回 mock 响应
 */
export async function submitForgotPassword(
   data: ForgotPasswordRequest
): Promise<ApplicationResponse> {
   // TODO: 等待服务端实现找回密码申请接口
   // 模拟提交成功
   await new Promise(resolve => setTimeout(resolve, 1000))
   return {
      success: true,
      applicationId: 'FP-' + Date.now(),
      message: '密码重置申请已提交，管理员将在24小时内审核',
   }
}

/**
 * 提交注册申请
 * 直接调用服务端注册接口
 */
export async function submitRegister(data: RegisterRequest): Promise<ApplicationResponse> {
   return register(data)
}

/**
 * 用户登出
 */
export async function logout(): Promise<void> {
   // 清除本地存储的 token
   localStorage.removeItem('access_token')
   localStorage.removeItem('refresh_token')
   localStorage.removeItem('auth_user')
}

/**
 * 检查 token 是否有效
 */
export async function validateToken(): Promise<boolean> {
   try {
      // 尝试调用需要认证的接口来验证 token
      await httpClient.get('/_/v1/dashboard/stats')
      return true
   } catch {
      return false
   }
}
