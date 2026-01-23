<script setup lang="ts">
import Avatar from 'primevue/avatar';
import type { User } from '@/types';

defineProps<{ user: User }>();
</script>

<template>
   <div class="flex items-center gap-3.5">
      <Avatar :image="user.avatar" shape="circle" size="normal" />
      <div class="flex flex-col gap-0.5">
         <span
            class="font-semibold text-surface-900 dark:text-surface-100">
            {{ user.name }}
         </span>
         <span class="text-[0.8125rem] text-surface-500">
            {{ user.email }}
         </span>
      </div>
   </div>
</template>
