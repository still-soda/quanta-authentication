<script setup lang="ts">
import { ref, computed } from 'vue';
import Card from 'primevue/card';
import Button from 'primevue/button';
import Chart from 'primevue/chart';
import { useThemeStore } from '@/stores/theme';

const themeStore = useThemeStore();

const userDistOptions = computed(() => ({
   maintainAspectRatio: false,
   plugins: {
      legend: {
         position: 'right',
         labels: {
            color: themeStore.isDark ? '#a1a1aa' : '#71717a',
            usePointStyle: true,
            padding: 20,
         },
      },
   },
}));

const userDistData = ref({
   labels: ['管理员', '普通用户', '开发者', '访客'],
   datasets: [
      {
         data: [12, 847, 234, 156],
         backgroundColor: ['#f97316', '#3b82f6', '#10b981', '#8b5cf6'],
         borderWidth: 0,
      },
   ],
});
</script>

<template>
   <Card
      class="rounded-2xl border border-surface-100 dark:border-surface-800 overflow-hidden">
      <template #title>
         <div
            class="flex items-center justify-between text-base font-semibold text-surface-900 dark:text-surface-100">
            <span>用户分布</span>
            <div>
               <Button
                  icon="pi pi-ellipsis-h"
                  text
                  rounded
                  severity="secondary" />
            </div>
         </div>
      </template>
      <template #content>
         <div class="h-70">
            <Chart
               type="doughnut"
               :data="userDistData"
               :options="userDistOptions"
               class="h-full" />
         </div>
      </template>
   </Card>
</template>
