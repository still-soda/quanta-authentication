<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import Dialog from 'primevue/dialog';
import Button from 'primevue/button';

const props = defineProps<{
   visible: boolean;
}>();

const emit = defineEmits<{
   'update:visible': [value: boolean];
}>();

const router = useRouter();
const searchQuery = ref('');
const selectedIndex = ref(0);
const searchInputRef = ref<HTMLInputElement | null>(null);

// 搜索数据源
interface SearchItem {
   id: string;
   label: string;
   description?: string;
   icon: string;
   category: 'navigation' | 'action' | 'user' | 'app' | 'recent';
   action: () => void;
   keywords?: string[];
}

// 页面导航项
const navigationItems: SearchItem[] = [
   {
      id: 'nav-dashboard',
      label: '仪表盘',
      description: '查看系统概览和统计数据',
      icon: 'pi pi-chart-line',
      category: 'navigation',
      action: () => navigateTo('/'),
      keywords: ['首页', 'home', 'dashboard', '统计', '概览'],
   },
   {
      id: 'nav-users',
      label: '用户管理',
      description: '管理系统用户账号',
      icon: 'pi pi-users',
      category: 'navigation',
      action: () => navigateTo('/users'),
      keywords: ['user', '账号', '成员', 'member'],
   },
   {
      id: 'nav-roles',
      label: '角色权限',
      description: '配置角色和权限策略',
      icon: 'pi pi-shield',
      category: 'navigation',
      action: () => navigateTo('/roles'),
      keywords: ['role', 'permission', '权限', 'rbac'],
   },
   {
      id: 'nav-oauth',
      label: 'OAuth 应用',
      description: '管理 OAuth2.0 客户端应用',
      icon: 'pi pi-box',
      category: 'navigation',
      action: () => navigateTo('/oauth'),
      keywords: ['client', '客户端', '应用', 'app', 'oauth2'],
   },
   {
      id: 'nav-organizations',
      label: '组织架构',
      description: '管理组织和部门结构',
      icon: 'pi pi-building',
      category: 'navigation',
      action: () => navigateTo('/organizations'),
      keywords: ['org', '部门', 'department', '公司'],
   },
   {
      id: 'nav-audit',
      label: '审计日志',
      description: '查看系统操作记录',
      icon: 'pi pi-history',
      category: 'navigation',
      action: () => navigateTo('/audit'),
      keywords: ['log', '日志', '记录', 'history'],
   },
   {
      id: 'nav-settings',
      label: '系统设置',
      description: '配置系统参数和选项',
      icon: 'pi pi-cog',
      category: 'navigation',
      action: () => navigateTo('/settings'),
      keywords: ['config', '配置', '参数', 'options'],
   },
   {
      id: 'nav-profile',
      label: '个人资料',
      description: '管理您的账号信息',
      icon: 'pi pi-user',
      category: 'navigation',
      action: () => navigateTo('/profile'),
      keywords: ['account', '账户', 'profile', '信息'],
   },
   {
      id: 'nav-notifications',
      label: '通知中心',
      description: '查看所有通知消息',
      icon: 'pi pi-bell',
      category: 'navigation',
      action: () => navigateTo('/notifications'),
      keywords: ['message', '消息', 'notification'],
   },
];

// 快捷操作项
const actionItems: SearchItem[] = [
   {
      id: 'action-add-user',
      label: '添加用户',
      description: '创建新的系统用户',
      icon: 'pi pi-user-plus',
      category: 'action',
      action: () => {
         navigateTo('/users');
         // 这里可以触发打开添加用户对话框
      },
      keywords: ['create', '新建', '创建', 'add'],
   },
   {
      id: 'action-add-role',
      label: '新建角色',
      description: '创建新的角色权限组',
      icon: 'pi pi-plus-circle',
      category: 'action',
      action: () => {
         navigateTo('/roles');
      },
      keywords: ['create', '新建', '创建', 'add'],
   },
   {
      id: 'action-add-app',
      label: '注册应用',
      description: '注册新的 OAuth 应用',
      icon: 'pi pi-plus',
      category: 'action',
      action: () => {
         navigateTo('/oauth');
      },
      keywords: ['create', '新建', '创建', 'add', 'register'],
   },
   {
      id: 'action-export-users',
      label: '导出用户',
      description: '导出用户数据到 Excel',
      icon: 'pi pi-download',
      category: 'action',
      action: () => {
         console.log('Export users');
         closeDialog();
      },
      keywords: ['export', '下载', 'download', 'excel'],
   },
   {
      id: 'action-refresh',
      label: '刷新页面',
      description: '重新加载当前页面数据',
      icon: 'pi pi-refresh',
      category: 'action',
      action: () => {
         window.location.reload();
      },
      keywords: ['reload', '重新加载', 'refresh'],
   },
];

// 模拟用户数据
const mockUsers: SearchItem[] = [
   {
      id: 'user-1',
      label: 'admin@example.com',
      description: '系统管理员',
      icon: 'pi pi-user',
      category: 'user',
      action: () => {
         navigateTo('/users');
         // 可以传递查询参数
      },
      keywords: ['admin', '管理员'],
   },
   {
      id: 'user-2',
      label: 'john.doe@example.com',
      description: 'John Doe - 普通用户',
      icon: 'pi pi-user',
      category: 'user',
      action: () => navigateTo('/users'),
      keywords: ['john', 'doe'],
   },
   {
      id: 'user-3',
      label: 'jane.smith@example.com',
      description: 'Jane Smith - 运维人员',
      icon: 'pi pi-user',
      category: 'user',
      action: () => navigateTo('/users'),
      keywords: ['jane', 'smith', '运维'],
   },
];

// 模拟应用数据
const mockApps: SearchItem[] = [
   {
      id: 'app-1',
      label: 'Quanta Dashboard',
      description: 'Client ID: qd-dashboard-001',
      icon: 'pi pi-desktop',
      category: 'app',
      action: () => navigateTo('/oauth'),
      keywords: ['dashboard', '仪表盘'],
   },
   {
      id: 'app-2',
      label: 'Mobile App',
      description: 'Client ID: qm-mobile-002',
      icon: 'pi pi-mobile',
      category: 'app',
      action: () => navigateTo('/oauth'),
      keywords: ['mobile', '移动端', '手机'],
   },
   {
      id: 'app-3',
      label: 'API Gateway',
      description: 'Client ID: qg-api-003',
      icon: 'pi pi-server',
      category: 'app',
      action: () => navigateTo('/oauth'),
      keywords: ['api', 'gateway', '网关'],
   },
];

// 最近搜索历史
const recentSearches = ref<SearchItem[]>([]);
const MAX_RECENT_SEARCHES = 5;

// 加载最近搜索
const loadRecentSearches = () => {
   try {
      const saved = localStorage.getItem('qauth-recent-searches');
      if (saved) {
         const ids = JSON.parse(saved) as string[];
         const allItems = [
            ...navigationItems,
            ...actionItems,
            ...mockUsers,
            ...mockApps,
         ];
         recentSearches.value = ids
            .map((id) => allItems.find((item) => item.id === id))
            .filter(Boolean) as SearchItem[];
      }
   } catch {
      recentSearches.value = [];
   }
};

// 保存最近搜索
const saveRecentSearch = (item: SearchItem) => {
   const ids = recentSearches.value.map((i) => i.id);
   const newIds = [
      item.id,
      ...ids.filter((id) => id !== item.id),
   ].slice(0, MAX_RECENT_SEARCHES);

   localStorage.setItem('qauth-recent-searches', JSON.stringify(newIds));

   const allItems = [
      ...navigationItems,
      ...actionItems,
      ...mockUsers,
      ...mockApps,
   ];
   recentSearches.value = newIds
      .map((id) => allItems.find((i) => i.id === id))
      .filter(Boolean) as SearchItem[];
};

// 清除搜索历史
const clearRecentSearches = () => {
   localStorage.removeItem('qauth-recent-searches');
   recentSearches.value = [];
};

// 模糊搜索函数
const fuzzyMatch = (text: string, query: string): boolean => {
   const lowerText = text.toLowerCase();
   const lowerQuery = query.toLowerCase();

   // 简单包含匹配
   if (lowerText.includes(lowerQuery)) return true;

   // 首字母匹配
   const words = lowerText.split(/\s+/);
   const initials = words.map((w) => w[0]).join('');
   if (initials.includes(lowerQuery)) return true;

   return false;
};

// 搜索结果
const searchResults = computed(() => {
   const query = searchQuery.value.trim();

   if (!query) {
      // 无搜索词时显示最近搜索
      if (recentSearches.value.length > 0) {
         return [
            {
               category: 'recent' as const,
               label: '最近搜索',
               items: recentSearches.value.map((item) => ({
                  ...item,
                  category: 'recent' as const,
               })),
            },
         ];
      }
      // 显示默认推荐
      return [
         {
            category: 'navigation' as const,
            label: '快速导航',
            items: navigationItems.slice(0, 5),
         },
         {
            category: 'action' as const,
            label: '快捷操作',
            items: actionItems.slice(0, 3),
         },
      ];
   }

   // 搜索所有数据源
   const allItems = [
      ...navigationItems,
      ...actionItems,
      ...mockUsers,
      ...mockApps,
   ];

   const matched = allItems.filter((item) => {
      if (fuzzyMatch(item.label, query)) return true;
      if (item.description && fuzzyMatch(item.description, query)) return true;
      if (item.keywords?.some((kw) => fuzzyMatch(kw, query))) return true;
      return false;
   });

   // 按类别分组
   const groups: {
      category: 'navigation' | 'action' | 'user' | 'app';
      label: string;
      items: SearchItem[];
   }[] = [];

   const navResults = matched.filter((i) => i.category === 'navigation');
   const actionResults = matched.filter((i) => i.category === 'action');
   const userResults = matched.filter((i) => i.category === 'user');
   const appResults = matched.filter((i) => i.category === 'app');

   if (navResults.length)
      groups.push({ category: 'navigation', label: '页面', items: navResults });
   if (actionResults.length)
      groups.push({ category: 'action', label: '操作', items: actionResults });
   if (userResults.length)
      groups.push({ category: 'user', label: '用户', items: userResults });
   if (appResults.length)
      groups.push({ category: 'app', label: '应用', items: appResults });

   return groups;
});

// 扁平化的搜索结果列表
const flatResults = computed(() => {
   return searchResults.value.flatMap((group) => group.items);
});

// 高亮匹配文字
const highlightMatch = (text: string, query: string): string => {
   if (!query.trim()) return text;

   const regex = new RegExp(`(${escapeRegExp(query)})`, 'gi');
   return text.replace(
      regex,
      '<mark class="bg-primary-100 dark:bg-primary-900/40 text-primary-700 dark:text-primary-300 rounded px-0.5">$1</mark>',
   );
};

const escapeRegExp = (str: string): string => {
   return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
};

// 导航函数
const navigateTo = (path: string) => {
   router.push(path);
   closeDialog();
};

// 关闭弹窗
const closeDialog = () => {
   emit('update:visible', false);
   searchQuery.value = '';
   selectedIndex.value = 0;
};

// 执行选中项
const executeSelected = () => {
   const item = flatResults.value[selectedIndex.value];
   if (item) {
      saveRecentSearch(item);
      item.action();
   }
};

// 处理点击项
const handleItemClick = (item: SearchItem) => {
   saveRecentSearch(item);
   item.action();
};

// 键盘导航
const handleKeydown = (e: KeyboardEvent) => {
   const total = flatResults.value.length;
   if (!total) return;

   switch (e.key) {
      case 'ArrowDown':
         e.preventDefault();
         selectedIndex.value = (selectedIndex.value + 1) % total;
         scrollToSelected();
         break;
      case 'ArrowUp':
         e.preventDefault();
         selectedIndex.value = (selectedIndex.value - 1 + total) % total;
         scrollToSelected();
         break;
      case 'Enter':
         e.preventDefault();
         executeSelected();
         break;
      case 'Escape':
         closeDialog();
         break;
   }
};

// 滚动到选中项
const scrollToSelected = () => {
   nextTick(() => {
      const selected = document.querySelector('.search-item-selected');
      selected?.scrollIntoView({ block: 'nearest', behavior: 'smooth' });
   });
};

// 重置选中索引
watch(searchQuery, () => {
   selectedIndex.value = 0;
});

// 打开时聚焦输入框
watch(
   () => props.visible,
   (visible) => {
      if (visible) {
         loadRecentSearches();
         nextTick(() => {
            searchInputRef.value?.focus();
         });
      }
   },
);

// 全局快捷键
const handleGlobalKeydown = (e: KeyboardEvent) => {
   if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
      e.preventDefault();
      emit('update:visible', true);
   }
};

onMounted(() => {
   loadRecentSearches();
   window.addEventListener('keydown', handleGlobalKeydown);
});

onUnmounted(() => {
   window.removeEventListener('keydown', handleGlobalKeydown);
});

// 获取类别图标
const getCategoryIcon = (
   category: 'navigation' | 'action' | 'user' | 'app' | 'recent',
): string => {
   const icons = {
      navigation: 'pi pi-compass',
      action: 'pi pi-bolt',
      user: 'pi pi-users',
      app: 'pi pi-box',
      recent: 'pi pi-clock',
   };
   return icons[category];
};

// 获取类别颜色类
const getCategoryColorClass = (
   category: 'navigation' | 'action' | 'user' | 'app' | 'recent',
): string => {
   const colors = {
      navigation: 'text-blue-500 bg-blue-50 dark:bg-blue-900/20',
      action: 'text-amber-500 bg-amber-50 dark:bg-amber-900/20',
      user: 'text-emerald-500 bg-emerald-50 dark:bg-emerald-900/20',
      app: 'text-violet-500 bg-violet-50 dark:bg-violet-900/20',
      recent: 'text-surface-500 bg-surface-100 dark:bg-surface-700',
   };
   return colors[category];
};
</script>

<template>
   <Dialog
      :visible="visible"
      modal
      dismissable-mask
      :draggable="false"
      class="w-[95vw] max-w-2xl! m-4"
      :pt="{
         root: { class: 'border-0! rounded-2xl! overflow-hidden!' },
         header: { class: 'hidden!' },
         content: { class: 'p-0!' },
         mask: {
            class: 'backdrop-blur-sm! bg-surface-900/20! dark:bg-surface-950/40!',
         },
      }"
      @update:visible="emit('update:visible', $event)">
      <div
         class="bg-white dark:bg-surface-900 rounded-2xl overflow-hidden shadow-2xl border border-surface-200 dark:border-surface-700">
         <!-- 搜索头部 -->
         <div
            class="flex items-center gap-3 px-5 py-4 border-b border-surface-100 dark:border-surface-800">
            <i class="pi pi-search text-lg text-surface-400"></i>
            <input
               ref="searchInputRef"
               v-model="searchQuery"
               type="text"
               placeholder="搜索页面、用户、应用、操作..."
               class="flex-1 bg-transparent border-none outline-none text-base text-surface-900 dark:text-surface-100 placeholder:text-surface-400"
               @keydown="handleKeydown" />
            <div class="flex items-center gap-1.5">
               <kbd
                  class="py-1 px-2 rounded-lg bg-surface-100 dark:bg-surface-800 text-xs text-surface-500 border border-surface-200 dark:border-surface-700">
                  ESC
               </kbd>
            </div>
         </div>

         <!-- 搜索结果 -->
         <div
            class="max-h-[60vh] overflow-y-auto overscroll-contain"
            style="scrollbar-gutter: stable">
            <template v-if="searchResults.length">
               <div
                  v-for="(group, groupIndex) in searchResults"
                  :key="group.category + groupIndex"
                  class="py-2">
                  <!-- 分组标题 -->
                  <div
                     class="flex items-center justify-between px-5 py-2 text-xs font-semibold text-surface-500 dark:text-surface-400 uppercase tracking-wider">
                     <div class="flex items-center gap-2">
                        <i
                           :class="getCategoryIcon(group.category)"
                           class="text-[0.7rem]"></i>
                        <span>{{ group.label }}</span>
                     </div>
                     <Button
                        v-if="group.category === 'recent'"
                        label="清除"
                        text
                        size="small"
                        class="text-xs! py-0! px-2! h-auto!"
                        @click.stop="clearRecentSearches" />
                  </div>

                  <!-- 搜索项 -->
                  <div
                     v-for="item in group.items"
                     :key="item.id"
                     class="search-item flex items-center gap-3 mx-2 px-3 py-2.5 rounded-xl cursor-pointer transition-all duration-150"
                     :class="[
                        flatResults.findIndex((i) => i.id === item.id) ===
                        selectedIndex
                           ? 'search-item-selected bg-primary-50 dark:bg-primary-900/20 border border-primary-200 dark:border-primary-800'
                           : 'border border-transparent hover:bg-surface-50 dark:hover:bg-surface-800',
                     ]"
                     @click="handleItemClick(item)"
                     @mouseenter="
                        selectedIndex = flatResults.findIndex(
                           (i) => i.id === item.id,
                        )
                     ">
                     <!-- 图标 -->
                     <div
                        class="flex items-center justify-center w-10 h-10 rounded-xl"
                        :class="getCategoryColorClass(item.category)">
                        <i :class="item.icon" class="text-base"></i>
                     </div>

                     <!-- 内容 -->
                     <div class="flex-1 min-w-0">
                        <div
                           class="text-sm font-medium text-surface-900 dark:text-surface-100 truncate"
                           v-html="highlightMatch(item.label, searchQuery)">
                        </div>
                        <div
                           v-if="item.description"
                           class="text-xs text-surface-500 dark:text-surface-400 truncate mt-0.5"
                           v-html="
                              highlightMatch(item.description, searchQuery)
                           ">
                        </div>
                     </div>

                     <!-- 快捷键提示 -->
                     <div
                        v-if="
                           flatResults.findIndex((i) => i.id === item.id) ===
                           selectedIndex
                        "
                        class="flex items-center gap-1">
                        <kbd
                           class="py-0.5 px-1.5 rounded bg-surface-100 dark:bg-surface-700 text-[0.65rem] text-surface-500 border border-surface-200 dark:border-surface-600">
                           ↵
                        </kbd>
                     </div>
                  </div>
               </div>
            </template>

            <!-- 无结果 -->
            <div
               v-else
               class="flex flex-col items-center justify-center py-12 text-center">
               <div
                  class="w-16 h-16 rounded-2xl bg-surface-100 dark:bg-surface-800 flex items-center justify-center mb-4">
                  <i
                     class="pi pi-search text-2xl text-surface-400 dark:text-surface-500"></i>
               </div>
               <p
                  class="text-sm font-medium text-surface-700 dark:text-surface-300">
                  未找到相关结果
               </p>
               <p class="text-xs text-surface-500 mt-1">
                  尝试使用不同的关键词搜索
               </p>
            </div>
         </div>

         <!-- 底部提示 -->
         <div
            class="flex items-center justify-between px-5 py-3 bg-surface-50 dark:bg-surface-800/50 border-t border-surface-100 dark:border-surface-800">
            <div class="flex items-center gap-4 text-xs text-surface-500">
               <span class="flex items-center gap-1.5">
                  <kbd
                     class="py-0.5 px-1 rounded bg-surface-200 dark:bg-surface-700 text-[0.65rem]">
                     ↑
                  </kbd>
                  <kbd
                     class="py-0.5 px-1 rounded bg-surface-200 dark:bg-surface-700 text-[0.65rem]">
                     ↓
                  </kbd>
                  <span>选择</span>
               </span>
               <span class="flex items-center gap-1.5">
                  <kbd
                     class="py-0.5 px-1 rounded bg-surface-200 dark:bg-surface-700 text-[0.65rem]">
                     ↵
                  </kbd>
                  <span>执行</span>
               </span>
               <span class="flex items-center gap-1.5">
                  <kbd
                     class="py-0.5 px-1 rounded bg-surface-200 dark:bg-surface-700 text-[0.65rem]">
                     esc
                  </kbd>
                  <span>关闭</span>
               </span>
            </div>
            <div class="text-xs text-surface-400">
               {{ flatResults.length }} 个结果
            </div>
         </div>
      </div>
   </Dialog>
</template>
