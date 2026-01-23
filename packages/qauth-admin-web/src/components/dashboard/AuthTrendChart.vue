<script setup lang="ts">
import { computed } from 'vue';
import Card from 'primevue/card';
import Button from 'primevue/button';
import Chart from 'primevue/chart';
import { useThemeStore } from '@/stores/theme';

const themeStore = useThemeStore();

const authTrendOptions = computed(() => ({
   maintainAspectRatio: false,
   plugins: {
      legend: {
         display: false,
      },
      tooltip: {
         mode: 'index',
         intersect: false,
         backgroundColor: themeStore.isDark ? '#27272a' : '#ffffff',
         titleColor: themeStore.isDark ? '#fafafa' : '#18181b',
         bodyColor: themeStore.isDark ? '#a1a1aa' : '#71717a',
         borderColor: themeStore.isDark ? '#3f3f46' : '#e4e4e7',
         borderWidth: 1,
         padding: 12,
         cornerRadius: 8,
      },
   },
   scales: {
      x: {
         grid: {
            display: false,
         },
         ticks: {
            color: themeStore.isDark ? '#71717a' : '#a1a1aa',
         },
      },
      y: {
         grid: {
            color: themeStore.isDark ? '#27272a' : '#f4f4f5',
         },
         ticks: {
            color: themeStore.isDark ? '#71717a' : '#a1a1aa',
         },
      },
   },
   interaction: {
      intersect: false,
      mode: 'index',
   },
}));

const authTrendData = computed(() => ({
   labels: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
   datasets: [
      {
         label: '认证请求',
         data: [1200, 1900, 1500, 2100, 1800, 900, 1100],
         fill: true,
         borderColor: '#f97316',
         backgroundColor: themeStore.isDark
            ? 'rgba(249, 115, 22, 0.1)'
            : 'rgba(249, 115, 22, 0.08)',
         tension: 0.4,
         pointRadius: 0,
         pointHoverRadius: 6,
         pointHoverBackgroundColor: '#f97316',
         pointHoverBorderColor: '#ffffff',
         pointHoverBorderWidth: 2,
      },
   ],
}));
</script>

<template>
   <Card
      class="rounded-2xl border border-surface-100 dark:border-surface-800 overflow-hidden">
      <template #title>
         <div
            class="flex items-center justify-between text-base font-semibold text-surface-900 dark:text-surface-100">
            <span>认证趋势</span>
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
               type="line"
               :data="authTrendData"
               :options="authTrendOptions"
               class="h-full" />
         </div>
      </template>
   </Card>
</template>
