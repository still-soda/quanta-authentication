/**
 * Dashboard API - 仪表盘相关接口
 */
import { httpClient, mockResponse } from './index'
import type {
   Activity,
   TopApp,
   DashboardStats,
   AuthTrendData,
   UserDistData,
   StatCardData,
   DashboardStatsResponse,
   AuditStatsResponse,
   RecentActivityResponse,
   TopClientResponse,
} from '@/types'

/**
 * 计算环比变化率
 */
function calculateChange(
   current: number,
   previous: number
): { change: string; changeType: 'increase' | 'decrease' } {
   if (previous === 0) {
      return { change: current > 0 ? '+100%' : '0%', changeType: 'increase' }
   }
   const changeRate = ((current - previous) / previous) * 100
   const isIncrease = changeRate >= 0
   const formattedChange = `${isIncrease ? '+' : ''}${changeRate.toFixed(1)}%`
   return { change: formattedChange, changeType: isIncrease ? 'increase' : 'decrease' }
}

/**
 * 格式化数字（添加千分位分隔符）
 */
function formatNumber(num: number): string {
   return num.toLocaleString('zh-CN')
}

/**
 * 安全获取数组的倒数第二个元素
 */
function getSecondLast(arr: number[]): number {
   if (arr.length >= 2) {
      return arr[arr.length - 2] ?? 0
   }
   return 0
}

/**
 * 获取仪表盘统计卡片数据
 */
export async function getDashboardStats(): Promise<DashboardStats> {
   const response = await httpClient.get<{ code: number; data: DashboardStatsResponse }>(
      '/_/v1/dashboard/stats'
   )
   const data = response.data.data

   // 计算各项指标的变化率（对比上一期数据）
   const userChange = calculateChange(data.user_count, getSecondLast(data.user_trend))
   const authChange = calculateChange(data.user_auth_count, getSecondLast(data.user_auth_trend))
   const oauthChange = calculateChange(data.oauth_app_count, getSecondLast(data.oauth_app_trend))
   const activeChange = calculateChange(
      data.active_user_count,
      getSecondLast(data.active_user_trend)
   )

   const cards: StatCardData[] = [
      {
         title: '总用户数',
         value: formatNumber(data.user_count),
         change: userChange.change,
         changeType: userChange.changeType,
         icon: 'pi pi-users',
         color: 'blue',
         trendData: data.user_trend,
      },
      {
         title: 'OAuth 应用',
         value: formatNumber(data.oauth_app_count),
         change: oauthChange.change,
         changeType: oauthChange.changeType,
         icon: 'pi pi-key',
         color: 'orange',
         trendData: data.oauth_app_trend,
      },
      {
         title: '今日认证',
         value: formatNumber(data.user_auth_count),
         change: authChange.change,
         changeType: authChange.changeType,
         icon: 'pi pi-shield',
         color: 'green',
         trendData: data.user_auth_trend,
      },
      {
         title: '活跃会话',
         value: formatNumber(data.active_user_count),
         change: activeChange.change,
         changeType: activeChange.changeType,
         icon: 'pi pi-bolt',
         color: 'purple',
         trendData: data.active_user_trend,
      },
   ]

   return { cards }
}

/**
 * 获取认证趋势数据
 * 从审计统计 API 获取数据
 */
export async function getAuthTrendData(): Promise<AuthTrendData> {
   try {
      const response = await httpClient.get<{ code: number; data: AuditStatsResponse }>(
         '/_/v1/audit/stats'
      )
      const actionStats = response.data.data.action_stats

      // 将操作类型统计转换为趋势图表数据
      const actionLabels: string[] = []
      const actionData: number[] = []

      // 定义操作类型的中文映射
      const actionNameMap: Record<string, string> = {
         LOGIN: '登录',
         LOGOUT: '登出',
         REGISTER: '注册',
         OAUTH_AUTHORIZE: 'OAuth授权',
         TOKEN_REFRESH: '令牌刷新',
         PASSWORD_RESET: '密码重置',
         PASSWORD_CHANGE: '密码修改',
      }

      for (const [action, count] of Object.entries(actionStats)) {
         if (count > 0) {
            actionLabels.push(actionNameMap[action] || action)
            actionData.push(count)
         }
      }

      return {
         labels: actionLabels.length > 0 ? actionLabels : ['暂无数据'],
         data: actionData.length > 0 ? actionData : [0],
      }
   } catch (error) {
      // 如果 API 调用失败，返回默认数据
      return mockResponse({
         labels: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
         data: [0, 0, 0, 0, 0, 0, 0],
      })
   }
}

/**
 * 获取用户分布数据
 * 从审计统计 API 获取数据
 */
export async function getUserDistData(): Promise<UserDistData> {
   try {
      const response = await httpClient.get<{ code: number; data: AuditStatsResponse }>(
         '/_/v1/audit/stats'
      )
      const moduleStats = response.data.data.module_stats

      // 将模块统计转换为饼图数据
      const moduleLabels: string[] = []
      const moduleData: number[] = []

      // 定义模块的中文映射
      const moduleNameMap: Record<string, string> = {
         AUTH: '认证模块',
         OAUTH: 'OAuth模块',
         USER: '用户管理',
         ROLE: '角色管理',
         PERMISSION: '权限管理',
         CLIENT: '客户端管理',
         SYSTEM: '系统管理',
      }

      // 定义颜色
      const colors = ['#f97316', '#3b82f6', '#10b981', '#8b5cf6', '#06b6d4', '#ec4899', '#84cc16']

      let colorIndex = 0
      const usedColors: string[] = []

      for (const [module, count] of Object.entries(moduleStats)) {
         if (count > 0) {
            moduleLabels.push(moduleNameMap[module] || module)
            moduleData.push(count)
            const color = colors[colorIndex % colors.length]
            if (color) {
               usedColors.push(color)
            }
            colorIndex++
         }
      }

      return {
         labels: moduleLabels.length > 0 ? moduleLabels : ['暂无数据'],
         data: moduleData.length > 0 ? moduleData : [1],
         colors: usedColors.length > 0 ? usedColors : ['#e5e7eb'],
      }
   } catch (error) {
      // 如果 API 调用失败，返回默认数据
      return mockResponse({
         labels: ['暂无数据'],
         data: [1],
         colors: ['#e5e7eb'],
      })
   }
}

/**
 * 格式化时间为相对时间
 */
function formatRelativeTime(dateStr: string): string {
   const date = new Date(dateStr)
   const now = new Date()
   const diffMs = now.getTime() - date.getTime()
   const diffSeconds = Math.floor(diffMs / 1000)
   const diffMinutes = Math.floor(diffSeconds / 60)
   const diffHours = Math.floor(diffMinutes / 60)
   const diffDays = Math.floor(diffHours / 24)

   if (diffSeconds < 60) {
      return '刚刚'
   } else if (diffMinutes < 60) {
      return `${diffMinutes}分钟前`
   } else if (diffHours < 24) {
      return `${diffHours}小时前`
   } else if (diffDays < 7) {
      return `${diffDays}天前`
   } else {
      return date.toLocaleDateString('zh-CN')
   }
}

/**
 * 映射操作类型到显示文本
 */
function mapActionToDisplay(action: string): string {
   const actionMap: Record<string, string> = {
      LOGIN: '登录成功',
      LOGIN_FAILED: '登录失败',
      LOGOUT: '退出登录',
      REGISTER: '新用户注册',
      PASSWORD_RESET: '密码重置',
      PASSWORD_CHANGE: '密码修改',
      TOKEN_REFRESH: '令牌刷新',
      OAUTH_AUTHORIZE: 'OAuth 授权',
      OAUTH_TOKEN: 'OAuth 令牌',
      OAUTH_REVOKE: 'OAuth 撤销',
      CLIENT_CREATE: '创建客户端',
      CLIENT_UPDATE: '更新客户端',
      CLIENT_DELETE: '删除客户端',
      ROLE_CREATE: '创建角色',
      ROLE_UPDATE: '更新角色',
      ROLE_DELETE: '删除角色',
      PERMISSION_GRANT: '授予权限',
      PERMISSION_REVOKE: '撤销权限',
      USER_CREATE: '创建用户',
      USER_UPDATE: '更新用户',
      USER_DELETE: '删除用户',
   }
   return actionMap[action] || action
}

/**
 * 映射状态到前端状态
 */
function mapStatusToFrontend(status: string): string {
   const statusMap: Record<string, string> = {
      success: 'success',
      failure: 'danger',
      warning: 'warning',
      pending: 'info',
   }
   return statusMap[status.toLowerCase()] || 'info'
}

/**
 * 获取最近活动列表
 */
export async function getRecentActivities(): Promise<Activity[]> {
   try {
      const response = await httpClient.get<{ code: number; data: RecentActivityResponse[] }>(
         '/_/v1/audit/activities'
      )
      const activities = response.data.data

      return activities.map(activity => ({
         user: activity.user || '未知用户',
         avatar:
            activity.avatar ||
            `https://api.dicebear.com/7.x/avataaars/svg?seed=${activity.operator_id}`,
         action: mapActionToDisplay(activity.action),
         client: activity.client || 'System',
         time: formatRelativeTime(activity.time),
         status: mapStatusToFrontend(activity.status),
      }))
   } catch (error) {
      // 如果 API 调用失败，返回空数组
      return []
   }
}

/**
 * 获取热门应用列表
 */
export async function getTopApps(): Promise<TopApp[]> {
   try {
      const response = await httpClient.get<{ code: number; data: TopClientResponse[] }>(
         '/_/v1/audit/top-clients'
      )
      const clients = response.data.data

      if (!clients || clients.length === 0) {
         return []
      }

      // 找到最大值用于计算百分比
      const maxCount = Math.max(...clients.map(c => c.count))

      return clients.map(client => ({
         name: client.client_name || client.client_id,
         users: client.count,
         percentage: maxCount > 0 ? Math.round((client.count / maxCount) * 100) : 0,
      }))
   } catch (error) {
      // 如果 API 调用失败，返回空数组
      return []
   }
}
