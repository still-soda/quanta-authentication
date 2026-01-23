<script setup lang="ts">
export interface SimpleStatData {
   title: string;
   value: number | string;
   icon: string;
   /** 卡片主题色 */
   color?: 'orange' | 'blue' | 'green' | 'purple' | 'cyan' | 'red' | 'gray';
}

withDefaults(
   defineProps<{
      stat: SimpleStatData;
   }>(),
   {},
);
</script>

<template>
   <div class="simple-stat-card" :class="`theme-${stat.color || 'orange'}`">
      <div class="stat-icon-wrapper">
         <i :class="stat.icon" class="stat-icon"></i>
      </div>
      <div class="stat-content">
         <span class="stat-value">{{ stat.value }}</span>
         <span class="stat-title">{{ stat.title }}</span>
      </div>
   </div>
</template>

<style scoped>
.simple-stat-card {
   display: flex;
   align-items: center;
   gap: 1rem;
   padding: 1.25rem 1.5rem;
   background: var(--p-surface-0);
   border: 1px solid var(--p-surface-200);
   border-radius: 14px;
   transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.simple-stat-card:hover {
   transform: translateY(-2px);
   box-shadow: 0 8px 20px -6px rgba(0, 0, 0, 0.08);
   border-color: var(--theme-color);
}

:global(.app-dark) .simple-stat-card {
   background: var(--p-surface-900);
   border-color: var(--p-surface-800);
}

:global(.app-dark) .simple-stat-card:hover {
   box-shadow: 0 8px 20px -6px rgba(0, 0, 0, 0.35);
}

/* Theme Colors */
.theme-orange {
   --theme-color: var(--p-orange-500);
   --theme-bg: var(--p-orange-50);
   --theme-bg-dark: rgba(251, 146, 60, 0.15);
}

.theme-blue {
   --theme-color: var(--p-blue-500);
   --theme-bg: var(--p-blue-50);
   --theme-bg-dark: rgba(59, 130, 246, 0.15);
}

.theme-green {
   --theme-color: var(--p-green-500);
   --theme-bg: var(--p-green-50);
   --theme-bg-dark: rgba(34, 197, 94, 0.15);
}

.theme-purple {
   --theme-color: var(--p-purple-500);
   --theme-bg: var(--p-purple-50);
   --theme-bg-dark: rgba(168, 85, 247, 0.15);
}

.theme-cyan {
   --theme-color: var(--p-cyan-500);
   --theme-bg: var(--p-cyan-50);
   --theme-bg-dark: rgba(6, 182, 212, 0.15);
}

.theme-red {
   --theme-color: var(--p-red-500);
   --theme-bg: var(--p-red-50);
   --theme-bg-dark: rgba(239, 68, 68, 0.15);
}

.theme-gray {
   --theme-color: var(--p-surface-500);
   --theme-bg: var(--p-surface-100);
   --theme-bg-dark: rgba(107, 114, 128, 0.15);
}

/* Icon */
.stat-icon-wrapper {
   width: 3rem;
   height: 3rem;
   display: flex;
   align-items: center;
   justify-content: center;
   border-radius: 12px;
   background: var(--theme-bg);
   color: var(--theme-color);
   flex-shrink: 0;
   transition: all 0.25s ease;
}

:global(.app-dark) .stat-icon-wrapper {
   background: var(--theme-bg-dark);
}

.simple-stat-card:hover .stat-icon-wrapper {
   background: var(--theme-color);
   color: #ffffff;
   transform: scale(1.05);
}

.stat-icon {
   font-size: 1.25rem;
}

/* Content */
.stat-content {
   display: flex;
   flex-direction: column;
   gap: 0.125rem;
}

.stat-value {
   font-size: 1.5rem;
   font-weight: 700;
   color: var(--p-surface-900);
   letter-spacing: -0.02em;
   line-height: 1.2;
}

:global(.app-dark) .stat-value {
   color: var(--p-surface-0);
}

.stat-title {
   font-size: 0.8125rem;
   font-weight: 500;
   color: var(--p-surface-500);
}

:global(.app-dark) .stat-title {
   color: var(--p-surface-400);
}
</style>
