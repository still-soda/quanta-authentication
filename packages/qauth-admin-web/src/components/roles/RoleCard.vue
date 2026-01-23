<script setup lang="ts">
import Card from 'primevue/card';
import Button from 'primevue/button';
import Tag from 'primevue/tag';
import type { Role } from '@/types';

defineProps<{ role: Role }>();

const emit = defineEmits<{
   (e: 'edit', role: Role): void;
   (e: 'configPermissions', role: Role): void;
}>();

const getStatusSeverity = (status: string) => (status === 'active' ? 'success' : 'secondary');
const getStatusLabel = (status: string) => (status === 'active' ? '启用' : '禁用');
</script>

<template>
   <Card
      class="rounded-2xl border border-surface-100 dark:border-surface-800 transition-all duration-300 ease overflow-hidden hover:-translate-y-0.5 hover:shadow-[0_12px_24px_-8px_rgba(0,0,0,0.08)] dark:hover:shadow-[0_12px_24px_-8px_rgba(0,0,0,0.3)]"
      :class="{
         'border-primary-200 dark:border-[rgba(251,146,60,0.3)] bg-linear-to-br from-primary-50/50 to-transparent dark:from-[rgba(251,146,60,0.06)] dark:to-transparent':
            role.isSystem,
      }">
      <template #content>
         <!-- Header: Title + Status -->
         <div class="flex items-center justify-between gap-3 mb-2">
            <div class="flex items-center gap-2.5 min-w-0">
               <h3
                  class="text-base font-semibold text-surface-900 dark:text-surface-100 m-0 truncate">
                  {{ role.name }}
               </h3>
               <span
                  v-if="role.isSystem"
                  class="shrink-0 inline-flex items-center gap-1 py-0.5 px-2 bg-primary-100 dark:bg-[rgba(251,146,60,0.12)] text-primary-600 dark:text-primary-400 rounded text-[0.6875rem] font-medium">
                  <i class="pi pi-lock text-[0.625rem]"></i>
                  内置
               </span>
            </div>
            <Tag
               :severity="getStatusSeverity(role.status)"
               class="shrink-0 text-xs">
               {{ getStatusLabel(role.status) }}
            </Tag>
         </div>

         <!-- Code -->
         <code
            class="inline-block text-xs text-surface-400 dark:text-surface-500 font-mono mb-3">
            {{ role.code }}
         </code>

         <!-- Description -->
         <p
            class="text-sm text-surface-600 dark:text-surface-400 m-0 leading-relaxed line-clamp-2 mb-4">
            {{ role.description }}
         </p>

         <!-- Stats Row -->
         <div
            class="flex items-center gap-4 py-3 border-y border-surface-100 dark:border-surface-800">
            <div class="flex items-center gap-2">
               <div
                  class="w-7 h-7 flex items-center justify-center rounded-md bg-surface-100 dark:bg-surface-800">
                  <i
                     class="pi pi-users text-xs text-surface-500 dark:text-surface-400"></i>
               </div>
               <div class="flex flex-col">
                  <span
                     class="text-sm font-semibold text-surface-800 dark:text-surface-200 leading-tight">
                     {{ role.userCount }}
                  </span>
                  <span
                     class="text-[0.6875rem] text-surface-400 dark:text-surface-500 leading-tight">
                     用户
                  </span>
               </div>
            </div>
            <div class="w-px h-8 bg-surface-100 dark:bg-surface-800"></div>
            <div class="flex items-center gap-2">
               <div
                  class="w-7 h-7 flex items-center justify-center rounded-md bg-surface-100 dark:bg-surface-800">
                  <i
                     class="pi pi-key text-xs text-surface-500 dark:text-surface-400"></i>
               </div>
               <div class="flex flex-col">
                  <span
                     class="text-sm font-semibold text-surface-800 dark:text-surface-200 leading-tight">
                     {{ role.permissions }}
                  </span>
                  <span
                     class="text-[0.6875rem] text-surface-400 dark:text-surface-500 leading-tight">
                     权限
                  </span>
               </div>
            </div>
            <span
               class="ml-auto text-[0.6875rem] text-surface-400 dark:text-surface-500">
               {{ role.createdAt }}
            </span>
         </div>

         <!-- Actions -->
         <div class="flex items-center gap-2 mt-3">
            <Button
               icon="pi pi-key"
               label="配置权限"
               severity="secondary"
               outlined
               size="small"
               class="flex-1"
               @click="emit('configPermissions', role)" />
            <Button
               icon="pi pi-pencil"
               severity="secondary"
               text
               size="small"
               :disabled="role.isSystem"
               v-tooltip.top="role.isSystem ? '系统内置角色不可编辑' : '编辑角色'"
               @click="emit('edit', role)" />
         </div>
      </template>
   </Card>
</template>
