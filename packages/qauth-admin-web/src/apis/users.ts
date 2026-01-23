/**
 * Users API - 用户管理相关接口
 */
import { mockResponse } from './index';
import type { User, UserFormData } from '@/types';

/**
 * 获取用户列表
 */
export async function getUsers(): Promise<User[]> {
  return mockResponse([
    {
      id: 1,
      name: '张伟',
      email: 'zhang.wei@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhang',
      role: '管理员',
      status: 'active',
      lastLogin: '2026-01-23 10:30',
      createdAt: '2024-06-15',
    },
    {
      id: 2,
      name: '李明',
      email: 'li.ming@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=li',
      role: '开发者',
      status: 'active',
      lastLogin: '2026-01-22 18:45',
      createdAt: '2024-08-22',
    },
    {
      id: 3,
      name: '王芳',
      email: 'wang.fang@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=wang',
      role: '普通用户',
      status: 'inactive',
      lastLogin: '2026-01-10 09:15',
      createdAt: '2024-09-03',
    },
    {
      id: 4,
      name: '陈红',
      email: 'chen.hong@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=chen',
      role: '普通用户',
      status: 'active',
      lastLogin: '2026-01-23 08:20',
      createdAt: '2024-11-18',
    },
    {
      id: 5,
      name: '赵阳',
      email: 'zhao.yang@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhao',
      role: '开发者',
      status: 'locked',
      lastLogin: '2026-01-05 14:30',
      createdAt: '2025-01-02',
    },
    {
      id: 6,
      name: '刘洋',
      email: 'liu.yang@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=liu',
      role: '普通用户',
      status: 'active',
      lastLogin: '2026-01-21 16:45',
      createdAt: '2025-03-10',
    },
    {
      id: 7,
      name: '孙静',
      email: 'sun.jing@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=sun',
      role: '管理员',
      status: 'active',
      lastLogin: '2026-01-23 11:00',
      createdAt: '2024-05-20',
    },
    {
      id: 8,
      name: '周杰',
      email: 'zhou.jie@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhou',
      role: '开发者',
      status: 'pending',
      lastLogin: '-',
      createdAt: '2026-01-20',
    },
  ]);
}

/**
 * 获取单个用户
 */
export async function getUser(id: number): Promise<User | null> {
  const users = await getUsers();
  return users.find(u => u.id === id) || null;
}

/**
 * 创建用户
 */
export async function createUser(data: UserFormData): Promise<User> {
  return mockResponse({
    id: Date.now(),
    name: data.name,
    email: data.email,
    avatar: `https://api.dicebear.com/7.x/avataaars/svg?seed=${Date.now()}`,
    role: data.role,
    status: data.status ? 'active' : 'inactive',
    lastLogin: '-',
    createdAt: new Date().toISOString().split('T')[0],
  });
}

/**
 * 更新用户
 */
export async function updateUser(id: number, data: Partial<UserFormData>): Promise<User> {
  const user = await getUser(id);
  if (!user) throw new Error('User not found');
  
  return mockResponse({
    ...user,
    ...data,
    status: data.status !== undefined ? (data.status ? 'active' : 'inactive') : user.status,
  });
}

/**
 * 删除用户
 */
export async function deleteUser(id: number): Promise<void> {
  return mockResponse(undefined);
}

/**
 * 重置用户密码
 */
export async function resetUserPassword(id: number): Promise<{ success: boolean }> {
  return mockResponse({ success: true });
}

/**
 * 禁用用户
 */
export async function disableUser(id: number): Promise<User> {
  const user = await getUser(id);
  if (!user) throw new Error('User not found');
  
  return mockResponse({
    ...user,
    status: 'inactive',
  });
}

/**
 * 启用用户
 */
export async function enableUser(id: number): Promise<User> {
  const user = await getUser(id);
  if (!user) throw new Error('User not found');
  
  return mockResponse({
    ...user,
    status: 'active',
  });
}
