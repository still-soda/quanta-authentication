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
   <header
      class="fixed top-0 right-0 h-16 flex items-center justify-between px-6 bg-white/80 dark:bg-[rgba(24,24,27,0.85)] backdrop-blur-xl border-b border-surface-100 dark:border-surface-800 z-[900] transition-[left] duration-300 ease-[cubic-bezier(0.4,0,0.2,1)] max-lg:left-0!"
      :style="{ left: sidebarWidth }">
      <div class="flex items-center gap-3">
         <!-- 移动端菜单按钮 -->
         <Button
            icon="pi pi-bars"
            severity="secondary"
            text
            rounded
            class="hidden! max-lg:flex!"
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
            class="text-surface-500 max-lg:hidden!"
            @click="sidebarStore.toggleCollapsed"
            v-tooltip.bottom="
               sidebarStore.isCollapsed ? '展开菜单' : '折叠菜单'
            " />

         <!-- 面包屑 -->
         <nav
            class="flex items-center gap-2 text-sm text-surface-500 max-lg:hidden!">
            <span class="flex items-center text-surface-400">
               <i class="pi pi-home"></i>
            </span>
            <span class="text-surface-300">/</span>
            <span
               class="text-surface-900 dark:text-surface-100 font-medium">
               仪表盘
            </span>
         </nav>
      </div>

      <div class="flex-1 max-w-128 mx-6 max-lg:hidden!">
         <!-- 搜索框 -->
         <div class="relative flex items-center">
            <span
               class="absolute left-4 text-surface-400 text-sm pointer-events-none">
               <i class="pi pi-search"></i>
            </span>
            <InputText
               v-model="searchQuery"
               placeholder="搜索用户、应用、设置..."
               class="w-full pl-11 pr-18 h-10 rounded-xl bg-surface-50 dark:bg-surface-800 border border-surface-200 dark:border-surface-700 text-sm transition-all duration-200 ease focus:bg-surface-0! dark:focus:bg-surface-900! focus:border-primary-300! dark:focus:border-primary-500! focus:shadow-[0_0_0_3px_rgba(249,115,22,0.1)]! dark:focus:shadow-[0_0_0_3px_rgba(251,146,60,0.15)]!" />
            <span
               class="absolute right-3 flex gap-1 pointer-events-none">
               <kbd
                  class="py-0.5 px-1.5 rounded bg-surface-100 dark:bg-surface-700 border border-surface-200 dark:border-surface-600 text-xs font-inherit text-surface-500 dark:text-surface-400">
                  ⌘
               </kbd>
               <kbd
                  class="py-0.5 px-1.5 rounded bg-surface-100 dark:bg-surface-700 border border-surface-200 dark:border-surface-600 text-xs font-inherit text-surface-500 dark:text-surface-400">
                  K
               </kbd>
            </span>
         </div>
      </div>

      <div class="flex items-center gap-3">
         <!-- 主题切换 -->
         <Button
            :icon="themeStore.isDark ? 'pi pi-sun' : 'pi pi-moon'"
            severity="secondary"
            text
            rounded
            class="text-surface-600! dark:text-surface-400!"
            @click="themeStore.toggleTheme"
            v-tooltip.bottom="
               themeStore.isDark ? '切换到亮色模式' : '切换到暗色模式'
            " />

         <!-- 通知 -->
         <div class="relative">
            <Button
               icon="pi pi-bell"
               severity="secondary"
               text
               rounded
               class="text-surface-600! dark:text-surface-400!"
               @click="toggleNotifications"
               v-tooltip.bottom="'通知'" />
            <Badge
               v-if="unreadCount > 0"
               :value="unreadCount"
               severity="danger"
               class="absolute! top-0! right-0! translate-x-1/4! -translate-y-1/4! min-w-4.5! h-4.5! text-[0.625rem]!" />
         </div>

         <!-- 用户头像 -->
         <button
            class="flex items-center gap-2.5 py-1.5 pr-3 pl-1.5 border-none rounded-full bg-surface-50 dark:bg-surface-800 cursor-pointer transition-all duration-200 ease hover:bg-surface-100 dark:hover:bg-surface-700"
            @click="toggleUserMenu">
            <Avatar
               image="https://api.dicebear.com/7.x/avataaars/svg?seed=admin"
               shape="circle"
               class="w-8! h-8! border-2! border-primary-200! dark:border-primary-700!" />
            <span
               class="text-sm font-medium text-surface-700 dark:text-surface-200 max-lg:hidden!">
               管理员
            </span>
            <i class="pi pi-chevron-down text-xs text-surface-400"></i>
         </button>

         <Menu ref="userMenu" :model="userMenuItems" popup class="min-w-48" />

         <!-- 通知面板 -->
         <Menu
            ref="notificationMenu"
            popup
            class="w-88! max-w-[calc(100vw-2rem)]!">
            <template #start>
               <div
                  class="flex items-center justify-between p-4 border-b border-surface-100 dark:border-surface-700">
                  <span
                     class="font-semibold text-surface-900 dark:text-surface-100">
                     通知
                  </span>
                  <Button label="全部已读" text size="small" />
               </div>
            </template>
            <template #end>
               <div class="max-h-80 overflow-y-auto">
                  <div
                     v-for="notification in notifications"
                     :key="notification.id"
                     class="flex gap-3 p-4 border-b border-surface-50 dark:border-surface-800 cursor-pointer transition-colors duration-200 ease hover:bg-surface-50 dark:hover:bg-surface-800"
                     :class="{
                        'bg-primary-50! dark:bg-[rgba(251,146,60,0.08)]!':
                           notification.unread,
                     }">
                     <div
                        v-if="notification.unread"
                        class="shrink-0 w-2 h-2 mt-1.5 rounded-full bg-primary-500"></div>
                     <div class="flex flex-col gap-1">
                        <span
                           class="text-sm font-semibold text-surface-900 dark:text-surface-100">
                           {{ notification.title }}
                        </span>
                        <span
                           class="text-[0.8125rem] text-surface-600 dark:text-surface-400">
                           {{ notification.message }}
                        </span>
                        <span class="text-xs text-surface-400">
                           {{ notification.time }}
                        </span>
                     </div>
                  </div>
               </div>
               <div
                  class="p-3 border-t border-surface-100 dark:border-surface-700">
                  <Button label="查看全部通知" text class="w-full!" />
               </div>
            </template>
         </Menu>
      </div>
   </header>
</template>
