<script setup lang="ts">
import { computed, ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import Card from 'primevue/card'
import Button from 'primevue/button'
import Chart from 'primevue/chart'
import Popover from 'primevue/popover'
import { useThemeStore } from '@/stores/theme'
import { getAuthTrend } from '@/apis/dashboard'
import type { AuthTrendRange } from '@/types'

const themeStore = useThemeStore()

// 时间范围选项
const rangeOptions = [
   { label: '最近一周', value: 'weekly' as AuthTrendRange, days: 7 },
   { label: '最近半月', value: 'half-weekly' as AuthTrendRange, days: 15 },
   { label: '最近一月', value: 'monthly' as AuthTrendRange, days: 30 },
]

// 当前选中的时间范围
const selectedRange = ref<AuthTrendRange>('weekly')

// 弹窗引用
const popoverRef = ref()

// 获取认证趋势数据
const { data: authTrendData, isLoading } = useQuery({
   queryKey: ['dashboard', 'authTrend', selectedRange],
   queryFn: () => getAuthTrend({ range: selectedRange.value }),
})

// 计算趋势数据和标签
const trendData = computed(() => {
   if (isLoading.value || !authTrendData.value) {
      return {
         labels: [],
         data: [],
      }
   }

   const data = authTrendData.value
   const currentOption = rangeOptions.find(opt => opt.value === selectedRange.value)
   const days = currentOption?.days || 7

   // 生成日期标签
   const labels = []
   const today = new Date()

   // 根据时间范围决定标签格式和间隔
   if (days === 7) {
      // 一周：显示星期，每天一个标签
      const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
      for (let i = days - 1; i >= 0; i--) {
         const date = new Date(today)
         date.setDate(today.getDate() - i)
         labels.push(weekdays[date.getDay()])
      }
   } else if (days === 15) {
      // 半月：每天一个标签，格式为 MM/DD
      for (let i = days - 1; i >= 0; i--) {
         const date = new Date(today)
         date.setDate(today.getDate() - i)
         const month = date.getMonth() + 1
         const day = date.getDate()
         labels.push(`${month}/${day}`)
      }
   } else if (days === 30) {
      // 一个月：每天一个标签，格式为 MM/DD
      for (let i = days - 1; i >= 0; i--) {
         const date = new Date(today)
         date.setDate(today.getDate() - i)
         const month = date.getMonth() + 1
         const day = date.getDate()
         labels.push(`${month}/${day}`)
      }
   }

   // 确保数据长度正确
   let trendArray: number[]
   if (data.length < days) {
      // 数据不足时，前面补0
      trendArray = Array(days - data.length)
         .fill(0)
         .concat(data)
   } else {
      // 数据过多时，取最后 days 天的数据
      trendArray = data.slice(-days)
   }

   return { labels, data: trendArray }
})

const authTrendOptions = computed(() => ({
   maintainAspectRatio: false,
   plugins: {
      legend: {
         display: false,
      },
      tooltip: {
         mode: 'index',
         intersect: false,
         backgroundColor: themeStore.isDark ? '#27272a' : '#ffffff',
         titleColor: themeStore.isDark ? '#fafafa' : '#18181b',
         bodyColor: themeStore.isDark ? '#a1a1aa' : '#71717a',
         borderColor: themeStore.isDark ? '#3f3f46' : '#e4e4e7',
         borderWidth: 1,
         padding: 12,
         cornerRadius: 8,
      },
   },
   scales: {
      x: {
         grid: {
            display: false,
         },
         ticks: {
            color: themeStore.isDark ? '#71717a' : '#a1a1aa',
         },
      },
      y: {
         grid: {
            color: themeStore.isDark ? '#27272a' : '#f4f4f5',
         },
         ticks: {
            color: themeStore.isDark ? '#71717a' : '#a1a1aa',
         },
      },
   },
   interaction: {
      intersect: false,
      mode: 'index',
   },
}))

const authTrendChartData = computed(() => {
   if (!trendData.value) return null
   return {
      labels: trendData.value.labels,
      datasets: [
         {
            label: '认证请求',
            data: trendData.value.data,
            fill: true,
            borderColor: '#f97316',
            backgroundColor: themeStore.isDark
               ? 'rgba(249, 115, 22, 0.1)'
               : 'rgba(249, 115, 22, 0.08)',
            tension: 0.4,
            pointRadius: 0,
            pointHoverRadius: 6,
            pointHoverBackgroundColor: '#f97316',
            pointHoverBorderColor: '#ffffff',
            pointHoverBorderWidth: 2,
         },
      ],
   }
})

// 选择时间范围
function selectRange(range: AuthTrendRange) {
   selectedRange.value = range
   popoverRef.value?.hide()
}

// 显示弹窗
function showPopover(event: Event) {
   popoverRef.value?.toggle(event)
}
</script>

<template>
   <Card class="rounded-2xl border border-surface-100 dark:border-surface-800 overflow-hidden">
      <template #title>
         <div
            class="flex items-center justify-between text-base font-semibold text-surface-900 dark:text-surface-100"
         >
            <span>认证趋势</span>
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
               v-else-if="authTrendChartData"
               type="line"
               :data="authTrendChartData"
               :options="authTrendOptions"
               class="h-full"
            />
         </div>
      </template>
   </Card>

   <!-- 时间范围选择弹窗 -->
   <Popover ref="popoverRef">
      <div class="w-48 p-1">
         <div class="text-xs text-surface-500 dark:text-surface-400 mb-2 px-3 pt-2">时间范围</div>
         <div class="space-y-1">
            <button
               v-for="option in rangeOptions"
               :key="option.value"
               class="w-full flex items-center justify-between px-3 py-2 rounded-lg cursor-pointer hover:bg-surface-100 dark:hover:bg-surface-700 transition-colors text-left"
               :class="{
                  'bg-primary-50 dark:bg-primary-900/20': selectedRange === option.value,
               }"
               @click="selectRange(option.value)"
            >
               <span
                  class="text-sm font-medium"
                  :class="
                     selectedRange === option.value
                        ? 'text-primary-600 dark:text-primary-400'
                        : 'text-surface-700 dark:text-surface-200'
                  "
               >
                  {{ option.label }}
               </span>
               <i
                  v-if="selectedRange === option.value"
                  class="pi pi-check text-xs text-primary-600 dark:text-primary-400"
               ></i>
            </button>
         </div>
      </div>
   </Popover>
</template>
