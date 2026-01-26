/**
 * 应用组权限相关类型定义
 */

// 应用组管理员类型
export type AppGroupAdminType = 'owner' | 'admin' | 'role_manager' | 'permission_manager'

// 应用组管理员信息
export interface AppGroupAdmin {
   id: string
   client_id: string
   user_id: string
   user_name: string
   user_email: string
   admin_type: AppGroupAdminType
   granted_at: string
   granted_by: string
   granter_name?: string
   created_at: string
}

// 应用组权限
export interface AppGroupPermission {
   id: string
   client_id: string
   resource: string
   action: number // 1=Create, 2=Read, 3=Update, 4=Delete
   code: string
   name: string
   description: string
   created_at: string
   updated_at: string
}

// 应用组权限表单数据
export interface AppGroupPermissionFormData {
   resource: string
   action: number
   code: string
   name: string
   description?: string
}

// 应用组角色
export interface AppGroupRole {
   id: string
   client_id: string
   code: string
   name: string
   description: string
   is_default: boolean
   user_count: number
   permission_count: number
   created_at: string
   updated_at: string
}

// 应用组角色表单数据
export interface AppGroupRoleFormData {
   code: string
   name: string
   description?: string
   is_default?: boolean
}

// 应用组用户信息（简化版）
export interface AppGroupUserInfo {
   id: string
   name: string
   email: string
   student_id: string
}

// 添加应用组管理员请求
export interface AddAppGroupAdminRequest {
   user_id: string
   admin_type: AppGroupAdminType
}

// 设置角色权限请求
export interface SetAppGroupRolePermissionsRequest {
   permission_ids: string[]
}

// 分配用户角色请求
export interface AssignAppGroupRoleUserRequest {
   user_id: string
}

// 管理员类型显示名称
export const AppGroupAdminTypeLabels: Record<AppGroupAdminType, string> = {
   owner: '应用所有者',
   admin: '应用组管理员',
   role_manager: '角色管理员',
   permission_manager: '权限管理员',
}

// 管理员类型描述
export const AppGroupAdminTypeDescriptions: Record<AppGroupAdminType, string> = {
   owner: '拥有应用组的所有管理权限，不可被移除',
   admin: '可以管理应用组的权限、角色和管理员',
   role_manager: '可以创建、编辑、删除角色，并为用户分配角色',
   permission_manager: '可以创建、编辑、删除权限',
}

// 权限操作类型显示名称
export const PermissionActionLabels: Record<number, string> = {
   1: '创建',
   2: '读取',
   3: '更新',
   4: '删除',
}

// 权限操作类型选项
export const PermissionActionOptions = [
   { label: '创建 (Create)', value: 1 },
   { label: '读取 (Read)', value: 2 },
   { label: '更新 (Update)', value: 3 },
   { label: '删除 (Delete)', value: 4 },
]
