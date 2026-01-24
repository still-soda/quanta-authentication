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
   | 'PERMISSION_GRANT'
   | 'PERMISSION_REVOKE'
   | 'USER_CREATE'
   | 'USER_UPDATE'
   | 'USER_DELETE'
   | 'KEY_ROTATION'
   | 'SETTINGS_CHANGE'

// 审计状态枚举
export type AuditStatus = 'SUCCESS' | 'WARNING' | 'ERROR'

// 审计日志
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
   time: string
   durationMs: number
   status: Status
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
}

// 审计日志响应
export interface AuditLogResponse {
   items: AuditLog[]
   total: number
   page: number
   page_size: number
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
