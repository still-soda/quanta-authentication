// 主题色类型
export type ThemeColor = 'orange' | 'blue' | 'green' | 'purple' | 'cyan' | 'red';
export type ExtendedThemeColor = ThemeColor | 'gray';

// 通用状态类型
export type Status = 'success' | 'warning' | 'error';

// 仪表盘相关
export interface StatCardData {
  title: string;
  value: string;
  change: string;
  changeType: 'increase' | 'decrease';
  icon: string;
  color?: ThemeColor;
  trendData?: number[];
}

export interface SimpleStatData {
  title: string;
  value: number | string;
  icon: string;
  color?: ExtendedThemeColor;
}

export interface MiniStatItem {
  label: string;
  value: number | string;
  colorClass?: string;
}

export interface Activity {
  user: string;
  avatar: string;
  action: string;
  client: string;
  time: string;
  status: string;
}

export interface TopApp {
  name: string;
  users: number;
  percentage: number;
}

// 用户相关
export interface User {
  id: number;
  name: string;
  email: string;
  avatar: string;
  role: string;
  status: string;
  lastLogin: string;
  createdAt: string;
}

export interface UserFormData {
  name: string;
  email: string;
  role: string;
  status: boolean;
}

// 角色相关
export interface Role {
  id: number;
  name: string;
  code: string;
  description: string;
  userCount: number;
  permissions: number;
  status: string;
  isSystem: boolean;
  createdAt: string;
}

export interface RoleFormData {
  name: string;
  code: string;
  description: string;
}

export interface Permission {
  id: string;
  name: string;
  checked: boolean;
}

export interface PermissionGroup {
  name: string;
  icon: string;
  permissions: Permission[];
}

// OAuth 相关
export interface OAuthApp {
  id: number;
  name: string;
  clientId: string;
  description: string;
  icon: string;
  iconBg: string;
  redirectUris: string[];
  scopes: string[];
  grantTypes: string[];
  status: string;
  trusted: boolean;
  createdAt: string;
  lastUsed: string;
  requestCount: number;
}

export interface OAuthAppFormData {
  name: string;
  description: string;
  redirectUris: string;
  scopes: string[];
  trusted: boolean;
}

// 审计日志
export interface AuditLog {
  id: string;
  operatorId: string;
  operatorName: string;
  operatorAvatar: string;
  module: string;
  action: string;
  targetId: string;
  targetName?: string;
  detail: Record<string, unknown>;
  ip: string;
  time: string;
  durationMs: number;
  status: Status;
}

// 通知
export type NotificationType = 'system' | 'security' | 'user' | 'oauth' | 'alert';

export interface Notification {
  id: string;
  type: NotificationType;
  title: string;
  message: string;
  time: string;
  read: boolean;
  actionUrl?: string;
  actionLabel?: string;
}

// 组织架构
export interface OrgNodeData {
  id: string;
  name: string;
  displayName?: string;
  avatar?: string;
  orgRole: string;
  class?: string;
  email?: string;
  depth: number;
}

export interface OrgNode {
  key: string;
  type?: string;
  styleClass?: string;
  data: OrgNodeData;
  children?: OrgNode[];
  selectable?: boolean;
  expanded?: boolean;
}

// 搜索
export type SearchCategory = 'navigation' | 'action' | 'user' | 'app' | 'recent';

export interface SearchItem {
  id: string;
  label: string;
  description?: string;
  icon: string;
  category: SearchCategory;
  action: () => void;
  keywords?: string[];
}

// 设置
export interface SettingGroup {
  id: string;
  label: string;
  icon: string;
}

// 动画数字
export interface AnimatedNumberOptions {
  duration?: number;
  easing?: (t: number) => number;
  decimals?: number;
  formatter?: (value: number) => string;
  autoStart?: boolean;
  delay?: number;
  padStart?: boolean;
  useGrouping?: boolean;
}
