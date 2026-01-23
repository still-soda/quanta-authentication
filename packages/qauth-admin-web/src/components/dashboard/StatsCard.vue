<script setup lang="ts">
import { computed } from 'vue';
import { useThemeStore } from '@/stores/theme';
import { useAnimatedNumber } from '@/composables/useAnimatedNumber';

export interface StatCardData {
   title: string;
   value: string;
   change: string;
   changeType: 'increase' | 'decrease';
   icon: string;
   /** 卡片的主题色 (支持 PrimeVue 预设颜色) */
   color?: 'orange' | 'blue' | 'green' | 'purple' | 'cyan' | 'red';
   /** 可选的趋势数据，用于绘制迷你图 */
   trendData?: number[];
}

const props = withDefaults(
   defineProps<{
      stat: StatCardData;
   }>(),
   {},
);

const themeStore = useThemeStore();

// 解析 value 字符串中的数值（支持千分位格式，如 "1,234" -> 1234）
const numericValue = computed(() => {
   const cleaned = props.stat.value.replace(/[^0-9.-]/g, '');
   return parseFloat(cleaned) || 0;
});

// 判断 value 是否包含前缀或后缀（如 "¥1,234" 或 "1,234人"）
const valuePrefix = computed(() => {
   const match = props.stat.value.match(/^([^0-9]*)/);
   return match?.[1] || '';
});

const valueSuffix = computed(() => {
   const match = props.stat.value.match(/[0-9]([^0-9]*)$/);
   return match?.[1] || '';
});

// 使用动画数字
const { formattedValue: animatedValue } = useAnimatedNumber(numericValue, {
   duration: 1200,
   delay: 100,
   padStart: true,
});

// 完整的显示值（带前缀后缀）
const displayValue = computed(() => {
   return `${valuePrefix.value}${animatedValue.value}${valueSuffix.value}`;
});

// 主题配色 - 亮色模式柔和优雅，暗色模式深邃有质感
const themeColors = {
   orange: {
      // 亮色模式 - 温暖的橙色渐变
      bgLight: 'linear-gradient(145deg, #fff7ed 0%, #fed7aa 30%, #fdba74 80%, #fff7ed 120%)',
      textLight: '#9a3412',
      badgeLight: 'rgba(251, 146, 60, 0.25)',
      lineLight: '#ea580c',
      // 暗色模式 - 深邃的琥珀色调配发光边缘
      bgDark:
         'linear-gradient(145deg, #1c1917 0%, #292524 50%, #44403c 100%)',
      glowDark: 'inset 0 1px 0 rgba(251, 146, 60, 0.15)',
      textDark: '#fdba74',
      badgeDark: 'rgba(251, 146, 60, 0.2)',
      lineDark: '#fb923c',
      accentDark: '#f97316',
   },
   blue: {
      bgLight: 'linear-gradient(145deg, #eff6ff 0%, #bfdbfe 30%, #93c5fd 80%, #eff6ff 120%)',
      textLight: '#1e40af',
      badgeLight: 'rgba(59, 130, 246, 0.2)',
      lineLight: '#2563eb',
      bgDark:
         'linear-gradient(145deg, #0c1929 0%, #1e293b 50%, #334155 100%)',
      glowDark: 'inset 0 1px 0 rgba(59, 130, 246, 0.2)',
      textDark: '#93c5fd',
      badgeDark: 'rgba(59, 130, 246, 0.25)',
      lineDark: '#60a5fa',
      accentDark: '#3b82f6',
   },
   green: {
      bgLight: 'linear-gradient(145deg, #f0fdf4 0%, #bbf7d0 30%, #86efac 80%, #f0fdf4 120%)',
      textLight: '#166534',
      badgeLight: 'rgba(34, 197, 94, 0.2)',
      lineLight: '#16a34a',
      bgDark:
         'linear-gradient(145deg, #052e16 0%, #14532d 50%, #166534 100%)',
      glowDark: 'inset 0 1px 0 rgba(34, 197, 94, 0.2)',
      textDark: '#86efac',
      badgeDark: 'rgba(34, 197, 94, 0.25)',
      lineDark: '#4ade80',
      accentDark: '#22c55e',
   },
   purple: {
      bgLight: 'linear-gradient(145deg, #faf5ff 0%, #e9d5ff 30%, #d8b4fe 80%, #faf5ff 120%)',
      textLight: '#7e22ce',
      badgeLight: 'rgba(168, 85, 247, 0.2)',
      lineLight: '#9333ea',
      bgDark:
         'linear-gradient(145deg, #1a0a2e 0%, #2e1065 50%, #4c1d95 100%)',
      glowDark: 'inset 0 1px 0 rgba(168, 85, 247, 0.2)',
      textDark: '#d8b4fe',
      badgeDark: 'rgba(168, 85, 247, 0.25)',
      lineDark: '#c084fc',
      accentDark: '#a855f7',
   },
   cyan: {
      bgLight: 'linear-gradient(145deg, #ecfeff 0%, #a5f3fc 30%, #67e8f9 80%, #ecfeff 120%)',
      textLight: '#0e7490',
      badgeLight: 'rgba(6, 182, 212, 0.2)',
      lineLight: '#06b6d4',
      bgDark:
         'linear-gradient(145deg, #042f2e 0%, #134e4a 50%, #115e59 100%)',
      glowDark: 'inset 0 1px 0 rgba(6, 182, 212, 0.25)',
      textDark: '#67e8f9',
      badgeDark: 'rgba(6, 182, 212, 0.25)',
      lineDark: '#22d3ee',
      accentDark: '#06b6d4',
   },
   red: {
      bgLight: 'linear-gradient(145deg, #fef2f2 0%, #fecaca 30%, #fca5a5 80%, #fef2f2 120%)',
      textLight: '#b91c1c',
      badgeLight: 'rgba(239, 68, 68, 0.2)',
      lineLight: '#dc2626',
      bgDark:
         'linear-gradient(145deg, #1c0a0a 0%, #450a0a 50%, #7f1d1d 100%)',
      glowDark: 'inset 0 1px 0 rgba(239, 68, 68, 0.15)',
      textDark: '#fca5a5',
      badgeDark: 'rgba(239, 68, 68, 0.25)',
      lineDark: '#f87171',
      accentDark: '#ef4444',
   },
};

// 根据当前主题动态计算颜色
const colors = computed(() => {
   const colorScheme = themeColors[props.stat.color || 'orange'];
   const isDark = themeStore.isDark;

   return {
      bg: isDark ? colorScheme.bgDark : colorScheme.bgLight,
      glow: isDark ? colorScheme.glowDark : 'none',
      text: isDark ? colorScheme.textDark : colorScheme.textLight,
      badge: isDark ? colorScheme.badgeDark : colorScheme.badgeLight,
      line: isDark ? colorScheme.lineDark : colorScheme.lineLight,
      accent: isDark ? colorScheme.accentDark : colorScheme.lineLight,
   };
});

// 生成默认趋势数据
const trendData = computed(() => {
   if (props.stat.trendData && props.stat.trendData.length > 0) {
      return props.stat.trendData;
   }
   // 根据变化类型生成模拟数据
   const isIncrease = props.stat.changeType === 'increase';
   if (isIncrease) {
      return [30, 45, 35, 55, 40, 60, 50, 70, 55, 75];
   } else {
      return [70, 55, 65, 50, 60, 45, 55, 40, 50, 35];
   }
});

// 计算 SVG 路径
const trendPath = computed(() => {
   const data = trendData.value;
   const width = 120;
   const height = 36;
   const padding = 4;

   const min = Math.min(...data);
   const max = Math.max(...data);
   const range = max - min || 1;

   const points = data.map((val, i) => {
      const x = padding + (i / (data.length - 1)) * (width - padding * 2);
      const y =
         height - padding - ((val - min) / range) * (height - padding * 2);
      return { x, y };
   });

   // 生成平滑曲线
   let path = `M ${points[0]!.x} ${points[0]!.y}`;
   for (let i = 1; i < points.length; i++) {
      const prev = points[i - 1]!;
      const curr = points[i]!;
      const cpx = (prev.x + curr.x) / 2;
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
