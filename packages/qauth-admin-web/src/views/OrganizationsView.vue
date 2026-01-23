<script setup lang="ts">
import { ref, computed } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import Button from 'primevue/button';
import OrganizationChart from 'primevue/organizationchart';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Select from 'primevue/select';
import PageHeader from '@/components/shared/PageHeader.vue';
import SimpleStatCard from '@/components/shared/SimpleStatCard.vue';
import type { OrgNode, SimpleStatData } from '@/types';
import { ORG_CLASS_OPTIONS, ORG_CLASS_COLORS } from '@/config';
import {
   getOrganizationTree,
   addOrgMember,
   updateOrgMember,
} from '@/apis/organizations';

const queryClient = useQueryClient();

// 使用 TanStack Query 获取组织架构数据
const { data: orgData, isLoading } = useQuery({
   queryKey: ['organization'],
   queryFn: getOrganizationTree,
});

// 添加成员 mutation
const addMemberMutation = useMutation({
   mutationFn: addOrgMember,
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['organization'] });
      memberDialog.value = false;
   },
});

// 更新成员 mutation
const updateMemberMutation = useMutation({
   mutationFn: ({ id, data }: { id: string; data: any }) => updateOrgMember(id, data),
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['organization'] });
      memberDialog.value = false;
   },
});

// 统计人数的辅助函数
const countNodes = (node: OrgNode | undefined): number => {
   if (!node) return 0;
   let count = 1;
   if (node.children) {
      node.children.forEach((child) => {
         count += countNodes(child);
      });
   }
   return count;
};

// 统计特定级别的人数
const countByClass = (node: OrgNode | undefined, targetClass: string): number => {
   if (!node) return 0;
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
   if (isEditing.value && selectedNode.value) {
      updateMemberMutation.mutate({
         id: selectedNode.value.data.id,
         data: memberForm.value,
      });
   } else {
      addMemberMutation.mutate(memberForm.value);
   }
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
         <template v-if="isLoading">
            <div
               v-for="i in 4"
               :key="i"
               class="h-20 bg-surface-100 dark:bg-surface-800 rounded-xl animate-pulse" />
         </template>
         <template v-else>
            <SimpleStatCard v-for="stat in stats" :key="stat.title" :stat="stat" />
         </template>
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
         <div
            v-if="isLoading"
            class="flex items-center justify-center py-24">
            <i class="pi pi-spin pi-spinner text-3xl text-surface-400"></i>
         </div>
         <OrganizationChart
            v-else-if="orgData"
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
