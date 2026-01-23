<script setup lang="ts">
import { computed } from 'vue';
import { useQuery } from '@tanstack/vue-query';
import Card from 'primevue/card';
import Button from 'primevue/button';
import Chart from 'primevue/chart';
import { useThemeStore } from '@/stores/theme';
import { getUserDistData } from '@/apis/dashboard';

const themeStore = useThemeStore();

// 获取用户分布数据
const { data: distData, isLoading } = useQuery({
   queryKey: ['dashboard', 'userDist'],
   queryFn: getUserDistData,
});

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

const userDistChartData = computed(() => {
   if (!distData.value) return null;
   return {
      labels: distData.value.labels,
      datasets: [
         {
            data: distData.value.data,
            backgroundColor: distData.value.colors,
            borderWidth: 0,
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
            <div
               v-if="isLoading"
               class="h-full flex items-center justify-center">
               <i class="pi pi-spin pi-spinner text-2xl text-surface-400"></i>
            </div>
            <Chart
               v-else-if="userDistChartData"
               type="doughnut"
               :data="userDistChartData"
               :options="userDistOptions"
               class="h-full" />
         </div>
      </template>
   </Card>
</template>
