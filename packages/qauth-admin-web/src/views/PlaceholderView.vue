<script setup lang="ts">
import { useRoute } from 'vue-router';
import Card from 'primevue/card';

const route = useRoute();
const title = (route.meta?.title as string) || '页面';
</script>

<template>
   <div class="placeholder-page">
      <div class="page-header">
         <div class="header-content">
            <h1 class="page-title">{{ title }}</h1>
            <p class="page-subtitle">此页面正在开发中</p>
         </div>
      </div>

      <Card class="placeholder-card">
         <template #content>
            <div class="placeholder-content">
               <div class="placeholder-icon">
                  <i class="pi pi-wrench"></i>
               </div>
               <h2>功能开发中</h2>
               <p>{{ title }}功能正在紧张开发中，敬请期待！</p>
            </div>
         </template>
      </Card>
   </div>
</template>

<style scoped>
.placeholder-page {
   display: flex;
   flex-direction: column;
   gap: 1.5rem;
}

.page-header {
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

.placeholder-card {
   border-radius: 16px;
   border: 1px solid var(--p-surface-100);
}

:global(.app-dark) .placeholder-card {
   border-color: var(--p-surface-800);
}

.placeholder-content {
   display: flex;
   flex-direction: column;
   align-items: center;
   justify-content: center;
   padding: 4rem 2rem;
   text-align: center;
}

.placeholder-icon {
   width: 5rem;
   height: 5rem;
   display: flex;
   align-items: center;
   justify-content: center;
   background: linear-gradient(
      135deg,
      var(--p-orange-100) 0%,
      var(--p-orange-200) 100%
   );
   border-radius: 50%;
   margin-bottom: 1.5rem;
}

:global(.app-dark) .placeholder-icon {
   background: linear-gradient(
      135deg,
      rgba(251, 146, 60, 0.15) 0%,
      rgba(251, 146, 60, 0.25) 100%
   );
}

.placeholder-icon i {
   font-size: 2rem;
   color: var(--p-orange-600);
}

:global(.app-dark) .placeholder-icon i {
   color: var(--p-orange-400);
}

.placeholder-content h2 {
   font-size: 1.25rem;
   font-weight: 600;
   color: var(--p-surface-900);
   margin: 0 0 0.5rem 0;
}

:global(.app-dark) .placeholder-content h2 {
   color: var(--p-surface-100);
}

.placeholder-content p {
   font-size: 0.9375rem;
   color: var(--p-surface-500);
   margin: 0;
}
</style>
