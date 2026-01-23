/**
 * Dashboard API - 仪表盘相关接口
 */
import { mockResponse } from './index';
import type { StatCardData, Activity, TopApp } from '@/types';

// 统计卡片数据
export interface DashboardStats {
  cards: StatCardData[];
}

// 认证趋势数据
export interface AuthTrendData {
  labels: string[];
  data: number[];
}

// 用户分布数据
export interface UserDistData {
  labels: string[];
  data: number[];
  colors: string[];
}

/**
 * 获取仪表盘统计卡片数据
 */
export async function getDashboardStats(): Promise<DashboardStats> {
  return mockResponse({
    cards: [
      {
        title: '总用户数',
        value: '12,847',
        change: '+12.5%',
        changeType: 'increase',
        icon: 'pi pi-users',
        color: 'blue',
        trendData: [30, 42, 38, 52, 45, 58, 50, 65, 55, 72],
      },
      {
        title: 'OAuth 应用',
        value: '86',
        change: '+3',
        changeType: 'increase',
        icon: 'pi pi-key',
        color: 'orange',
        trendData: [40, 45, 42, 50, 48, 55, 52, 60, 58, 65],
      },
      {
        title: '今日认证',
        value: '2,451',
        change: '+8.2%',
        changeType: 'increase',
        icon: 'pi pi-shield',
        color: 'green',
        trendData: [35, 48, 40, 55, 45, 62, 52, 68, 58, 75],
      },
      {
        title: '活跃会话',
        value: '847',
        change: '-2.1%',
        changeType: 'decrease',
        icon: 'pi pi-bolt',
        color: 'purple',
        trendData: [70, 62, 68, 55, 60, 50, 58, 45, 52, 40],
      },
    ],
  });
}

/**
 * 获取认证趋势数据
 */
export async function getAuthTrendData(): Promise<AuthTrendData> {
  return mockResponse({
    labels: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
    data: [1200, 1900, 1500, 2100, 1800, 900, 1100],
  });
}

/**
 * 获取用户分布数据
 */
export async function getUserDistData(): Promise<UserDistData> {
  return mockResponse({
    labels: ['管理员', '普通用户', '开发者', '访客'],
    data: [12, 847, 234, 156],
    colors: ['#f97316', '#3b82f6', '#10b981', '#8b5cf6'],
  });
}

/**
 * 获取最近活动列表
 */
export async function getRecentActivities(): Promise<Activity[]> {
  return mockResponse([
    {
      user: 'zhang.wei@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhang',
      action: '登录成功',
      client: 'Web Dashboard',
      time: '2分钟前',
      status: 'success',
    },
    {
      user: 'li.ming@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=li',
      action: '密码重置',
      client: 'Mobile App',
      time: '5分钟前',
      status: 'warning',
    },
    {
      user: 'wang.fang@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=wang',
      action: 'OAuth 授权',
      client: 'Third Party App',
      time: '12分钟前',
      status: 'success',
    },
    {
      user: 'chen.hong@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=chen',
      action: '登录失败',
      client: 'API Client',
      time: '15分钟前',
      status: 'danger',
    },
    {
      user: 'zhao.yang@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhao',
      action: '新用户注册',
      client: 'Web Portal',
      time: '23分钟前',
      status: 'info',
    },
  ]);
}

/**
 * 获取热门应用列表
 */
export async function getTopApps(): Promise<TopApp[]> {
  return mockResponse([
    { name: 'Web Dashboard', users: 4521, percentage: 85 },
    { name: 'Mobile App', users: 3842, percentage: 72 },
    { name: 'API Gateway', users: 2156, percentage: 54 },
    { name: 'Admin Portal', users: 1234, percentage: 38 },
  ]);
}
