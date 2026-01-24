/**
 * OAuth 相关类型定义
 */

// OAuth 应用
export interface OAuthApp {
   id: number
   name: string
   clientId: string
   description: string
   icon: string
   iconBg: string
   redirectUris: string[]
   scopes: string[]
   grantTypes: string[]
   status: string
   trusted: boolean
   createdAt: string
   lastUsed: string
   requestCount: number
}

// OAuth 应用表单数据
export interface OAuthAppFormData {
   name: string
   description: string
   redirectUris: string
   scopes: string[]
   trusted: boolean
}

// OAuth 设置
export interface OAuthSettings {
   accessTokenLifetime: number
   refreshTokenLifetime: number
   authCodeLifetime: number
   allowImplicitGrant: boolean
   allowPasswordGrant: boolean
   allowClientCredentials: boolean
   requirePkce: boolean
   allowedScopes: string
   jwksRotationDays: number
}
