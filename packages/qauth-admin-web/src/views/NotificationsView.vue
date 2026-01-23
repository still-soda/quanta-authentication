<script setup lang="ts">
import { ref, computed } from 'vue';
import Button from 'primevue/button';
import Tag from 'primevue/tag';
import Checkbox from 'primevue/checkbox';
import Menu from 'primevue/menu';
import PageHeader from '@/components/shared/PageHeader.vue';
import SimpleStatCard from '@/components/shared/SimpleStatCard.vue';
import type { Notification, NotificationType, SimpleStatData } from '@/types';
import { NOTIFICATION_TYPE_CONFIG } from '@/config';

// 通知数据
const notifications = ref<Notification[]>([
   {
      id: '1',
      type: 'security',
      title: '新设备登录提醒',
      message: '检测到您的账号在新设备上登录：Chrome / Windows，IP: 192.168.1.100，位置：北京',
      time: '2026-01-23 14:30:00',
      read: false,
      actionUrl: '/audit',
      actionLabel: '查看详情',
   },
   {
      id: '2',
      type: 'system',
      title: '系统维护通知',
      message: '系统将于 2026年1月25日 02:00-06:00 进行例行维护，届时服务将暂停访问',
      time: '2026-01-23 10:00:00',
      read: false,
   },
   {
      id: '3',
      type: 'oauth',
      title: 'OAuth 应用审核通过',
      message: '您提交的应用「内部管理系统」已通过审核，现在可以正常使用',
      time: '2026-01-22 16:45:00',
      read: false,
      actionUrl: '/oauth',
      actionLabel: '前往管理',
   },
   {
      id: '4',
      type: 'user',
      title: '新用户注册通知',
      message: '新用户 李四 (lisi@example.com) 已完成注册并通过邮箱验证',
      time: '2026-01-22 14:20:00',
      read: true,
      actionUrl: '/users',
      actionLabel: '查看用户',
   },
   {
      id: '5',
      type: 'alert',
      title: '异常登录警告',
      message: '用户 赵阳 在短时间内多次登录失败，已触发账号锁定机制',
      time: '2026-01-22 11:30:00',
      read: true,
      actionUrl: '/audit',
      actionLabel: '查看日志',
   },
   {
      id: '6',
      type: 'security',
      title: '密码修改成功',
      message: '您的账号密码已成功修改，如非本人操作请立即联系管理员',
      time: '2026-01-21 18:00:00',
      read: true,
   },
   {
      id: '7',
      type: 'system',
      title: '系统更新完成',
      message: '系统已升级至 v2.5.0 版本，新增组织架构管理等功能',
      time: '2026-01-20 09:00:00',
      read: true,
   },
   {
      id: '8',
      type: 'oauth',
      title: 'OAuth 密钥即将过期',
      message: '应用「客户端 App」的密钥将于 30 天后过期，请及时更新',
      time: '2026-01-19 14:00:00',
      read: true,
      actionUrl: '/oauth',
      actionLabel: '立即更新',
   },
   {
      id: '9',
      type: 'user',
      title: '用户权限变更',
      message: '管理员为您分配了新角色：审计员，您现在可以查看系统审计日志',
      time: '2026-01-18 11:20:00',
      read: true,
   },
   {
      id: '10',
      type: 'alert',
      title: '存储空间预警',
      message: '系统存储空间使用率已达 85%，建议清理无用文件或扩容',
      time: '2026-01-17 16:30:00',
      read: true,
   },
]);

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
const markAsRead = (notification: Notification) => {
   notification.read = true;
};

const markSelectedAsRead = () => {
   notifications.value.forEach((n) => {
      if (selectedIds.value.includes(n.id)) {
         n.read = true;
      }
   });
   selectedIds.value = [];
};

const markAllAsRead = () => {
   notifications.value.forEach((n) => {
      n.read = true;
   });
};

// 删除通知
const deleteNotification = (id: string) => {
   notifications.value = notifications.value.filter((n) => n.id !== id);
   selectedIds.value = selectedIds.value.filter((i) => i !== id);
};

const deleteSelected = () => {
   notifications.value = notifications.value.filter(
      (n) => !selectedIds.value.includes(n.id),
   );
   selectedIds.value = [];
};

// 更多操作菜单
const moreMenu = ref();
const moreMenuItems = ref([
   {
      label: '全部标记已读',
      icon: 'pi pi-check-circle',
      command: markAllAsRead,
   },
   {
      label: '删除所有已读',
      icon: 'pi pi-trash',
      command: () => {
         notifications.value = notifications.value.filter((n) => !n.read);
      },
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
         <SimpleStatCard v-for="stat in stats" :key="stat.title" :stat="stat" />
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

         <!-- Notification Items -->
         <div
            v-if="filteredNotifications.length > 0"
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
                           @click="markAsRead(notification)" />
                        <Button
                           icon="pi pi-trash"
                           severity="danger"
                           text
                           rounded
                           size="small"
                           v-tooltip.top="'删除'"
                           @click="deleteNotification(notification.id)" />
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
