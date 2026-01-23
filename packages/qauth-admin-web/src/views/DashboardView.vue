<script setup lang="ts">
import { ref } from 'vue';
import Button from 'primevue/button';
import PageHeader from '@/components/shared/PageHeader.vue';
import StatsCard, {
   type StatCardData,
} from '@/components/dashboard/StatsCard.vue';
import AuthTrendChart from '@/components/dashboard/AuthTrendChart.vue';
import UserDistChart from '@/components/dashboard/UserDistChart.vue';
import RecentActivities from '@/components/dashboard/RecentActivities.vue';
import TopApps from '@/components/dashboard/TopApps.vue';

// 统计卡片数据
const statsCards = ref<StatCardData[]>([
   {
      title: '总用户数',
      value: '12,847',
      change: '+12.5%',
      changeType: 'increase',
      icon: 'pi pi-users',
      color: 'blue',
      trendData: [30, 42, 38, 52, 45, 58, 50, 65, 55, 72],
   },
   {
      title: 'OAuth 应用',
      value: '86',
      change: '+3',
      changeType: 'increase',
      icon: 'pi pi-key',
      color: 'orange',
      trendData: [40, 45, 42, 50, 48, 55, 52, 60, 58, 65],
   },
   {
      title: '今日认证',
      value: '2,451',
      change: '+8.2%',
      changeType: 'increase',
      icon: 'pi pi-shield',
      color: 'green',
      trendData: [35, 48, 40, 55, 45, 62, 52, 68, 58, 75],
   },
   {
      title: '活跃会话',
      value: '847',
      change: '-2.1%',
      changeType: 'decrease',
      icon: 'pi pi-bolt',
      color: 'purple',
      trendData: [70, 62, 68, 55, 60, 50, 58, 45, 52, 40],
   },
]);
</script>

<template>
   <div class="dashboard">
      <!-- Page Header -->
      <PageHeader title="仪表盘" subtitle="欢迎回来，这是您的认证中心概览">
         <template #actions>
            <Button
               label="导出报表"
               icon="pi pi-download"
               severity="secondary"
               outlined />
            <Button label="刷新数据" icon="pi pi-refresh" />
         </template>
      </PageHeader>

      <!-- Stats Cards -->
      <div class="stats-grid">
         <StatsCard v-for="stat in statsCards" :key="stat.title" :stat="stat" />
      </div>

      <!-- Charts Row -->
      <div class="charts-row">
         <AuthTrendChart class="trend-chart" />
         <UserDistChart class="distribution-chart" />
      </div>

      <!-- Bottom Row -->
      <div class="bottom-row">
         <RecentActivities />
         <TopApps />
      </div>
   </div>
</template>

<style scoped>
.dashboard {
   display: flex;
   flex-direction: column;
   gap: 1.5rem;
}

/* Stats Grid */
.stats-grid {
   display: grid;
   grid-template-columns: repeat(4, 1fr);
   gap: 1.25rem;
}

/* Charts Row */
.charts-row {
   display: grid;
   grid-template-columns: 2fr 1fr;
   gap: 1.25rem;
}

/* Bottom Row */
.bottom-row {
   display: grid;
   grid-template-columns: 2fr 1fr;
   gap: 1.25rem;
}

/* Responsive */
@media (max-width: 1280px) {
   .stats-grid {
      grid-template-columns: repeat(2, 1fr);
   }
}

@media (max-width: 1024px) {
   .charts-row,
   .bottom-row {
      grid-template-columns: 1fr;
   }
}

@media (max-width: 640px) {
   .stats-grid {
      grid-template-columns: 1fr;
   }
}
</style>
