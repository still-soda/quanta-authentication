<script setup lang="ts">
import Card from 'primevue/card';
import Button from 'primevue/button';
import Tag from 'primevue/tag';

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

defineProps<{
   role: Role;
}>();

const emit = defineEmits<{
   (e: 'edit', role: Role): void;
   (e: 'configPermissions', role: Role): void;
}>();

const getStatusSeverity = (status: string) => {
   return status === 'active' ? 'success' : 'secondary';
};

const getStatusLabel = (status: string) => {
   return status === 'active' ? '启用' : '禁用';
};
</script>

<template>
   <Card
      class="rounded-2xl border border-surface-100 dark:border-surface-800 transition-all duration-300 ease overflow-hidden hover:-translate-y-0.5 hover:shadow-[0_12px_24px_-8px_rgba(0,0,0,0.08)] dark:hover:shadow-[0_12px_24px_-8px_rgba(0,0,0,0.3)]"
      :class="{
         'border-primary-200 dark:border-[rgba(251,146,60,0.3)] bg-linear-to-br from-primary-50 to-transparent dark:from-[rgba(251,146,60,0.08)] dark:to-transparent':
            role.isSystem,
      }">
      <template #content>
         <div class="flex justify-between items-start mb-4">
            <div
               class="w-12 h-12 flex items-center justify-center rounded-xl text-xl"
               :class="
                  role.isSystem
                     ? 'bg-linear-to-br from-primary-400 to-primary-600 text-white shadow-[0_4px_12px_rgba(249,115,22,0.3)]'
                     : 'bg-surface-100 dark:bg-surface-800 text-surface-600 dark:text-surface-400'
               ">
               <i class="pi pi-shield"></i>
            </div>
            <Tag :severity="getStatusSeverity(role.status)">
               {{ getStatusLabel(role.status) }}
            </Tag>
         </div>

         <div class="flex flex-col gap-3">
            <h3
               class="text-lg font-semibold text-surface-900 dark:text-surface-100 m-0">
               {{ role.name }}
            </h3>
            <p class="text-sm text-surface-500 m-0 leading-relaxed">
               {{ role.description }}
            </p>

            <div class="flex gap-5">
               <div
                  class="flex items-center gap-2 text-[0.8125rem] text-surface-600 dark:text-surface-400">
                  <i class="pi pi-users text-sm"></i>
                  <span>{{ role.userCount }} 用户</span>
               </div>
               <div
                  class="flex items-center gap-2 text-[0.8125rem] text-surface-600 dark:text-surface-400">
                  <i class="pi pi-key text-sm"></i>
                  <span>{{ role.permissions }} 权限</span>
               </div>
            </div>

            <div
               class="flex items-center gap-4 pt-3 border-t border-surface-100 dark:border-surface-800">
               <span
                  v-if="role.isSystem"
                  class="inline-flex items-center gap-1.5 py-1 px-2.5 bg-primary-100 dark:bg-[rgba(251,146,60,0.15)] text-primary-700 dark:text-primary-400 rounded-full text-xs font-medium">
                  <i class="pi pi-lock"></i>
                  系统内置
               </span>
               <span class="text-xs text-surface-400">
                  创建于 {{ role.createdAt }}
               </span>
            </div>
         </div>

         <div
            class="flex gap-2 mt-4 pt-4 border-t border-surface-100 dark:border-surface-800">
            <Button
               icon="pi pi-key"
               label="权限"
               severity="secondary"
               outlined
               size="small"
               @click="emit('configPermissions', role)" />
            <Button
               icon="pi pi-pencil"
               label="编辑"
               severity="secondary"
               text
               size="small"
               :disabled="role.isSystem"
               @click="emit('edit', role)" />
         </div>
      </template>
   </Card>
</template>
