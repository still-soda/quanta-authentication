/**
 * Audit API - 审计日志相关接口
 */
import { mockResponse } from './index';
import type { AuditLog } from '@/types';

export interface AuditLogFilter {
  search?: string;
  module?: string | null;
  action?: string | null;
  dateRange?: [Date, Date] | null;
  page?: number;
  pageSize?: number;
}

export interface AuditLogResponse {
  data: AuditLog[];
  total: number;
}

/**
 * 获取审计日志列表
 */
export async function getAuditLogs(filter?: AuditLogFilter): Promise<AuditLogResponse> {
  const allLogs: AuditLog[] = [
    {
      id: '1',
      operatorId: 'u1',
      operatorName: '张伟',
      operatorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhang',
      module: '用户管理',
      action: '创建用户',
      targetId: 'u15',
      targetName: '新用户 - 李四',
      detail: { email: 'lisi@example.com', role: '普通用户' },
      ip: '192.168.1.100',
      time: '2026-01-23 14:30:25',
      durationMs: 125,
      status: 'success',
    },
    {
      id: '2',
      operatorId: 'u1',
      operatorName: '张伟',
      operatorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhang',
      module: 'OAuth应用',
      action: '重置密钥',
      targetId: 'app3',
      targetName: '内部管理系统',
      detail: { reason: '安全检查' },
      ip: '192.168.1.100',
      time: '2026-01-23 14:15:00',
      durationMs: 89,
      status: 'success',
    },
    {
      id: '3',
      operatorId: 'u2',
      operatorName: '李明',
      operatorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=li',
      module: '角色权限',
      action: '修改权限',
      targetId: 'role2',
      targetName: '开发者角色',
      detail: { added: ['oauth:create'], removed: [] },
      ip: '10.0.0.55',
      time: '2026-01-23 13:45:12',
      durationMs: 234,
      status: 'success',
    },
    {
      id: '4',
      operatorId: 'u3',
      operatorName: '王芳',
      operatorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=wang',
      module: '用户管理',
      action: '禁用用户',
      targetId: 'u8',
      targetName: '异常账户 - 测试',
      detail: { reason: '违规操作' },
      ip: '192.168.2.88',
      time: '2026-01-23 12:20:45',
      durationMs: 67,
      status: 'warning',
    },
    {
      id: '5',
      operatorId: 'u1',
      operatorName: '张伟',
      operatorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhang',
      module: '系统设置',
      action: '修改配置',
      targetId: 'config',
      targetName: '登录安全策略',
      detail: { maxAttempts: 5, lockDuration: 30 },
      ip: '192.168.1.100',
      time: '2026-01-23 11:00:00',
      durationMs: 156,
      status: 'success',
    },
    {
      id: '6',
      operatorId: 'u4',
      operatorName: '陈红',
      operatorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=chen',
      module: 'OAuth应用',
      action: '创建应用',
      targetId: 'app12',
      targetName: '新财务系统',
      detail: { type: 'web', redirectUri: 'https://finance.example.com/callback' },
      ip: '172.16.0.22',
      time: '2026-01-23 10:30:18',
      durationMs: 312,
      status: 'success',
    },
    {
      id: '7',
      operatorId: 'u2',
      operatorName: '李明',
      operatorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=li',
      module: '用户管理',
      action: '删除用户',
      targetId: 'u99',
      targetName: '已离职 - 张三',
      detail: { reason: '离职清退' },
      ip: '10.0.0.55',
      time: '2026-01-22 18:45:30',
      durationMs: 198,
      status: 'error',
    },
    {
      id: '8',
      operatorId: 'u5',
      operatorName: '赵阳',
      operatorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhaoyang',
      module: '认证登录',
      action: '登录失败',
      targetId: 'u5',
      targetName: '赵阳',
      detail: { error: '密码错误', attempts: 3 },
      ip: '203.0.113.45',
      time: '2026-01-22 16:20:00',
      durationMs: 45,
      status: 'error',
    },
    {
      id: '9',
      operatorId: 'u1',
      operatorName: '张伟',
      operatorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhang',
      module: '角色权限',
      action: '创建角色',
      targetId: 'role8',
      targetName: '访客角色',
      detail: { permissions: ['read:basic'] },
      ip: '192.168.1.100',
      time: '2026-01-22 14:10:22',
      durationMs: 178,
      status: 'success',
    },
    {
      id: '10',
      operatorId: 'u3',
      operatorName: '王芳',
      operatorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=wang',
      module: '数据导出',
      action: '导出用户列表',
      targetId: 'export1',
      targetName: '用户数据导出',
      detail: { format: 'xlsx', count: 1024 },
      ip: '192.168.2.88',
      time: '2026-01-22 11:30:00',
      durationMs: 2456,
      status: 'success',
    },
  ];

  return mockResponse({
    data: allLogs,
    total: allLogs.length,
  });
}

/**
 * 获取单条审计日志详情
 */
export async function getAuditLogDetail(id: string): Promise<AuditLog | null> {
  const { data } = await getAuditLogs();
  return data.find(log => log.id === id) || null;
}

/**
 * 导出审计日志
 */
export async function exportAuditLogs(filter?: AuditLogFilter): Promise<Blob> {
  // 模拟导出，返回空 Blob
  return mockResponse(new Blob([''], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' }));
}
