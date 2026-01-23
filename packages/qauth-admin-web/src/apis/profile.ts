/**
 * Profile API - 个人资料相关接口
 */
import { mockResponse } from './index';

export interface ProfileData {
  id: string;
  studentId: string;
  name: string;
  displayName: string;
  email: string;
  phone: string;
  avatar: string;
  status: 'ACTIVE' | 'LOCKED' | 'BANNED';
  emailVerified: boolean;
  createdAt: string;
  lastLogin: string;
  bio: string;
}

export interface ProfileUpdateData {
  displayName?: string;
  phone?: string;
  bio?: string;
}

export interface UserRole {
  name: string;
  code: string;
  isSystem: boolean;
}

export interface LoginRecord {
  id: number;
  time: string;
  ip: string;
  location: string;
  device: string;
  status: 'success' | 'failed';
}

export interface SecuritySettings {
  mfaEnabled: boolean;
  emailNotifications: boolean;
  loginAlerts: boolean;
}

export interface PasswordChangeData {
  currentPassword: string;
  newPassword: string;
}

/**
 * 获取当前用户资料
 */
export async function getProfile(): Promise<ProfileData> {
  return mockResponse({
    id: 'u1',
    studentId: '20210001',
    name: '张伟',
    displayName: '管理员',
    email: 'zhang.wei@example.com',
    phone: '138****8888',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhang',
    status: 'ACTIVE',
    emailVerified: true,
    createdAt: '2024-06-15 10:30:00',
    lastLogin: '2026-01-23 14:30:00',
    bio: '系统管理员，负责平台日常运维和用户管理工作。',
  });
}

/**
 * 更新用户资料
 */
export async function updateProfile(data: ProfileUpdateData): Promise<ProfileData> {
  const profile = await getProfile();
  return mockResponse({
    ...profile,
    ...data,
  });
}

/**
 * 获取用户角色
 */
export async function getUserRoles(): Promise<UserRole[]> {
  return mockResponse([
    { name: '超级管理员', code: 'super_admin', isSystem: true },
    { name: '管理员', code: 'admin', isSystem: true },
  ]);
}

/**
 * 获取登录历史
 */
export async function getLoginHistory(): Promise<LoginRecord[]> {
  return mockResponse([
    {
      id: 1,
      time: '2026-01-23 14:30:00',
      ip: '192.168.1.100',
      location: '北京',
      device: 'Chrome / Windows',
      status: 'success',
    },
    {
      id: 2,
      time: '2026-01-22 09:15:00',
      ip: '192.168.1.100',
      location: '北京',
      device: 'Chrome / Windows',
      status: 'success',
    },
    {
      id: 3,
      time: '2026-01-21 18:45:00',
      ip: '10.0.0.55',
      location: '上海',
      device: 'Safari / macOS',
      status: 'success',
    },
    {
      id: 4,
      time: '2026-01-20 11:20:00',
      ip: '203.0.113.45',
      location: '未知',
      device: 'Firefox / Linux',
      status: 'failed',
    },
  ]);
}

/**
 * 获取安全设置
 */
export async function getSecuritySettings(): Promise<SecuritySettings> {
  return mockResponse({
    mfaEnabled: false,
    emailNotifications: true,
    loginAlerts: true,
  });
}

/**
 * 更新安全设置
 */
export async function updateSecuritySettings(data: Partial<SecuritySettings>): Promise<SecuritySettings> {
  const settings = await getSecuritySettings();
  return mockResponse({
    ...settings,
    ...data,
  });
}

/**
 * 修改密码
 */
export async function changePassword(data: PasswordChangeData): Promise<{ success: boolean }> {
  return mockResponse({ success: true });
}

/**
 * 上传头像
 */
export async function uploadAvatar(file: File): Promise<{ url: string }> {
  return mockResponse({
    url: `https://api.dicebear.com/7.x/avataaars/svg?seed=${Date.now()}`,
  });
}
