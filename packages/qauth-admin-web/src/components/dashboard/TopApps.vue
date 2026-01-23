<script setup lang="ts">
import { ref } from 'vue';
import Card from 'primevue/card';
import Button from 'primevue/button';
import ProgressBar from 'primevue/progressbar';

export interface TopApp {
   name: string;
   users: number;
   percentage: number;
}

const topApps = ref<TopApp[]>([
   { name: 'Web Dashboard', users: 4521, percentage: 85 },
   { name: 'Mobile App', users: 3842, percentage: 72 },
   { name: 'API Gateway', users: 2156, percentage: 54 },
   { name: 'Admin Portal', users: 1234, percentage: 38 },
]);
</script>

<template>
   <Card class="top-apps-card">
      <template #title>
         <div class="card-header">
            <span>热门应用</span>
            <Button label="管理" text size="small" icon="pi pi-cog" />
         </div>
      </template>
      <template #content>
         <div class="apps-list">
            <div v-for="app in topApps" :key="app.name" class="app-item">
               <div class="app-info">
                  <div class="app-icon">
                     <i class="pi pi-box"></i>
                  </div>
                  <div class="app-details">
                     <span class="app-name">{{ app.name }}</span>
                     <span class="app-users">
                        {{ app.users.toLocaleString() }} 用户
                     </span>
                  </div>
               </div>
               <div class="app-progress">
                  <ProgressBar :value="app.percentage" :showValue="false" />
                  <span class="app-percentage">{{ app.percentage }}%</span>
               </div>
            </div>
         </div>
      </template>
   </Card>
</template>

<style scoped>
.top-apps-card {
   border-radius: 16px;
   border: 1px solid var(--p-surface-100);
}

:global(.app-dark) .top-apps-card {
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

.apps-list {
   display: flex;
   flex-direction: column;
   gap: 1.25rem;
}

.app-item {
   display: flex;
   flex-direction: column;
   gap: 0.75rem;
}

.app-info {
   display: flex;
   align-items: center;
   gap: 0.75rem;
}

.app-icon {
   width: 2.5rem;
   height: 2.5rem;
   display: flex;
   align-items: center;
   justify-content: center;
   background: var(--p-orange-50);
   color: var(--p-orange-600);
   border-radius: 10px;
}

:global(.app-dark) .app-icon {
   background: rgba(251, 146, 60, 0.15);
   color: var(--p-orange-400);
}

.app-details {
   display: flex;
   flex-direction: column;
   gap: 0.125rem;
}

.app-name {
   font-weight: 600;
   color: var(--p-surface-900);
   font-size: 0.875rem;
}

:global(.app-dark) .app-name {
   color: var(--p-surface-100);
}

.app-users {
   font-size: 0.75rem;
   color: var(--p-surface-500);
}

.app-progress {
   display: flex;
   align-items: center;
   gap: 0.75rem;
}

.app-progress :deep(.p-progressbar) {
   flex: 1;
   height: 0.5rem;
   border-radius: 9999px;
}

.app-percentage {
   font-size: 0.75rem;
   font-weight: 600;
   color: var(--p-surface-600);
   min-width: 2.5rem;
   text-align: right;
}

:global(.app-dark) .app-percentage {
   color: var(--p-surface-400);
}
</style>
