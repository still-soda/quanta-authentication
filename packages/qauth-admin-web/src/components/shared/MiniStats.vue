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
   <div
      class="flex gap-8 p-4 px-6 bg-surface-0 dark:bg-surface-900 rounded-xl border border-surface-100 dark:border-surface-800">
      <div v-for="stat in items" :key="stat.label" class="flex flex-col gap-1">
         <span
            class="text-2xl font-bold text-surface-900 dark:text-surface-100 tracking-tight"
            :class="stat.colorClass">
            {{ stat.value }}
         </span>
         <span class="text-[0.8125rem] text-surface-500">
            {{ stat.label }}
         </span>
      </div>
   </div>
</template>
