<script setup lang="ts">
import { useRoute } from 'vue-router';
import Card from 'primevue/card';

const route = useRoute();
const title = (route.meta?.title as string) || '页面';
</script>

<template>
   <div class="flex flex-col gap-6">
      <div class="flex flex-col gap-1">
         <h1
            class="text-[1.75rem] font-bold text-surface-900 dark:text-surface-50 m-0 tracking-tight">
            {{ title }}
         </h1>
         <p class="text-[0.9375rem] text-surface-500 m-0">
            此页面正在开发中
         </p>
      </div>

      <Card
         class="rounded-2xl border border-surface-100 dark:border-surface-800">
         <template #content>
            <div
               class="flex flex-col items-center justify-center py-16 px-8 text-center">
               <div
                  class="w-20 h-20 flex items-center justify-center bg-linear-to-br from-primary-100 to-primary-200 dark:from-[rgba(251,146,60,0.15)] dark:to-[rgba(251,146,60,0.25)] rounded-full mb-6">
                  <i
                     class="pi pi-wrench text-[2rem] text-primary-600 dark:text-primary-400"></i>
               </div>
               <h2
                  class="text-xl font-semibold text-surface-900 dark:text-surface-100 m-0 mb-2">
                  功能开发中
               </h2>
               <p class="text-[0.9375rem] text-surface-500 m-0">
                  {{ title }}功能正在紧张开发中，敬请期待！
               </p>
            </div>
         </template>
      </Card>
   </div>
</template>
