/**
 * 系统设置相关类型定义
 */

import type { OAuthSettings } from './oauth'

// 设置组
export interface SettingGroup {
   id: string
   label: string
   icon: string
}

// 通用设置
export interface GeneralSettings {
   siteName: string
   siteDescription: string
   adminEmail: string
   defaultLanguage: string
   timezone: string
   maintenanceMode: boolean
   allowRegistration: boolean
   requireEmailVerification: boolean
}

// 安全设置配置
export interface SecuritySettingsConfig {
   passwordMinLength: number
   passwordRequireUppercase: boolean
   passwordRequireLowercase: boolean
   passwordRequireNumbers: boolean
   passwordRequireSpecial: boolean
   maxLoginAttempts: number
   lockoutDuration: number
   sessionTimeout: number
   enableMfa: boolean
   enforceHttps: boolean
   allowedIpRanges: string
   csrfProtection: boolean
}

// 邮件设置
export interface EmailSettings {
   smtpHost: string
   smtpPort: number
   smtpUser: string
   smtpPassword: string
   smtpEncryption: string
   fromAddress: string
   fromName: string
   enableEmailNotifications: boolean
}

// 存储设置
export interface StorageSettings {
   storageType: string
   localPath: string
   s3Bucket: string
   s3Region: string
   s3AccessKey: string
   s3SecretKey: string
   maxFileSize: number
   allowedFileTypes: string
}

// 所有设置
export interface AllSettings {
   general: GeneralSettings
   security: SecuritySettingsConfig
   oauth: OAuthSettings
   email: EmailSettings
   storage: StorageSettings
}
