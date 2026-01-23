<script setup lang="ts">
defineProps<{
   title: string;
   subtitle?: string;
}>();
</script>

<template>
   <div class="page-header">
      <div class="header-content">
         <h1 class="page-title">{{ title }}</h1>
         <p v-if="subtitle" class="page-subtitle">{{ subtitle }}</p>
      </div>
      <div class="header-actions">
         <slot name="actions" />
      </div>
   </div>
</template>

<style scoped>
.page-header {
   display: flex;
   justify-content: space-between;
   align-items: flex-start;
   flex-wrap: wrap;
   gap: 1rem;
}

.header-content {
   display: flex;
   flex-direction: column;
   gap: 0.25rem;
}

.page-title {
   font-size: 1.75rem;
   font-weight: 700;
   color: var(--p-surface-900);
   margin: 0;
   letter-spacing: -0.02em;
}

:global(.app-dark) .page-title {
   color: var(--p-surface-50);
}

.page-subtitle {
   font-size: 0.9375rem;
   color: var(--p-surface-500);
   margin: 0;
}

.header-actions {
   display: flex;
   gap: 0.75rem;
}

@media (max-width: 640px) {
   .page-header {
      flex-direction: column;
   }

   .header-actions {
      width: 100%;
   }

   .header-actions :deep(.p-button) {
      flex: 1;
   }
}
</style>
