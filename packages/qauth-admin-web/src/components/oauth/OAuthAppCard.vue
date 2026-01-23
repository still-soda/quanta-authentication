<script setup lang="ts">
import Card from 'primevue/card';
import Button from 'primevue/button';
import Tag from 'primevue/tag';
import Chip from 'primevue/chip';

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

defineProps<{
   app: OAuthApp;
}>();

const emit = defineEmits<{
   (e: 'view', app: OAuthApp): void;
   (e: 'edit', app: OAuthApp): void;
   (e: 'regenerateSecret', app: OAuthApp): void;
}>();

const getStatusSeverity = (status: string) => {
   const map: Record<string, 'success' | 'warn' | 'danger' | 'secondary'> = {
      active: 'success',
      development: 'warn',
      deprecated: 'secondary',
   };
   return map[status] || 'info';
};

const getStatusLabel = (status: string) => {
   const map: Record<string, string> = {
      active: '生产环境',
      development: '开发中',
      deprecated: '已弃用',
   };
   return map[status] || status;
};

const formatNumber = (num: number) => {
   if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
   if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
   return num.toString();
};

const copyToClipboard = (text: string) => {
   navigator.clipboard.writeText(text);
};
</script>

<template>
   <Card class="app-card">
      <template #content>
         <div class="app-header">
            <div class="app-icon" :style="{ background: app.iconBg }">
               <i :class="app.icon"></i>
            </div>
            <div class="app-badges">
               <Tag
                  v-if="app.trusted"
                  severity="info"
                  class="trusted-badge"
                  rounded>
                  <i class="pi pi-verified"></i>
                  可信
               </Tag>
               <Tag :severity="getStatusSeverity(app.status)">
                  {{ getStatusLabel(app.status) }}
               </Tag>
            </div>
         </div>

         <div class="app-body">
            <h3 class="app-name">{{ app.name }}</h3>
            <p class="app-description">{{ app.description }}</p>

            <div class="app-client-id">
               <label>Client ID</label>
               <div class="client-id-value">
                  <code>{{ app.clientId }}</code>
                  <Button
                     icon="pi pi-copy"
                     text
                     rounded
                     severity="secondary"
                     size="small"
                     @click="copyToClipboard(app.clientId)"
                     v-tooltip.top="'复制'" />
               </div>
            </div>

            <div class="app-scopes">
               <label>授权范围</label>
               <div class="scope-tags">
                  <Chip
                     v-for="scope in app.scopes"
                     :key="scope"
                     :label="scope" />
               </div>
            </div>

            <div class="app-stats">
               <div class="stat-item">
                  <i class="pi pi-chart-line"></i>
                  <span>{{ formatNumber(app.requestCount) }} 请求</span>
               </div>
               <div class="stat-item">
                  <i class="pi pi-clock"></i>
                  <span>{{ app.lastUsed }}</span>
               </div>
            </div>
         </div>

         <div class="app-actions">
            <Button
               icon="pi pi-eye"
               label="查看"
               severity="secondary"
               outlined
               size="small"
               @click="emit('view', app)" />
            <Button
               icon="pi pi-refresh"
               label="重置密钥"
               severity="secondary"
               text
               size="small"
               @click="emit('regenerateSecret', app)" />
            <Button
               icon="pi pi-pencil"
               text
               rounded
               severity="secondary"
               @click="emit('edit', app)" />
         </div>
      </template>
   </Card>
</template>

<style scoped>
.app-card {
   border-radius: 16px;
   border: 1px solid var(--p-surface-100);
   transition: all 0.3s ease;
}

.app-card:hover {
   transform: translateY(-2px);
   box-shadow: 0 12px 24px -8px rgba(0, 0, 0, 0.08);
}

:global(.app-dark) .app-card {
   border-color: var(--p-surface-800);
}

:global(.app-dark) .app-card:hover {
   box-shadow: 0 12px 24px -8px rgba(0, 0, 0, 0.3);
}

.app-header {
   display: flex;
   justify-content: space-between;
   align-items: flex-start;
   margin-bottom: 1rem;
}

.app-icon {
   width: 3rem;
   height: 3rem;
   display: flex;
   align-items: center;
   justify-content: center;
   border-radius: 12px;
   color: white;
   font-size: 1.25rem;
   box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.app-badges {
   display: flex;
   gap: 0.5rem;
}

.trusted-badge {
   display: flex;
   align-items: center;
   gap: 0.25rem;
}

.app-body {
   display: flex;
   flex-direction: column;
   gap: 1rem;
}

.app-name {
   font-size: 1.125rem;
   font-weight: 600;
   color: var(--p-surface-900);
   margin: 0;
}

:global(.app-dark) .app-name {
   color: var(--p-surface-100);
}

.app-description {
   font-size: 0.875rem;
   color: var(--p-surface-500);
   margin: 0;
   line-height: 1.5;
}

.app-client-id {
   display: flex;
   flex-direction: column;
   gap: 0.375rem;
}

.app-client-id label,
.app-scopes label {
   font-size: 0.75rem;
   font-weight: 500;
   color: var(--p-surface-500);
   text-transform: uppercase;
   letter-spacing: 0.05em;
}

.client-id-value {
   display: flex;
   align-items: center;
   gap: 0.5rem;
   padding: 0.5rem 0.75rem;
   background: var(--p-surface-50);
   border-radius: 8px;
}

:global(.app-dark) .client-id-value {
   background: var(--p-surface-800);
}

.client-id-value code {
   flex: 1;
   font-size: 0.8125rem;
   font-family: monospace;
   color: var(--p-surface-700);
}

:global(.app-dark) .client-id-value code {
   color: var(--p-surface-300);
}

.app-scopes {
   display: flex;
   flex-direction: column;
   gap: 0.5rem;
}

.scope-tags {
   display: flex;
   flex-wrap: wrap;
   gap: 0.375rem;
}

.scope-tags :deep(.p-chip) {
   font-size: 0.75rem;
   padding: 0.25rem 0.5rem;
}

.app-stats {
   display: flex;
   gap: 1.25rem;
}

.stat-item {
   display: flex;
   align-items: center;
   gap: 0.5rem;
   font-size: 0.8125rem;
   color: var(--p-surface-500);
}

.stat-item i {
   font-size: 0.875rem;
}

.app-actions {
   display: flex;
   gap: 0.5rem;
   margin-top: 1rem;
   padding-top: 1rem;
   border-top: 1px solid var(--p-surface-100);
}

:global(.app-dark) .app-actions {
   border-top-color: var(--p-surface-800);
}
</style>
