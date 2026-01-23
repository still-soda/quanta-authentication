<script setup lang="ts">
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useSidebarStore } from '@/stores/sidebar';
import { useThemeStore } from '@/stores/theme';

const route = useRoute();
const router = useRouter();
const sidebarStore = useSidebarStore();
const themeStore = useThemeStore();

interface MenuItem {
   label: string;
   icon: string;
   to?: string;
   items?: MenuItem[];
   badge?: string | number;
}

const menuItems: MenuItem[] = [
   {
      label: '仪表盘',
      icon: 'pi pi-th-large',
      to: '/',
   },
   {
      label: '用户管理',
      icon: 'pi pi-users',
      to: '/users',
      badge: '128',
   },
   {
      label: '角色权限',
      icon: 'pi pi-shield',
      to: '/roles',
   },
   {
      label: 'OAuth 应用',
      icon: 'pi pi-key',
      to: '/oauth',
      badge: '12',
   },
   {
      label: '组织管理',
      icon: 'pi pi-building',
      to: '/organizations',
   },
   {
      label: '审计日志',
      icon: 'pi pi-history',
      to: '/audit',
   },
   {
      label: '系统设置',
      icon: 'pi pi-cog',
      to: '/settings',
   },
];

const isActive = (item: MenuItem) => {
   if (item.to === '/') {
      return route.path === '/';
   }
   return route.path.startsWith(item.to || '');
};

const navigateTo = (item: MenuItem) => {
   if (item.to) {
      router.push(item.to);
      sidebarStore.closeMobile();
   }
};

const sidebarWidth = computed(() =>
   sidebarStore.isCollapsed ? '5rem' : '17rem',
);
</script>

<template>
   <aside
      class="sidebar"
      :class="{
         collapsed: sidebarStore.isCollapsed,
         'mobile-open': sidebarStore.isMobileOpen,
      }">
      <!-- 遮罩层 (移动端) -->
      <div
         class="sidebar-overlay"
         :class="{ active: sidebarStore.isMobileOpen }"
         @click="sidebarStore.closeMobile"></div>

      <div class="sidebar-inner">
         <!-- Logo 区域 -->
         <div class="sidebar-header">
            <div class="logo-container">
               <div class="logo-icon">
                  <img src="/quantacenter.jpg" alt="Quanta Logo" />
               </div>
               <Transition name="fade">
                  <div v-if="!sidebarStore.isCollapsed" class="logo-text">
                     <span class="logo-title">Quanta</span>
                     <span class="logo-subtitle">认证中心</span>
                  </div>
               </Transition>
            </div>
         </div>

         <!-- 导航菜单 -->
         <nav class="sidebar-nav">
            <ul class="nav-list">
               <li
                  v-for="item in menuItems"
                  :key="item.label"
                  class="nav-item"
                  :class="{ active: isActive(item) }">
                  <button
                     class="nav-link"
                     @click="navigateTo(item)"
                     v-tooltip.right="{
                        value: item.label,
                        disabled: !sidebarStore.isCollapsed,
                     }">
                     <span class="nav-icon">
                        <i :class="item.icon"></i>
                     </span>
                     <Transition name="fade">
                        <span
                           v-if="!sidebarStore.isCollapsed"
                           class="nav-label">
                           {{ item.label }}
                        </span>
                     </Transition>
                     <Transition name="fade">
                        <span
                           v-if="item.badge && !sidebarStore.isCollapsed"
                           class="nav-badge">
                           {{ item.badge }}
                        </span>
                     </Transition>
                  </button>
               </li>
            </ul>
         </nav>

         <!-- 底部用户区域 -->
         <div class="sidebar-footer">
            <div class="user-card">
               <div class="user-avatar">
                  <img
                     src="https://api.dicebear.com/7.x/avataaars/svg?seed=admin"
                     alt="User Avatar" />
               </div>
               <Transition name="fade">
                  <div v-if="!sidebarStore.isCollapsed" class="user-info">
                     <span class="user-name">管理员</span>
                     <span class="user-role">超级管理员</span>
                  </div>
               </Transition>
            </div>
         </div>
      </div>
   </aside>
</template>

<style scoped>
.sidebar {
   position: fixed;
   left: 0;
   top: 0;
   bottom: 0;
   width: v-bind(sidebarWidth);
   z-index: 1000;
   transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.sidebar-overlay {
   display: none;
   position: fixed;
   inset: 0;
   background: rgba(0, 0, 0, 0.5);
   backdrop-filter: blur(4px);
   z-index: -1;
   opacity: 0;
   transition: opacity 0.3s ease;
}

.sidebar-inner {
   display: flex;
   flex-direction: column;
   height: 100%;
   background: linear-gradient(
      180deg,
      var(--p-surface-0) 0%,
      var(--p-surface-50) 100%
   );
   border-right: 1px solid var(--p-surface-200);
   box-shadow: 4px 0 24px rgba(0, 0, 0, 0.04);
}

:global(.app-dark) .sidebar-inner {
   background: linear-gradient(
      180deg,
      var(--p-surface-900) 0%,
      var(--p-surface-950) 100%
   );
   border-right-color: var(--p-surface-700);
   box-shadow: 4px 0 24px rgba(0, 0, 0, 0.3);
}

/* Header / Logo */
.sidebar-header {
   padding: 1.5rem 1rem;
   border-bottom: 1px solid var(--p-surface-100);
}

:global(.app-dark) .sidebar-header {
   border-bottom-color: var(--p-surface-800);
}

.logo-container {
   display: flex;
   align-items: center;
   gap: 0.75rem;
}

.logo-icon {
   flex-shrink: 0;
   width: 3rem;
   height: 3rem;
   display: flex;
   align-items: center;
   justify-content: center;
   background: linear-gradient(
      135deg,
      var(--p-orange-400) 0%,
      var(--p-orange-600) 100%
   );
   border-radius: 12px;
   color: white;
   font-size: 1.25rem;
   box-shadow: 0 4px 12px rgba(51, 21, 0, 0.3);
   transition: transform 0.2s ease;
}

.logo-icon img {
   width: 100%;
   height: 100%;
   object-fit: cover;
   border-radius: 12px;
}

.logo-icon:hover {
   transform: scale(1.05);
}

.logo-text {
   display: flex;
   flex-direction: column;
   line-height: 1.2;
   overflow: hidden;
   white-space: nowrap;
}

.logo-title {
   font-size: 1.125rem;
   font-weight: 700;
   background: linear-gradient(
      135deg,
      var(--p-orange-500) 0%,
      var(--p-orange-700) 100%
   );
   -webkit-background-clip: text;
   -webkit-text-fill-color: transparent;
   background-clip: text;
}

:global(.app-dark) .logo-title {
   background: linear-gradient(
      135deg,
      var(--p-orange-300) 0%,
      var(--p-orange-500) 100%
   );
   -webkit-background-clip: text;
   background-clip: text;
}

.logo-subtitle {
   font-size: 0.75rem;
   color: var(--p-surface-500);
   letter-spacing: 0.05em;
}

/* Navigation */
.sidebar-nav {
   flex: 1;
   padding: 1rem 0.75rem;
   overflow-y: auto;
   overflow-x: hidden;
}

.nav-list {
   list-style: none;
   margin: 0;
   padding: 0;
   display: flex;
   flex-direction: column;
   gap: 0.25rem;
}

.nav-item {
   position: relative;
}

.nav-link {
   display: flex;
   align-items: center;
   gap: 0.875rem;
   width: 100%;
   padding: 0.875rem 1rem;
   border: none;
   border-radius: 10px;
   background: transparent;
   color: var(--p-surface-600);
   cursor: pointer;
   transition: all 0.2s ease;
   text-align: left;
   font-size: 0.9375rem;
   font-weight: 500;
}

.nav-link:hover {
   background: var(--p-surface-100);
   color: var(--p-surface-900);
}

:global(.app-dark) .nav-link {
   color: var(--p-surface-400);
}

:global(.app-dark) .nav-link:hover {
   background: var(--p-surface-800);
   color: var(--p-surface-100);
}

.nav-item.active .nav-link {
   background: linear-gradient(
      135deg,
      var(--p-orange-50) 0%,
      var(--p-orange-100) 100%
   );
   color: var(--p-orange-700);
   box-shadow: 0 2px 8px rgba(249, 115, 22, 0.12);
}

:global(.app-dark) .nav-item.active .nav-link {
   background: linear-gradient(
      135deg,
      rgba(251, 146, 60, 0.15) 0%,
      rgba(251, 146, 60, 0.1) 100%
   );
   color: var(--p-orange-400);
   box-shadow: 0 2px 8px rgba(251, 146, 60, 0.15);
}

.nav-icon {
   flex-shrink: 0;
   width: 1.5rem;
   display: flex;
   align-items: center;
   justify-content: center;
   font-size: 1.125rem;
}

.nav-label {
   flex: 1;
   white-space: nowrap;
   overflow: hidden;
}

.nav-badge {
   padding: 0.125rem 0.5rem;
   border-radius: 9999px;
   background: var(--p-orange-100);
   color: var(--p-orange-700);
   font-size: 0.75rem;
   font-weight: 600;
}

:global(.app-dark) .nav-badge {
   background: rgba(251, 146, 60, 0.2);
   color: var(--p-orange-400);
}

/* Footer / User */
.sidebar-footer {
   padding: 1rem;
   border-top: 1px solid var(--p-surface-100);
}

:global(.app-dark) .sidebar-footer {
   border-top-color: var(--p-surface-800);
}

.user-card {
   display: flex;
   align-items: center;
   gap: 0.75rem;
   padding: 0.75rem;
   border-radius: 12px;
   background: var(--p-surface-50);
   transition: background 0.2s ease;
}

:global(.app-dark) .user-card {
   background: var(--p-surface-800);
}

.user-avatar {
   flex-shrink: 0;
   width: 2.5rem;
   height: 2.5rem;
   border-radius: 10px;
   overflow: hidden;
   background: linear-gradient(
      135deg,
      var(--p-orange-200) 0%,
      var(--p-orange-300) 100%
   );
   box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.user-avatar img {
   width: 100%;
   height: 100%;
   object-fit: cover;
}

.user-info {
   display: flex;
   flex-direction: column;
   overflow: hidden;
   white-space: nowrap;
}

.user-name {
   font-size: 0.875rem;
   font-weight: 600;
   color: var(--p-surface-900);
}

:global(.app-dark) .user-name {
   color: var(--p-surface-100);
}

.user-role {
   font-size: 0.75rem;
   color: var(--p-surface-500);
}

/* Collapsed state */
.sidebar.collapsed .logo-container {
   justify-content: center;
}

.sidebar.collapsed .nav-link {
   justify-content: center;
   padding: 0.875rem;
}

.sidebar.collapsed .user-card {
   justify-content: center;
   padding: 0.75rem;
}

/* Transitions */
.fade-enter-active,
.fade-leave-active {
   transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
   opacity: 0;
}

/* Mobile */
@media (max-width: 1024px) {
   .sidebar {
      width: 17rem !important;
      transform: translateX(-100%);
   }

   .sidebar.mobile-open {
      transform: translateX(0);
   }

   .sidebar-overlay {
      display: block;
   }

   .sidebar-overlay.active {
      opacity: 1;
      z-index: -1;
   }
}
</style>
