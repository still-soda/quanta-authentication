/**
 * Users API - 用户管理相关接口
 */
import { httpClient } from './index'
import type {
   User,
   CreateUserFormData,
   UpdateUserFormData,
   ListUsersParams,
   ListUsersResult,
   UserStatusCounts,
} from '@/types'
import type { Role } from '@/types/role'

/**
 * 获取用户列表（带分页）
 */
export async function getUsers(params: ListUsersParams = {}): Promise<ListUsersResult> {
   const response = await httpClient.get('/_/v1/users', { params })
   return response.data.data
}

/**
 * 获取单个用户
 */
export async function getUser(id: string): Promise<User> {
   const response = await httpClient.get(`/_/v1/users/${id}`)
   return response.data.data
}

/**
 * 创建用户
 */
export async function createUser(data: CreateUserFormData): Promise<User> {
   const response = await httpClient.post('/_/v1/users', data)
   return response.data.data
}

/**
 * 更新用户
 */
export async function updateUser(id: string, data: UpdateUserFormData): Promise<User> {
   const response = await httpClient.put(`/_/v1/users/${id}`, data)
   return response.data.data
}

/**
 * 删除用户
 */
export async function deleteUser(id: string): Promise<void> {
   await httpClient.delete(`/_/v1/users/${id}`)
}

/**
 * 获取用户角色
 */
export async function getUserRoles(userId: string): Promise<Role[]> {
   const response = await httpClient.get(`/_/v1/users/${userId}/roles`)
   return response.data.data || []
}

/**
 * 设置用户角色（替换现有角色）
 */
export async function setUserRoles(userId: string, roleIds: string[]): Promise<Role[]> {
   const response = await httpClient.put(`/_/v1/users/${userId}/roles`, {
      role_ids: roleIds,
   })
   return response.data.data || []
}

/**
 * 为用户分配角色（追加）
 */
export async function assignRolesToUser(userId: string, roleIds: string[]): Promise<Role[]> {
   const response = await httpClient.post(`/_/v1/users/${userId}/roles/assign`, {
      role_ids: roleIds,
   })
   return response.data.data || []
}

/**
 * 从用户撤销角色
 */
export async function revokeRolesFromUser(userId: string, roleIds: string[]): Promise<Role[]> {
   const response = await httpClient.post(`/_/v1/users/${userId}/roles/revoke`, {
      role_ids: roleIds,
   })
   return response.data.data || []
}

/**
 * 重置用户密码
 */
export async function resetUserPassword(
   userId: string,
   newPassword?: string
): Promise<{ success: boolean; new_password?: string }> {
   const response = await httpClient.post(`/_/v1/users/${userId}/reset-password`, {
      new_password: newPassword,
   })
   return response.data.data
}

/**
 * 获取用户状态统计
 */
export async function getUserStatusCounts(): Promise<UserStatusCounts> {
   const response = await httpClient.get('/_/v1/users/stats')
   return response.data.data || {}
}

/**
 * 禁用用户（更改状态为 LOCKED）
 */
export async function disableUser(id: string): Promise<User> {
   return updateUser(id, { status: 'LOCKED' })
}

/**
 * 启用用户（更改状态为 ACTIVE）
 */
export async function enableUser(id: string): Promise<User> {
   return updateUser(id, { status: 'ACTIVE' })
}
