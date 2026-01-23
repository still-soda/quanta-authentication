<script setup lang="ts">
import InputText from 'primevue/inputtext';

const searchQuery = defineModel<string>('modelValue', { default: '' });
</script>

<template>
   <div class="flex gap-4">
      <div class="relative flex items-center flex-1 max-w-96">
         <i
            class="pi pi-search absolute left-3.5 text-surface-400 text-sm pointer-events-none"></i>
         <InputText
            v-model="searchQuery"
            placeholder="搜索应用名称或 Client ID..."
            class="w-full pl-10 h-10 rounded-[10px]" />
      </div>
   </div>
</template>
