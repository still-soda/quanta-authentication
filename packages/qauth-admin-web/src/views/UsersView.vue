<script setup lang="ts">
import { ref, computed } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import Button from 'primevue/button';
import PageHeader from '@/components/shared/PageHeader.vue';
import SimpleStatCard from '@/components/shared/SimpleStatCard.vue';
import UsersTable from '@/components/users/UsersTable.vue';
import UserDialog from '@/components/users/UserDialog.vue';
import type { User } from '@/components/users/UserCell.vue';
import type { SimpleStatData } from '@/types';
import {
   getUsers,
   createUser,
   updateUser,
   deleteUser,
   resetUserPassword,
   disableUser,
} from '@/apis/users';

const queryClient = useQueryClient();

// 获取用户数据
const { data: users, isLoading } = useQuery({
   queryKey: ['users'],
   queryFn: getUsers,
});

// 创建用户 mutation
const createUserMutation = useMutation({
   mutationFn: createUser,
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      userDialog.value = false;
   },
});

// 更新用户 mutation
const updateUserMutation = useMutation({
   mutationFn: ({ id, data }: { id: number; data: any }) => updateUser(id, data),
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      userDialog.value = false;
   },
});

// 删除用户 mutation
const deleteUserMutation = useMutation({
   mutationFn: deleteUser,
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
   },
});

// 重置密码 mutation
const resetPasswordMutation = useMutation({
   mutationFn: resetUserPassword,
});

// 禁用用户 mutation
const disableUserMutation = useMutation({
   mutationFn: disableUser,
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
   },
});

const selectedUsers = ref<User[]>([]);
const userDialog = ref(false);
const isEditing = ref(false);
const currentUser = ref<User | null>(null);

// 统计数据
const stats = computed<SimpleStatData[]>(() => {
   const userList = users.value || [];
   return [
      {
         title: '总用户',
         value: userList.length,
         icon: 'pi pi-users',
         color: 'blue',
      },
      {
         title: '活跃用户',
         value: userList.filter((u) => u.status === 'active').length,
         icon: 'pi pi-check-circle',
         color: 'green',
      },
      {
         title: '未激活',
         value: userList.filter((u) => u.status === 'inactive').length,
         icon: 'pi pi-clock',
         color: 'gray',
      },
      {
         title: '已锁定',
         value: userList.filter((u) => u.status === 'locked').length,
         icon: 'pi pi-lock',
         color: 'red',
      },
   ];
});

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
   if (isEditing.value && currentUser.value) {
      updateUserMutation.mutate({ id: currentUser.value.id, data });
   } else {
      createUserMutation.mutate(data);
   }
};

const handleDelete = (user: User) => {
   deleteUserMutation.mutate(user.id);
};

const handleResetPassword = (user: User) => {
   resetPasswordMutation.mutate(user.id);
};

const handleDisable = (user: User) => {
   disableUserMutation.mutate(user.id);
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
         <template v-if="isLoading">
            <div
               v-for="i in 4"
               :key="i"
               class="h-20 bg-surface-100 dark:bg-surface-800 rounded-xl animate-pulse" />
         </template>
         <template v-else>
            <SimpleStatCard v-for="stat in stats" :key="stat.title" :stat="stat" />
         </template>
      </div>

      <!-- Users Table -->
      <div
         v-if="isLoading"
         class="h-96 bg-surface-100 dark:bg-surface-800 rounded-xl animate-pulse" />
      <UsersTable
         v-else
         :users="users || []"
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
