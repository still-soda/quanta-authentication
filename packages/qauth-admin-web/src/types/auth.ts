/**
 * 认证相关类型定义
 */

// 认证用户信息（服务端返回格式）
export interface AuthUser {
   id: string
   student_id: string
   name: string
   email: string
   avatar?: string
   display_name?: string
   role: string
   status: 'ACTIVE' | 'LOCKED' | 'BANNED'
   created_at?: string
   updated_at?: string
}

// 登录请求（适配服务端）
export interface LoginRequest {
   student_id: string
   password: string
}

// 登录响应（服务端返回格式）
export interface LoginResponse {
   user: AuthUser
   access_token: string
   refresh_token: string
}

// 刷新 Token 请求
export interface RefreshTokenRequest {
   refresh_token: string
}

// 刷新 Token 响应
export interface RefreshTokenResponse {
   access_token: string
   refresh_token: string
}

// 注册请求（适配服务端）
export interface RegisterRequest {
   student_id: string
   password: string
   email: string
   name: string
}

// 忘记密码申请请求
export interface ForgotPasswordRequest {
   email: string
   student_id: string
   reason: string
}

// 通用申请响应
export interface ApplicationResponse {
   success: boolean
   applicationId?: string
   message: string
}
