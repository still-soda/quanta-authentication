<script setup lang="ts">
import Avatar from 'primevue/avatar';
import Tag from 'primevue/tag';

export interface User {
   id: number;
   name: string;
   email: string;
   avatar: string;
   role: string;
   status: string;
   lastLogin: string;
   createdAt: string;
}

defineProps<{
   user: User;
}>();

const getRoleSeverity = (role: string) => {
   const map: Record<string, 'danger' | 'warn' | 'info'> = {
      管理员: 'danger',
      开发者: 'warn',
      普通用户: 'info',
   };
   return map[role] || 'info';
};
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
