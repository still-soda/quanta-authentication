/**
 * 组织架构相关类型定义
 */

// 组织节点数据
export interface OrgNodeData {
   id: string
   name: string
   displayName?: string
   avatar?: string
   orgRole: string
   class?: string
   email?: string
   depth: number
}

// 组织节点
export interface OrgNode {
   key: string
   type?: string
   styleClass?: string
   data: OrgNodeData
   children?: OrgNode[]
   selectable?: boolean
   expanded?: boolean
}

// 组织成员表单数据
export interface OrgMemberFormData {
   name: string
   orgRole: string
   class: string
   email?: string
   parentId?: string
}
