/**
 * OAuth API - OAuth 应用相关接口
 */
import { httpClient } from './index'
import type {
   OAuthApp,
   OAuthAppFormData,
   CreateOAuthAppResponse,
   ListOAuthAppsParams,
   ListOAuthAppsResult,
   OAuthClientStats,
   OAuthClientOptions,
} from '@/types'

/**
 * 获取 OAuth 应用列表
 */
export async function getOAuthApps(params: ListOAuthAppsParams = {}): Promise<ListOAuthAppsResult> {
   const response = await httpClient.get('/_/v1/clients', { params })
   return response.data.data
}

/**
 * 获取单个 OAuth 应用
 */
export async function getOAuthApp(id: string): Promise<OAuthApp> {
   const response = await httpClient.get(`/_/v1/clients/${id}`)
   return response.data.data
}

/**
 * 创建 OAuth 应用
 */
export async function createOAuthApp(data: OAuthAppFormData): Promise<CreateOAuthAppResponse> {
   const response = await httpClient.post('/_/v1/clients', data)
   return response.data.data
}

/**
 * 更新 OAuth 应用
 */
export async function updateOAuthApp(
   id: string,
   data: Partial<OAuthAppFormData>
): Promise<OAuthApp> {
   const response = await httpClient.put(`/_/v1/clients/${id}`, data)
   return response.data.data
}

/**
 * 删除 OAuth 应用
 */
export async function deleteOAuthApp(id: string): Promise<void> {
   await httpClient.delete(`/_/v1/clients/${id}`)
}

/**
 * 重新生成客户端密钥
 */
export async function regenerateClientSecret(id: string): Promise<{ secret: string }> {
   const response = await httpClient.post(`/_/v1/clients/${id}/regenerate-secret`)
   return response.data.data
}

/**
 * 获取 OAuth 客户端统计
 */
export async function getOAuthClientStats(): Promise<OAuthClientStats> {
   const response = await httpClient.get('/_/v1/clients/stats')
   return response.data.data
}

/**
 * 获取 OAuth 客户端可选配置项（授权类型、授权范围）
 */
export async function getOAuthClientOptions(): Promise<OAuthClientOptions> {
   const response = await httpClient.get('/_/v1/clients/options')
   return response.data.data
}
