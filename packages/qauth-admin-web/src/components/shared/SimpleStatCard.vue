<script setup lang="ts">
import { computed } from 'vue';

export interface SimpleStatData {
   title: string;
   value: number | string;
   icon: string;
   /** 卡片主题色 */
   color?: 'orange' | 'blue' | 'green' | 'purple' | 'cyan' | 'red' | 'gray';
}

const props = withDefaults(
   defineProps<{
      stat: SimpleStatData;
   }>(),
   {},
);

const themeColors = computed(() => {
   const colorMap = {
      orange: {
         color: 'var(--p-orange-500)',
         bg: 'var(--p-orange-50)',
         bgDark: 'rgba(251, 146, 60, 0.15)',
      },
      blue: {
         color: 'var(--p-blue-500)',
         bg: 'var(--p-blue-50)',
         bgDark: 'rgba(59, 130, 246, 0.15)',
      },
      green: {
         color: 'var(--p-green-500)',
         bg: 'var(--p-green-50)',
         bgDark: 'rgba(34, 197, 94, 0.15)',
      },
      purple: {
         color: 'var(--p-purple-500)',
         bg: 'var(--p-purple-50)',
         bgDark: 'rgba(168, 85, 247, 0.15)',
      },
      cyan: {
         color: 'var(--p-cyan-500)',
         bg: 'var(--p-cyan-50)',
         bgDark: 'rgba(6, 182, 212, 0.15)',
      },
      red: {
         color: 'var(--p-red-500)',
         bg: 'var(--p-red-50)',
         bgDark: 'rgba(239, 68, 68, 0.15)',
      },
      gray: {
         color: 'var(--p-surface-500)',
         bg: 'var(--p-surface-100)',
         bgDark: 'rgba(107, 114, 128, 0.15)',
      },
   };
   return colorMap[props.stat.color || 'orange'];
});
</script>

<template>
   <div
      class="group flex items-center gap-4 py-5 px-6 bg-surface-0 dark:bg-surface-900 border border-surface-200 dark:border-surface-800 rounded-[14px] transition-all duration-250 ease-in-out hover:-translate-y-0.5 hover:shadow-[0_8px_20px_-6px_rgba(0,0,0,0.08)] dark:hover:shadow-[0_8px_20px_-6px_rgba(0,0,0,0.35)]"
      :style="{
         '--theme-color': themeColors.color,
         '--theme-bg': themeColors.bg,
         '--theme-bg-dark': themeColors.bgDark,
      }">
      <div
         class="w-12 h-12 flex items-center justify-center rounded-xl bg-(--theme-bg) dark:bg-(--theme-bg-dark) text-(--theme-color) shrink-0 transition-all duration-250 ease group-hover:bg-(--theme-color) group-hover:text-white group-hover:scale-105">
         <i :class="stat.icon" class="text-xl"></i>
      </div>
      <div class="flex flex-col gap-0.5">
         <span
            class="text-2xl font-bold text-surface-900 dark:text-surface-0 tracking-tight leading-tight">
            {{ stat.value }}
         </span>
         <span
            class="text-[0.8125rem] font-medium text-surface-500 dark:text-surface-400">
            {{ stat.title }}
         </span>
      </div>
   </div>
</template>
