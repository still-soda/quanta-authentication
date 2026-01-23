<script setup lang="ts">
import { computed } from 'vue';
import { useAnimatedNumber } from '@/composables/use-animated-number';
import type { SimpleStatData } from '@/types';
import { SIMPLE_STAT_COLORS } from '@/config';

const props = defineProps<{ stat: SimpleStatData }>();

const numericValue = computed(() => {
   if (typeof props.stat.value === 'number') return props.stat.value;
   return parseFloat(props.stat.value.replace(/[^0-9.-]/g, '')) || 0;
});

const { formattedValue: animatedValue } = useAnimatedNumber(numericValue, { duration: 1000, delay: 50, padStart: true });
const themeColors = computed(() => SIMPLE_STAT_COLORS[props.stat.color || 'orange']);
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
            class="text-2xl font-bold text-surface-900 dark:text-surface-0 tracking-tight leading-tight tabular-nums">
            {{ animatedValue }}
         </span>
         <span
            class="text-[0.8125rem] font-medium text-surface-500 dark:text-surface-400">
            {{ stat.title }}
         </span>
      </div>
   </div>
</template>
