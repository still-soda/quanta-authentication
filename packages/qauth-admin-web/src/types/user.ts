/**
 * 用户相关类型定义
 */

// 用户信息
export interface User {
   id: number
   name: string
   email: string
   avatar: string
   role: string
   status: string
   lastLogin: string
   createdAt: string
}

// 用户表单数据
export interface UserFormData {
   name: string
   email: string
   role: string
   status: boolean
}
