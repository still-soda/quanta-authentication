/**
 * 用户相关类型定义
 */

// 用户状态枚举
export type UserStatus = 'ACTIVE' | 'LOCKED' | 'BANNED'

// 用户状态显示名称
export const UserStatusNames: Record<UserStatus, string> = {
   ACTIVE: '正常',
   LOCKED: '已锁定',
   BANNED: '已禁用',
}

// 用户状态标签颜色
export const UserStatusSeverity: Record<UserStatus, 'success' | 'warn' | 'danger' | 'secondary'> = {
   ACTIVE: 'success',
   LOCKED: 'warn',
   BANNED: 'danger',
}

// 用户角色信息（简化版）
export interface UserRole {
   id: string
   name: string
   code: string
   is_system: boolean
}

// 用户信息（来自后端）
export interface User {
   id: string
   student_id: string
   email: string
   name: string
   phone?: string
   display_name?: string
   avatar_id?: string
   status: UserStatus
   email_verified: boolean
   roles: UserRole[]
   role_names: string[]
   last_login_at?: string
   created_at: string
   updated_at: string
}

// 用户表单数据（创建）
export interface CreateUserFormData {
   student_id: string
   name: string
   email: string
   password?: string
   phone?: string
   display_name?: string
   role_ids?: string[]
}

// 用户表单数据（更新）
export interface UpdateUserFormData {
   name?: string
   email?: string
   phone?: string
   display_name?: string
   status?: UserStatus
}

// 用户列表查询参数
export interface ListUsersParams {
   page?: number
   page_size?: number
   search?: string
   status?: UserStatus | ''
   role_id?: string
   sort_by?: string
   sort_desc?: boolean
}

// 用户列表响应
export interface ListUsersResult {
   users: User[]
   total: number
   page: number
   page_size: number
   total_pages: number
}

// 用户状态统计
export interface UserStatusCounts {
   ACTIVE?: number
   LOCKED?: number
   BANNED?: number
}
