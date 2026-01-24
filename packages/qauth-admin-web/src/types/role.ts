/**
 * 角色和权限相关类型定义
 */

// 权限操作类型
export type PermissionAction = 1 | 2 | 3 | 4 // 1=Create, 2=Read, 3=Update, 4=Delete

// 权限信息（来自后端）
export interface Permission {
   id: string
   resource: string
   action: PermissionAction
   code: string
   description: string
   created_at: string
   updated_at: string
}

// 权限表单数据
export interface PermissionFormData {
   resource: string
   action: PermissionAction
   code: string
   description: string
}

// 权限组（用于 UI 展示）
export interface PermissionGroup {
   name: string
   resource: string
   icon: string
   permissions: PermissionItem[]
}

// 权限组中的权限项
export interface PermissionItem {
   id: string
   code: string
   name: string
   description: string
   checked: boolean
}

// 角色信息（来自后端）
export interface Role {
   id: string
   name: string
   code: string
   description: string
   is_system: boolean
   user_count: number
   permission_count: number
   created_at: string
   updated_at: string
}

// 角色表单数据
export interface RoleFormData {
   name: string
   code: string
   description: string
}

// 用户角色
export interface UserRole {
   name: string
   code: string
   isSystem: boolean
}

// 资源名称映射
export const ResourceNames: Record<string, string> = {
   oauth_clients: 'OAuth 应用',
   system: '系统管理',
   roles: '角色管理',
   audit: '审计日志',
   users: '用户管理',
   permissions: '权限管理',
}

// 资源图标映射
export const ResourceIcons: Record<string, string> = {
   oauth_clients: 'pi pi-key',
   system: 'pi pi-cog',
   roles: 'pi pi-shield',
   audit: 'pi pi-history',
   users: 'pi pi-users',
   permissions: 'pi pi-lock',
}

// 操作类型名称映射
export const ActionNames: Record<PermissionAction, string> = {
   1: '创建',
   2: '查看',
   3: '更新',
   4: '删除',
}
