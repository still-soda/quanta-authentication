<script setup lang="ts">
import InputText from 'primevue/inputtext';

const searchQuery = defineModel<string>('modelValue', { default: '' });
</script>

<template>
   <div class="search-section">
      <div class="search-box">
         <i class="pi pi-search"></i>
         <InputText
            v-model="searchQuery"
            placeholder="搜索应用名称或 Client ID..."
            class="search-input" />
      </div>
   </div>
</template>

<style scoped>
.search-section {
   display: flex;
   gap: 1rem;
}

.search-box {
   position: relative;
   display: flex;
   align-items: center;
   flex: 1;
   max-width: 24rem;
}

.search-box i {
   position: absolute;
   left: 0.875rem;
   color: var(--p-surface-400);
   font-size: 0.875rem;
}

.search-input {
   width: 100%;
   padding-left: 2.5rem;
   height: 2.5rem;
   border-radius: 10px;
}
</style>
