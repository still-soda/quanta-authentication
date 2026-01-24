/**
 * Settings API - 系统设置相关接口
 */
import { mockResponse } from './index'
import type {
   GeneralSettings,
   SecuritySettingsConfig,
   OAuthSettings,
   EmailSettings,
   StorageSettings,
   AllSettings,
} from '@/types'

/**
 * 获取所有设置
 */
export async function getAllSettings(): Promise<AllSettings> {
   return mockResponse({
      general: {
         siteName: 'Quanta 认证中心',
         siteDescription: '企业级统一身份认证平台',
         adminEmail: 'admin@example.com',
         defaultLanguage: 'zh-CN',
         timezone: 'Asia/Shanghai',
         maintenanceMode: false,
         allowRegistration: true,
         requireEmailVerification: true,
      },
      security: {
         passwordMinLength: 8,
         passwordRequireUppercase: true,
         passwordRequireLowercase: true,
         passwordRequireNumbers: true,
         passwordRequireSpecial: false,
         maxLoginAttempts: 5,
         lockoutDuration: 30,
         sessionTimeout: 60,
         enableMfa: false,
         enforceHttps: true,
         allowedIpRanges: '',
         csrfProtection: true,
      },
      oauth: {
         accessTokenLifetime: 3600,
         refreshTokenLifetime: 604800,
         authCodeLifetime: 600,
         allowImplicitGrant: false,
         allowPasswordGrant: false,
         allowClientCredentials: true,
         requirePkce: true,
         allowedScopes: 'openid profile email',
         jwksRotationDays: 90,
      },
      email: {
         smtpHost: 'smtp.example.com',
         smtpPort: 587,
         smtpUser: 'noreply@example.com',
         smtpPassword: '********',
         smtpEncryption: 'tls',
         fromAddress: 'noreply@example.com',
         fromName: 'Quanta Auth',
         enableEmailNotifications: true,
      },
      storage: {
         storageType: 'local',
         localPath: '/var/data/uploads',
         s3Bucket: '',
         s3Region: '',
         s3AccessKey: '',
         s3SecretKey: '',
         maxFileSize: 10,
         allowedFileTypes: 'jpg,png,gif,pdf,doc,docx',
      },
   })
}

/**
 * 获取基本设置
 */
export async function getGeneralSettings(): Promise<GeneralSettings> {
   const all = await getAllSettings()
   return all.general
}

/**
 * 更新基本设置
 */
export async function updateGeneralSettings(
   data: Partial<GeneralSettings>
): Promise<GeneralSettings> {
   const current = await getGeneralSettings()
   return mockResponse({ ...current, ...data })
}

/**
 * 获取安全设置
 */
export async function getSecuritySettingsConfig(): Promise<SecuritySettingsConfig> {
   const all = await getAllSettings()
   return all.security
}

/**
 * 更新安全设置
 */
export async function updateSecuritySettingsConfig(
   data: Partial<SecuritySettingsConfig>
): Promise<SecuritySettingsConfig> {
   const current = await getSecuritySettingsConfig()
   return mockResponse({ ...current, ...data })
}

/**
 * 获取 OAuth 设置
 */
export async function getOAuthSettings(): Promise<OAuthSettings> {
   const all = await getAllSettings()
   return all.oauth
}

/**
 * 更新 OAuth 设置
 */
export async function updateOAuthSettings(data: Partial<OAuthSettings>): Promise<OAuthSettings> {
   const current = await getOAuthSettings()
   return mockResponse({ ...current, ...data })
}

/**
 * 获取邮件设置
 */
export async function getEmailSettings(): Promise<EmailSettings> {
   const all = await getAllSettings()
   return all.email
}

/**
 * 更新邮件设置
 */
export async function updateEmailSettings(data: Partial<EmailSettings>): Promise<EmailSettings> {
   const current = await getEmailSettings()
   return mockResponse({ ...current, ...data })
}

/**
 * 获取存储设置
 */
export async function getStorageSettings(): Promise<StorageSettings> {
   const all = await getAllSettings()
   return all.storage
}

/**
 * 更新存储设置
 */
export async function updateStorageSettings(
   data: Partial<StorageSettings>
): Promise<StorageSettings> {
   const current = await getStorageSettings()
   return mockResponse({ ...current, ...data })
}

/**
 * 保存所有设置
 */
export async function saveAllSettings(data: Partial<AllSettings>): Promise<AllSettings> {
   const current = await getAllSettings()
   return mockResponse({
      general: { ...current.general, ...data.general },
      security: { ...current.security, ...data.security },
      oauth: { ...current.oauth, ...data.oauth },
      email: { ...current.email, ...data.email },
      storage: { ...current.storage, ...data.storage },
   })
}

/**
 * 重置设置为默认值
 */
export async function resetSettings(): Promise<AllSettings> {
   return getAllSettings()
}

/**
 * 测试邮件发送
 */
export async function testEmailSend(): Promise<{ success: boolean; message: string }> {
   return mockResponse({
      success: true,
      message: '测试邮件已发送',
   })
}
