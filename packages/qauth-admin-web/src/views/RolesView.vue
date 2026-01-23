<script setup lang="ts">
import { ref } from 'vue';
import Button from 'primevue/button';
import PageHeader from '@/components/shared/PageHeader.vue';
import RoleCard, { type Role } from '@/components/roles/RoleCard.vue';
import RoleDialog, {
   type RoleFormData,
} from '@/components/roles/RoleDialog.vue';
import PermissionDialog, {
   type PermissionGroup,
} from '@/components/roles/PermissionDialog.vue';

// 角色数据
const roles = ref<Role[]>([
   {
      id: 1,
      name: '超级管理员',
      code: 'super_admin',
      description: '拥有系统所有权限',
      userCount: 3,
      permissions: 42,
      status: 'active',
      isSystem: true,
      createdAt: '2024-01-01',
   },
   {
      id: 2,
      name: '管理员',
      code: 'admin',
      description: '管理用户和基本设置',
      userCount: 12,
      permissions: 28,
      status: 'active',
      isSystem: true,
      createdAt: '2024-01-01',
   },
   {
      id: 3,
      name: '开发者',
      code: 'developer',
      description: '管理 OAuth 应用和 API',
      userCount: 45,
      permissions: 15,
      status: 'active',
      isSystem: false,
      createdAt: '2024-03-15',
   },
   {
      id: 4,
      name: '审计员',
      code: 'auditor',
      description: '查看审计日志和报表',
      userCount: 8,
      permissions: 8,
      status: 'active',
      isSystem: false,
      createdAt: '2024-05-20',
   },
   {
      id: 5,
      name: '普通用户',
      code: 'user',
      description: '基本账号访问权限',
      userCount: 847,
      permissions: 5,
      status: 'active',
      isSystem: true,
      createdAt: '2024-01-01',
   },
   {
      id: 6,
      name: '访客',
      code: 'guest',
      description: '只读访问权限',
      userCount: 156,
      permissions: 2,
      status: 'inactive',
      isSystem: false,
      createdAt: '2024-08-10',
   },
]);

// 权限组
const permissionGroups = ref<PermissionGroup[]>([
   {
      name: '用户管理',
      icon: 'pi pi-users',
      permissions: [
         { id: 'user:read', name: '查看用户', checked: true },
         { id: 'user:create', name: '创建用户', checked: true },
         { id: 'user:update', name: '编辑用户', checked: true },
         { id: 'user:delete', name: '删除用户', checked: false },
         { id: 'user:export', name: '导出用户', checked: true },
      ],
   },
   {
      name: '角色权限',
      icon: 'pi pi-shield',
      permissions: [
         { id: 'role:read', name: '查看角色', checked: true },
         { id: 'role:create', name: '创建角色', checked: false },
         { id: 'role:update', name: '编辑角色', checked: false },
         { id: 'role:delete', name: '删除角色', checked: false },
         { id: 'role:assign', name: '分配角色', checked: true },
      ],
   },
   {
      name: 'OAuth 应用',
      icon: 'pi pi-key',
      permissions: [
         { id: 'oauth:read', name: '查看应用', checked: true },
         { id: 'oauth:create', name: '创建应用', checked: true },
         { id: 'oauth:update', name: '编辑应用', checked: true },
         { id: 'oauth:delete', name: '删除应用', checked: false },
         { id: 'oauth:secret', name: '查看密钥', checked: false },
      ],
   },
   {
      name: '系统设置',
      icon: 'pi pi-cog',
      permissions: [
         { id: 'system:read', name: '查看设置', checked: true },
         { id: 'system:update', name: '修改设置', checked: false },
         { id: 'system:audit', name: '审计日志', checked: true },
         { id: 'system:backup', name: '数据备份', checked: false },
      ],
   },
]);

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
   console.log('Saving role:', data);
   roleDialog.value = false;
};

const savePermissions = (groups: PermissionGroup[]) => {
   console.log('Saving permissions for:', selectedRole.value?.name, groups);
   permissionDialog.value = false;
};
</script>

<template>
   <div class="roles-page">
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
      <div class="roles-grid">
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
         :permissionGroups="permissionGroups"
         @save="savePermissions" />
   </div>
</template>

<style scoped>
.roles-page {
   display: flex;
   flex-direction: column;
   gap: 1.5rem;
}

.roles-grid {
   display: grid;
   grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
   gap: 1.25rem;
}

@media (max-width: 768px) {
   .roles-grid {
      grid-template-columns: 1fr;
   }
}
</style>
