/**
 * Organizations API - 组织架构相关接口
 */
import { mockResponse } from './index';
import type { OrgNode } from '@/types';

export interface OrgMemberFormData {
  name: string;
  orgRole: string;
  class: string;
  email?: string;
  parentId?: string;
}

/**
 * 获取组织架构树
 */
export async function getOrganizationTree(): Promise<OrgNode> {
  return mockResponse({
    key: '0',
    type: 'person',
    data: {
      id: '1',
      name: '张伟',
      displayName: '张伟',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhang',
      orgRole: 'CEO',
      class: '高管层',
      email: 'zhang.wei@company.com',
      depth: 0,
    },
    expanded: true,
    children: [
      {
        key: '0-0',
        type: 'person',
        data: {
          id: '2',
          name: '李明',
          displayName: '李明',
          avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=li',
          orgRole: 'CTO',
          class: '高管层',
          email: 'li.ming@company.com',
          depth: 1,
        },
        expanded: true,
        children: [
          {
            key: '0-0-0',
            type: 'person',
            data: {
              id: '5',
              name: '王强',
              displayName: '王强',
              avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=wangq',
              orgRole: '技术总监',
              class: '管理层',
              email: 'wang.qiang@company.com',
              depth: 2,
            },
            children: [
              {
                key: '0-0-0-0',
                type: 'person',
                data: {
                  id: '9',
                  name: '刘洋',
                  displayName: '刘洋',
                  avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=liu',
                  orgRole: '前端工程师',
                  class: '员工',
                  email: 'liu.yang@company.com',
                  depth: 3,
                },
              },
              {
                key: '0-0-0-1',
                type: 'person',
                data: {
                  id: '10',
                  name: '陈浩',
                  displayName: '陈浩',
                  avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=chenh',
                  orgRole: '后端工程师',
                  class: '员工',
                  email: 'chen.hao@company.com',
                  depth: 3,
                },
              },
            ],
          },
          {
            key: '0-0-1',
            type: 'person',
            data: {
              id: '6',
              name: '赵娜',
              displayName: '赵娜',
              avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhao',
              orgRole: '产品总监',
              class: '管理层',
              email: 'zhao.na@company.com',
              depth: 2,
            },
            children: [
              {
                key: '0-0-1-0',
                type: 'person',
                data: {
                  id: '11',
                  name: '孙静',
                  displayName: '孙静',
                  avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=sun',
                  orgRole: '产品经理',
                  class: '员工',
                  email: 'sun.jing@company.com',
                  depth: 3,
                },
              },
            ],
          },
        ],
      },
      {
        key: '0-1',
        type: 'person',
        data: {
          id: '3',
          name: '王芳',
          displayName: '王芳',
          avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=wang',
          orgRole: 'CFO',
          class: '高管层',
          email: 'wang.fang@company.com',
          depth: 1,
        },
        expanded: true,
        children: [
          {
            key: '0-1-0',
            type: 'person',
            data: {
              id: '7',
              name: '周杰',
              displayName: '周杰',
              avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhou',
              orgRole: '财务经理',
              class: '管理层',
              email: 'zhou.jie@company.com',
              depth: 2,
            },
            children: [
              {
                key: '0-1-0-0',
                type: 'person',
                data: {
                  id: '12',
                  name: '吴敏',
                  displayName: '吴敏',
                  avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=wu',
                  orgRole: '会计',
                  class: '员工',
                  email: 'wu.min@company.com',
                  depth: 3,
                },
              },
            ],
          },
        ],
      },
      {
        key: '0-2',
        type: 'person',
        data: {
          id: '4',
          name: '陈红',
          displayName: '陈红',
          avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=chen',
          orgRole: 'COO',
          class: '高管层',
          email: 'chen.hong@company.com',
          depth: 1,
        },
        expanded: true,
        children: [
          {
            key: '0-2-0',
            type: 'person',
            data: {
              id: '8',
              name: '郑磊',
              displayName: '郑磊',
              avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zheng',
              orgRole: '运营总监',
              class: '管理层',
              email: 'zheng.lei@company.com',
              depth: 2,
            },
            children: [
              {
                key: '0-2-0-0',
                type: 'person',
                data: {
                  id: '13',
                  name: '林涛',
                  displayName: '林涛',
                  avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=lin',
                  orgRole: '运营专员',
                  class: '员工',
                  email: 'lin.tao@company.com',
                  depth: 3,
                },
              },
              {
                key: '0-2-0-1',
                type: 'person',
                data: {
                  id: '14',
                  name: '黄梅',
                  displayName: '黄梅',
                  avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=huang',
                  orgRole: '客服专员',
                  class: '员工',
                  email: 'huang.mei@company.com',
                  depth: 3,
                },
              },
            ],
          },
        ],
      },
    ],
  });
}

/**
 * 添加组织成员
 */
export async function addOrgMember(data: OrgMemberFormData): Promise<OrgNode> {
  const id = Date.now().toString();
  return mockResponse({
    key: id,
    type: 'person',
    data: {
      id,
      name: data.name,
      displayName: data.name,
      avatar: `https://api.dicebear.com/7.x/avataaars/svg?seed=${id}`,
      orgRole: data.orgRole,
      class: data.class,
      email: data.email,
      depth: 0,
    },
  });
}

/**
 * 更新组织成员
 */
export async function updateOrgMember(id: string, data: Partial<OrgMemberFormData>): Promise<OrgNode> {
  return mockResponse({
    key: id,
    type: 'person',
    data: {
      id,
      name: data.name || '',
      displayName: data.name || '',
      avatar: `https://api.dicebear.com/7.x/avataaars/svg?seed=${id}`,
      orgRole: data.orgRole || '',
      class: data.class || '',
      email: data.email,
      depth: 0,
    },
  });
}

/**
 * 删除组织成员
 */
export async function deleteOrgMember(id: string): Promise<void> {
  return mockResponse(undefined);
}

/**
 * 移动组织成员
 */
export async function moveOrgMember(id: string, newParentId: string): Promise<void> {
  return mockResponse(undefined);
}
