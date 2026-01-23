<script setup lang="ts">
import { computed } from 'vue';
import { useThemeStore } from '@/stores/theme';
import { useAnimatedNumber } from '@/composables/use-animated-number';
import type { StatCardData } from '@/types';
import { THEME_COLORS } from '@/config';

const props = defineProps<{ stat: StatCardData }>();

const themeStore = useThemeStore();

const numericValue = computed(() => parseFloat(props.stat.value.replace(/[^0-9.-]/g, '')) || 0);
const valuePrefix = computed(() => props.stat.value.match(/^([^0-9]*)/)?.[1] || '');
const valueSuffix = computed(() => props.stat.value.match(/[0-9]([^0-9]*)$/)?.[1] || '');

const { formattedValue: animatedValue } = useAnimatedNumber(numericValue, { duration: 1200, delay: 100, padStart: true });
const displayValue = computed(() => `${valuePrefix.value}${animatedValue.value}${valueSuffix.value}`);

const colors = computed(() => {
   const scheme = THEME_COLORS[props.stat.color || 'orange'];
   const isDark = themeStore.isDark;
   return {
      bg: isDark ? scheme.bgDark : scheme.bgLight,
      glow: isDark ? scheme.glowDark : 'none',
      text: isDark ? scheme.textDark : scheme.textLight,
      badge: isDark ? scheme.badgeDark : scheme.badgeLight,
      line: isDark ? scheme.lineDark : scheme.lineLight,
      accent: isDark ? scheme.accentDark : scheme.lineLight,
   };
});

const trendData = computed(() => {
   if (props.stat.trendData?.length) return props.stat.trendData;
   return props.stat.changeType === 'increase'
      ? [30, 45, 35, 55, 40, 60, 50, 70, 55, 75]
      : [70, 55, 65, 50, 60, 45, 55, 40, 50, 35];
});

const trendPath = computed(() => {
   const data = trendData.value;
   const width = 120;
   const height = 36;
   const padding = 4;

   const min = Math.min(...data);
   const max = Math.max(...data);
   const range = max - min || 1;

   const points = data.map((val, i) => ({
      x: padding + (i / (data.length - 1)) * (width - padding * 2),
      y: height - padding - ((val - min) / range) * (height - padding * 2),
   }));

   let path = `M ${points[0]!.x} ${points[0]!.y}`;
   for (let i = 1; i < points.length; i++) {
      const prev = points[i - 1]!, curr = points[i]!, cpx = (prev.x + curr.x) / 2;
      path += ` Q ${prev.x + (cpx - prev.x) * 0.5} ${prev.y}, ${cpx} ${(prev.y + curr.y) / 2}`;
      path += ` Q ${cpx + (curr.x - cpx) * 0.5} ${curr.y}, ${curr.x} ${curr.y}`;
   }
   return path;
});
</script>

<template>
   <div
      class="stats-card group relative flex flex-col py-5 px-6 rounded-xl overflow-hidden transition-all duration-300 ease-out min-h-40"
      :style="{
         background: colors.bg,
         boxShadow: colors.glow,
      }">
      <!-- 暗色模式下的发光边框效果 -->
      <div
         v-if="themeStore.isDark"
         class="absolute inset-0 rounded-xl pointer-events-none opacity-0 group-hover:opacity-100 transition-opacity duration-500"
         :style="{
            boxShadow: `0 0 30px -5px ${colors.accent}40, inset 0 0 20px -10px ${colors.accent}20`,
         }"></div>

      <!-- Header: Icon + Title / Change Badge -->
      <div class="relative z-10 flex items-center justify-between mb-3">
         <div class="flex items-center gap-2.5">
            <div
               class="flex items-center justify-center w-8 h-8 rounded-lg transition-transform duration-300 group-hover:scale-110"
               :style="{
                  background: colors.badge,
               }">
               <i
                  :class="stat.icon"
                  class="text-sm"
                  :style="{ color: colors.text }"></i>
            </div>
            <span
               class="text-[0.9375rem] font-semibold tracking-[0.01em]"
               :style="{ color: colors.text }">
               {{ stat.title }}
            </span>
         </div>
         <span
            class="inline-flex items-center py-1.5 px-3 rounded-full text-xs font-bold backdrop-blur-sm"
            :style="{
               background: colors.badge,
               color: colors.text,
            }">
            <i
               :class="
                  stat.changeType === 'increase'
                     ? 'pi pi-arrow-up-right'
                     : 'pi pi-arrow-down-right'
               "
               class="text-[10px] mr-1"></i>
            {{ stat.change }}
         </span>
      </div>

      <!-- Value -->
      <div
         class="relative z-10 font-mono text-4xl font-bold tracking-tight leading-none mb-auto transition-transform duration-300 group-hover:translate-x-1"
         :style="{ color: colors.text }">
         {{ displayValue }}
      </div>

      <!-- Mini Trend Chart -->
      <div class="relative z-10 mt-3 h-9 w-full">
         <svg
            viewBox="0 0 120 36"
            preserveAspectRatio="none"
            class="w-full h-full overflow-visible">
            <!-- 渐变定义 -->
            <defs>
               <linearGradient
                  :id="`line-gradient-${stat.color || 'orange'}`"
                  x1="0%"
                  y1="0%"
                  x2="100%"
                  y2="0%">
                  <stop
                     offset="0%"
                     :stop-color="colors.line"
                     stop-opacity="0.4" />
                  <stop
                     offset="50%"
                     :stop-color="colors.line"
                     stop-opacity="1" />
                  <stop
                     offset="100%"
                     :stop-color="colors.accent"
                     stop-opacity="0.8" />
               </linearGradient>
            </defs>
            <!-- 发光效果（暗色模式） -->
            <path
               v-if="themeStore.isDark"
               :d="trendPath"
               fill="none"
               :stroke="colors.accent"
               stroke-width="6"
               stroke-linecap="round"
               stroke-linejoin="round"
               class="opacity-20 blur-[2px]" />
            <!-- 主线条 -->
            <path
               :d="trendPath"
               fill="none"
               :stroke="`url(#line-gradient-${stat.color || 'orange'})`"
               stroke-width="2.5"
               stroke-linecap="round"
               stroke-linejoin="round"
               class="transition-all duration-300" />
         </svg>
      </div>
   </div>
</template>

<style scoped>
.stats-card {
   /* 悬浮效果 */
   &:hover {
      transform: translateY(-2px);
   }
}

/* 暗色模式下的特殊效果 */
:global(.app-dark) .stats-card {
   /* 添加微妙的边框 */
   border: 1px solid rgba(255, 255, 255, 0.06);
}

:global(.app-dark) .stats-card:hover {
   border-color: rgba(255, 255, 255, 0.1);
}
</style>
