<script setup lang="ts">
import { ref, computed } from 'vue';
import Button from 'primevue/button';
import OrganizationChart from 'primevue/organizationchart';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Select from 'primevue/select';
import PageHeader from '@/components/shared/PageHeader.vue';
import SimpleStatCard from '@/components/shared/SimpleStatCard.vue';
import type { OrgNode, SimpleStatData } from '@/types';
import { ORG_CLASS_OPTIONS, ORG_CLASS_COLORS } from '@/config';

// 模拟组织架构数据
const orgData = ref<OrgNode>({
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

// 统计人数的辅助函数
const countNodes = (node: OrgNode): number => {
   let count = 1;
   if (node.children) {
      node.children.forEach((child) => {
         count += countNodes(child);
      });
   }
   return count;
};

// 统计特定级别的人数
const countByClass = (node: OrgNode, targetClass: string): number => {
   let count = node.data.class === targetClass ? 1 : 0;
   if (node.children) {
      node.children.forEach((child) => {
         count += countByClass(child, targetClass);
      });
   }
   return count;
};

// 统计数据
const stats = computed<SimpleStatData[]>(() => [
   {
      title: '总人数',
      value: countNodes(orgData.value),
      icon: 'pi pi-users',
      color: 'blue',
   },
   {
      title: '高管层',
      value: countByClass(orgData.value, '高管层'),
      icon: 'pi pi-star',
      color: 'orange',
   },
   {
      title: '管理层',
      value: countByClass(orgData.value, '管理层'),
      icon: 'pi pi-briefcase',
      color: 'purple',
   },
   {
      title: '员工',
      value: countByClass(orgData.value, '员工'),
      icon: 'pi pi-user',
      color: 'green',
   },
]);

// 选中的节点
const selectedNode = ref<OrgNode | null>(null);

// 对话框状态
const memberDialog = ref(false);
const isEditing = ref(false);
const memberForm = ref({
   name: '',
   orgRole: '',
   class: '',
   email: '',
});

const classOptions = ORG_CLASS_OPTIONS;

// 节点选择处理
const onNodeSelect = (node: OrgNode) => {
   selectedNode.value = node;
};

// 打开新增成员对话框
const openAddMemberDialog = () => {
   isEditing.value = false;
   memberForm.value = {
      name: '',
      orgRole: '',
      class: '员工',
      email: '',
   };
   memberDialog.value = true;
};

// 编辑当前选中节点
const editSelectedMember = () => {
   if (!selectedNode.value) return;
   isEditing.value = true;
   memberForm.value = {
      name: selectedNode.value.data.name,
      orgRole: selectedNode.value.data.orgRole,
      class: selectedNode.value.data.class || '员工',
      email: selectedNode.value.data.email || '',
   };
   memberDialog.value = true;
};

// 保存成员
const saveMember = () => {
   console.log('Saving member:', memberForm.value);
   memberDialog.value = false;
};

const getClassColor = (classType?: string) => ORG_CLASS_COLORS[classType || ''] || 'bg-surface-100 text-surface-700 dark:bg-surface-800 dark:text-surface-400';
</script>

<template>
   <div class="flex flex-col gap-6">
      <!-- Page Header -->
      <PageHeader title="组织架构" subtitle="可视化管理企业组织结构">
         <template #actions>
            <Button
               v-if="selectedNode"
               label="编辑成员"
               icon="pi pi-pencil"
               severity="secondary"
               outlined
               @click="editSelectedMember" />
            <Button
               label="添加成员"
               icon="pi pi-plus"
               @click="openAddMemberDialog" />
         </template>
      </PageHeader>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4">
         <SimpleStatCard v-for="stat in stats" :key="stat.title" :stat="stat" />
      </div>

      <!-- Selected Member Info -->
      <Transition name="slide-fade">
         <div
            v-if="selectedNode"
            class="flex items-center gap-4 p-4 bg-surface-0 dark:bg-surface-900 border border-primary-200 dark:border-primary-800 rounded-xl">
            <div
               class="w-14 h-14 rounded-xl overflow-hidden bg-linear-to-br from-primary-100 to-primary-200 dark:from-primary-900/50 dark:to-primary-800/50 shadow-md">
               <img
                  :src="selectedNode.data.avatar"
                  :alt="selectedNode.data.name"
                  class="w-full h-full object-cover" />
            </div>
            <div class="flex-1">
               <div class="flex items-center gap-3">
                  <h3
                     class="text-lg font-bold text-surface-900 dark:text-surface-100 m-0">
                     {{ selectedNode.data.name }}
                  </h3>
                  <span
                     :class="getClassColor(selectedNode.data.class)"
                     class="text-xs font-medium px-2 py-0.5 rounded-full">
                     {{ selectedNode.data.class }}
                  </span>
               </div>
               <p class="text-sm text-surface-500 m-0 mt-0.5">
                  {{ selectedNode.data.orgRole }}
                  <span v-if="selectedNode.data.email" class="ml-2">
                     · {{ selectedNode.data.email }}
                  </span>
               </p>
            </div>
            <Button
               icon="pi pi-times"
               severity="secondary"
               text
               rounded
               @click="selectedNode = null" />
         </div>
      </Transition>

      <!-- Organization Chart -->
      <div
         class="bg-surface-0 dark:bg-surface-900 border border-surface-200 dark:border-surface-800 rounded-xl p-6 overflow-auto">
         <OrganizationChart
            :value="orgData"
            selectionMode="single"
            v-model:selectionKeys="selectedNode!"
            @node-select="onNodeSelect!"
            :pt="{
               root: { class: 'org-chart-root' },
               table: { class: 'org-chart-table' },
               node: { class: 'org-chart-node' },
               nodeToggleButton: { class: 'org-chart-toggle' },
            }">
            <template #default="slotProps">
               <div
                  class="org-node group cursor-pointer transition-all duration-200"
                  :class="{
                     'ring-2 ring-primary-500 ring-offset-2 dark:ring-offset-surface-900':
                        selectedNode?.key === slotProps.node.key,
                  }">
                  <!-- Avatar -->
                  <div
                     class="w-16 h-16 mx-auto mb-3 rounded-xl overflow-hidden bg-linear-to-br from-surface-100 to-surface-200 dark:from-surface-700 dark:to-surface-800 shadow-lg transition-transform duration-200 group-hover:scale-105">
                     <img
                        :src="slotProps.node.data.avatar"
                        :alt="slotProps.node.data.name"
                        class="w-full h-full object-cover" />
                  </div>
                  <!-- Name -->
                  <div
                     class="text-sm font-bold text-surface-900 dark:text-surface-100 mb-0.5">
                     {{ slotProps.node.data.name }}
                  </div>
                  <!-- Role -->
                  <div class="text-xs text-surface-500 mb-2">
                     {{ slotProps.node.data.orgRole }}
                  </div>
                  <!-- Class Badge -->
                  <span
                     :class="getClassColor(slotProps.node.data.class)"
                     class="text-[10px] font-medium px-2 py-0.5 rounded-full">
                     {{ slotProps.node.data.class }}
                  </span>
               </div>
            </template>
         </OrganizationChart>
      </div>

      <!-- Member Dialog -->
      <Dialog
         v-model:visible="memberDialog"
         :header="isEditing ? '编辑成员' : '添加成员'"
         modal
         :style="{ width: '480px' }"
         :pt="{
            root: {
               class: 'border-none shadow-2xl rounded-2xl',
            },
            header: {
               class: 'border-b border-surface-200 dark:border-surface-700 px-6 py-4',
            },
            content: { class: 'p-6' },
         }">
         <div class="flex flex-col gap-5">
            <div class="flex flex-col gap-2">
               <label class="text-sm font-medium text-surface-700 dark:text-surface-300">
                  姓名 <span class="text-red-500">*</span>
               </label>
               <InputText
                  v-model="memberForm.name"
                  placeholder="请输入姓名"
                  class="w-full" />
            </div>
            <div class="flex flex-col gap-2">
               <label class="text-sm font-medium text-surface-700 dark:text-surface-300">
                  职位 <span class="text-red-500">*</span>
               </label>
               <InputText
                  v-model="memberForm.orgRole"
                  placeholder="请输入职位"
                  class="w-full" />
            </div>
            <div class="flex flex-col gap-2">
               <label class="text-sm font-medium text-surface-700 dark:text-surface-300">
                  级别 <span class="text-red-500">*</span>
               </label>
               <Select
                  v-model="memberForm.class"
                  :options="classOptions"
                  optionLabel="label"
                  optionValue="value"
                  placeholder="选择级别"
                  class="w-full" />
            </div>
            <div class="flex flex-col gap-2">
               <label class="text-sm font-medium text-surface-700 dark:text-surface-300">
                  邮箱
               </label>
               <InputText
                  v-model="memberForm.email"
                  type="email"
                  placeholder="请输入邮箱"
                  class="w-full" />
            </div>
         </div>
         <template #footer>
            <div class="flex justify-end gap-3">
               <Button
                  label="取消"
                  severity="secondary"
                  outlined
                  @click="memberDialog = false" />
               <Button :label="isEditing ? '保存' : '添加'" @click="saveMember" />
            </div>
         </template>
      </Dialog>
   </div>
</template>

<style scoped>
/* Organization Chart 样式定制 */
:deep(.org-chart-root) {
   padding: 0;
}

:deep(.org-chart-table) {
   border-collapse: separate;
   border-spacing: 0 16px;
}

/* 使用 PrimeVue 默认节点样式 */
:deep(.p-organizationchart-node-content) {
   padding: 1rem;
   text-align: center;
}

:deep(.p-organizationchart-connectors) {
   color: var(--p-surface-300);
}

:deep(.p-organizationchart-line-down),
:deep(.p-organizationchart-line-left),
:deep(.p-organizationchart-line-right) {
   border-color: var(--p-surface-300);
}

.app-dark :deep(.p-organizationchart-connectors) {
   color: var(--p-surface-600);
}

.app-dark :deep(.p-organizationchart-line-down),
.app-dark :deep(.p-organizationchart-line-left),
.app-dark :deep(.p-organizationchart-line-right) {
   border-color: var(--p-surface-600);
}

.org-node {
   text-align: center;
}

/* Transition */
.slide-fade-enter-active,
.slide-fade-leave-active {
   transition: all 0.3s ease;
}

.slide-fade-enter-from,
.slide-fade-leave-to {
   opacity: 0;
   transform: translateY(-10px);
}
</style>
