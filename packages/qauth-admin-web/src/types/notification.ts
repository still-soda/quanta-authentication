/**
 * 通知相关类型定义
 */

// 通知类型
export type NotificationType = 'system' | 'security' | 'user' | 'oauth' | 'alert'

// 通知信息
export interface Notification {
   id: string
   type: NotificationType
   title: string
   message: string
   time: string
   read: boolean
   actionUrl?: string
   actionLabel?: string
}

// 通知过滤器
export interface NotificationFilter {
   type?: NotificationType | 'all'
   read?: boolean
}
