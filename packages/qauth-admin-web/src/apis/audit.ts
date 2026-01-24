/**
 * Audit API - 审计日志相关接口
 */
import { httpClient } from './index'
import type {
   AuditLog,
   AuditModule,
   AuditAction,
   AuditStatus,
   AuditLogFilter,
   AuditLogResponse,
   AuditActivity,
   AuditStats,
   TopClient,
   Status,
} from '@/types'

/**
 * 获取审计日志列表
 */
export async function getAuditLogs(filter?: AuditLogFilter): Promise<AuditLogResponse> {
   const params = new URLSearchParams()
   if (filter) {
      Object.entries(filter).forEach(([key, value]) => {
         if (value !== undefined && value !== null && value !== '') {
            params.append(key, String(value))
         }
      })
   }

   const response = await httpClient.get<AuditLogResponse>(`/_/v1/audit/logs?${params.toString()}`)
   return response.data
}

/**
 * 获取单条审计日志详情
 */
export async function getAuditLogDetail(id: string): Promise<AuditLog> {
   const response = await httpClient.get<AuditLog>(`/_/v1/audit/logs/${id}`)
   return response.data
}

/**
 * 获取审计活动列表
 */
export async function getAuditActivities(): Promise<AuditActivity[]> {
   const response = await httpClient.get<AuditActivity[]>('/_/v1/audit/activities')
   return response.data
}

/**
 * 获取审计统计数据
 */
export async function getAuditStats(): Promise<AuditStats> {
   const response = await httpClient.get<AuditStats>('/_/v1/audit/stats')
   return response.data
}

/**
 * 获取热门客户端
 */
export async function getTopClients(): Promise<TopClient[]> {
   const response = await httpClient.get<TopClient[]>('/_/v1/audit/top-clients')
   return response.data
}

/**
 * 导出审计日志
 */
export async function exportAuditLogs(filter?: AuditLogFilter): Promise<Blob> {
   const params = new URLSearchParams()
   if (filter) {
      Object.entries(filter).forEach(([key, value]) => {
         if (value !== undefined && value !== null && value !== '') {
            params.append(key, String(value))
         }
      })
   }

   const response = await httpClient.get(`/_/v1/audit/export?${params.toString()}`, {
      responseType: 'blob',
   })
   return response.data
}

// 辅助函数：将后端状态映射为前端状态
export function mapAuditStatus(status: AuditStatus): Status {
   switch (status) {
      case 'SUCCESS':
         return 'success'
      case 'WARNING':
         return 'warning'
      case 'ERROR':
         return 'error'
      default:
         return 'success'
   }
}

// 辅助函数：获取模块显示名称
export function getModuleDisplayName(module: AuditModule): string {
   const moduleNames: Record<AuditModule, string> = {
      AUTH: '认证登录',
      OAUTH: 'OAuth授权',
      USER: '用户管理',
      ROLE: '角色管理',
      PERMISSION: '权限管理',
      CLIENT: '客户端管理',
      SYSTEM: '系统管理',
   }
   return moduleNames[module] || module
}

// 辅助函数：获取操作显示名称
export function getActionDisplayName(action: AuditAction): string {
   const actionNames: Record<AuditAction, string> = {
      LOGIN: '登录',
      LOGOUT: '登出',
      REGISTER: '注册',
      PASSWORD_RESET: '密码重置',
      PASSWORD_CHANGE: '密码修改',
      TOKEN_REFRESH: '令牌刷新',
      OAUTH_AUTHORIZE: 'OAuth授权',
      OAUTH_TOKEN: 'OAuth令牌',
      OAUTH_REVOKE: 'OAuth撤销',
      CLIENT_CREATE: '创建客户端',
      CLIENT_UPDATE: '更新客户端',
      CLIENT_DELETE: '删除客户端',
      ROLE_CREATE: '创建角色',
      ROLE_UPDATE: '更新角色',
      ROLE_DELETE: '删除角色',
      PERMISSION_GRANT: '授予权限',
      PERMISSION_REVOKE: '撤销权限',
      USER_CREATE: '创建用户',
      USER_UPDATE: '更新用户',
      USER_DELETE: '删除用户',
      KEY_ROTATION: '密钥轮换',
      SETTINGS_CHANGE: '设置变更',
   }
   return actionNames[action] || action
}
