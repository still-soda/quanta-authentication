<script setup lang="ts">
import { ref } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import Button from 'primevue/button';
import PageHeader from '@/components/shared/PageHeader.vue';
import RoleCard from '@/components/roles/RoleCard.vue';
import RoleDialog from '@/components/roles/RoleDialog.vue';
import PermissionDialog from '@/components/roles/PermissionDialog.vue';
import type { Role, RoleFormData, PermissionGroup } from '@/types';
import {
   getRoles,
   createRole,
   updateRole,
   getPermissionGroups,
   saveRolePermissions,
} from '@/apis/roles';

const queryClient = useQueryClient();

// 使用 TanStack Query 获取角色数据
const { data: roles, isLoading: isLoadingRoles } = useQuery({
   queryKey: ['roles'],
   queryFn: getRoles,
});

// 使用 TanStack Query 获取权限组数据
const { data: permissionGroups, isLoading: isLoadingPermissions } = useQuery({
   queryKey: ['permissions'],
   queryFn: getPermissionGroups,
});

// 创建角色 mutation
const createRoleMutation = useMutation({
   mutationFn: createRole,
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['roles'] });
      roleDialog.value = false;
   },
});

// 更新角色 mutation
const updateRoleMutation = useMutation({
   mutationFn: ({ id, data }: { id: number; data: Partial<RoleFormData> }) =>
      updateRole(id, data),
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['roles'] });
      roleDialog.value = false;
   },
});

// 保存权限 mutation
const savePermissionsMutation = useMutation({
   mutationFn: ({ roleId, permissions }: { roleId: number; permissions: PermissionGroup[] }) =>
      saveRolePermissions(roleId, permissions),
   onSuccess: () => {
      permissionDialog.value = false;
   },
});

const selectedRole = ref<Role | null>(null);
const roleDialog = ref(false);
const permissionDialog = ref(false);
const isEditing = ref(false);

const openNewRoleDialog = () => {
   isEditing.value = false;
   selectedRole.value = null;
   roleDialog.value = true;
};

const editRole = (role: Role) => {
   isEditing.value = true;
   selectedRole.value = role;
   roleDialog.value = true;
};

const openPermissionDialog = (role: Role) => {
   selectedRole.value = role;
   permissionDialog.value = true;
};

const saveRole = (data: RoleFormData) => {
   if (isEditing.value && selectedRole.value) {
      updateRoleMutation.mutate({ id: selectedRole.value.id, data });
   } else {
      createRoleMutation.mutate(data);
   }
};

const savePermissions = (groups: PermissionGroup[]) => {
   if (selectedRole.value) {
      savePermissionsMutation.mutate({
         roleId: selectedRole.value.id,
         permissions: groups,
      });
   }
};
</script>

<template>
   <div class="flex flex-col gap-6">
      <!-- Page Header -->
      <PageHeader title="角色权限" subtitle="管理系统角色和权限配置">
         <template #actions>
            <Button
               label="新建角色"
               icon="pi pi-plus"
               @click="openNewRoleDialog" />
         </template>
      </PageHeader>

      <!-- Roles Grid -->
      <div
         v-if="isLoadingRoles"
         class="grid grid-cols-1 md:grid-cols-[repeat(auto-fill,minmax(320px,1fr))] gap-5">
         <div
            v-for="i in 6"
            :key="i"
            class="h-48 bg-surface-100 dark:bg-surface-800 rounded-xl animate-pulse" />
      </div>
      <div
         v-else
         class="grid grid-cols-1 md:grid-cols-[repeat(auto-fill,minmax(320px,1fr))] gap-5">
         <RoleCard
            v-for="role in roles"
            :key="role.id"
            :role="role"
            @edit="editRole"
            @configPermissions="openPermissionDialog" />
      </div>

      <!-- Role Dialog -->
      <RoleDialog
         v-model:visible="roleDialog"
         :isEditing="isEditing"
         :initialData="
            selectedRole
               ? {
                    name: selectedRole.name,
                    code: selectedRole.code,
                    description: selectedRole.description,
                 }
               : undefined
         "
         @save="saveRole" />

      <!-- Permission Dialog -->
      <PermissionDialog
         v-model:visible="permissionDialog"
         :role="selectedRole"
         :permissionGroups="permissionGroups || []"
         @save="savePermissions" />
   </div>
</template>
