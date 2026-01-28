<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useQuery, keepPreviousData } from '@tanstack/vue-query'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import DatePicker from 'primevue/datepicker'
import Dialog from 'primevue/dialog'
import PageHeader from '@/components/shared/PageHeader.vue'
import SimpleStatCard from '@/components/shared/SimpleStatCard.vue'
import type {
   AuditLog,
   SimpleStatData,
   AuditLogBackend,
   Status,
   AuditAction,
   AuditLogFilter,
} from '@/types'
import { AUDIT_MODULE_ICONS } from '@/config'
import { getAuditLogs } from '@/apis/audit'
import { actionNamesMapping } from '@/mappings/action'

// 分页参数
const page = ref(1)
const pageSize = ref(10)

// 排序参数
const sortField = ref<string | undefined>(undefined)
const sortOrder = ref<1 | -1 | 0>(0)

// 筛选参数
const filters = ref({
   search: '',
   module: null as string | null,
   action: null as string | null,
   dateRange: null as [Date, Date] | null,
})

// 构建查询参数
const queryParams = computed<AuditLogFilter>(() => {
   const params: AuditLogFilter = {
      page: page.value,
      page_size: pageSize.value,
   }

   if (filters.value.module) {
      params.module = filters.value.module as any
   }
   if (filters.value.action) {
      params.action = filters.value.action as any
   }
   if (filters.value.dateRange && filters.value.dateRange[0] && filters.value.dateRange[1]) {
      params.start_time = filters.value.dateRange[0].toISOString()
      params.end_time = filters.value.dateRange[1].toISOString()
   }
   if (sortField.value && sortOrder.value !== 0) {
      params.sort_by = sortField.value
      params.sort_desc = sortOrder.value === -1
   }

   return params
})

// 获取审计日志数据
const {
   data: auditLogsData,
   isLoading,
   isFetching,
} = useQuery({
   queryKey: ['audit-logs', queryParams],
   queryFn: () => getAuditLogs(queryParams.value),
   placeholderData: keepPreviousData,
})

// 监听 filters 变化,重置分页
watch(
   () => [filters.value.module, filters.value.action, filters.value.dateRange],
   () => {
      page.value = 1
   },
   { deep: true }
)

// 映射后端数据到前端格式
const auditLogs = computed<AuditLog[]>(() => {
   const items = auditLogsData.value?.data?.items || []
   return items.map((item: AuditLogBackend): AuditLog => {
      const status = (item.status?.toLowerCase() as Status) || 'success'
      return {
         id: item.id,
         operatorId: item.operator_id,
         operatorName: item.operator_name,
         operatorAvatar: item.operator?.avatar?.file?.storage_key
            ? `/uploads/${item.operator.avatar.file.storage_key}`
            : `https://api.dicebear.com/7.x/avataaars/svg?seed=${item.operator_id}`,
         module: item.module,
         action: item.action,
         targetId: item.target_id,
         targetName: item.target_name,
         detail: item.detail,
         ip: item.ip,
         userAgent: item.user_agent,
         location: item.location,
         status,
         errorMessage: item.error_message,
         durationMs: item.duration_ms || 0,
         time: item.created_at,
      }
   })
})

const total = computed(() => auditLogsData.value?.data?.total || 0)

// 模块选项
const moduleOptions = [
   { label: '全部模块', value: null },
   { label: '认证', value: 'AUTH' },
   { label: 'OAuth', value: 'OAUTH' },
   { label: '用户', value: 'USER' },
   { label: '角色', value: 'ROLE' },
   { label: '权限', value: 'PERMISSION' },
   { label: '客户端', value: 'CLIENT' },
   { label: '系统', value: 'SYSTEM' },
]

// 操作选项
const actionOptions = [
   { label: '全部操作', value: null },
   { label: '登录', value: 'LOGIN' },
   { label: '登出', value: 'LOGOUT' },
   { label: '注册', value: 'REGISTER' },
   { label: '创建客户端', value: 'CLIENT_CREATE' },
   { label: '更新客户端', value: 'CLIENT_UPDATE' },
   { label: '删除客户端', value: 'CLIENT_DELETE' },
   { label: '创建角色', value: 'ROLE_CREATE' },
   { label: '更新角色', value: 'ROLE_UPDATE' },
   { label: '删除角色', value: 'ROLE_DELETE' },
]

// 统计数据 (基于当前页)
const stats = computed<SimpleStatData[]>(() => {
   return [
      {
         title: '总记录数',
         value: total.value,
         icon: 'pi pi-history',
         color: 'blue',
      },
      {
         title: '成功操作',
         value: auditLogs.value.filter(log => log.status === 'success').length,
         icon: 'pi pi-check-circle',
         color: 'green',
      },
      {
         title: '警告事件',
         value: auditLogs.value.filter(log => log.status === 'warning').length,
         icon: 'pi pi-exclamation-triangle',
         color: 'orange',
      },
      {
         title: '失败操作',
         value: auditLogs.value.filter(log => log.status === 'error').length,
         icon: 'pi pi-times-circle',
         color: 'red',
      },
   ]
})

// 前端搜索过滤 (仅对当前页数据)
const filteredLogs = computed(() => {
   if (!filters.value.search) return auditLogs.value

   const search = filters.value.search.toLowerCase()
   return auditLogs.value.filter(log => {
      return (
         log.operatorName.toLowerCase().includes(search) ||
         log.targetName?.toLowerCase().includes(search) ||
         log.action.toLowerCase().includes(search) ||
         log.ip.includes(search)
      )
   })
})

const STATUS_CONFIG = {
   success: { label: '成功', severity: 'success' },
   warning: { label: '警告', severity: 'warn' },
   error: { label: '失败', severity: 'danger' },
} as const
const getStatusSeverity = (status: string) =>
   STATUS_CONFIG[status as keyof typeof STATUS_CONFIG]?.severity || 'secondary'
const getStatusLabel = (status: string) =>
   STATUS_CONFIG[status as keyof typeof STATUS_CONFIG]?.label || status
const getModuleIcon = (module: string) => AUDIT_MODULE_ICONS[module] || 'pi pi-circle'

// 模块中文名称映射
const MODULE_LABELS: Record<string, string> = {
   AUTH: '认证',
   OAUTH: 'OAuth',
   USER: '用户',
   ROLE: '角色',
   PERMISSION: '权限',
   CLIENT: '客户端',
   SYSTEM: '系统',
}
const getModuleLabel = (module: string) => MODULE_LABELS[module] || module

// 操作中文名称映射
const getActionLabel = (action: string) => actionNamesMapping[action as AuditAction] || action

// 格式化时间显示
const formatTime = (time: string) => {
   if (!time) return '-'
   const date = new Date(time)
   return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false,
   })
}

const detailDialog = ref(false)
const selectedLog = ref<AuditLog | null>(null)
const showDetail = (log: AuditLog) => {
   selectedLog.value = log
   detailDialog.value = true
}
const exportLogs = () => console.log('Exporting logs...')
const clearFilters = () => {
   filters.value = { search: '', module: null, action: null, dateRange: null }
   sortField.value = undefined
   sortOrder.value = 0
}

// 分页事件处理
const onPage = (event: any) => {
   page.value = event.page + 1
   pageSize.value = event.rows
}

// 排序事件处理
const onSort = (event: any) => {
   sortField.value = event.sortField
   sortOrder.value = event.sortOrder
   page.value = 1
}
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
               @click="exportLogs"
            />
            <Button
               label="清除筛选"
               icon="pi pi-filter-slash"
               :disabled="
                  !filters.search && !filters.module && !filters.action && !filters.dateRange
               "
               @click="clearFilters"
            />
         </template>
      </PageHeader>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4">
         <template v-if="isLoading">
            <div
               v-for="i in 4"
               :key="i"
               class="h-20 bg-surface-100 dark:bg-surface-800 rounded-xl animate-pulse"
            />
         </template>
         <template v-else>
            <SimpleStatCard v-for="stat in stats" :key="stat.title" :stat="stat" />
         </template>
      </div>

      <!-- Filters -->
      <div
         class="flex flex-wrap items-center gap-4 p-4 bg-surface-0 dark:bg-surface-900 border border-surface-200 dark:border-surface-800 rounded-xl"
      >
         <div class="flex-1 min-w-60 relative">
            <i
               class="pi pi-search absolute left-3 top-1/2 -translate-y-1/2 text-surface-400 z-10"
            ></i>
            <InputText
               v-model="filters.search"
               placeholder="搜索操作者、目标、IP..."
               class="w-full pl-10"
            />
         </div>
         <Select
            v-model="filters.module"
            :options="moduleOptions"
            optionLabel="label"
            optionValue="value"
            placeholder="选择模块"
            :pt="{ root: { style: 'min-width: 140px' } }"
         />
         <Select
            v-model="filters.action"
            :options="actionOptions"
            optionLabel="label"
            optionValue="value"
            placeholder="操作类型"
            :pt="{ root: { style: 'min-width: 130px' } }"
         />
         <DatePicker
            v-model="filters.dateRange"
            selectionMode="range"
            :manualInput="false"
            placeholder="选择日期范围"
            :pt="{ root: { style: 'min-width: 220px' } }"
         />
      </div>

      <!-- Audit Logs Table -->
      <div
         class="bg-surface-0 dark:bg-surface-900 border border-surface-200 dark:border-surface-800 rounded-xl overflow-hidden"
      >
         <DataTable
            :loading="isFetching"
            :value="filteredLogs"
            :paginator="true"
            :rows="pageSize"
            :totalRecords="total"
            :lazy="true"
            :sortField="sortField"
            :sortOrder="sortOrder"
            removableSort
            @page="onPage"
            @sort="onSort"
            :rowsPerPageOptions="[10, 20, 50]"
            paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown CurrentPageReport"
            currentPageReportTemplate="显示 {first} 到 {last} 条，共 {totalRecords} 条"
            responsiveLayout="scroll"
            :pt="{
               root: { class: 'border-none' },
               header: { class: 'bg-transparent border-none px-5 py-4' },
               thead: { class: 'bg-surface-50 dark:bg-surface-800' },
               tbody: { class: 'divide-y divide-surface-100 dark:divide-surface-800' },
               paginator: {
                  class: 'border-t border-surface-200 dark:border-surface-700 px-4 py-3',
               },
            }"
         >
            <template #empty>
               <div class="flex flex-col items-center justify-center py-12 text-surface-500">
                  <i class="pi pi-inbox text-4xl mb-3 opacity-50"></i>
                  <span>暂无审计日志</span>
               </div>
            </template>

            <Column field="operator_name" header="操作者" sortable style="min-width: 180px">
               <template #body="{ data }">
                  <div class="flex items-center gap-3">
                     <div
                        class="w-9 h-9 rounded-lg overflow-hidden bg-surface-100 dark:bg-surface-700 shrink-0"
                     >
                        <img
                           :src="data.operatorAvatar"
                           :alt="data.operatorName"
                           class="w-full h-full object-cover"
                        />
                     </div>
                     <div class="flex flex-col">
                        <span class="font-medium text-surface-900 dark:text-surface-100 text-sm">
                           {{ data.operatorName }}
                        </span>
                        <span class="text-xs text-surface-500">{{ data.ip }}</span>
                     </div>
                  </div>
               </template>
            </Column>

            <Column field="module" header="模块" sortable style="min-width: 120px">
               <template #body="{ data }">
                  <div class="flex items-center gap-2">
                     <i :class="getModuleIcon(data.module)" class="text-surface-500"></i>
                     <span class="text-sm text-surface-700 dark:text-surface-300">
                        {{ getModuleLabel(data.module) }}
                     </span>
                  </div>
               </template>
            </Column>

            <Column field="action" header="操作" sortable style="min-width: 100px">
               <template #body="{ data }">
                  <span class="font-medium text-sm text-surface-900 dark:text-surface-100">
                     {{ getActionLabel(data.action) }}
                  </span>
               </template>
            </Column>

            <Column field="target_name" header="目标" sortable style="min-width: 180px">
               <template #body="{ data }">
                  <span
                     v-if="data.targetName"
                     class="text-sm text-surface-600 dark:text-surface-400"
                  >
                     {{ data.targetName }}
                  </span>
                  <span v-else class="text-sm text-surface-400">-</span>
               </template>
            </Column>

            <Column field="status" header="状态" sortable style="min-width: 90px">
               <template #body="{ data }">
                  <Tag
                     :severity="getStatusSeverity(data.status)"
                     :value="getStatusLabel(data.status)"
                     :pt="{
                        root: { class: 'text-xs px-2 py-1' },
                     }"
                  />
               </template>
            </Column>

            <Column field="duration_ms" header="耗时" sortable style="min-width: 80px">
               <template #body="{ data }">
                  <span class="text-sm text-surface-500 tabular-nums">
                     {{ data.durationMs }}ms
                  </span>
               </template>
            </Column>

            <Column field="created_at" header="时间" sortable style="min-width: 160px">
               <template #body="{ data }">
                  <span class="text-sm text-surface-500 tabular-nums">
                     {{ formatTime(data.time) }}
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
                     @click="showDetail(data)"
                  />
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
         }"
      >
         <div v-if="selectedLog" class="flex flex-col gap-5">
            <!-- Operator Info -->
            <div class="flex items-center gap-4 p-4 bg-surface-50 dark:bg-surface-800 rounded-xl">
               <div class="w-12 h-12 rounded-xl overflow-hidden bg-surface-200 dark:bg-surface-700">
                  <img
                     :src="selectedLog.operatorAvatar"
                     :alt="selectedLog.operatorName"
                     class="w-full h-full object-cover"
                  />
               </div>
               <div class="flex-1">
                  <div class="font-semibold text-surface-900 dark:text-surface-100">
                     {{ selectedLog.operatorName }}
                  </div>
                  <div class="text-sm text-surface-500">IP: {{ selectedLog.ip }}</div>
               </div>
               <Tag
                  :severity="getStatusSeverity(selectedLog.status)"
                  :value="getStatusLabel(selectedLog.status)"
               />
            </div>

            <!-- Details Grid -->
            <div class="grid grid-cols-2 gap-4">
               <div>
                  <div class="text-xs text-surface-500 uppercase tracking-wide mb-1">模块</div>
                  <div class="text-sm font-medium text-surface-900 dark:text-surface-100">
                     {{ getModuleLabel(selectedLog.module) }}
                  </div>
               </div>
               <div>
                  <div class="text-xs text-surface-500 uppercase tracking-wide mb-1">操作</div>
                  <div class="text-sm font-medium text-surface-900 dark:text-surface-100">
                     {{ getActionLabel(selectedLog.action) }}
                  </div>
               </div>
               <div>
                  <div class="text-xs text-surface-500 uppercase tracking-wide mb-1">目标</div>
                  <div class="text-sm font-medium text-surface-900 dark:text-surface-100">
                     {{ selectedLog.targetName || '-' }}
                  </div>
               </div>
               <div>
                  <div class="text-xs text-surface-500 uppercase tracking-wide mb-1">耗时</div>
                  <div class="text-sm font-medium text-surface-900 dark:text-surface-100">
                     {{ selectedLog.durationMs }}ms
                  </div>
               </div>
               <div class="col-span-2">
                  <div class="text-xs text-surface-500 uppercase tracking-wide mb-1">时间</div>
                  <div class="text-sm font-medium text-surface-900 dark:text-surface-100">
                     {{ formatTime(selectedLog.time) }}
                  </div>
               </div>
            </div>

            <!-- Detail JSON -->
            <div>
               <div class="text-xs text-surface-500 uppercase tracking-wide mb-2">详细信息</div>
               <pre
                  class="p-4 bg-surface-900 dark:bg-surface-950 text-green-400 text-xs rounded-lg overflow-auto max-h-48 font-mono"
                  >{{ JSON.stringify(selectedLog.detail, null, 2) }}</pre
               >
            </div>
         </div>
      </Dialog>
   </div>
</template>
