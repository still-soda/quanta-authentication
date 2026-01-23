<script setup lang="ts">
import { computed } from 'vue';
import { useQuery } from '@tanstack/vue-query';
import Card from 'primevue/card';
import Button from 'primevue/button';
import Chart from 'primevue/chart';
import { useThemeStore } from '@/stores/theme';
import { getAuthTrendData } from '@/apis/dashboard';

const themeStore = useThemeStore();

// 获取认证趋势数据
const { data: trendData, isLoading } = useQuery({
   queryKey: ['dashboard', 'authTrend'],
   queryFn: getAuthTrendData,
});

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

const authTrendChartData = computed(() => {
   if (!trendData.value) return null;
   return {
      labels: trendData.value.labels,
      datasets: [
         {
            label: '认证请求',
            data: trendData.value.data,
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
   };
});
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
            <div
               v-if="isLoading"
               class="h-full flex items-center justify-center">
               <i class="pi pi-spin pi-spinner text-2xl text-surface-400"></i>
            </div>
            <Chart
               v-else-if="authTrendChartData"
               type="line"
               :data="authTrendChartData"
               :options="authTrendOptions"
               class="h-full" />
         </div>
      </template>
   </Card>
</template>
