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
   <Card
      class="rounded-2xl border border-surface-100 dark:border-surface-800 transition-all duration-300 ease hover:-translate-y-0.5 hover:shadow-[0_12px_24px_-8px_rgba(0,0,0,0.08)] dark:hover:shadow-[0_12px_24px_-8px_rgba(0,0,0,0.3)]">
      <template #content>
         <div class="flex justify-between items-start mb-4">
            <div
               class="w-12 h-12 flex items-center justify-center rounded-xl text-white text-xl shadow-[0_4px_12px_rgba(0,0,0,0.15)]"
               :style="{ background: app.iconBg }">
               <i :class="app.icon"></i>
            </div>
            <div class="flex gap-2">
               <Tag
                  v-if="app.trusted"
                  severity="info"
                  class="flex items-center gap-1"
                  rounded>
                  <i class="pi pi-verified"></i>
                  可信
               </Tag>
               <Tag :severity="getStatusSeverity(app.status)">
                  {{ getStatusLabel(app.status) }}
               </Tag>
            </div>
         </div>

         <div class="flex flex-col gap-4">
            <h3
               class="text-lg font-semibold text-surface-900 dark:text-surface-100 m-0">
               {{ app.name }}
            </h3>
            <p class="text-sm text-surface-500 m-0 leading-relaxed">
               {{ app.description }}
            </p>

            <div class="flex flex-col gap-1.5">
               <label
                  class="text-xs font-medium text-surface-500 uppercase tracking-wider">
                  Client ID
               </label>
               <div
                  class="flex items-center gap-2 py-2 px-3 bg-surface-50 dark:bg-surface-800 rounded-lg">
                  <code
                     class="flex-1 text-[0.8125rem] font-mono text-surface-700 dark:text-surface-300">
                     {{ app.clientId }}
                  </code>
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

            <div class="flex flex-col gap-2">
               <label
                  class="text-xs font-medium text-surface-500 uppercase tracking-wider">
                  授权范围
               </label>
               <div class="flex flex-wrap gap-1.5">
                  <Chip
                     v-for="scope in app.scopes"
                     :key="scope"
                     :label="scope"
                     class="text-xs py-1 px-2" />
               </div>
            </div>

            <div class="flex gap-5">
               <div
                  class="flex items-center gap-2 text-[0.8125rem] text-surface-500">
                  <i class="pi pi-chart-line text-sm"></i>
                  <span>{{ formatNumber(app.requestCount) }} 请求</span>
               </div>
               <div
                  class="flex items-center gap-2 text-[0.8125rem] text-surface-500">
                  <i class="pi pi-clock text-sm"></i>
                  <span>{{ app.lastUsed }}</span>
               </div>
            </div>
         </div>

         <div
            class="flex gap-2 mt-4 pt-4 border-t border-surface-100 dark:border-surface-800">
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
