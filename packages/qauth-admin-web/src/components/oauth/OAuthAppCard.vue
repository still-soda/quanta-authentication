<script setup lang="ts">
import Card from 'primevue/card'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Chip from 'primevue/chip'
import { useToast } from 'primevue/usetoast'
import type { OAuthApp } from '@/types'

defineProps<{ app: OAuthApp }>()

const emit = defineEmits<{
   (e: 'view', app: OAuthApp): void
   (e: 'edit', app: OAuthApp): void
   (e: 'delete', app: OAuthApp): void
   (e: 'regenerateSecret', app: OAuthApp): void
   (e: 'managePermissions', app: OAuthApp): void
}>()

const toast = useToast()

const STATUS_MAP = {
   active: { label: '生产环境', severity: 'success' },
   development: { label: '开发中', severity: 'warn' },
   deprecated: { label: '已弃用', severity: 'secondary' },
} as const
const getStatusSeverity = (status: string) =>
   (STATUS_MAP[status as keyof typeof STATUS_MAP]?.severity || 'info') as
      | 'success'
      | 'warn'
      | 'danger'
      | 'secondary'
const getStatusLabel = (status: string) =>
   STATUS_MAP[status as keyof typeof STATUS_MAP]?.label || status

const formatNumber = (num: number) =>
   num >= 1e6
      ? (num / 1e6).toFixed(1) + 'M'
      : num >= 1e3
        ? (num / 1e3).toFixed(1) + 'K'
        : String(num)

const formatDate = (dateStr: string | null) => {
   if (!dateStr) return '从未使用'
   const date = new Date(dateStr)
   return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
   })
}

const copyToClipboard = async (text: string) => {
   await navigator.clipboard.writeText(text)
   toast.add({
      severity: 'success',
      summary: '已复制',
      detail: 'Client ID 已复制到剪贴板',
      life: 2000,
   })
}
</script>

<template>
   <Card
      class="rounded-2xl border border-surface-100 dark:border-surface-800 transition-all duration-300 ease hover:-translate-y-0.5 hover:shadow-[0_12px_24px_-8px_rgba(0,0,0,0.08)] dark:hover:shadow-[0_12px_24px_-8px_rgba(0,0,0,0.3)]"
   >
      <template #content>
         <div class="flex justify-between items-start mb-4">
            <div
               class="w-12 h-12 flex items-center justify-center rounded-xl text-white text-xl shadow-[0_4px_12px_rgba(0,0,0,0.15)] overflow-hidden"
               :style="{ background: app.logo ? 'transparent' : app.icon_bg }"
            >
               <img
                  v-if="app.logo"
                  :src="`/api/uploads/${app.logo}`"
                  alt="App Logo"
                  class="w-full h-full object-cover"
               />
               <i v-else :class="app.icon"></i>
            </div>
            <div class="flex gap-2">
               <Tag v-if="app.trusted" severity="info" class="flex items-center gap-1" rounded>
                  <i class="pi pi-verified"></i>
                  可信
               </Tag>
               <Tag :severity="getStatusSeverity(app.status)">
                  {{ getStatusLabel(app.status) }}
               </Tag>
            </div>
         </div>

         <div class="flex flex-col gap-4">
            <h3 class="text-lg font-semibold text-surface-900 dark:text-surface-100 m-0">
               {{ app.name }}
            </h3>
            <p class="text-sm text-surface-500 m-0 leading-relaxed line-clamp-2 min-h-10">
               {{ app.description || '暂无描述' }}
            </p>

            <div class="flex flex-col gap-1.5">
               <label class="text-xs font-medium text-surface-500 uppercase tracking-wider">
                  Client ID
               </label>
               <div
                  class="flex items-center gap-2 py-2 px-3 bg-surface-50 dark:bg-surface-800 rounded-lg"
               >
                  <code
                     class="flex-1 text-[0.8125rem] font-mono text-surface-700 dark:text-surface-300 truncate"
                  >
                     {{ app.client_id }}
                  </code>
                  <Button
                     icon="pi pi-copy"
                     text
                     rounded
                     severity="secondary"
                     size="small"
                     @click="copyToClipboard(app.client_id)"
                     v-tooltip.top="'复制'"
                  />
               </div>
            </div>

            <div class="flex flex-col gap-2">
               <label class="text-xs font-medium text-surface-500 uppercase tracking-wider">
                  授权范围
               </label>
               <div class="flex flex-wrap gap-1.5">
                  <Chip
                     v-for="scope in app.scopes?.slice(0, 4)"
                     :key="scope"
                     :label="scope"
                     class="text-xs py-1 px-2"
                  />
                  <Chip
                     v-if="app.scopes && app.scopes.length > 4"
                     :label="`+${app.scopes.length - 4}`"
                     class="text-xs py-1 px-2 bg-surface-200 dark:bg-surface-700"
                  />
               </div>
            </div>

            <div class="flex gap-5">
               <div class="flex items-center gap-2 text-[0.8125rem] text-surface-500">
                  <i class="pi pi-chart-line text-sm"></i>
                  <span>{{ formatNumber(app.request_count) }} 请求</span>
               </div>
               <div class="flex items-center gap-2 text-[0.8125rem] text-surface-500">
                  <i class="pi pi-clock text-sm"></i>
                  <span>{{ formatDate(app.last_used_at) }}</span>
               </div>
            </div>
         </div>

         <div
            class="flex gap-2 mt-4 pt-4 border-t border-surface-100 whitespace-nowrap dark:border-surface-800"
         >
            <Button
               icon="pi pi-eye"
               label="查看"
               severity="secondary"
               outlined
               size="small"
               @click="emit('view', app)"
            />
            <Button
               icon="pi pi-shield"
               label="权限"
               severity="secondary"
               text
               size="small"
               @click="emit('managePermissions', app)"
               v-tooltip.top="'管理应用组权限'"
            />
            <Button
               icon="pi pi-refresh"
               label="重置密钥"
               severity="secondary"
               text
               size="small"
               @click="emit('regenerateSecret', app)"
            />
            <Button
               icon="pi pi-pencil"
               text
               rounded
               severity="secondary"
               @click="emit('edit', app)"
               v-tooltip.top="'编辑'"
            />
            <Button
               icon="pi pi-trash"
               text
               rounded
               severity="danger"
               @click="emit('delete', app)"
               v-tooltip.top="'删除'"
            />
         </div>
      </template>
   </Card>
</template>
