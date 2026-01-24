/**
 * 仪表盘相关类型定义
 */

import type { ThemeColor, ExtendedThemeColor } from './common'

// 统计卡片数据
export interface StatCardData {
   title: string
   value: string
   change: string
   changeType: 'increase' | 'decrease'
   icon: string
   color?: ThemeColor
   trendData?: number[]
}

export interface SimpleStatData {
   title: string
   value: number | string
   icon: string
   color?: ExtendedThemeColor
}

export interface MiniStatItem {
   label: string
   value: number | string
   colorClass?: string
}

// 活动记录
export interface Activity {
   user: string
   avatar: string
   action: string
   client: string
   time: string
   status: string
}

// 热门应用
export interface TopApp {
   name: string
   users: number
   percentage: number
}

// Dashboard API 响应类型
export interface DashboardStats {
   cards: StatCardData[]
}

export interface AuthTrendData {
   labels: string[]
   data: number[]
}

// 认证趋势时间范围类型
export type AuthTrendRange = 'weekly' | 'half-weekly' | 'monthly'

// 认证趋势请求参数
export interface AuthTrendParams {
   range?: AuthTrendRange
}

// 认证趋势响应数据
export interface AuthTrendResponse {
   auth_trend: number[]
}

// 用户分布数据 - 按角色分类
export interface UserDistData {
   labels: string[] // 角色名称
   data: number[] // 每个角色的用户数量
   colors: string[] // 对应的颜色
}

// 后端返回的 Dashboard 统计数据类型
export interface DashboardStatsResponse {
   user_count: number
   user_trend: number[]
   user_auth_count: number
   user_auth_trend: number[]
   oauth_app_count: number
   oauth_app_trend: number[]
   active_user_count: number
   active_user_trend: number[]
}
