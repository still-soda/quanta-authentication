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
   <Card class="role-card" :class="{ 'system-role': role.isSystem }">
      <template #content>
         <div class="role-header">
            <div class="role-icon" :class="{ system: role.isSystem }">
               <i class="pi pi-shield"></i>
            </div>
            <Tag :severity="getStatusSeverity(role.status)" class="role-status">
               {{ getStatusLabel(role.status) }}
            </Tag>
         </div>

         <div class="role-body">
            <h3 class="role-name">{{ role.name }}</h3>
            <p class="role-description">{{ role.description }}</p>

            <div class="role-stats">
               <div class="stat-item">
                  <i class="pi pi-users"></i>
                  <span>{{ role.userCount }} 用户</span>
               </div>
               <div class="stat-item">
                  <i class="pi pi-key"></i>
                  <span>{{ role.permissions }} 权限</span>
               </div>
            </div>

            <div class="role-meta">
               <span v-if="role.isSystem" class="system-badge">
                  <i class="pi pi-lock"></i>
                  系统内置
               </span>
               <span class="created-at">创建于 {{ role.createdAt }}</span>
            </div>
         </div>

         <div class="role-actions">
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

<style scoped>
.role-card {
   border-radius: 16px;
   border: 1px solid var(--p-surface-100);
   transition: all 0.3s ease;
   overflow: hidden;
}

.role-card:hover {
   transform: translateY(-2px);
   box-shadow: 0 12px 24px -8px rgba(0, 0, 0, 0.08);
}

:global(.app-dark) .role-card {
   border-color: var(--p-surface-800);
}

:global(.app-dark) .role-card:hover {
   box-shadow: 0 12px 24px -8px rgba(0, 0, 0, 0.3);
}

.role-card.system-role {
   border-color: var(--p-orange-200);
   background: linear-gradient(135deg, var(--p-orange-50) 0%, transparent 100%);
}

:global(.app-dark) .role-card.system-role {
   border-color: rgba(251, 146, 60, 0.3);
   background: linear-gradient(
      135deg,
      rgba(251, 146, 60, 0.08) 0%,
      transparent 100%
   );
}

.role-header {
   display: flex;
   justify-content: space-between;
   align-items: flex-start;
   margin-bottom: 1rem;
}

.role-icon {
   width: 3rem;
   height: 3rem;
   display: flex;
   align-items: center;
   justify-content: center;
   background: var(--p-surface-100);
   color: var(--p-surface-600);
   border-radius: 12px;
   font-size: 1.25rem;
}

:global(.app-dark) .role-icon {
   background: var(--p-surface-800);
   color: var(--p-surface-400);
}

.role-icon.system {
   background: linear-gradient(
      135deg,
      var(--p-orange-400) 0%,
      var(--p-orange-600) 100%
   );
   color: white;
   box-shadow: 0 4px 12px rgba(249, 115, 22, 0.3);
}

.role-body {
   display: flex;
   flex-direction: column;
   gap: 0.75rem;
}

.role-name {
   font-size: 1.125rem;
   font-weight: 600;
   color: var(--p-surface-900);
   margin: 0;
}

:global(.app-dark) .role-name {
   color: var(--p-surface-100);
}

.role-description {
   font-size: 0.875rem;
   color: var(--p-surface-500);
   margin: 0;
   line-height: 1.5;
}

.role-stats {
   display: flex;
   gap: 1.25rem;
}

.stat-item {
   display: flex;
   align-items: center;
   gap: 0.5rem;
   font-size: 0.8125rem;
   color: var(--p-surface-600);
}

:global(.app-dark) .stat-item {
   color: var(--p-surface-400);
}

.stat-item i {
   font-size: 0.875rem;
}

.role-meta {
   display: flex;
   align-items: center;
   gap: 1rem;
   padding-top: 0.75rem;
   border-top: 1px solid var(--p-surface-100);
}

:global(.app-dark) .role-meta {
   border-top-color: var(--p-surface-800);
}

.system-badge {
   display: inline-flex;
   align-items: center;
   gap: 0.375rem;
   padding: 0.25rem 0.625rem;
   background: var(--p-orange-100);
   color: var(--p-orange-700);
   border-radius: 9999px;
   font-size: 0.75rem;
   font-weight: 500;
}

:global(.app-dark) .system-badge {
   background: rgba(251, 146, 60, 0.15);
   color: var(--p-orange-400);
}

.created-at {
   font-size: 0.75rem;
   color: var(--p-surface-400);
}

.role-actions {
   display: flex;
   gap: 0.5rem;
   margin-top: 1rem;
   padding-top: 1rem;
   border-top: 1px solid var(--p-surface-100);
}

:global(.app-dark) .role-actions {
   border-top-color: var(--p-surface-800);
}
</style>
