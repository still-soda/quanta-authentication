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
   <Card
      class="rounded-2xl border border-surface-100 dark:border-surface-800">
      <template #title>
         <div
            class="flex items-center justify-between text-base font-semibold text-surface-900 dark:text-surface-100">
            <span>热门应用</span>
            <Button label="管理" text size="small" icon="pi pi-cog" />
         </div>
      </template>
      <template #content>
         <div class="flex flex-col gap-5">
            <div
               v-for="app in topApps"
               :key="app.name"
               class="flex flex-col gap-3">
               <div class="flex items-center gap-3">
                  <div
                     class="w-10 h-10 flex items-center justify-center bg-primary-50 dark:bg-[rgba(251,146,60,0.15)] text-primary-600 dark:text-primary-400 rounded-[10px]">
                     <i class="pi pi-box"></i>
                  </div>
                  <div class="flex flex-col gap-0.5">
                     <span
                        class="font-semibold text-surface-900 dark:text-surface-100 text-sm">
                        {{ app.name }}
                     </span>
                     <span class="text-xs text-surface-500">
                        {{ app.users.toLocaleString() }} 用户
                     </span>
                  </div>
               </div>
               <div class="flex items-center gap-3">
                  <ProgressBar
                     :value="app.percentage"
                     :showValue="false"
                     class="flex-1 h-2 rounded-full [&_.p-progressbar-value]:rounded-full" />
                  <span
                     class="text-xs font-semibold text-surface-600 dark:text-surface-400 min-w-10 text-right">
                     {{ app.percentage }}%
                  </span>
               </div>
            </div>
         </div>
      </template>
   </Card>
</template>
