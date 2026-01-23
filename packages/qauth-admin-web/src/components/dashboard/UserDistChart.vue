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
   <Card class="chart-card distribution-chart">
      <template #title>
         <div class="card-header">
            <span>用户分布</span>
            <div class="card-actions">
               <Button
                  icon="pi pi-ellipsis-h"
                  text
                  rounded
                  severity="secondary" />
            </div>
         </div>
      </template>
      <template #content>
         <div class="chart-container">
            <Chart
               type="doughnut"
               :data="userDistData"
               :options="userDistOptions"
               class="dist-chart" />
         </div>
      </template>
   </Card>
</template>

<style scoped>
.chart-card {
   border-radius: 16px;
   border: 1px solid var(--p-surface-100);
   overflow: hidden;
}

:global(.app-dark) .chart-card {
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

.chart-container {
   height: 280px;
}

.dist-chart {
   height: 100%;
}
</style>
