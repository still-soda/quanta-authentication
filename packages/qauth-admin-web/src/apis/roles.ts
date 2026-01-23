/**
 * Roles API - 角色权限相关接口
 */
import { mockResponse } from './index';
import type { Role, RoleFormData, PermissionGroup } from '@/types';

/**
 * 获取角色列表
 */
export async function getRoles(): Promise<Role[]> {
  return mockResponse([
    {
      id: 1,
      name: '超级管理员',
      code: 'super_admin',
      description: '拥有系统所有权限',
      userCount: 3,
      permissions: 42,
      status: 'active',
      isSystem: true,
      createdAt: '2024-01-01',
    },
    {
      id: 2,
      name: '管理员',
      code: 'admin',
      description: '管理用户和基本设置',
      userCount: 12,
      permissions: 28,
      status: 'active',
      isSystem: true,
      createdAt: '2024-01-01',
    },
    {
      id: 3,
      name: '开发者',
      code: 'developer',
      description: '管理 OAuth 应用和 API',
      userCount: 45,
      permissions: 15,
      status: 'active',
      isSystem: false,
      createdAt: '2024-03-15',
    },
    {
      id: 4,
      name: '审计员',
      code: 'auditor',
      description: '查看审计日志和报表',
      userCount: 8,
      permissions: 8,
      status: 'active',
      isSystem: false,
      createdAt: '2024-05-20',
    },
    {
      id: 5,
      name: '普通用户',
      code: 'user',
      description: '基本账号访问权限',
      userCount: 847,
      permissions: 5,
      status: 'active',
      isSystem: true,
      createdAt: '2024-01-01',
    },
    {
      id: 6,
      name: '访客',
      code: 'guest',
      description: '只读访问权限',
      userCount: 156,
      permissions: 2,
      status: 'inactive',
      isSystem: false,
      createdAt: '2024-08-10',
    },
  ].sort((a, b) => Number(b.isSystem) - Number(a.isSystem)));
}

/**
 * 获取单个角色
 */
export async function getRole(id: number): Promise<Role | null> {
  const roles = await getRoles();
  return roles.find(r => r.id === id) || null;
}

/**
 * 创建角色
 */
export async function createRole(data: RoleFormData): Promise<Role> {
  return mockResponse({
    id: Date.now(),
    name: data.name,
    code: data.code,
    description: data.description,
    userCount: 0,
    permissions: 0,
    status: 'active',
    isSystem: false,
    createdAt: new Date().toISOString().split('T')[0],
  });
}

/**
 * 更新角色
 */
export async function updateRole(id: number, data: Partial<RoleFormData>): Promise<Role> {
  const role = await getRole(id);
  if (!role) throw new Error('Role not found');
  
  return mockResponse({
    ...role,
    ...data,
  });
}

/**
 * 删除角色
 */
export async function deleteRole(id: number): Promise<void> {
  return mockResponse(undefined);
}

/**
 * 获取权限组列表
 */
export async function getPermissionGroups(): Promise<PermissionGroup[]> {
  return mockResponse([
    {
      name: '用户管理',
      icon: 'pi pi-users',
      permissions: [
        { id: 'user:read', name: '查看用户', checked: true },
        { id: 'user:create', name: '创建用户', checked: true },
        { id: 'user:update', name: '编辑用户', checked: true },
        { id: 'user:delete', name: '删除用户', checked: false },
        { id: 'user:export', name: '导出用户', checked: true },
      ],
    },
    {
      name: '角色权限',
      icon: 'pi pi-shield',
      permissions: [
        { id: 'role:read', name: '查看角色', checked: true },
        { id: 'role:create', name: '创建角色', checked: false },
        { id: 'role:update', name: '编辑角色', checked: false },
        { id: 'role:delete', name: '删除角色', checked: false },
        { id: 'role:assign', name: '分配角色', checked: true },
      ],
    },
    {
      name: 'OAuth 应用',
      icon: 'pi pi-key',
      permissions: [
        { id: 'oauth:read', name: '查看应用', checked: true },
        { id: 'oauth:create', name: '创建应用', checked: true },
        { id: 'oauth:update', name: '编辑应用', checked: true },
        { id: 'oauth:delete', name: '删除应用', checked: false },
        { id: 'oauth:secret', name: '查看密钥', checked: false },
      ],
    },
    {
      name: '系统设置',
      icon: 'pi pi-cog',
      permissions: [
        { id: 'system:read', name: '查看设置', checked: true },
        { id: 'system:update', name: '修改设置', checked: false },
        { id: 'system:audit', name: '审计日志', checked: true },
        { id: 'system:backup', name: '数据备份', checked: false },
      ],
    },
  ]);
}

/**
 * 获取角色的权限
 */
export async function getRolePermissions(roleId: number): Promise<PermissionGroup[]> {
  // 返回默认权限组，实际应用中会根据角色ID返回不同的权限配置
  return getPermissionGroups();
}

/**
 * 保存角色权限
 */
export async function saveRolePermissions(roleId: number, permissions: PermissionGroup[]): Promise<void> {
  return mockResponse(undefined);
}
