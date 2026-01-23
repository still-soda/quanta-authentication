<script setup lang="ts">
import { computed } from 'vue';

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

// 主题配色 - 使用柔和的背景色搭配深色文字
const themeColors = {
   orange: {
      bg: 'linear-gradient(135deg, #fef3e2 0%, #ffd7aa 80%, #fef3e2 120%)',
      bgDark: 'linear-gradient(135deg, #431407 0%, #7c2d12 100%)',
      text: '#9a3412',
      textDark: '#fed7aa',
      badge: '#fdba74',
      badgeDark: '#c2410c',
      line: '#ea580c',
      lineDark: '#fb923c',
   },
   blue: {
      bg: 'linear-gradient(135deg, #eff6ff 0%, #bfebfe 80%, #eff6ff 120%)',
      bgDark: 'linear-gradient(135deg, #172554 0%, #1e3a8a 100%)',
      text: '#1e40af',
      textDark: '#bfdbfe',
      badge: '#93c5fd',
      badgeDark: '#1d4ed8',
      line: '#2563eb',
      lineDark: '#60a5fa',
   },
   green: {
      bg: 'linear-gradient(135deg, #f0fdf4 0%, #bbf7d0 80%, #f0fdf4 120%)',
      bgDark: 'linear-gradient(135deg, #052e16 0%, #14532d 100%)',
      text: '#166534',
      textDark: '#bbf7d0',
      badge: '#86efac',
      badgeDark: '#15803d',
      line: '#16a34a',
      lineDark: '#4ade80',
   },
   purple: {
      bg: 'linear-gradient(135deg, #faf5ff 0%, #e9d5ff 80%, #faf5ff 120%)',
      bgDark: 'linear-gradient(135deg, #3b0764 0%, #581c87 100%)',
      text: '#7e22ce',
      textDark: '#e9d5ff',
      badge: '#d8b4fe',
      badgeDark: '#7e22ce',
      line: '#9333ea',
      lineDark: '#c084fc',
   },
   cyan: {
      bg: 'linear-gradient(135deg, #ecfeff 0%, #a5f3fc 80%, #ecfeff 120%)',
      bgDark: 'linear-gradient(135deg, #083344 0%, #164e63 100%)',
      text: '#0e7490',
      textDark: '#a5f3fc',
      badge: '#67e8f9',
      badgeDark: '#0891b2',
      line: '#06b6d4',
      lineDark: '#22d3ee',
   },
   red: {
      bg: 'linear-gradient(135deg, #fef2f2 0%, #fecaca 100%)',
      bgDark: 'linear-gradient(135deg, #450a0a 0%, #7f1d1d 100%)',
      text: '#b91c1c',
      textDark: '#fecaca',
      badge: '#fca5a5',
      badgeDark: '#dc2626',
      line: '#dc2626',
      lineDark: '#f87171',
   },
};

const colors = computed(() => themeColors[props.stat.color || 'orange']);

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
      class="stat-card"
      :style="{
         '--card-bg': colors.bg,
         '--card-bg-dark': colors.bgDark,
         '--card-text': colors.text,
         '--card-text-dark': colors.textDark,
         '--badge-bg': colors.badge,
         '--badge-bg-dark': colors.badgeDark,
         '--line-color': colors.line,
         '--line-color-dark': colors.lineDark,
      }">
      <!-- Header: Icon + Title / Change Badge -->
      <div class="stat-header">
         <div class="stat-label">
            <i :class="stat.icon" class="stat-icon"></i>
            <span class="stat-title">{{ stat.title }}</span>
         </div>
         <span class="stat-badge">
            {{ stat.change }}
         </span>
      </div>

      <!-- Value -->
      <div class="stat-value">{{ stat.value }}</div>

      <!-- Mini Trend Chart -->
      <div class="stat-trend">
         <svg viewBox="0 0 120 36" preserveAspectRatio="none">
            <path
               :d="trendPath"
               fill="none"
               stroke-width="2.5"
               stroke-linecap="round"
               stroke-linejoin="round"
               class="trend-line" />
         </svg>
      </div>
   </div>
</template>

<style scoped>
.stat-card {
   position: relative;
   display: flex;
   flex-direction: column;
   padding: 1.25rem 1.5rem;
   background: var(--card-bg);
   border-radius: 20px;
   overflow: hidden;
   transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
   min-height: 160px;
}

:global(.app-dark) .stat-card {
   background: var(--card-bg-dark);
}

:global(.app-dark) .stat-card:hover {
   box-shadow: 0 20px 40px -12px rgba(0, 0, 0, 0.5);
}

/* Header */
.stat-header {
   display: flex;
   align-items: center;
   justify-content: space-between;
   margin-bottom: 0.75rem;
}

.stat-label {
   display: flex;
   align-items: center;
   gap: 0.5rem;
}

.stat-icon {
   font-size: 1rem;
   color: var(--card-text);
   opacity: 0.9;
}

:global(.app-dark) .stat-icon {
   color: var(--card-text-dark);
}

.stat-title {
   font-size: 0.9375rem;
   font-weight: 600;
   color: var(--card-text);
   letter-spacing: 0.01em;
}

:global(.app-dark) .stat-title {
   color: var(--card-text-dark);
}

/* Badge */
.stat-badge {
   display: inline-flex;
   align-items: center;
   padding: 0.3rem 0.625rem;
   border-radius: 8px;
   font-size: 0.75rem;
   font-weight: 700;
   background: var(--badge-bg);
   color: var(--card-text);
}

:global(.app-dark) .stat-badge {
   background: var(--badge-bg-dark);
   color: var(--card-text-dark);
}

/* Value */
.stat-value {
   font-family: monospace;
   font-size: 2.25rem;
   font-weight: 700;
   color: var(--card-text);
   letter-spacing: -0.03em;
   line-height: 1.1;
   margin-bottom: auto;
}

:global(.app-dark) .stat-value {
   color: var(--card-text-dark);
}

/* Mini Trend Chart */
.stat-trend {
   margin-top: 0.75rem;
   height: 36px;
   width: 100%;
}

.stat-trend svg {
   width: 100%;
   height: 100%;
   overflow: visible;
}

.trend-line {
   stroke: var(--line-color);
   opacity: 0.65;
}

:global(.app-dark) .trend-line {
   stroke: var(--line-color-dark);
   opacity: 0.75;
}
</style>
