/**
 * 个人资料相关类型定义
 */

// 个人资料数据
export interface ProfileData {
   id: string
   studentId: string
   name: string
   displayName: string
   email: string
   phone: string
   avatar: string
   status: 'ACTIVE' | 'LOCKED' | 'BANNED'
   emailVerified: boolean
   createdAt: string
   lastLogin: string
   bio: string
}

// 个人资料更新数据
export interface ProfileUpdateData {
   displayName?: string
   phone?: string
   bio?: string
}

// 登录记录
export interface LoginRecord {
   id: number
   time: string
   ip: string
   location: string
   device: string
   status: 'success' | 'failed'
}

// 安全设置
export interface SecuritySettings {
   mfaEnabled: boolean
   emailNotifications: boolean
   loginAlerts: boolean
}

// 密码修改数据
export interface PasswordChangeData {
   currentPassword: string
   newPassword: string
}
