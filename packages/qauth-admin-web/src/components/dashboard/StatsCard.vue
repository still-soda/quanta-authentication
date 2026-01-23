<script setup lang="ts">
export interface StatCardData {
   title: string;
   value: string;
   change: string;
   changeType: 'increase' | 'decrease';
   icon: string;
   iconBg: string;
}

defineProps<{
   stat: StatCardData;
}>();
</script>

<template>
   <div class="stat-card">
      <div class="stat-icon" :style="{ background: stat.iconBg }">
         <i :class="stat.icon"></i>
      </div>
      <div class="stat-content">
         <span class="stat-title">{{ stat.title }}</span>
         <span class="stat-value">{{ stat.value }}</span>
         <span
            class="stat-change"
            :class="stat.changeType === 'increase' ? 'positive' : 'negative'">
            <i
               :class="
                  stat.changeType === 'increase'
                     ? 'pi pi-arrow-up'
                     : 'pi pi-arrow-down'
               "></i>
            {{ stat.change }}
         </span>
      </div>
      <div class="stat-decoration"></div>
   </div>
</template>

<style scoped>
.stat-card {
   position: relative;
   display: flex;
   align-items: flex-start;
   gap: 1rem;
   padding: 1.5rem;
   background: var(--p-surface-0);
   border-radius: 16px;
   border: 1px solid var(--p-surface-100);
   overflow: hidden;
   transition: all 0.3s ease;
}

.stat-card:hover {
   transform: translateY(-2px);
   box-shadow: 0 12px 24px -8px rgba(0, 0, 0, 0.08);
}

:global(.app-dark) .stat-card {
   background: var(--p-surface-900);
   border-color: var(--p-surface-800);
}

:global(.app-dark) .stat-card:hover {
   box-shadow: 0 12px 24px -8px rgba(0, 0, 0, 0.3);
}

.stat-icon {
   flex-shrink: 0;
   width: 3rem;
   height: 3rem;
   display: flex;
   align-items: center;
   justify-content: center;
   border-radius: 12px;
   color: white;
   font-size: 1.25rem;
   box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.stat-content {
   display: flex;
   flex-direction: column;
   gap: 0.25rem;
   z-index: 1;
}

.stat-title {
   font-size: 0.8125rem;
   color: var(--p-surface-500);
   font-weight: 500;
}

.stat-value {
   font-size: 1.5rem;
   font-weight: 700;
   color: var(--p-surface-900);
   letter-spacing: -0.02em;
}

:global(.app-dark) .stat-value {
   color: var(--p-surface-50);
}

.stat-change {
   display: inline-flex;
   align-items: center;
   gap: 0.25rem;
   font-size: 0.75rem;
   font-weight: 600;
}

.stat-change.positive {
   color: #10b981;
}

.stat-change.negative {
   color: #ef4444;
}

.stat-change i {
   font-size: 0.625rem;
}

.stat-decoration {
   position: absolute;
   right: -2rem;
   bottom: -2rem;
   width: 6rem;
   height: 6rem;
   background: linear-gradient(
      135deg,
      var(--p-orange-100) 0%,
      transparent 100%
   );
   border-radius: 50%;
   opacity: 0.5;
}

:global(.app-dark) .stat-decoration {
   background: linear-gradient(
      135deg,
      rgba(251, 146, 60, 0.1) 0%,
      transparent 100%
   );
}
</style>
