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
   <div class="user-cell">
      <Avatar :image="user.avatar" shape="circle" size="normal" />
      <div class="user-info">
         <span class="user-name">{{ user.name }}</span>
         <span class="user-email">{{ user.email }}</span>
      </div>
   </div>
</template>

<style scoped>
.user-cell {
   display: flex;
   align-items: center;
   gap: 0.875rem;
}

.user-info {
   display: flex;
   flex-direction: column;
   gap: 0.125rem;
}

.user-name {
   font-weight: 600;
   color: var(--p-surface-900);
}

:global(.app-dark) .user-name {
   color: var(--p-surface-100);
}

.user-email {
   font-size: 0.8125rem;
   color: var(--p-surface-500);
}
</style>
