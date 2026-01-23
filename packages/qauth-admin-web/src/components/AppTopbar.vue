<script setup lang="ts">
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useSidebarStore } from '@/stores/sidebar';
import { useThemeStore } from '@/stores/theme';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Avatar from 'primevue/avatar';
import Menu from 'primevue/menu';
import Badge from 'primevue/badge';

const router = useRouter();
const sidebarStore = useSidebarStore();
const themeStore = useThemeStore();

const searchQuery = ref('');
const userMenu = ref();
const notificationMenu = ref();

const userMenuItems = ref([
   {
      label: '个人资料',
      icon: 'pi pi-user',
      command: () => router.push('/profile'),
   },
   {
      label: '账号设置',
      icon: 'pi pi-cog',
      command: () => router.push('/settings'),
   },
   { separator: true },
   {
      label: '退出登录',
      icon: 'pi pi-sign-out',
      command: () => console.log('Logout'),
   },
]);

const notifications = ref([
   {
      id: 1,
      title: '新用户注册',
      message: 'john.doe@example.com 注册了新账号',
      time: '2分钟前',
      unread: true,
   },
   {
      id: 2,
      title: 'OAuth 授权请求',
      message: 'MyApp 请求用户授权',
      time: '15分钟前',
      unread: true,
   },
   {
      id: 3,
      title: '系统更新',
      message: '系统已成功更新至 v2.1.0',
      time: '1小时前',
      unread: false,
   },
]);

const unreadCount = computed(
   () => notifications.value.filter((n) => n.unread).length,
);

const toggleUserMenu = (event: Event) => {
   userMenu.value.toggle(event);
};

const toggleNotifications = (event: Event) => {
   notificationMenu.value.toggle(event);
};

const sidebarWidth = computed(() =>
   sidebarStore.isCollapsed ? '5rem' : '17rem',
);
</script>

<template>
   <header class="topbar" :style="{ '--sidebar-width': sidebarWidth }">
      <div class="topbar-left">
         <!-- 移动端菜单按钮 -->
         <Button
            icon="pi pi-bars"
            severity="secondary"
            text
            rounded
            class="mobile-menu-btn"
            @click="sidebarStore.toggleMobile" />

         <!-- 折叠按钮 -->
         <Button
            :icon="
               sidebarStore.isCollapsed
                  ? 'pi pi-angle-right'
                  : 'pi pi-angle-left'
            "
            severity="secondary"
            text
            rounded
            class="collapse-btn"
            @click="sidebarStore.toggleCollapsed"
            v-tooltip.bottom="
               sidebarStore.isCollapsed ? '展开菜单' : '折叠菜单'
            " />

         <!-- 面包屑 -->
         <nav class="breadcrumb">
            <span class="breadcrumb-item home">
               <i class="pi pi-home"></i>
            </span>
            <span class="breadcrumb-separator">/</span>
            <span class="breadcrumb-item current">仪表盘</span>
         </nav>
      </div>

      <div class="topbar-center">
         <!-- 搜索框 -->
         <div class="search-container">
            <span class="search-icon">
               <i class="pi pi-search"></i>
            </span>
            <InputText
               v-model="searchQuery"
               placeholder="搜索用户、应用、设置..."
               class="search-input" />
            <span class="search-shortcut">
               <kbd>⌘</kbd>
               <kbd>K</kbd>
            </span>
         </div>
      </div>

      <div class="topbar-right">
         <!-- 主题切换 -->
         <Button
            :icon="themeStore.isDark ? 'pi pi-sun' : 'pi pi-moon'"
            severity="secondary"
            text
            rounded
            class="theme-toggle"
            @click="themeStore.toggleTheme"
            v-tooltip.bottom="
               themeStore.isDark ? '切换到亮色模式' : '切换到暗色模式'
            " />

         <!-- 通知 -->
         <div class="notification-wrapper">
            <Button
               icon="pi pi-bell"
               severity="secondary"
               text
               rounded
               class="notification-btn"
               @click="toggleNotifications"
               v-tooltip.bottom="'通知'" />
            <Badge
               v-if="unreadCount > 0"
               :value="unreadCount"
               severity="danger"
               class="notification-badge" />
         </div>

         <!-- 用户头像 -->
         <button class="user-avatar-btn" @click="toggleUserMenu">
            <Avatar
               image="https://api.dicebear.com/7.x/avataaars/svg?seed=admin"
               shape="circle"
               class="user-avatar" />
            <span class="user-name">管理员</span>
            <i class="pi pi-chevron-down"></i>
         </button>

         <Menu ref="userMenu" :model="userMenuItems" popup class="user-menu" />

         <!-- 通知面板 -->
         <Menu ref="notificationMenu" popup class="notification-panel">
            <template #start>
               <div class="notification-header">
                  <span class="notification-title">通知</span>
                  <Button label="全部已读" text size="small" />
               </div>
            </template>
            <template #end>
               <div class="notification-list">
                  <div
                     v-for="notification in notifications"
                     :key="notification.id"
                     class="notification-item"
                     :class="{ unread: notification.unread }">
                     <div
                        class="notification-dot"
                        v-if="notification.unread"></div>
                     <div class="notification-content">
                        <span class="notification-item-title">{{
                           notification.title
                        }}</span>
                        <span class="notification-message">{{
                           notification.message
                        }}</span>
                        <span class="notification-time">{{
                           notification.time
                        }}</span>
                     </div>
                  </div>
               </div>
               <div class="notification-footer">
                  <Button label="查看全部通知" text class="view-all-btn" />
               </div>
            </template>
         </Menu>
      </div>
   </header>
</template>

<style scoped>
.topbar {
   position: fixed;
   top: 0;
   left: var(--sidebar-width);
   right: 0;
   height: 4rem;
   display: flex;
   align-items: center;
   justify-content: space-between;
   padding: 0 1.5rem;
   background: rgba(255, 255, 255, 0.8);
   backdrop-filter: blur(12px);
   border-bottom: 1px solid var(--p-surface-100);
   z-index: 900;
   transition: left 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

:global(.app-dark) .topbar {
   background: rgba(24, 24, 27, 0.85);
   border-bottom-color: var(--p-surface-800);
}

.topbar-left,
.topbar-right {
   display: flex;
   align-items: center;
   gap: 0.75rem;
}

.topbar-center {
   flex: 1;
   max-width: 32rem;
   margin: 0 1.5rem;
}

/* Mobile menu button */
.mobile-menu-btn {
   display: none;
}

/* Collapse button */
.collapse-btn {
   color: var(--p-surface-500);
}

/* Breadcrumb */
.breadcrumb {
   display: flex;
   align-items: center;
   gap: 0.5rem;
   font-size: 0.875rem;
   color: var(--p-surface-500);
}

.breadcrumb-item {
   display: flex;
   align-items: center;
}

.breadcrumb-item.home {
   color: var(--p-surface-400);
}

.breadcrumb-item.current {
   color: var(--p-surface-900);
   font-weight: 500;
}

:global(.app-dark) .breadcrumb-item.current {
   color: var(--p-surface-100);
}

.breadcrumb-separator {
   color: var(--p-surface-300);
}

/* Search */
.search-container {
   position: relative;
   display: flex;
   align-items: center;
}

.search-icon {
   position: absolute;
   left: 1rem;
   color: var(--p-surface-400);
   font-size: 0.875rem;
   pointer-events: none;
}

.search-input {
   width: 100%;
   padding-left: 2.75rem;
   padding-right: 4.5rem;
   height: 2.5rem;
   border-radius: 12px;
   background: var(--p-surface-50);
   border: 1px solid var(--p-surface-200);
   font-size: 0.875rem;
   transition: all 0.2s ease;
}

.search-input:focus {
   background: var(--p-surface-0);
   border-color: var(--p-orange-300);
   box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.1);
}

:global(.app-dark) .search-input {
   background: var(--p-surface-800);
   border-color: var(--p-surface-700);
}

:global(.app-dark) .search-input:focus {
   background: var(--p-surface-900);
   border-color: var(--p-orange-500);
   box-shadow: 0 0 0 3px rgba(251, 146, 60, 0.15);
}

.search-shortcut {
   position: absolute;
   right: 0.75rem;
   display: flex;
   gap: 0.25rem;
   pointer-events: none;
}

.search-shortcut kbd {
   padding: 0.125rem 0.375rem;
   border-radius: 4px;
   background: var(--p-surface-100);
   border: 1px solid var(--p-surface-200);
   font-size: 0.75rem;
   font-family: inherit;
   color: var(--p-surface-500);
}

:global(.app-dark) .search-shortcut kbd {
   background: var(--p-surface-700);
   border-color: var(--p-surface-600);
   color: var(--p-surface-400);
}

/* Theme toggle */
.theme-toggle {
   color: var(--p-surface-600);
}

:global(.app-dark) .theme-toggle {
   color: var(--p-surface-400);
}

/* Notifications */
.notification-wrapper {
   position: relative;
}

.notification-btn {
   color: var(--p-surface-600);
}

:global(.app-dark) .notification-btn {
   color: var(--p-surface-400);
}

.notification-badge {
   position: absolute;
   top: 0;
   right: 0;
   transform: translate(25%, -25%);
   min-width: 1.125rem;
   height: 1.125rem;
   font-size: 0.625rem;
}

/* User avatar button */
.user-avatar-btn {
   display: flex;
   align-items: center;
   gap: 0.625rem;
   padding: 0.375rem 0.75rem 0.375rem 0.375rem;
   border: none;
   border-radius: 9999px;
   background: var(--p-surface-50);
   cursor: pointer;
   transition: all 0.2s ease;
}

.user-avatar-btn:hover {
   background: var(--p-surface-100);
}

:global(.app-dark) .user-avatar-btn {
   background: var(--p-surface-800);
}

:global(.app-dark) .user-avatar-btn:hover {
   background: var(--p-surface-700);
}

.user-avatar {
   width: 2rem;
   height: 2rem;
   border: 2px solid var(--p-orange-200);
}

:global(.app-dark) .user-avatar {
   border-color: var(--p-orange-700);
}

.user-name {
   font-size: 0.875rem;
   font-weight: 500;
   color: var(--p-surface-700);
}

:global(.app-dark) .user-name {
   color: var(--p-surface-200);
}

.user-avatar-btn i {
   font-size: 0.75rem;
   color: var(--p-surface-400);
}

/* Menu styles */
:deep(.user-menu) {
   min-width: 12rem;
}

/* Notification panel */
:deep(.notification-panel) {
   width: 22rem;
   max-width: calc(100vw - 2rem);
}

.notification-header {
   display: flex;
   align-items: center;
   justify-content: space-between;
   padding: 1rem;
   border-bottom: 1px solid var(--p-surface-100);
}

:global(.app-dark) .notification-header {
   border-bottom-color: var(--p-surface-700);
}

.notification-title {
   font-weight: 600;
   color: var(--p-surface-900);
}

:global(.app-dark) .notification-title {
   color: var(--p-surface-100);
}

.notification-list {
   max-height: 20rem;
   overflow-y: auto;
}

.notification-item {
   display: flex;
   gap: 0.75rem;
   padding: 1rem;
   border-bottom: 1px solid var(--p-surface-50);
   cursor: pointer;
   transition: background 0.2s ease;
}

.notification-item:hover {
   background: var(--p-surface-50);
}

:global(.app-dark) .notification-item {
   border-bottom-color: var(--p-surface-800);
}

:global(.app-dark) .notification-item:hover {
   background: var(--p-surface-800);
}

.notification-item.unread {
   background: var(--p-orange-50);
}

:global(.app-dark) .notification-item.unread {
   background: rgba(251, 146, 60, 0.08);
}

.notification-dot {
   flex-shrink: 0;
   width: 0.5rem;
   height: 0.5rem;
   margin-top: 0.375rem;
   border-radius: 9999px;
   background: var(--p-orange-500);
}

.notification-content {
   display: flex;
   flex-direction: column;
   gap: 0.25rem;
}

.notification-item-title {
   font-size: 0.875rem;
   font-weight: 600;
   color: var(--p-surface-900);
}

:global(.app-dark) .notification-item-title {
   color: var(--p-surface-100);
}

.notification-message {
   font-size: 0.8125rem;
   color: var(--p-surface-600);
}

:global(.app-dark) .notification-message {
   color: var(--p-surface-400);
}

.notification-time {
   font-size: 0.75rem;
   color: var(--p-surface-400);
}

.notification-footer {
   padding: 0.75rem;
   border-top: 1px solid var(--p-surface-100);
}

:global(.app-dark) .notification-footer {
   border-top-color: var(--p-surface-700);
}

.view-all-btn {
   width: 100%;
}

/* Mobile */
@media (max-width: 1024px) {
   .topbar {
      left: 0;
   }

   .mobile-menu-btn {
      display: flex;
   }

   .collapse-btn {
      display: none;
   }

   .topbar-center {
      display: none;
   }

   .breadcrumb {
      display: none;
   }

   .user-name {
      display: none;
   }
}
</style>
