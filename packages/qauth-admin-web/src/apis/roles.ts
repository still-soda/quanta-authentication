/**
 * Roles API - 角色相关接口
 */
import { httpClient } from './index'
import type { Role, RoleFormData, Permission, PermissionGroup } from '@/types'

/**
 * 获取角色列表
 * @param params 查询参数（可选）
 */
export async function getRoles(params?: {
   page?: number
   page_size?: number
   search?: string
   all?: boolean
}): Promise<Role[]> {
   // 如果明确请求所有数据或没有传入分页参数，则请求所有
   if (params?.all === true || (!params?.page && !params?.page_size)) {
      const response = await httpClient.get('/_/v1/roles', { params: { all: true } })
      return response.data.data || []
   }

   // 否则使用分页
   const response = await httpClient.get('/_/v1/roles', { params })
   return response.data.data?.items || response.data.data || []
}

/**
 * 获取角色列表（带分页）
 */
export async function listRoles(params: { page: number; page_size: number; search?: string }) {
   const response = await httpClient.get('/_/v1/roles', { params })
   return response.data.data
}

/**
 * 获取单个角色
 */
export async function getRole(id: string): Promise<Role> {
   const response = await httpClient.get(`/_/v1/roles/${id}`)
   return response.data.data
}

/**
 * 创建角色
 */
export async function createRole(data: RoleFormData & { permissions?: string[] }): Promise<Role> {
   const response = await httpClient.post('/_/v1/roles', {
      name: data.name,
      code: data.code,
      description: data.description,
      permissions: data.permissions || [],
   })
   return response.data.data
}

/**
 * 更新角色
 */
export async function updateRole(id: string, data: RoleFormData): Promise<Role> {
   const response = await httpClient.put(`/_/v1/roles/${id}`, data)
   return response.data.data
}

/**
 * 删除角色
 */
export async function deleteRole(id: string): Promise<void> {
   await httpClient.delete(`/_/v1/roles/${id}`)
}

/**
 * 获取角色的权限
 */
export async function getRolePermissions(roleId: string): Promise<Permission[]> {
   const response = await httpClient.get(`/_/v1/roles/${roleId}/permissions`)
   return response.data.data || []
}

/**
 * 设置角色的权限
 */
export async function setRolePermissions(roleId: string, permissionCodes: string[]): Promise<void> {
   await httpClient.put(`/_/v1/roles/${roleId}/permissions`, {
      permission_codes: permissionCodes,
   })
}

/**
 * 获取所有权限
 * @param params 查询参数（可选）
 */
export async function getPermissions(params?: {
   page?: number
   page_size?: number
   search?: string
   resource?: string
   all?: boolean
}): Promise<Permission[]> {
   // 如果明确请求所有数据或没有传入分页参数，则请求所有
   if (params?.all === true || (!params?.page && !params?.page_size)) {
      const response = await httpClient.get('/_/v1/permissions', { params: { all: true } })
      return response.data.data || []
   }

   // 否则使用分页
   const response = await httpClient.get('/_/v1/permissions', { params })
   return response.data.data?.items || response.data.data || []
}

/**
 * 获取权限列表（带分页）
 */
export async function listPermissions(params: {
   page: number
   page_size: number
   search?: string
   resource?: string
   sort_field?: string
   sort_order?: 'asc' | 'desc'
}) {
   const response = await httpClient.get('/_/v1/permissions', { params })
   return response.data.data
}

/**
 * 获取按资源分组的权限
 */
export async function getPermissionsGrouped(): Promise<Record<string, Permission[]>> {
   const response = await httpClient.get('/_/v1/permissions/grouped')
   return response.data.data || {}
}

/**
 * 创建权限
 */
export async function createPermission(data: {
   resource: string
   action: number
   code: string
   description: string
}): Promise<Permission> {
   const response = await httpClient.post('/_/v1/permissions', data)
   return response.data.data
}

/**
 * 更新权限
 */
export async function updatePermission(
   id: string,
   data: {
      resource: string
      action: number
      code: string
      description: string
   }
): Promise<Permission> {
   const response = await httpClient.put(`/_/v1/permissions/${id}`, data)
   return response.data.data
}

/**
 * 删除权限
 */
export async function deletePermission(id: string): Promise<void> {
   await httpClient.delete(`/_/v1/permissions/${id}`)
}

// 资源名称映射
const resourceNamesMap: Record<string, string> = {
   oauth_clients: 'OAuth 应用',
   system: '系统管理',
   roles: '角色管理',
   audit: '审计日志',
   users: '用户管理',
   permissions: '权限管理',
}

// 资源图标映射
const resourceIconsMap: Record<string, string> = {
   oauth_clients: 'pi pi-key',
   system: 'pi pi-cog',
   roles: 'pi pi-shield',
   audit: 'pi pi-history',
   users: 'pi pi-users',
   permissions: 'pi pi-lock',
}

/**
 * 将权限列表转换为权限组格式
 */
export function transformPermissionsToGroups(
   permissions: Permission[],
   checkedCodes: string[] = []
): PermissionGroup[] {
   const grouped: Record<string, Permission[]> = {}

   // 按资源分组
   for (const perm of permissions) {
      if (!grouped[perm.resource]) {
         grouped[perm.resource] = []
      }
      grouped[perm.resource].push(perm)
   }

   // 转换为 PermissionGroup 格式
   return Object.entries(grouped).map(([resource, perms]) => ({
      name: resourceNamesMap[resource] || resource,
      resource,
      icon: resourceIconsMap[resource] || 'pi pi-key',
      permissions: perms.map(p => ({
         id: p.id,
         code: p.code,
         name: p.description || p.code,
         description: p.description,
         checked: checkedCodes.includes(p.code),
      })),
   }))
}

/**
 * 从权限组中提取选中的权限代码
 */
export function extractCheckedPermissionCodes(groups: PermissionGroup[]): string[] {
   return groups.flatMap(group => group.permissions.filter(p => p.checked).map(p => p.code))
}
