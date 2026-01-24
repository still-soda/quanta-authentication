/**
 * 角色和权限相关类型定义
 */

// 角色信息
export interface Role {
   id: number
   name: string
   code: string
   description: string
   userCount: number
   permissions: number
   status: string
   isSystem: boolean
   createdAt: string
}

// 角色表单数据
export interface RoleFormData {
   name: string
   code: string
   description: string
}

// 权限
export interface Permission {
   id: string
   name: string
   checked: boolean
}

// 权限组
export interface PermissionGroup {
   name: string
   icon: string
   permissions: Permission[]
}

// 用户角色
export interface UserRole {
   name: string
   code: string
   isSystem: boolean
}
