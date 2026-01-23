<script setup lang="ts">
import { ref } from 'vue';
import Card from 'primevue/card';
import Button from 'primevue/button';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Tag from 'primevue/tag';
import Avatar from 'primevue/avatar';
import type { Activity } from '@/types';

const recentActivities = ref<Activity[]>([
   {
      user: 'zhang.wei@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhang',
      action: '登录成功',
      client: 'Web Dashboard',
      time: '2分钟前',
      status: 'success',
   },
   {
      user: 'li.ming@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=li',
      action: '密码重置',
      client: 'Mobile App',
      time: '5分钟前',
      status: 'warning',
   },
   {
      user: 'wang.fang@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=wang',
      action: 'OAuth 授权',
      client: 'Third Party App',
      time: '12分钟前',
      status: 'success',
   },
   {
      user: 'chen.hong@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=chen',
      action: '登录失败',
      client: 'API Client',
      time: '15分钟前',
      status: 'danger',
   },
   {
      user: 'zhao.yang@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhao',
      action: '新用户注册',
      client: 'Web Portal',
      time: '23分钟前',
      status: 'info',
   },
]);

const getStatusSeverity = (status: string) => {
   const map: Record<string, 'success' | 'warn' | 'danger' | 'info'> = {
      success: 'success',
      warning: 'warn',
      danger: 'danger',
      info: 'info',
   };
   return map[status] || 'info';
};
</script>

<template>
   <Card
      class="rounded-2xl border border-surface-100 dark:border-surface-800">
      <template #title>
         <div
            class="flex items-center justify-between text-base font-semibold text-surface-900 dark:text-surface-100">
            <span>最近活动</span>
            <Button label="查看全部" text size="small" />
         </div>
      </template>
      <template #content>
         <DataTable
            :value="recentActivities"
            :rows="5"
            class="text-sm"
            :pt="{
               table: { style: 'min-width: 40rem' },
            }">
            <Column field="user" header="用户">
               <template #body="{ data }">
                  <div class="flex items-center gap-3">
                     <Avatar :image="data.avatar" shape="circle" size="small" />
                     <span
                        class="text-surface-700 dark:text-surface-200 font-medium">
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
                  <Tag :severity="getStatusSeverity(data.status)">
                     {{ data.action.includes('失败') ? '失败' : '成功' }}
                  </Tag>
               </template>
            </Column>
         </DataTable>
      </template>
   </Card>
</template>
