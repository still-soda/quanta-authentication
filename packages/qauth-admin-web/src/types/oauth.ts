/**
 * OAuth 相关类型定义
 */

// OAuth 客户端状态
export type OAuthClientStatus = 'active' | 'development' | 'deprecated'

// OAuth 应用（API 响应格式）
export interface OAuthApp {
   id: string
   client_id: string
   name: string
   description: string
   domain: string
   redirect_uris: string[]
   scopes: string[]
   grant_types: string[]
   status: OAuthClientStatus
   trusted: boolean
   logo: string
   icon: string
   icon_bg: string
   last_used_at: string | null
   request_count: number
   created_at: string
   updated_at: string
}

// OAuth 应用表单数据
export interface OAuthAppFormData {
   name: string
   description: string
   domain: string
   redirect_uris: string[]
   scopes: string[]
   grant_types: string[]
   status: OAuthClientStatus
   trusted: boolean
   logo?: string
   icon?: string
   icon_bg?: string
}

// 创建 OAuth 应用响应
export interface CreateOAuthAppResponse {
   client: OAuthApp
   client_secret: string
}

// 列表查询参数
export interface ListOAuthAppsParams {
   page?: number
   page_size?: number
   search?: string
   status?: OAuthClientStatus
   sort_by?: string
   sort_desc?: boolean
}

// 列表响应
export interface ListOAuthAppsResult {
   items: OAuthApp[]
   total: number
   page: number
   size: number
}

// OAuth 客户端统计
export interface OAuthClientStats {
   total: number
   active: number
   development: number
   deprecated: number
}

// OAuth 选项（授权类型、授权范围）
export interface OAuthOption {
   label: string
   value: string
}

// OAuth 客户端可选配置项
export interface OAuthClientOptions {
   scopes: OAuthOption[]
   grant_types: OAuthOption[]
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
