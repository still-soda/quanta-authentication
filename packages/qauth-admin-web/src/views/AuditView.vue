<script setup lang="ts">
import { ref, computed } from 'vue';
import Button from 'primevue/button';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Tag from 'primevue/tag';
import InputText from 'primevue/inputtext';
import Select from 'primevue/select';
import DatePicker from 'primevue/datepicker';
import Dialog from 'primevue/dialog';
import PageHeader from '@/components/shared/PageHeader.vue';
import SimpleStatCard from '@/components/shared/SimpleStatCard.vue';
import type { AuditLog, SimpleStatData } from '@/types';
import { AUDIT_MODULE_ICONS } from '@/config';

// 模拟审计日志数据
const auditLogs = ref<AuditLog[]>([
   {
      id: '1',
      operatorId: 'u1',
      operatorName: '张伟',
      operatorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhang',
      module: '用户管理',
      action: '创建用户',
      targetId: 'u15',
      targetName: '新用户 - 李四',
      detail: { email: 'lisi@example.com', role: '普通用户' },
      ip: '192.168.1.100',
      time: '2026-01-23 14:30:25',
      durationMs: 125,
      status: 'success',
   },
   {
      id: '2',
      operatorId: 'u1',
      operatorName: '张伟',
      operatorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhang',
      module: 'OAuth应用',
      action: '重置密钥',
      targetId: 'app3',
      targetName: '内部管理系统',
      detail: { reason: '安全检查' },
      ip: '192.168.1.100',
      time: '2026-01-23 14:15:00',
      durationMs: 89,
      status: 'success',
   },
   {
      id: '3',
      operatorId: 'u2',
      operatorName: '李明',
      operatorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=li',
      module: '角色权限',
      action: '修改权限',
      targetId: 'role2',
      targetName: '开发者角色',
      detail: { added: ['oauth:create'], removed: [] },
      ip: '10.0.0.55',
      time: '2026-01-23 13:45:12',
      durationMs: 234,
      status: 'success',
   },
   {
      id: '4',
      operatorId: 'u3',
      operatorName: '王芳',
      operatorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=wang',
      module: '用户管理',
      action: '禁用用户',
      targetId: 'u8',
      targetName: '异常账户 - 测试',
      detail: { reason: '违规操作' },
      ip: '192.168.2.88',
      time: '2026-01-23 12:20:45',
      durationMs: 67,
      status: 'warning',
   },
   {
      id: '5',
      operatorId: 'u1',
      operatorName: '张伟',
      operatorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhang',
      module: '系统设置',
      action: '修改配置',
      targetId: 'config',
      targetName: '登录安全策略',
      detail: { maxAttempts: 5, lockDuration: 30 },
      ip: '192.168.1.100',
      time: '2026-01-23 11:00:00',
      durationMs: 156,
      status: 'success',
   },
   {
      id: '6',
      operatorId: 'u4',
      operatorName: '陈红',
      operatorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=chen',
      module: 'OAuth应用',
      action: '创建应用',
      targetId: 'app12',
      targetName: '新财务系统',
      detail: { type: 'web', redirectUri: 'https://finance.example.com/callback' },
      ip: '172.16.0.22',
      time: '2026-01-23 10:30:18',
      durationMs: 312,
      status: 'success',
   },
   {
      id: '7',
      operatorId: 'u2',
      operatorName: '李明',
      operatorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=li',
      module: '用户管理',
      action: '删除用户',
      targetId: 'u99',
      targetName: '已离职 - 张三',
      detail: { reason: '离职清退' },
      ip: '10.0.0.55',
      time: '2026-01-22 18:45:30',
      durationMs: 198,
      status: 'error',
   },
   {
      id: '8',
      operatorId: 'u5',
      operatorName: '赵阳',
      operatorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhaoyang',
      module: '认证登录',
      action: '登录失败',
      targetId: 'u5',
      targetName: '赵阳',
      detail: { error: '密码错误', attempts: 3 },
      ip: '203.0.113.45',
      time: '2026-01-22 16:20:00',
      durationMs: 45,
      status: 'error',
   },
   {
      id: '9',
      operatorId: 'u1',
      operatorName: '张伟',
      operatorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhang',
      module: '角色权限',
      action: '创建角色',
      targetId: 'role8',
      targetName: '访客角色',
      detail: { permissions: ['read:basic'] },
      ip: '192.168.1.100',
      time: '2026-01-22 14:10:22',
      durationMs: 178,
      status: 'success',
   },
   {
      id: '10',
      operatorId: 'u3',
      operatorName: '王芳',
      operatorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=wang',
      module: '数据导出',
      action: '导出用户列表',
      targetId: 'export1',
      targetName: '用户数据导出',
      detail: { format: 'xlsx', count: 1024 },
      ip: '192.168.2.88',
      time: '2026-01-22 11:30:00',
      durationMs: 2456,
      status: 'success',
   },
]);

// 筛选状态
const filters = ref({
   search: '',
   module: null as string | null,
   action: null as string | null,
   dateRange: null as [Date, Date] | null,
});

// 模块选项
const moduleOptions = [
   { label: '全部模块', value: null },
   { label: '用户管理', value: '用户管理' },
   { label: 'OAuth应用', value: 'OAuth应用' },
   { label: '角色权限', value: '角色权限' },
   { label: '系统设置', value: '系统设置' },
   { label: '认证登录', value: '认证登录' },
   { label: '数据导出', value: '数据导出' },
];

// 操作选项
const actionOptions = [
   { label: '全部操作', value: null },
   { label: '创建', value: '创建' },
   { label: '修改', value: '修改' },
   { label: '删除', value: '删除' },
   { label: '禁用', value: '禁用' },
   { label: '登录', value: '登录' },
   { label: '导出', value: '导出' },
];

// 统计数据
const stats = computed<SimpleStatData[]>(() => [
   {
      title: '今日操作',
      value: auditLogs.value.filter((log) =>
         log.time.startsWith('2026-01-23'),
      ).length,
      icon: 'pi pi-history',
      color: 'blue',
   },
   {
      title: '成功操作',
      value: auditLogs.value.filter((log) => log.status === 'success').length,
      icon: 'pi pi-check-circle',
      color: 'green',
   },
   {
      title: '警告事件',
      value: auditLogs.value.filter((log) => log.status === 'warning').length,
      icon: 'pi pi-exclamation-triangle',
      color: 'orange',
   },
   {
      title: '失败操作',
      value: auditLogs.value.filter((log) => log.status === 'error').length,
      icon: 'pi pi-times-circle',
      color: 'red',
   },
]);

// 过滤后的日志
const filteredLogs = computed(() => {
   return auditLogs.value.filter((log) => {
      // 搜索过滤
      if (filters.value.search) {
         const search = filters.value.search.toLowerCase();
         const matchesSearch =
            log.operatorName.toLowerCase().includes(search) ||
            log.targetName?.toLowerCase().includes(search) ||
            log.action.toLowerCase().includes(search) ||
            log.ip.includes(search);
         if (!matchesSearch) return false;
      }
      // 模块过滤
      if (filters.value.module && log.module !== filters.value.module) {
         return false;
      }
      // 操作类型过滤
      if (filters.value.action && !log.action.includes(filters.value.action)) {
         return false;
      }
      return true;
   });
});

const STATUS_CONFIG = { success: { label: '成功', severity: 'success' }, warning: { label: '警告', severity: 'warn' }, error: { label: '失败', severity: 'danger' } } as const;
const getStatusSeverity = (status: string) => STATUS_CONFIG[status as keyof typeof STATUS_CONFIG]?.severity || 'secondary';
const getStatusLabel = (status: string) => STATUS_CONFIG[status as keyof typeof STATUS_CONFIG]?.label || status;
const getModuleIcon = (module: string) => AUDIT_MODULE_ICONS[module] || 'pi pi-circle';

const detailDialog = ref(false);
const selectedLog = ref<AuditLog | null>(null);
const showDetail = (log: AuditLog) => { selectedLog.value = log; detailDialog.value = true; };
const exportLogs = () => console.log('Exporting logs...');
const clearFilters = () => { filters.value = { search: '', module: null, action: null, dateRange: null }; };
</script>

<template>
   <div class="flex flex-col gap-6">
      <!-- Page Header -->
      <PageHeader title="审计日志" subtitle="追踪和审查系统中的所有操作记录">
         <template #actions>
            <Button
               label="导出日志"
               icon="pi pi-download"
               severity="secondary"
               outlined
               @click="exportLogs" />
            <Button
               label="清除筛选"
               icon="pi pi-filter-slash"
               :disabled="
                  !filters.search &&
                  !filters.module &&
                  !filters.action &&
                  !filters.dateRange
               "
               @click="clearFilters" />
         </template>
      </PageHeader>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4">
         <SimpleStatCard v-for="stat in stats" :key="stat.title" :stat="stat" />
      </div>

      <!-- Filters -->
      <div
         class="flex flex-wrap items-center gap-4 p-4 bg-surface-0 dark:bg-surface-900 border border-surface-200 dark:border-surface-800 rounded-xl">
         <div class="flex-1 min-w-60 relative">
            <i class="pi pi-search absolute left-3 top-1/2 -translate-y-1/2 text-surface-400 z-10"></i>
            <InputText
               v-model="filters.search"
               placeholder="搜索操作者、目标、IP..."
               class="w-full pl-10" />
         </div>
         <Select
            v-model="filters.module"
            :options="moduleOptions"
            optionLabel="label"
            optionValue="value"
            placeholder="选择模块"
            :pt="{ root: { style: 'min-width: 140px' } }" />
         <Select
            v-model="filters.action"
            :options="actionOptions"
            optionLabel="label"
            optionValue="value"
            placeholder="操作类型"
            :pt="{ root: { style: 'min-width: 130px' } }" />
         <DatePicker
            v-model="filters.dateRange"
            selectionMode="range"
            :manualInput="false"
            placeholder="选择日期范围"
            :pt="{ root: { style: 'min-width: 220px' } }" />
      </div>

      <!-- Audit Logs Table -->
      <div
         class="bg-surface-0 dark:bg-surface-900 border border-surface-200 dark:border-surface-800 rounded-xl overflow-hidden">
         <DataTable
            :value="filteredLogs"
            :paginator="true"
            :rows="10"
            :rowsPerPageOptions="[10, 20, 50]"
            paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
            responsiveLayout="scroll"
            :pt="{
               root: { class: 'border-none' },
               header: { class: 'bg-transparent border-none px-5 py-4' },
               thead: { class: 'bg-surface-50 dark:bg-surface-800' },
               tbody: { class: 'divide-y divide-surface-100 dark:divide-surface-800' },
               paginator: {
                  class: 'border-t border-surface-200 dark:border-surface-700 px-4 py-3',
               },
            }">
            <template #empty>
               <div class="flex flex-col items-center justify-center py-12 text-surface-500">
                  <i class="pi pi-inbox text-4xl mb-3 opacity-50"></i>
                  <span>暂无审计日志</span>
               </div>
            </template>

            <Column header="操作者" style="min-width: 180px">
               <template #body="{ data }">
                  <div class="flex items-center gap-3">
                     <div
                        class="w-9 h-9 rounded-lg overflow-hidden bg-surface-100 dark:bg-surface-700 shrink-0">
                        <img
                           :src="data.operatorAvatar"
                           :alt="data.operatorName"
                           class="w-full h-full object-cover" />
                     </div>
                     <div class="flex flex-col">
                        <span
                           class="font-medium text-surface-900 dark:text-surface-100 text-sm">
                           {{ data.operatorName }}
                        </span>
                        <span class="text-xs text-surface-500">{{
                           data.ip
                        }}</span>
                     </div>
                  </div>
               </template>
            </Column>

            <Column header="模块" style="min-width: 120px">
               <template #body="{ data }">
                  <div class="flex items-center gap-2">
                     <i
                        :class="getModuleIcon(data.module)"
                        class="text-surface-500"></i>
                     <span class="text-sm text-surface-700 dark:text-surface-300">
                        {{ data.module }}
                     </span>
                  </div>
               </template>
            </Column>

            <Column header="操作" style="min-width: 100px">
               <template #body="{ data }">
                  <span
                     class="font-medium text-sm text-surface-900 dark:text-surface-100">
                     {{ data.action }}
                  </span>
               </template>
            </Column>

            <Column header="目标" style="min-width: 180px">
               <template #body="{ data }">
                  <span
                     v-if="data.targetName"
                     class="text-sm text-surface-600 dark:text-surface-400">
                     {{ data.targetName }}
                  </span>
                  <span v-else class="text-sm text-surface-400">-</span>
               </template>
            </Column>

            <Column header="状态" style="min-width: 90px">
               <template #body="{ data }">
                  <Tag
                     :severity="getStatusSeverity(data.status)"
                     :value="getStatusLabel(data.status)"
                     :pt="{
                        root: { class: 'text-xs px-2 py-1' },
                     }" />
               </template>
            </Column>

            <Column header="耗时" style="min-width: 80px">
               <template #body="{ data }">
                  <span class="text-sm text-surface-500 tabular-nums">
                     {{ data.durationMs }}ms
                  </span>
               </template>
            </Column>

            <Column header="时间" style="min-width: 160px">
               <template #body="{ data }">
                  <span class="text-sm text-surface-500 tabular-nums">
                     {{ data.time }}
                  </span>
               </template>
            </Column>

            <Column header="" style="width: 60px">
               <template #body="{ data }">
                  <Button
                     icon="pi pi-eye"
                     severity="secondary"
                     text
                     rounded
                     size="small"
                     @click="showDetail(data)" />
               </template>
            </Column>
         </DataTable>
      </div>

      <!-- Detail Dialog -->
      <Dialog
         v-model:visible="detailDialog"
         header="操作详情"
         modal
         :style="{ width: '560px' }"
         :pt="{
            root: { class: 'border-none shadow-2xl rounded-2xl' },
            header: {
               class: 'border-b border-surface-200 dark:border-surface-700 px-6 py-4',
            },
            content: { class: 'p-6' },
         }">
         <div v-if="selectedLog" class="flex flex-col gap-5">
            <!-- Operator Info -->
            <div
               class="flex items-center gap-4 p-4 bg-surface-50 dark:bg-surface-800 rounded-xl">
               <div
                  class="w-12 h-12 rounded-xl overflow-hidden bg-surface-200 dark:bg-surface-700">
                  <img
                     :src="selectedLog.operatorAvatar"
                     :alt="selectedLog.operatorName"
                     class="w-full h-full object-cover" />
               </div>
               <div class="flex-1">
                  <div
                     class="font-semibold text-surface-900 dark:text-surface-100">
                     {{ selectedLog.operatorName }}
                  </div>
                  <div class="text-sm text-surface-500">
                     IP: {{ selectedLog.ip }}
                  </div>
               </div>
               <Tag
                  :severity="getStatusSeverity(selectedLog.status)"
                  :value="getStatusLabel(selectedLog.status)" />
            </div>

            <!-- Details Grid -->
            <div class="grid grid-cols-2 gap-4">
               <div>
                  <div
                     class="text-xs text-surface-500 uppercase tracking-wide mb-1">
                     模块
                  </div>
                  <div
                     class="text-sm font-medium text-surface-900 dark:text-surface-100">
                     {{ selectedLog.module }}
                  </div>
               </div>
               <div>
                  <div
                     class="text-xs text-surface-500 uppercase tracking-wide mb-1">
                     操作
                  </div>
                  <div
                     class="text-sm font-medium text-surface-900 dark:text-surface-100">
                     {{ selectedLog.action }}
                  </div>
               </div>
               <div>
                  <div
                     class="text-xs text-surface-500 uppercase tracking-wide mb-1">
                     目标
                  </div>
                  <div
                     class="text-sm font-medium text-surface-900 dark:text-surface-100">
                     {{ selectedLog.targetName || '-' }}
                  </div>
               </div>
               <div>
                  <div
                     class="text-xs text-surface-500 uppercase tracking-wide mb-1">
                     耗时
                  </div>
                  <div
                     class="text-sm font-medium text-surface-900 dark:text-surface-100">
                     {{ selectedLog.durationMs }}ms
                  </div>
               </div>
               <div class="col-span-2">
                  <div
                     class="text-xs text-surface-500 uppercase tracking-wide mb-1">
                     时间
                  </div>
                  <div
                     class="text-sm font-medium text-surface-900 dark:text-surface-100">
                     {{ selectedLog.time }}
                  </div>
               </div>
            </div>

            <!-- Detail JSON -->
            <div>
               <div
                  class="text-xs text-surface-500 uppercase tracking-wide mb-2">
                  详细信息
               </div>
               <pre
                  class="p-4 bg-surface-900 dark:bg-surface-950 text-green-400 text-xs rounded-lg overflow-auto max-h-48 font-mono">{{ JSON.stringify(selectedLog.detail, null, 2) }}</pre>
            </div>
         </div>
      </Dialog>
   </div>
</template>
