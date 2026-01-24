/**
 * 通用类型定义
 */

// 主题色类型
export type ThemeColor = 'orange' | 'blue' | 'green' | 'purple' | 'cyan' | 'red'
export type ExtendedThemeColor = ThemeColor | 'gray'

// 通用状态类型
export type Status = 'success' | 'warning' | 'error'

// 搜索
export type SearchCategory = 'navigation' | 'action' | 'user' | 'app' | 'recent'

export interface SearchItem {
   id: string
   label: string
   description?: string
   icon: string
   category: SearchCategory
   action: () => void
   keywords?: string[]
}

export interface SearchGroup {
   category: SearchCategory
   label: string
   items: SearchItem[]
}

// 动画数字
export interface AnimatedNumberOptions {
   duration?: number
   easing?: (t: number) => number
   decimals?: number
   formatter?: (value: number) => string
   autoStart?: boolean
   delay?: number
   padStart?: boolean
   useGrouping?: boolean
}
