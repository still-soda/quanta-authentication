<script setup lang="ts">
import { ref, computed } from 'vue';
import Button from 'primevue/button';
import PageHeader from '@/components/shared/PageHeader.vue';
import SimpleStatCard, {
   type SimpleStatData,
} from '@/components/shared/SimpleStatCard.vue';
import UsersTable from '@/components/users/UsersTable.vue';
import UserDialog from '@/components/users/UserDialog.vue';
import type { User } from '@/components/users/UserCell.vue';

// 用户数据
const users = ref<User[]>([
   {
      id: 1,
      name: '张伟',
      email: 'zhang.wei@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhang',
      role: '管理员',
      status: 'active',
      lastLogin: '2026-01-23 10:30',
      createdAt: '2024-06-15',
   },
   {
      id: 2,
      name: '李明',
      email: 'li.ming@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=li',
      role: '开发者',
      status: 'active',
      lastLogin: '2026-01-22 18:45',
      createdAt: '2024-08-22',
   },
   {
      id: 3,
      name: '王芳',
      email: 'wang.fang@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=wang',
      role: '普通用户',
      status: 'inactive',
      lastLogin: '2026-01-10 09:15',
      createdAt: '2024-09-03',
   },
   {
      id: 4,
      name: '陈红',
      email: 'chen.hong@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=chen',
      role: '普通用户',
      status: 'active',
      lastLogin: '2026-01-23 08:20',
      createdAt: '2024-11-18',
   },
   {
      id: 5,
      name: '赵阳',
      email: 'zhao.yang@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhao',
      role: '开发者',
      status: 'locked',
      lastLogin: '2026-01-05 14:30',
      createdAt: '2025-01-02',
   },
   {
      id: 6,
      name: '刘洋',
      email: 'liu.yang@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=liu',
      role: '普通用户',
      status: 'active',
      lastLogin: '2026-01-21 16:45',
      createdAt: '2025-03-10',
   },
   {
      id: 7,
      name: '孙静',
      email: 'sun.jing@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=sun',
      role: '管理员',
      status: 'active',
      lastLogin: '2026-01-23 11:00',
      createdAt: '2024-05-20',
   },
   {
      id: 8,
      name: '周杰',
      email: 'zhou.jie@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhou',
      role: '开发者',
      status: 'pending',
      lastLogin: '-',
      createdAt: '2026-01-20',
   },
]);

const selectedUsers = ref<User[]>([]);
const userDialog = ref(false);
const isEditing = ref(false);
const currentUser = ref<User | null>(null);

// 统计数据
const stats = computed<SimpleStatData[]>(() => [
   {
      title: '总用户',
      value: users.value.length,
      icon: 'pi pi-users',
      color: 'blue',
   },
   {
      title: '活跃用户',
      value: users.value.filter((u) => u.status === 'active').length,
      icon: 'pi pi-check-circle',
      color: 'green',
   },
   {
      title: '未激活',
      value: users.value.filter((u) => u.status === 'inactive').length,
      icon: 'pi pi-clock',
      color: 'gray',
   },
   {
      title: '已锁定',
      value: users.value.filter((u) => u.status === 'locked').length,
      icon: 'pi pi-lock',
      color: 'red',
   },
]);

const openNewUserDialog = () => {
   isEditing.value = false;
   currentUser.value = null;
   userDialog.value = true;
};

const editUser = (user: User) => {
   isEditing.value = true;
   currentUser.value = user;
   userDialog.value = true;
};

const saveUser = (data: any) => {
   console.log('Saving user:', data);
   userDialog.value = false;
};

const handleDelete = (user: User) => {
   console.log('Delete user:', user);
};

const handleResetPassword = (user: User) => {
   console.log('Reset password for:', user);
};

const handleDisable = (user: User) => {
   console.log('Disable user:', user);
};
</script>

<template>
   <div class="flex flex-col gap-6">
      <!-- Page Header -->
      <PageHeader title="用户管理" subtitle="管理系统用户账号和权限">
         <template #actions>
            <Button
               label="导入用户"
               icon="pi pi-upload"
               severity="secondary"
               outlined />
            <Button
               label="新建用户"
               icon="pi pi-plus"
               @click="openNewUserDialog" />
         </template>
      </PageHeader>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4">
         <SimpleStatCard v-for="stat in stats" :key="stat.title" :stat="stat" />
      </div>

      <!-- Users Table -->
      <UsersTable
         :users="users"
         v-model:selectedUsers="selectedUsers"
         @edit="editUser"
         @delete="handleDelete"
         @resetPassword="handleResetPassword"
         @disable="handleDisable" />

      <!-- User Dialog -->
      <UserDialog
         v-model:visible="userDialog"
         :isEditing="isEditing"
         :initialData="
            currentUser
               ? {
                    name: currentUser.name,
                    email: currentUser.email,
                    role: currentUser.role,
                    status: currentUser.status === 'active',
                 }
               : undefined
         "
         @save="saveUser" />
   </div>
</template>
