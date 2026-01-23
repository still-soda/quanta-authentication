<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query';
import Button from 'primevue/button';
import PageHeader from '@/components/shared/PageHeader.vue';
import StatsCard from '@/components/dashboard/StatsCard.vue';
import AuthTrendChart from '@/components/dashboard/AuthTrendChart.vue';
import UserDistChart from '@/components/dashboard/UserDistChart.vue';
import RecentActivities from '@/components/dashboard/RecentActivities.vue';
import TopApps from '@/components/dashboard/TopApps.vue';
import { getDashboardStats } from '@/apis/dashboard';

// 使用 TanStack Query 获取统计卡片数据
const { data: statsData, isLoading, refetch } = useQuery({
   queryKey: ['dashboard', 'stats'],
   queryFn: getDashboardStats,
});

const handleRefresh = () => {
   refetch();
};
</script>

<template>
   <div class="flex flex-col gap-6">
      <!-- Page Header -->
      <PageHeader title="仪表盘" subtitle="欢迎回来，这是您的认证中心概览">
         <template #actions>
            <Button
               label="导出报表"
               icon="pi pi-download"
               severity="secondary"
               outlined />
            <Button
               label="刷新数据"
               icon="pi pi-refresh"
               :loading="isLoading"
               @click="handleRefresh" />
         </template>
      </PageHeader>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-5">
         <template v-if="isLoading">
            <div
               v-for="i in 4"
               :key="i"
               class="h-32 bg-surface-100 dark:bg-surface-800 rounded-2xl animate-pulse" />
         </template>
         <template v-else-if="statsData">
            <StatsCard
               v-for="stat in statsData.cards"
               :key="stat.title"
               :stat="stat" />
         </template>
      </div>

      <!-- Charts Row -->
      <div class="grid grid-cols-1 lg:grid-cols-[2fr_1fr] gap-5">
         <AuthTrendChart />
         <UserDistChart />
      </div>

      <!-- Bottom Row -->
      <div class="grid grid-cols-1 lg:grid-cols-[2fr_1fr] gap-5">
         <RecentActivities />
         <TopApps />
      </div>
   </div>
</template>
