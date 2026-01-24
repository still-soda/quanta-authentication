/**
 * 认证相关 API
 */

import { mockResponse } from './index'

// 登录请求参数
export interface LoginRequest {
   username: string
   password: string
   rememberMe?: boolean
}

// 登录响应
export interface LoginResponse {
   token: string
   user: {
      id: string
      username: string
      email: string
      fullName: string
      avatar?: string
      roles: string[]
   }
}

// 忘记密码请求参数
export interface ForgotPasswordRequest {
   email: string
   reason: string
}

// 注册申请请求参数
export interface RegisterRequest {
   username: string
   email: string
   fullName: string
   department: string
   phone?: string
   purpose: string
}

// 申请响应
export interface ApplicationResponse {
   success: boolean
   applicationId: string
   message: string
}

/**
 * 用户登录
 */
export async function login(data: LoginRequest): Promise<LoginResponse> {
   // Mock: 模拟登录逻辑
   if (data.username === 'admin' && data.password === 'admin123') {
      return mockResponse<LoginResponse>(
         {
            token: 'mock-jwt-token-' + Date.now(),
            user: {
               id: '1',
               username: 'admin',
               email: 'admin@quanta.io',
               fullName: '系统管理员',
               avatar: undefined,
               roles: ['admin'],
            },
         },
         800
      )
   }

   // 模拟登录失败
   await mockResponse(null, 500)
   throw new Error('用户名或密码错误')
}

/**
 * 提交忘记密码申请
 */
export async function submitForgotPassword(
   data: ForgotPasswordRequest
): Promise<ApplicationResponse> {
   return mockResponse<ApplicationResponse>(
      {
         success: true,
         applicationId: 'FP-' + Date.now(),
         message: '密码重置申请已提交，管理员将在24小时内审核',
      },
      1000
   )
}

/**
 * 提交注册申请
 */
export async function submitRegister(data: RegisterRequest): Promise<ApplicationResponse> {
   return mockResponse<ApplicationResponse>(
      {
         success: true,
         applicationId: 'REG-' + Date.now(),
         message: '注册申请已提交，管理员将在1-3个工作日内审核',
      },
      1500
   )
}

/**
 * 用户登出
 */
export async function logout(): Promise<void> {
   return mockResponse(undefined, 200)
}

/**
 * 检查 token 是否有效
 */
export async function validateToken(token: string): Promise<boolean> {
   return mockResponse(token.startsWith('mock-jwt-token-'), 100)
}
