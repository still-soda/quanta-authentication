<script setup lang="ts">
export interface MiniStatItem {
   label: string;
   value: number | string;
   colorClass?: string;
}

defineProps<{
   items: MiniStatItem[];
}>();
</script>

<template>
   <div class="stats-row">
      <div v-for="stat in items" :key="stat.label" class="mini-stat">
         <span class="mini-stat-value" :class="stat.colorClass">
            {{ stat.value }}
         </span>
         <span class="mini-stat-label">{{ stat.label }}</span>
      </div>
   </div>
</template>

<style scoped>
.stats-row {
   display: flex;
   gap: 2rem;
   padding: 1rem 1.5rem;
   background: var(--p-surface-0);
   border-radius: 12px;
   border: 1px solid var(--p-surface-100);
}

:global(.app-dark) .stats-row {
   background: var(--p-surface-900);
   border-color: var(--p-surface-800);
}

.mini-stat {
   display: flex;
   flex-direction: column;
   gap: 0.25rem;
}

.mini-stat-value {
   font-size: 1.5rem;
   font-weight: 700;
   color: var(--p-surface-900);
   letter-spacing: -0.02em;
}

:global(.app-dark) .mini-stat-value {
   color: var(--p-surface-100);
}

.mini-stat-label {
   font-size: 0.8125rem;
   color: var(--p-surface-500);
}
</style>
