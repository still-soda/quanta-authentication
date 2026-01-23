/**
 * Notifications API - 通知相关接口
 */
import { mockResponse } from './index';
import type { Notification, NotificationType } from '@/types';

export interface NotificationFilter {
  type?: NotificationType | 'all';
  read?: boolean;
}

/**
 * 获取通知列表
 */
export async function getNotifications(filter?: NotificationFilter): Promise<Notification[]> {
  const allNotifications: Notification[] = [
    {
      id: '1',
      type: 'security',
      title: '新设备登录提醒',
      message: '检测到您的账号在新设备上登录：Chrome / Windows，IP: 192.168.1.100，位置：北京',
      time: '2026-01-23 14:30:00',
      read: false,
      actionUrl: '/audit',
      actionLabel: '查看详情',
    },
    {
      id: '2',
      type: 'system',
      title: '系统维护通知',
      message: '系统将于 2026年1月25日 02:00-06:00 进行例行维护，届时服务将暂停访问',
      time: '2026-01-23 10:00:00',
      read: false,
    },
    {
      id: '3',
      type: 'oauth',
      title: 'OAuth 应用审核通过',
      message: '您提交的应用「内部管理系统」已通过审核，现在可以正常使用',
      time: '2026-01-22 16:45:00',
      read: false,
      actionUrl: '/oauth',
      actionLabel: '前往管理',
    },
    {
      id: '4',
      type: 'user',
      title: '新用户注册通知',
      message: '新用户 李四 (lisi@example.com) 已完成注册并通过邮箱验证',
      time: '2026-01-22 14:20:00',
      read: true,
      actionUrl: '/users',
      actionLabel: '查看用户',
    },
    {
      id: '5',
      type: 'alert',
      title: '异常登录警告',
      message: '用户 赵阳 在短时间内多次登录失败，已触发账号锁定机制',
      time: '2026-01-22 11:30:00',
      read: true,
      actionUrl: '/audit',
      actionLabel: '查看日志',
    },
    {
      id: '6',
      type: 'security',
      title: '密码修改成功',
      message: '您的账号密码已成功修改，如非本人操作请立即联系管理员',
      time: '2026-01-21 18:00:00',
      read: true,
    },
    {
      id: '7',
      type: 'system',
      title: '系统更新完成',
      message: '系统已升级至 v2.5.0 版本，新增组织架构管理等功能',
      time: '2026-01-20 09:00:00',
      read: true,
    },
    {
      id: '8',
      type: 'oauth',
      title: 'OAuth 密钥即将过期',
      message: '应用「客户端 App」的密钥将于 30 天后过期，请及时更新',
      time: '2026-01-19 14:00:00',
      read: true,
      actionUrl: '/oauth',
      actionLabel: '立即更新',
    },
    {
      id: '9',
      type: 'user',
      title: '用户权限变更',
      message: '管理员为您分配了新角色：审计员，您现在可以查看系统审计日志',
      time: '2026-01-18 11:20:00',
      read: true,
    },
    {
      id: '10',
      type: 'alert',
      title: '存储空间预警',
      message: '系统存储空间使用率已达 85%，建议清理无用文件或扩容',
      time: '2026-01-17 16:30:00',
      read: true,
    },
  ];

  return mockResponse(allNotifications);
}

/**
 * 获取未读通知数量
 */
export async function getUnreadCount(): Promise<number> {
  const notifications = await getNotifications();
  return notifications.filter(n => !n.read).length;
}

/**
 * 标记通知为已读
 */
export async function markAsRead(id: string): Promise<void> {
  return mockResponse(undefined);
}

/**
 * 批量标记通知为已读
 */
export async function markMultipleAsRead(ids: string[]): Promise<void> {
  return mockResponse(undefined);
}

/**
 * 标记所有通知为已读
 */
export async function markAllAsRead(): Promise<void> {
  return mockResponse(undefined);
}

/**
 * 删除通知
 */
export async function deleteNotification(id: string): Promise<void> {
  return mockResponse(undefined);
}

/**
 * 批量删除通知
 */
export async function deleteMultipleNotifications(ids: string[]): Promise<void> {
  return mockResponse(undefined);
}

/**
 * 删除所有已读通知
 */
export async function deleteAllReadNotifications(): Promise<void> {
  return mockResponse(undefined);
}
