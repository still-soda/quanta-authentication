<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useSidebarStore } from '@/stores/sidebar'
import { useThemeStore } from '@/stores/theme'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'
import Menu from 'primevue/menu'
import Badge from 'primevue/badge'
import GlobalSearchDialog from '@/components/shared/GlobalSearchDialog.vue'
import type { MenuItem } from 'primevue/menuitem'

const router = useRouter()
const route = useRoute()
const sidebarStore = useSidebarStore()
const themeStore = useThemeStore()

const searchDialogVisible = ref(false)
const userMenu = ref()
const notificationMenu = ref()

// 当前页面标题
const currentPageTitle = computed(() => {
   return (route.meta?.title as string) || '仪表盘'
})

const userMenuItems = ref<MenuItem[]>([
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
      command: () => handleLogout(),
   },
])

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
])

const unreadCount = computed(() => notifications.value.filter(n => n.unread).length)

const handleLogout = () => {
   localStorage.clear()
   router.replace('/auth/login')
}

const toggleUserMenu = (event: Event) => {
   userMenu.value.toggle(event)
}

const toggleNotifications = (event: Event) => {
   notificationMenu.value.toggle(event)
}

const sidebarWidth = computed(() => (sidebarStore.isCollapsed ? '5rem' : '17rem'))
</script>

<template>
   <header
      class="fixed top-0 right-0 h-16 flex items-center justify-between px-6 bg-white/80 dark:bg-[rgba(24,24,27,0.85)] backdrop-blur-xl border-b border-surface-100 dark:border-surface-800 z-900 transition-[left] duration-300 ease-in-out max-lg:left-0!"
      :style="{ left: sidebarWidth }"
   >
      <div class="flex items-center gap-3">
         <!-- 移动端菜单按钮 -->
         <Button
            icon="pi pi-bars"
            severity="secondary"
            text
            rounded
            class="hidden! max-lg:flex!"
            @click="sidebarStore.toggleMobile"
         />

         <!-- 折叠按钮 -->
         <Button
            :icon="sidebarStore.isCollapsed ? 'pi pi-angle-right' : 'pi pi-angle-left'"
            severity="secondary"
            text
            rounded
            class="text-surface-500 max-lg:hidden!"
            @click="sidebarStore.toggleCollapsed"
            v-tooltip.bottom="sidebarStore.isCollapsed ? '展开菜单' : '折叠菜单'"
         />

         <!-- 面包屑 -->
         <nav class="flex items-center gap-2 text-sm text-surface-500 max-lg:hidden! shrink-0">
            <span class="flex items-center text-surface-400">
               <i class="pi pi-home"></i>
            </span>
            <span class="text-surface-300">/</span>
            <span class="text-surface-900 dark:text-surface-100 font-medium">
               {{ currentPageTitle }}
            </span>
         </nav>
      </div>

      <div class="flex-1 max-w-lg mx-6 max-lg:hidden!">
         <!-- 搜索触发按钮 -->
         <button
            class="w-full flex items-center gap-3 px-4 h-10 rounded-xl bg-surface-50 dark:bg-surface-800 border border-surface-200 dark:border-surface-700 text-sm cursor-pointer transition-all duration-200 ease group hover:bg-surface-100 dark:hover:bg-surface-700 hover:border-surface-300 dark:hover:border-surface-600"
            @click="searchDialogVisible = true"
         >
            <i class="pi pi-search text-surface-400 group-hover:text-surface-500"></i>
            <span class="flex-1 text-left text-surface-400 group-hover:text-surface-500">
               搜索用户、应用、设置...
            </span>
            <span class="flex gap-1">
               <kbd
                  class="py-0.5 px-1.5 rounded bg-surface-100 dark:bg-surface-700 border border-surface-200 dark:border-surface-600 text-xs font-inherit text-surface-500 dark:text-surface-400"
               >
                  ⌘
               </kbd>
               <kbd
                  class="py-0.5 px-1.5 rounded bg-surface-100 dark:bg-surface-700 border border-surface-200 dark:border-surface-600 text-xs font-inherit text-surface-500 dark:text-surface-400"
               >
                  K
               </kbd>
            </span>
         </button>
      </div>

      <!-- 移动端搜索按钮 -->
      <Button
         icon="pi pi-search"
         severity="secondary"
         text
         rounded
         class="text-surface-600! dark:text-surface-400! hidden! max-lg:flex!"
         @click="searchDialogVisible = true"
         v-tooltip.bottom="'搜索'"
      />

      <div class="flex items-center gap-3">
         <!-- 主题切换 -->
         <Button
            :icon="themeStore.isDark ? 'pi pi-sun' : 'pi pi-moon'"
            severity="secondary"
            text
            rounded
            class="text-surface-600! dark:text-surface-400!"
            @click="themeStore.toggleTheme"
            v-tooltip.bottom="themeStore.isDark ? '切换到亮色模式' : '切换到暗色模式'"
         />

         <!-- 通知 -->
         <div class="relative">
            <Button
               icon="pi pi-bell"
               severity="secondary"
               text
               rounded
               class="text-surface-600! dark:text-surface-400!"
               @click="toggleNotifications"
               v-tooltip.bottom="'通知'"
            />
            <Badge
               v-if="unreadCount > 0"
               :value="unreadCount"
               severity="danger"
               class="absolute! top-0! right-0! translate-x-1/4! -translate-y-1/4! min-w-4.5! h-4.5! text-[0.625rem]!"
            />
         </div>

         <!-- 用户头像 -->
         <button
            class="flex items-center gap-2.5 py-1.5 pr-3 pl-1.5 border-none rounded-full bg-surface-50 dark:bg-surface-800 cursor-pointer transition-all duration-200 ease hover:bg-surface-100 dark:hover:bg-surface-700"
            @click="toggleUserMenu"
         >
            <Avatar
               image="https://api.dicebear.com/7.x/avataaars/svg?seed=admin"
               shape="circle"
               class="w-8! h-8! border-2! border-primary-200! dark:border-primary-700!"
            />
            <span class="text-sm font-medium text-surface-700 dark:text-surface-200 max-lg:hidden!">
               管理员
            </span>
            <i class="pi pi-chevron-down text-xs text-surface-400"></i>
         </button>

         <Menu ref="userMenu" :model="userMenuItems" popup class="min-w-48" />

         <!-- 通知面板 -->
         <Menu ref="notificationMenu" popup class="w-88! max-w-[calc(100vw-2rem)]!">
            <template #start>
               <div
                  class="flex items-center justify-between p-4 border-b border-surface-100 dark:border-surface-700"
               >
                  <span class="font-semibold text-surface-900 dark:text-surface-100"> 通知 </span>
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
                        'bg-primary-50! dark:bg-[rgba(251,146,60,0.08)]!': notification.unread,
                     }"
                  >
                     <div
                        v-if="notification.unread"
                        class="shrink-0 w-2 h-2 mt-1.5 rounded-full bg-primary-500"
                     ></div>
                     <div class="flex flex-col gap-1">
                        <span class="text-sm font-semibold text-surface-900 dark:text-surface-100">
                           {{ notification.title }}
                        </span>
                        <span class="text-[0.8125rem] text-surface-600 dark:text-surface-400">
                           {{ notification.message }}
                        </span>
                        <span class="text-xs text-surface-400">
                           {{ notification.time }}
                        </span>
                     </div>
                  </div>
               </div>
               <div class="p-3 border-t border-surface-100 dark:border-surface-700">
                  <Button
                     label="查看全部通知"
                     text
                     class="w-full!"
                     @click="router.push('/notifications')"
                  />
               </div>
            </template>
         </Menu>
      </div>
   </header>

   <!-- 全局搜索弹窗 -->
   <GlobalSearchDialog v-model:visible="searchDialogVisible" />
</template>
