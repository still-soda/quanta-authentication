<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import Card from 'primevue/card'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Avatar from 'primevue/avatar'
import { getRecentActivities } from '@/apis/dashboard'

// 获取最近活动数据
const { data: recentActivities, isLoading } = useQuery({
   queryKey: ['dashboard', 'recentActivities'],
   queryFn: getRecentActivities,
})

const getStatusSeverity = (status: string) => {
   const map: Record<string, 'success' | 'warn' | 'danger' | 'info'> = {
      success: 'success',
      warning: 'warn',
      danger: 'danger',
      info: 'info',
   }
   return map[status] || 'info'
}
</script>

<template>
   <Card class="rounded-2xl border border-surface-100 dark:border-surface-800">
      <template #title>
         <div
            class="flex items-center justify-between text-base font-semibold text-surface-900 dark:text-surface-100"
         >
            <span>最近活动</span>
            <Button label="查看全部" text size="small" />
         </div>
      </template>
      <template #content>
         <div v-if="isLoading" class="flex items-center justify-center py-12">
            <i class="pi pi-spin pi-spinner text-2xl text-surface-400"></i>
         </div>
         <DataTable
            v-else
            :value="recentActivities"
            :rows="5"
            class="text-sm"
            :pt="{
               table: { style: 'min-width: 40rem' },
            }"
         >
            <Column field="user" header="用户">
               <template #body="{ data }">
                  <div class="flex items-center gap-3">
                     <Avatar :image="data.avatar" shape="circle" size="small" />
                     <span class="text-surface-700 dark:text-surface-200 font-medium">
                        {{ data.user }}
                     </span>
                  </div>
               </template>
            </Column>
            <Column field="action" header="操作">
               <template #body="{ data }">
                  <span class="text-surface-600 dark:text-surface-400">
                     {{ data.action }}
                  </span>
               </template>
            </Column>
            <Column field="client" header="来源">
               <template #body="{ data }">
                  <span class="text-surface-500 text-[0.8125rem]">
                     {{ data.client }}
                  </span>
               </template>
            </Column>
            <Column field="time" header="时间">
               <template #body="{ data }">
                  <span class="text-surface-400 text-[0.8125rem]">
                     {{ data.time }}
                  </span>
               </template>
            </Column>
            <Column field="status" header="状态">
               <template #body="{ data }">
                  <Tag :severity="getStatusSeverity(data.status)" class="whitespace-nowrap">
                     {{ data.action.includes('失败') ? '失败' : '成功' }}
                  </Tag>
               </template>
            </Column>
         </DataTable>
      </template>
   </Card>
</template>
