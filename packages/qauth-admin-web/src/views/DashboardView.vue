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
      iconBg: 'linear-gradient(135deg, #f97316 0%, #ea580c 100%)',
   },
   {
      title: 'OAuth 应用',
      value: '86',
      change: '+3',
      changeType: 'increase',
      icon: 'pi pi-key',
      iconBg: 'linear-gradient(135deg, #3b82f6 0%, #2563eb 100%)',
   },
   {
      title: '今日认证',
      value: '2,451',
      change: '+8.2%',
      changeType: 'increase',
      icon: 'pi pi-shield',
      iconBg: 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
   },
   {
      title: '活跃会话',
      value: '847',
      change: '-2.1%',
      changeType: 'decrease',
      icon: 'pi pi-bolt',
      iconBg: 'linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%)',
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
