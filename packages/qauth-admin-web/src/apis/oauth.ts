/**
 * OAuth API - OAuth 应用相关接口
 */
import { mockResponse } from './index';
import type { OAuthApp, OAuthAppFormData } from '@/types';

/**
 * 获取 OAuth 应用列表
 */
export async function getOAuthApps(): Promise<OAuthApp[]> {
  return mockResponse([
    {
      id: 1,
      name: 'Web Dashboard',
      clientId: 'web_dashboard_prod',
      description: '主要的 Web 管理后台应用',
      icon: 'pi pi-desktop',
      iconBg: 'linear-gradient(135deg, #f97316 0%, #ea580c 100%)',
      redirectUris: ['https://dashboard.example.com/callback'],
      scopes: ['openid', 'profile', 'email', 'admin'],
      grantTypes: ['authorization_code', 'refresh_token'],
      status: 'active',
      trusted: true,
      createdAt: '2024-01-15',
      lastUsed: '2026-01-23',
      requestCount: 125840,
    },
    {
      id: 2,
      name: 'Mobile App',
      clientId: 'mobile_app_v2',
      description: 'iOS 和 Android 移动应用',
      icon: 'pi pi-mobile',
      iconBg: 'linear-gradient(135deg, #3b82f6 0%, #2563eb 100%)',
      redirectUris: ['myapp://callback', 'https://mobile.example.com/callback'],
      scopes: ['openid', 'profile', 'email', 'offline_access'],
      grantTypes: ['authorization_code', 'refresh_token'],
      status: 'active',
      trusted: true,
      createdAt: '2024-03-20',
      lastUsed: '2026-01-23',
      requestCount: 89562,
    },
    {
      id: 3,
      name: 'Partner API',
      clientId: 'partner_api_client',
      description: '第三方合作伙伴 API 访问',
      icon: 'pi pi-link',
      iconBg: 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
      redirectUris: ['https://partner.example.com/oauth/callback'],
      scopes: ['openid', 'profile', 'read:users'],
      grantTypes: ['client_credentials'],
      status: 'active',
      trusted: false,
      createdAt: '2024-06-10',
      lastUsed: '2026-01-22',
      requestCount: 34521,
    },
    {
      id: 4,
      name: 'Internal Tools',
      clientId: 'internal_tools_dev',
      description: '内部开发工具和脚本',
      icon: 'pi pi-wrench',
      iconBg: 'linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%)',
      redirectUris: ['http://localhost:3000/callback'],
      scopes: ['openid', 'profile', 'admin'],
      grantTypes: ['authorization_code', 'client_credentials'],
      status: 'development',
      trusted: true,
      createdAt: '2024-09-05',
      lastUsed: '2026-01-20',
      requestCount: 8421,
    },
    {
      id: 5,
      name: 'Legacy System',
      clientId: 'legacy_system_v1',
      description: '旧版系统兼容接口',
      icon: 'pi pi-history',
      iconBg: 'linear-gradient(135deg, #6b7280 0%, #4b5563 100%)',
      redirectUris: ['https://legacy.example.com/auth'],
      scopes: ['openid', 'profile'],
      grantTypes: ['authorization_code'],
      status: 'deprecated',
      trusted: true,
      createdAt: '2023-05-01',
      lastUsed: '2026-01-15',
      requestCount: 2150,
    },
  ]);
}

/**
 * 获取单个 OAuth 应用
 */
export async function getOAuthApp(id: number): Promise<OAuthApp | null> {
  const apps = await getOAuthApps();
  return apps.find(a => a.id === id) || null;
}

/**
 * 创建 OAuth 应用
 */
export async function createOAuthApp(data: OAuthAppFormData): Promise<OAuthApp> {
  const clientId = `app_${Date.now()}_${Math.random().toString(36).substring(2, 8)}`;
  
  return mockResponse({
    id: Date.now(),
    name: data.name,
    clientId,
    description: data.description,
    icon: 'pi pi-box',
    iconBg: 'linear-gradient(135deg, #6366f1 0%, #4f46e5 100%)',
    redirectUris: data.redirectUris.split('\n').filter(Boolean),
    scopes: data.scopes,
    grantTypes: ['authorization_code', 'refresh_token'],
    status: 'development',
    trusted: data.trusted,
    createdAt: new Date().toISOString().split('T')[0],
    lastUsed: '-',
    requestCount: 0,
  });
}

/**
 * 更新 OAuth 应用
 */
export async function updateOAuthApp(id: number, data: Partial<OAuthAppFormData>): Promise<OAuthApp> {
  const app = await getOAuthApp(id);
  if (!app) throw new Error('OAuth App not found');
  
  return mockResponse({
    ...app,
    name: data.name ?? app.name,
    description: data.description ?? app.description,
    redirectUris: data.redirectUris ? data.redirectUris.split('\n').filter(Boolean) : app.redirectUris,
    scopes: data.scopes ?? app.scopes,
    trusted: data.trusted ?? app.trusted,
  });
}

/**
 * 删除 OAuth 应用
 */
export async function deleteOAuthApp(id: number): Promise<void> {
  return mockResponse(undefined);
}

/**
 * 重新生成客户端密钥
 */
export async function regenerateClientSecret(id: number): Promise<{ secret: string }> {
  const secret = 'sk_live_' + Math.random().toString(36).substring(2, 34);
  return mockResponse({ secret });
}
