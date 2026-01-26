<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import Card from 'primevue/card'
import Button from 'primevue/button'
import Chart from 'primevue/chart'
import Popover from 'primevue/popover'
import Checkbox from 'primevue/checkbox'
import Divider from 'primevue/divider'
import { useThemeStore } from '@/stores/theme'
import { getUserDistData, getAllRoles } from '@/apis/dashboard'
import {
   useRoleDistributionConfig,
   BUILTIN_ROLE_CODES,
} from '@/composables/use-role-distribution-config'
import { generateBeautifulColors } from '@/composables/use-color-wheel'

const themeStore = useThemeStore()

// 角色配置
const { selectedRoleCodes, setSelectedRoles } = useRoleDistributionConfig()

// 临时选中状态（用于弹窗编辑）
const tempSelectedRoles = ref<string[]>([])

// 弹窗引用
const popoverRef = ref()

// 获取所有角色列表
const { data: allRoles, isLoading: isLoadingRoles } = useQuery({
   queryKey: ['dashboard', 'allRoles'],
   queryFn: getAllRoles,
})

// 获取用户分布数据
const {
   data: distData,
   isLoading,
   refetch,
} = useQuery({
   queryKey: ['dashboard', 'userDist', selectedRoleCodes],
   queryFn: () => getUserDistData(selectedRoleCodes.value),
})

// 当选中的角色变化时，重新获取数据
watch(
   selectedRoleCodes,
   () => {
      refetch()
   },
   { deep: true }
)

const userDistOptions = computed(() => ({
   maintainAspectRatio: false,
   plugins: {
      legend: {
         position: 'right',
         labels: {
            color: themeStore.isDark ? '#a1a1aa' : '#71717a',
            usePointStyle: true,
            padding: 20,
         },
      },
   },
}))

// 根据选中的角色数量生成颜色
const chartColors = computed(() => {
   const count = distData.value?.labels.length || 0
   return generateBeautifulColors(count, themeStore.isDark)
})

const userDistChartData = computed(() => {
   if (!distData.value) return null

   // 构建角色名称到代码的映射
   const roleNameToCodeMap: Record<string, string> = {}
   allRoles.value?.forEach(role => {
      roleNameToCodeMap[role.name] = role.code
   })

   // 根据选中的角色代码过滤数据
   const filteredIndices: number[] = []
   distData.value.labels.forEach((label, index) => {
      const roleCode = roleNameToCodeMap[label]
      if (roleCode && selectedRoleCodes.value.includes(roleCode)) {
         filteredIndices.push(index)
      }
   })

   if (filteredIndices.length === 0) {
      return {
         labels: ['暂无数据'],
         datasets: [{ data: [1], backgroundColor: ['#E5E7EB'], borderWidth: 0 }],
      }
   }

   const labels = filteredIndices.map(i => distData.value!.labels[i]!)
   const data = filteredIndices.map(i => distData.value!.data[i]!)
   const colors = filteredIndices.map(i => distData.value!.colors?.[i] || chartColors.value[i])

   return {
      labels,
      datasets: [
         {
            data,
            backgroundColor: colors,
            borderWidth: 0,
         },
      ],
   }
})

// 打开弹窗时，同步临时状态
function onPopoverShow() {
   tempSelectedRoles.value = [...selectedRoleCodes.value]
}

// 切换临时选中状态
function toggleTempRole(roleCode: string) {
   const index = tempSelectedRoles.value.indexOf(roleCode)
   if (index === -1) {
      tempSelectedRoles.value.push(roleCode)
   } else {
      tempSelectedRoles.value.splice(index, 1)
   }
}

// 保存选择
function saveSelection() {
   setSelectedRoles(tempSelectedRoles.value)
   popoverRef.value?.hide()
}

// 重置为默认
function handleReset() {
   tempSelectedRoles.value = [...BUILTIN_ROLE_CODES]
}

// 取消
function handleCancel() {
   popoverRef.value?.hide()
}

// 显示弹窗
function showPopover(event: Event) {
   popoverRef.value?.toggle(event)
}

// 检查角色是否被临时选中
function isTempSelected(roleCode: string): boolean {
   return tempSelectedRoles.value.includes(roleCode)
}

// 分组角色：系统内置角色 vs 自定义角色
const systemRoles = computed(() => {
   return Array.isArray(allRoles.value) ? allRoles.value.filter(r => r.is_system) : []
})

const customRoles = computed(() => {
   return Array.isArray(allRoles.value) ? allRoles.value.filter(r => !r.is_system) : []
})
</script>

<template>
   <Card class="rounded-2xl border border-surface-100 dark:border-surface-800 overflow-hidden">
      <template #title>
         <div
            class="flex items-center justify-between text-base font-semibold text-surface-900 dark:text-surface-100"
         >
            <span>角色分布</span>
            <div>
               <Button
                  icon="pi pi-ellipsis-h"
                  text
                  rounded
                  severity="secondary"
                  @click="showPopover"
               />
            </div>
         </div>
      </template>
      <template #content>
         <div class="h-70">
            <div v-if="isLoading" class="h-full flex items-center justify-center">
               <i class="pi pi-spin pi-spinner text-2xl text-surface-400"></i>
            </div>
            <Chart
               v-else-if="userDistChartData"
               type="doughnut"
               :data="userDistChartData"
               :options="userDistOptions"
               class="h-full"
            />
         </div>
      </template>
   </Card>

   <!-- 角色选择弹窗 -->
   <Popover ref="popoverRef" @show="onPopoverShow">
      <div class="w-72 p-1">
         <div class="flex items-center justify-between mb-3">
            <span class="text-sm font-semibold text-surface-700 dark:text-surface-200">
               选择显示的角色
            </span>
            <Button
               label="重置"
               icon="pi pi-refresh"
               size="small"
               text
               severity="secondary"
               @click="handleReset"
            />
         </div>

         <div v-if="isLoadingRoles" class="flex items-center justify-center py-4">
            <i class="pi pi-spin pi-spinner text-lg text-surface-400"></i>
         </div>

         <div v-else class="max-h-80 overflow-y-auto">
            <!-- 系统内置角色 -->
            <div v-if="systemRoles.length > 0" class="mb-3">
               <div class="text-xs text-surface-500 dark:text-surface-400 mb-2 px-1">
                  系统内置角色
               </div>
               <div class="space-y-1">
                  <label
                     v-for="role in systemRoles"
                     :key="role.code"
                     class="flex items-center gap-3 px-3 py-1 rounded-lg cursor-pointer hover:bg-surface-100 dark:hover:bg-surface-700 transition-colors"
                  >
                     <Checkbox
                        :model-value="isTempSelected(role.code)"
                        :binary="true"
                        @update:model-value="toggleTempRole(role.code)"
                     />
                     <div class="flex-1 min-w-0">
                        <div
                           class="text-sm font-medium text-surface-700 dark:text-surface-200 truncate"
                        >
                           {{ role.name }}
                        </div>
                        <div class="text-xs text-surface-500 dark:text-surface-400 truncate">
                           {{ role.code }}
                        </div>
                     </div>
                     <span
                        class="shrink-0 text-xs px-2 py-0.5 rounded-full bg-primary-100 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400"
                     >
                        内置
                     </span>
                  </label>
               </div>
            </div>

            <!-- 自定义角色 -->
            <div v-if="customRoles.length > 0">
               <div class="text-xs text-surface-500 dark:text-surface-400 mb-2 px-1">
                  自定义角色
               </div>
               <div class="space-y-1">
                  <label
                     v-for="role in customRoles"
                     :key="role.code"
                     class="flex items-center gap-2 px-3 py-1 rounded-lg cursor-pointer hover:bg-surface-100 dark:hover:bg-surface-700 transition-colors"
                  >
                     <Checkbox
                        :model-value="isTempSelected(role.code)"
                        :binary="true"
                        @update:model-value="toggleTempRole(role.code)"
                     />
                     <div class="flex-1 min-w-0">
                        <div
                           class="text-sm font-medium text-surface-700 dark:text-surface-200 truncate"
                        >
                           {{ role.name }}
                        </div>
                        <div class="text-xs text-surface-500 dark:text-surface-400 truncate">
                           {{ role.code }}
                        </div>
                     </div>
                  </label>
               </div>
            </div>

            <!-- 无自定义角色提示 -->
            <div
               v-if="customRoles.length === 0 && !isLoadingRoles"
               class="text-center py-4 text-sm text-surface-500 dark:text-surface-400"
            >
               暂无自定义角色
            </div>
         </div>

         <Divider class="my-3" />

         <div class="flex items-center justify-end gap-2">
            <Button label="取消" size="small" text severity="secondary" @click="handleCancel" />
            <Button
               label="保存"
               size="small"
               :disabled="tempSelectedRoles.length === 0"
               @click="saveSelection"
            />
         </div>
      </div>
   </Popover>
</template>
