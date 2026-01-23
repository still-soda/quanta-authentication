<script setup lang="ts">
import { ref, computed } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import Button from 'primevue/button';
import Tag from 'primevue/tag';
import Checkbox from 'primevue/checkbox';
import Menu from 'primevue/menu';
import PageHeader from '@/components/shared/PageHeader.vue';
import SimpleStatCard from '@/components/shared/SimpleStatCard.vue';
import type { Notification, NotificationType, SimpleStatData } from '@/types';
import { NOTIFICATION_TYPE_CONFIG } from '@/config';
import {
   getNotifications,
   markAsRead,
   markMultipleAsRead,
   markAllAsRead as markAllAsReadApi,
   deleteNotification,
   deleteMultipleNotifications,
   deleteAllReadNotifications,
} from '@/apis/notifications';

const queryClient = useQueryClient();

// 使用 TanStack Query 获取通知数据
const { data: notificationsData, isLoading } = useQuery({
   queryKey: ['notifications'],
   queryFn: () => getNotifications(),
});

const notifications = computed(() => notificationsData.value || []);

// Mutations
const markAsReadMutation = useMutation({
   mutationFn: markAsRead,
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] });
   },
});

const markMultipleMutation = useMutation({
   mutationFn: markMultipleAsRead,
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] });
      selectedIds.value = [];
   },
});

const markAllMutation = useMutation({
   mutationFn: markAllAsReadApi,
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] });
   },
});

const deleteMutation = useMutation({
   mutationFn: deleteNotification,
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] });
   },
});

const deleteMultipleMutation = useMutation({
   mutationFn: deleteMultipleNotifications,
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] });
      selectedIds.value = [];
   },
});

const deleteAllReadMutation = useMutation({
   mutationFn: deleteAllReadNotifications,
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] });
   },
});

// 筛选类型
const filterType = ref<NotificationType | 'all'>('all');

// 选中的通知
const selectedIds = ref<string[]>([]);

// 统计数据
const stats = computed<SimpleStatData[]>(() => {
   const total = notifications.value.length;
   const unread = notifications.value.filter((n) => !n.read).length;
   const today = notifications.value.filter((n) =>
      n.time.startsWith('2026-01-23'),
   ).length;
   const alerts = notifications.value.filter((n) => n.type === 'alert').length;
   return [
      {
         title: '全部通知',
         value: total,
         icon: 'pi pi-bell',
         color: 'blue',
      },
      {
         title: '未读消息',
         value: unread,
         icon: 'pi pi-envelope',
         color: 'orange',
      },
      {
         title: '今日通知',
         value: today,
         icon: 'pi pi-calendar',
         color: 'green',
      },
      {
         title: '警告通知',
         value: alerts,
         icon: 'pi pi-exclamation-triangle',
         color: 'red',
      },
   ];
});

// 过滤后的通知
const filteredNotifications = computed(() => {
   if (filterType.value === 'all') {
      return notifications.value;
   }
   return notifications.value.filter((n) => n.type === filterType.value);
});

const typeConfig = NOTIFICATION_TYPE_CONFIG;

// 过滤选项
const filterOptions = [
   { label: '全部', value: 'all', icon: 'pi pi-inbox' },
   { label: '系统通知', value: 'system', icon: 'pi pi-cog' },
   { label: '安全通知', value: 'security', icon: 'pi pi-shield' },
   { label: '用户通知', value: 'user', icon: 'pi pi-user' },
   { label: 'OAuth 通知', value: 'oauth', icon: 'pi pi-key' },
   { label: '警告通知', value: 'alert', icon: 'pi pi-exclamation-triangle' },
];

// 全选
const isAllSelected = computed(() => {
   return (
      filteredNotifications.value.length > 0 &&
      filteredNotifications.value.every((n) => selectedIds.value.includes(n.id))
   );
});

const toggleSelectAll = () => {
   if (isAllSelected.value) {
      selectedIds.value = [];
   } else {
      selectedIds.value = filteredNotifications.value.map((n) => n.id);
   }
};

// 标记已读
const handleMarkAsRead = (notification: Notification) => {
   markAsReadMutation.mutate(notification.id);
};

const markSelectedAsRead = () => {
   markMultipleMutation.mutate(selectedIds.value);
};

const handleMarkAllAsRead = () => {
   markAllMutation.mutate();
};

// 删除通知
const handleDeleteNotification = (id: string) => {
   deleteMutation.mutate(id);
};

const deleteSelected = () => {
   deleteMultipleMutation.mutate(selectedIds.value);
};

// 更多操作菜单
const moreMenu = ref();
const moreMenuItems = ref([
   {
      label: '全部标记已读',
      icon: 'pi pi-check-circle',
      command: handleMarkAllAsRead,
   },
   {
      label: '删除所有已读',
      icon: 'pi pi-trash',
      command: () => deleteAllReadMutation.mutate(),
   },
   { separator: true },
   {
      label: '通知设置',
      icon: 'pi pi-cog',
      command: () => {
         console.log('Go to settings');
      },
   },
]);

const toggleMoreMenu = (event: Event) => {
   moreMenu.value.toggle(event);
};

// 格式化时间
const formatTime = (time: string) => {
   const date = new Date(time);
   const now = new Date();
   const diff = now.getTime() - date.getTime();
   const minutes = Math.floor(diff / 60000);
   const hours = Math.floor(diff / 3600000);
   const days = Math.floor(diff / 86400000);

   if (minutes < 60) {
      return minutes <= 0 ? '刚刚' : `${minutes} 分钟前`;
   } else if (hours < 24) {
      return `${hours} 小时前`;
   } else if (days < 7) {
      return `${days} 天前`;
   } else {
      return time.split(' ')[0];
   }
};
</script>

<template>
   <div class="flex flex-col gap-6">
      <!-- Page Header -->
      <PageHeader title="全部通知" subtitle="查看和管理系统通知消息">
         <template #actions>
            <Button
               v-if="selectedIds.length > 0"
               :label="`标记已读 (${selectedIds.length})`"
               icon="pi pi-check"
               severity="secondary"
               outlined
               @click="markSelectedAsRead" />
            <Button
               v-if="selectedIds.length > 0"
               :label="`删除 (${selectedIds.length})`"
               icon="pi pi-trash"
               severity="danger"
               outlined
               @click="deleteSelected" />
            <Button
               icon="pi pi-ellipsis-v"
               severity="secondary"
               text
               rounded
               @click="toggleMoreMenu" />
            <Menu ref="moreMenu" :model="moreMenuItems" :popup="true" />
         </template>
      </PageHeader>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4">
         <template v-if="isLoading">
            <div
               v-for="i in 4"
               :key="i"
               class="h-20 bg-surface-100 dark:bg-surface-800 rounded-xl animate-pulse" />
         </template>
         <template v-else>
            <SimpleStatCard v-for="stat in stats" :key="stat.title" :stat="stat" />
         </template>
      </div>

      <!-- Filter Tabs -->
      <div
         class="flex flex-wrap gap-2 p-4 bg-surface-0 dark:bg-surface-900 border border-surface-200 dark:border-surface-800 rounded-xl">
         <button
            v-for="option in filterOptions"
            :key="option.value"
            class="flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all border-none cursor-pointer"
            :class="
               filterType === option.value
                  ? 'bg-primary-100 dark:bg-primary-900/30 text-primary-700 dark:text-primary-400'
                  : 'bg-transparent text-surface-600 dark:text-surface-400 hover:bg-surface-100 dark:hover:bg-surface-800'
            "
            @click="filterType = option.value as NotificationType | 'all'">
            <i :class="option.icon"></i>
            <span>{{ option.label }}</span>
            <span
               v-if="option.value !== 'all'"
               class="text-xs px-1.5 py-0.5 rounded-full bg-surface-200 dark:bg-surface-700 text-surface-600 dark:text-surface-400">
               {{
                  notifications.filter(
                     (n) => n.type === option.value,
                  ).length
               }}
            </span>
         </button>
      </div>

      <!-- Notifications List -->
      <div
         class="bg-surface-0 dark:bg-surface-900 border border-surface-200 dark:border-surface-800 rounded-xl overflow-hidden">
         <!-- List Header -->
         <div
            class="flex items-center gap-4 px-5 py-3 border-b border-surface-200 dark:border-surface-700 bg-surface-50 dark:bg-surface-800">
            <Checkbox
               :modelValue="isAllSelected"
               :binary="true"
               @change="toggleSelectAll" />
            <span class="text-sm text-surface-500">
               共 {{ filteredNotifications.length }} 条通知
               <template
                  v-if="
                     filteredNotifications.filter((n) => !n.read).length > 0
                  ">
                  ，{{ filteredNotifications.filter((n) => !n.read).length }} 条未读
               </template>
            </span>
         </div>

         <!-- Loading State -->
         <div
            v-if="isLoading"
            class="flex items-center justify-center py-24">
            <i class="pi pi-spin pi-spinner text-3xl text-surface-400"></i>
         </div>

         <!-- Notification Items -->
         <div
            v-else-if="filteredNotifications.length > 0"
            class="divide-y divide-surface-100 dark:divide-surface-800">
            <div
               v-for="notification in filteredNotifications"
               :key="notification.id"
               class="group flex items-start gap-4 p-5 transition-colors"
               :class="{
                  'bg-primary-50/50 dark:bg-primary-900/10': !notification.read,
                  'hover:bg-surface-50 dark:hover:bg-surface-800/50':
                     notification.read,
               }">
               <!-- Checkbox -->
               <div class="pt-1">
                  <Checkbox v-model="selectedIds" :value="notification.id" />
               </div>

               <!-- Type Icon -->
               <div
                  class="shrink-0 w-10 h-10 rounded-xl flex items-center justify-center"
                  :class="[
                     typeConfig[notification.type].bgColor,
                     typeConfig[notification.type].color,
                  ]">
                  <i :class="typeConfig[notification.type].icon"></i>
               </div>

               <!-- Content -->
               <div class="flex-1 min-w-0">
                  <div class="flex items-start justify-between gap-4">
                     <div class="flex-1 min-w-0">
                        <div class="flex items-center gap-2 mb-1">
                           <h4
                              class="text-sm font-semibold text-surface-900 dark:text-surface-100 m-0"
                              :class="{ 'font-bold': !notification.read }">
                              {{ notification.title }}
                           </h4>
                           <span
                              v-if="!notification.read"
                              class="w-2 h-2 rounded-full bg-primary-500 shrink-0"></span>
                        </div>
                        <p
                           class="text-sm text-surface-600 dark:text-surface-400 m-0 line-clamp-2">
                           {{ notification.message }}
                        </p>
                        <div
                           class="flex items-center gap-3 mt-3 text-xs text-surface-500">
                           <Tag
                              :value="typeConfig[notification.type].label"
                              severity="secondary"
                              :pt="{ root: { class: 'text-xs px-2 py-0.5' } }" />
                           <span>{{ formatTime(notification.time) }}</span>
                        </div>
                     </div>

                     <!-- Actions -->
                     <div
                        class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity shrink-0">
                        <Button
                           v-if="notification.actionUrl"
                           :label="notification.actionLabel"
                           size="small"
                           text
                           class="text-xs" />
                        <Button
                           v-if="!notification.read"
                           icon="pi pi-check"
                           severity="secondary"
                           text
                           rounded
                           size="small"
                           v-tooltip.top="'标记已读'"
                           @click="handleMarkAsRead(notification)" />
                        <Button
                           icon="pi pi-trash"
                           severity="danger"
                           text
                           rounded
                           size="small"
                           v-tooltip.top="'删除'"
                           @click="handleDeleteNotification(notification.id)" />
                     </div>
                  </div>
               </div>
            </div>
         </div>

         <!-- Empty State -->
         <div
            v-else
            class="flex flex-col items-center justify-center py-16 text-surface-500">
            <i class="pi pi-inbox text-5xl mb-4 opacity-30"></i>
            <span class="text-lg font-medium mb-1">暂无通知</span>
            <span class="text-sm">当前筛选条件下没有通知消息</span>
         </div>
      </div>
   </div>
</template>

<style scoped>
.line-clamp-2 {
   display: -webkit-box;
   -webkit-line-clamp: 2;
   -webkit-box-orient: vertical;
   overflow: hidden;
}
</style>
