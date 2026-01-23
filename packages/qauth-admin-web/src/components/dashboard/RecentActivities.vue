<script setup lang="ts">
import { ref } from 'vue';
import Card from 'primevue/card';
import Button from 'primevue/button';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Tag from 'primevue/tag';
import Avatar from 'primevue/avatar';

export interface Activity {
   user: string;
   avatar: string;
   action: string;
   client: string;
   time: string;
   status: string;
}

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
   <Card class="activities-card">
      <template #title>
         <div class="card-header">
            <span>最近活动</span>
            <Button label="查看全部" text size="small" />
         </div>
      </template>
      <template #content>
         <DataTable
            :value="recentActivities"
            :rows="5"
            class="activities-table"
            :pt="{
               table: { style: 'min-width: 40rem' },
            }">
            <Column field="user" header="用户">
               <template #body="{ data }">
                  <div class="user-cell">
                     <Avatar :image="data.avatar" shape="circle" size="small" />
                     <span class="user-email">{{ data.user }}</span>
                  </div>
               </template>
            </Column>
            <Column field="action" header="操作">
               <template #body="{ data }">
                  <span class="action-text">{{ data.action }}</span>
               </template>
            </Column>
            <Column field="client" header="来源">
               <template #body="{ data }">
                  <span class="client-text">{{ data.client }}</span>
               </template>
            </Column>
            <Column field="time" header="时间">
               <template #body="{ data }">
                  <span class="time-text">{{ data.time }}</span>
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

<style scoped>
.activities-card {
   border-radius: 16px;
   border: 1px solid var(--p-surface-100);
}

:global(.app-dark) .activities-card {
   border-color: var(--p-surface-800);
}

.card-header {
   display: flex;
   align-items: center;
   justify-content: space-between;
   font-size: 1rem;
   font-weight: 600;
   color: var(--p-surface-900);
}

:global(.app-dark) .card-header {
   color: var(--p-surface-100);
}

.activities-table {
   font-size: 0.875rem;
}

.user-cell {
   display: flex;
   align-items: center;
   gap: 0.75rem;
}

.user-email {
   color: var(--p-surface-700);
   font-weight: 500;
}

:global(.app-dark) .user-email {
   color: var(--p-surface-200);
}

.action-text {
   color: var(--p-surface-600);
}

:global(.app-dark) .action-text {
   color: var(--p-surface-400);
}

.client-text {
   color: var(--p-surface-500);
   font-size: 0.8125rem;
}

.time-text {
   color: var(--p-surface-400);
   font-size: 0.8125rem;
}
</style>
