<script setup lang="ts">
defineProps<{
   title: string;
   subtitle?: string;
}>();
</script>

<template>
   <div
      class="flex justify-between items-start flex-wrap gap-4 max-sm:flex-col">
      <div class="flex flex-col gap-1">
         <h1
            class="text-[1.75rem] font-bold text-surface-900 dark:text-surface-50 m-0 tracking-tight">
            {{ title }}
         </h1>
         <p v-if="subtitle" class="text-[0.9375rem] text-surface-500 m-0">
            {{ subtitle }}
         </p>
      </div>
      <div class="flex gap-3 max-sm:w-full max-sm:*:flex-1">
         <slot name="actions" />
      </div>
   </div>
</template>
