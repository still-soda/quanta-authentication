/**
 * App Group API - 应用组权限相关接口
 */
import { httpClient } from './index'
import type {
   AppGroupAdmin,
   AppGroupPermission,
   AppGroupRole,
   AppGroupUserInfo,
   AppGroupPermissionFormData,
   AppGroupRoleFormData,
   AddAppGroupAdminRequest,
   SetAppGroupRolePermissionsRequest,
   AssignAppGroupRoleUserRequest,
   AppGroupAdminType,
} from '@/types'

// ======================== 应用组管理员相关 API ========================

/**
 * 获取应用组管理员列表
 */
export async function getAppGroupAdmins(clientId: string): Promise<AppGroupAdmin[]> {
   const response = await httpClient.get(`/_/v1/clients/${clientId}/admins`)
   return response.data.data || []
}

/**
 * 添加应用组管理员
 */
export async function addAppGroupAdmin(
   clientId: string,
   data: AddAppGroupAdminRequest
): Promise<void> {
   await httpClient.post(`/_/v1/clients/${clientId}/admins`, data)
}

/**
 * 移除应用组管理员
 */
export async function removeAppGroupAdmin(
   clientId: string,
   userId: string,
   adminType: AppGroupAdminType
): Promise<void> {
   await httpClient.delete(`/_/v1/clients/${clientId}/admins/${userId}`, {
      params: { admin_type: adminType },
   })
}

// ======================== 应用组权限相关 API ========================

/**
 * 获取应用组权限列表
 */
export async function getAppGroupPermissions(clientId: string): Promise<AppGroupPermission[]> {
   const response = await httpClient.get(`/_/v1/clients/${clientId}/permissions`)
   return response.data.data || []
}

/**
 * 获取应用组权限（按资源分组）
 */
export async function getAppGroupPermissionsGrouped(
   clientId: string
): Promise<Record<string, AppGroupPermission[]>> {
   const response = await httpClient.get(`/_/v1/clients/${clientId}/permissions/grouped`)
   return response.data.data || {}
}

/**
 * 创建应用组权限
 */
export async function createAppGroupPermission(
   clientId: string,
   data: AppGroupPermissionFormData
): Promise<AppGroupPermission> {
   const response = await httpClient.post(`/_/v1/clients/${clientId}/permissions`, data)
   return response.data.data
}

/**
 * 更新应用组权限
 */
export async function updateAppGroupPermission(
   clientId: string,
   permissionId: string,
   data: Pick<AppGroupPermissionFormData, 'name' | 'description'>
): Promise<AppGroupPermission> {
   const response = await httpClient.put(
      `/_/v1/clients/${clientId}/permissions/${permissionId}`,
      data
   )
   return response.data.data
}

/**
 * 删除应用组权限
 */
export async function deleteAppGroupPermission(
   clientId: string,
   permissionId: string
): Promise<void> {
   await httpClient.delete(`/_/v1/clients/${clientId}/permissions/${permissionId}`)
}

// ======================== 应用组角色相关 API ========================

/**
 * 获取应用组角色列表
 */
export async function getAppGroupRoles(clientId: string): Promise<AppGroupRole[]> {
   const response = await httpClient.get(`/_/v1/clients/${clientId}/roles`)
   return response.data.data || []
}

/**
 * 获取应用组角色详情
 */
export async function getAppGroupRole(clientId: string, roleId: string): Promise<AppGroupRole> {
   const response = await httpClient.get(`/_/v1/clients/${clientId}/roles/${roleId}`)
   return response.data.data
}

/**
 * 创建应用组角色
 */
export async function createAppGroupRole(
   clientId: string,
   data: AppGroupRoleFormData
): Promise<AppGroupRole> {
   const response = await httpClient.post(`/_/v1/clients/${clientId}/roles`, data)
   return response.data.data
}

/**
 * 更新应用组角色
 */
export async function updateAppGroupRole(
   clientId: string,
   roleId: string,
   data: Pick<AppGroupRoleFormData, 'name' | 'description' | 'is_default'>
): Promise<AppGroupRole> {
   const response = await httpClient.put(`/_/v1/clients/${clientId}/roles/${roleId}`, data)
   return response.data.data
}

/**
 * 删除应用组角色
 */
export async function deleteAppGroupRole(clientId: string, roleId: string): Promise<void> {
   await httpClient.delete(`/_/v1/clients/${clientId}/roles/${roleId}`)
}

/**
 * 获取应用组角色的权限
 */
export async function getAppGroupRolePermissions(
   clientId: string,
   roleId: string
): Promise<AppGroupPermission[]> {
   const response = await httpClient.get(`/_/v1/clients/${clientId}/roles/${roleId}/permissions`)
   return response.data.data || []
}

/**
 * 设置应用组角色的权限
 */
export async function setAppGroupRolePermissions(
   clientId: string,
   roleId: string,
   data: SetAppGroupRolePermissionsRequest
): Promise<void> {
   await httpClient.put(`/_/v1/clients/${clientId}/roles/${roleId}/permissions`, data)
}

/**
 * 获取应用组角色的用户列表
 */
export async function getAppGroupRoleUsers(
   clientId: string,
   roleId: string
): Promise<AppGroupUserInfo[]> {
   const response = await httpClient.get(`/_/v1/clients/${clientId}/roles/${roleId}/users`)
   return response.data.data || []
}

/**
 * 为用户分配应用组角色
 */
export async function assignAppGroupRoleToUser(
   clientId: string,
   roleId: string,
   data: AssignAppGroupRoleUserRequest
): Promise<void> {
   await httpClient.post(`/_/v1/clients/${clientId}/roles/${roleId}/users`, data)
}

/**
 * 从用户撤销应用组角色
 */
export async function revokeAppGroupRoleFromUser(
   clientId: string,
   roleId: string,
   userId: string
): Promise<void> {
   await httpClient.delete(`/_/v1/clients/${clientId}/roles/${roleId}/users/${userId}`)
}

// ======================== 应用组用户权限查询 API ========================

/**
 * 获取用户在应用组的角色
 */
export async function getUserAppGroupRoles(
   clientId: string,
   userId: string
): Promise<AppGroupRole[]> {
   const response = await httpClient.get(`/_/v1/clients/${clientId}/users/${userId}/roles`)
   return response.data.data || []
}

/**
 * 获取用户在应用组的权限
 */
export async function getUserAppGroupPermissions(
   clientId: string,
   userId: string
): Promise<AppGroupPermission[]> {
   const response = await httpClient.get(`/_/v1/clients/${clientId}/users/${userId}/permissions`)
   return response.data.data || []
}
