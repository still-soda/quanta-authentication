/**
 * 审计日志相关类型定义
 */

import type { Status } from './common'

// 审计模块枚举
export type AuditModule = 'AUTH' | 'OAUTH' | 'USER' | 'ROLE' | 'PERMISSION' | 'CLIENT' | 'SYSTEM'

// 审计操作枚举
export type AuditAction =
   | 'LOGIN'
   | 'LOGOUT'
   | 'REGISTER'
   | 'PASSWORD_RESET'
   | 'PASSWORD_CHANGE'
   | 'TOKEN_REFRESH'
   | 'OAUTH_AUTHORIZE'
   | 'OAUTH_TOKEN'
   | 'OAUTH_REVOKE'
   | 'CLIENT_CREATE'
   | 'CLIENT_UPDATE'
   | 'CLIENT_DELETE'
   | 'ROLE_CREATE'
   | 'ROLE_UPDATE'
   | 'ROLE_DELETE'
   | 'PERMISSION_CREATE'
   | 'PERMISSION_UPDATE'
   | 'PERMISSION_DELETE'
   | 'PERMISSION_GRANT'
   | 'PERMISSION_REVOKE'
   | 'USER_CREATE'
   | 'USER_UPDATE'
   | 'USER_DELETE'
   | 'KEY_ROTATION'
   | 'SETTINGS_CHANGE'
   | 'APP_GROUP_ROLE_CREATE'
   | 'APP_GROUP_ROLE_UPDATE'
   | 'APP_GROUP_ROLE_DELETE'
   | 'APP_GROUP_ROLE_ASSIGN_PERMISSIONS'
   | 'APP_GROUP_ROLE_REVOKE_PERMISSIONS'
   | 'APP_GROUP_USER_ASSIGN_ROLES'
   | 'APP_GROUP_USER_REVOKE_ROLES'
   | 'APP_GROUP_ADMIN_ADD'
   | 'APP_GROUP_ADMIN_REMOVE'
   | 'APP_GROUP_PERMISSION_CREATE'
   | 'APP_GROUP_PERMISSION_UPDATE'
   | 'APP_GROUP_PERMISSION_DELETE'

// 审计状态枚举
export type AuditStatus = 'SUCCESS' | 'WARNING' | 'ERROR'

// 后端返回的审计日志数据（snake_case）
export interface AuditLogBackend {
   id: string
   operator_id: string
   operator_name: string
   module: string
   action: string
   target_id: string
   target_type?: string
   target_name?: string
   detail: Record<string, unknown>
   ip: string
   user_agent?: string
   location?: string
   status: AuditStatus
   error_message?: string
   duration_ms: number
   client_id?: string
   session_id?: string
   request_id?: string
   created_at: string
   operator?: {
      id: string
      username: string
      avatar?: {
         file?: {
            storage_key: string
         }
      }
   }
}

// 审计日志（前端格式，camelCase）
export interface AuditLog {
   id: string
   operatorId: string
   operatorName: string
   operatorAvatar: string
   module: string
   action: string
   targetId: string
   targetName?: string
   detail: Record<string, unknown>
   ip: string
   userAgent?: string
   location?: string
   time: string
   durationMs: number
   status: Status
   errorMessage?: string
}

// 审计日志查询过滤器
export interface AuditLogFilter {
   operator_id?: string
   module?: AuditModule | null
   action?: AuditAction | null
   target_id?: string
   status?: AuditStatus | null
   client_id?: string
   start_time?: string
   end_time?: string
   page?: number
   page_size?: number
   sort_by?: string
   sort_desc?: boolean
}

// 审计日志响应
export interface AuditLogResponse {
   code: number
   msg: string
   data: {
      items: AuditLogBackend[]
      total: number
      page: number
      page_size: number
   }
}

// 审计活动
export interface AuditActivity {
   id: string
   user: string
   operator_id: string
   avatar: string
   action: AuditAction
   module: AuditModule
   target_id: string
   target_name: string
   client: string
   ip: string
   time: string
   status: AuditStatus
   duration_ms: number
}

// 审计统计
export interface AuditStats {
   module_stats: Record<string, number>
   action_stats: Record<string, number>
   status_stats: Record<string, number>
   login_stats: {
      success: number
      fail: number
   }
   start_time: string
   end_time: string
}

// 热门客户端
export interface TopClient {
   client_id: string
   client_name: string
   count: number
}

// 后端审计统计响应数据类型
export interface AuditStatsResponse {
   module_stats: Record<string, number>
   action_stats: Record<string, number>
   status_stats: Record<string, number>
   login_stats: {
      success: number
      fail: number
   }
   start_time: string
   end_time: string
}

// 后端最近活动响应类型
export interface RecentActivityResponse {
   id: string
   user: string
   operator_id: string
   action: string
   module: string
   target_id: string
   target_name: string
   ip: string
   time: string
   status: string
   duration_ms: number
   avatar: string
   client: string
}

// 后端热门客户端响应类型
export interface TopClientResponse {
   client_id: string
   client_name: string
   count: number
}
