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
      class="fixed left-0 top-0 bottom-0 z-[1000] transition-[width] duration-300 ease-[cubic-bezier(0.4,0,0.2,1)] max-lg:w-68! max-lg:-translate-x-full"
      :class="{
         'max-lg:translate-x-0': sidebarStore.isMobileOpen,
      }"
      :style="{ width: sidebarWidth }">
      <!-- 遮罩层 (移动端) -->
      <div
         class="hidden max-lg:block fixed inset-0 bg-black/50 backdrop-blur-sm -z-1 opacity-0 transition-opacity duration-300 ease"
         :class="{ 'opacity-100!': sidebarStore.isMobileOpen }"
         @click="sidebarStore.closeMobile"></div>

      <div
         class="flex flex-col h-full bg-linear-to-b from-surface-0 to-surface-50 dark:from-surface-900 dark:to-surface-950 border-r border-surface-200 dark:border-surface-700 shadow-[4px_0_24px_rgba(0,0,0,0.04)] dark:shadow-[4px_0_24px_rgba(0,0,0,0.3)]">
         <!-- Logo 区域 -->
         <div
            class="p-6 px-4 border-b border-surface-100 dark:border-surface-800">
            <div
               class="flex items-center gap-3"
               :class="{ 'justify-center': sidebarStore.isCollapsed }">
               <div
                  class="shrink-0 w-12 h-12 flex items-center justify-center bg-linear-to-br from-primary-400 to-primary-600 rounded-xl text-white text-xl shadow-[0_4px_12px_rgba(51,21,0,0.3)] transition-transform duration-200 ease hover:scale-105">
                  <img
                     src="/quantacenter.jpg"
                     alt="Quanta Logo"
                     class="w-full h-full object-cover rounded-xl" />
               </div>
               <Transition name="fade">
                  <div
                     v-if="!sidebarStore.isCollapsed"
                     class="flex flex-col leading-tight overflow-hidden whitespace-nowrap">
                     <span
                        class="text-lg font-bold bg-linear-to-br from-primary-500 to-primary-700 dark:from-primary-300 dark:to-primary-500 bg-clip-text text-transparent">
                        Quanta
                     </span>
                     <span
                        class="text-xs text-surface-500 tracking-wider">
                        认证中心
                     </span>
                  </div>
               </Transition>
            </div>
         </div>

         <!-- 导航菜单 -->
         <nav class="flex-1 p-4 px-3 overflow-y-auto overflow-x-hidden">
            <ul class="list-none m-0 p-0 flex flex-col gap-1">
               <li
                  v-for="item in menuItems"
                  :key="item.label"
                  class="relative group">
                  <button
                     class="flex items-center gap-3.5 w-full py-3.5 px-4 border-none rounded-[10px] bg-transparent text-surface-600 dark:text-surface-400 cursor-pointer transition-all duration-200 ease text-left text-[0.9375rem] font-medium hover:bg-surface-100 hover:text-surface-900 dark:hover:bg-surface-800 dark:hover:text-surface-100"
                     :class="{
                        'justify-center! px-3.5!': sidebarStore.isCollapsed,
                        'bg-linear-to-br! from-primary-50! to-primary-100! text-primary-700! shadow-[0_2px_8px_rgba(249,115,22,0.12)]! dark:from-[rgba(251,146,60,0.15)]! dark:to-[rgba(251,146,60,0.1)]! dark:text-primary-400! dark:shadow-[0_2px_8px_rgba(251,146,60,0.15)]!':
                           isActive(item),
                     }"
                     @click="navigateTo(item)"
                     v-tooltip.right="{
                        value: item.label,
                        disabled: !sidebarStore.isCollapsed,
                     }">
                     <span
                        class="shrink-0 w-6 flex items-center justify-center text-lg">
                        <i :class="item.icon"></i>
                     </span>
                     <Transition name="fade">
                        <span
                           v-if="!sidebarStore.isCollapsed"
                           class="flex-1 whitespace-nowrap overflow-hidden">
                           {{ item.label }}
                        </span>
                     </Transition>
                     <Transition name="fade">
                        <span
                           v-if="item.badge && !sidebarStore.isCollapsed"
                           class="py-0.5 px-2 rounded-full bg-primary-100 dark:bg-[rgba(251,146,60,0.2)] text-primary-700 dark:text-primary-400 text-xs font-semibold">
                           {{ item.badge }}
                        </span>
                     </Transition>
                  </button>
               </li>
            </ul>
         </nav>

         <!-- 底部用户区域 -->
         <div
            class="p-4 border-t border-surface-100 dark:border-surface-800">
            <div
               class="flex items-center gap-3 p-3 rounded-xl bg-surface-50 dark:bg-surface-800 transition-colors duration-200 ease"
               :class="{ 'justify-center! p-3!': sidebarStore.isCollapsed }">
               <div
                  class="shrink-0 w-10 h-10 rounded-[10px] overflow-hidden bg-linear-to-br from-primary-200 to-primary-300 shadow-[0_2px_8px_rgba(0,0,0,0.08)]">
                  <img
                     src="https://api.dicebear.com/7.x/avataaars/svg?seed=admin"
                     alt="User Avatar"
                     class="w-full h-full object-cover" />
               </div>
               <Transition name="fade">
                  <div
                     v-if="!sidebarStore.isCollapsed"
                     class="flex flex-col overflow-hidden whitespace-nowrap">
                     <span
                        class="text-sm font-semibold text-surface-900 dark:text-surface-100">
                        管理员
                     </span>
                     <span class="text-xs text-surface-500">
                        超级管理员
                     </span>
                  </div>
               </Transition>
            </div>
         </div>
      </div>
   </aside>
</template>

<style scoped>
/* Transitions */
.fade-enter-active,
.fade-leave-active {
   transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
   opacity: 0;
}
</style>
