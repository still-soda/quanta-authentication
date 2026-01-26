<script setup lang="ts">
import { ref, computed } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import PageHeader from '@/components/shared/PageHeader.vue'
import RoleCard from '@/components/roles/RoleCard.vue'
import RoleDialog from '@/components/roles/RoleDialog.vue'
import PermissionDialog from '@/components/roles/PermissionDialog.vue'
import DeleteRoleDialog from '@/components/roles/DeleteRoleDialog.vue'
import type { Role, RoleFormData, PermissionGroup } from '@/types'
import {
   getRoles,
   createRole,
   updateRole,
   deleteRole,
   getPermissions,
   getRolePermissions,
   setRolePermissions,
   transformPermissionsToGroups,
   extractCheckedPermissionCodes,
} from '@/apis/roles'

const queryClient = useQueryClient()
const toast = useToast()

// 获取角色数据（获取所有）
const {
   data: roles,
   isLoading: isLoadingRoles,
   error: rolesError,
} = useQuery({
   queryKey: ['roles', 'all'],
   queryFn: () => getRoles({ all: true }),
})

// 获取所有权限数据（用于权限配置,获取所有）
const { data: allPermissions, isLoading: isLoadingPermissions } = useQuery({
   queryKey: ['permissions', 'all'],
   queryFn: () => getPermissions({ all: true }),
})

// 状态管理
const selectedRole = ref<Role | null>(null)
const roleDialog = ref(false)
const permissionDialog = ref(false)
const deleteDialog = ref(false)
const isEditing = ref(false)
const isDeleting = ref(false)
const isSavingPermissions = ref(false)

// 当前选中角色的权限
const currentRolePermissionCodes = ref<string[]>([])

// 权限组（基于所有权限和当前角色权限）
const permissionGroups = computed<PermissionGroup[]>(() => {
   if (!allPermissions.value) return []
   return transformPermissionsToGroups(allPermissions.value, currentRolePermissionCodes.value)
})

// 创建角色 mutation
const createRoleMutation = useMutation({
   mutationFn: (data: RoleFormData & { permissions?: string[] }) => createRole(data),
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['roles'] })
      roleDialog.value = false
      toast.add({
         severity: 'success',
         summary: '创建成功',
         detail: '角色已成功创建',
         life: 3000,
      })
   },
   onError: (error: any) => {
      toast.add({
         severity: 'error',
         summary: '创建失败',
         detail: error.response?.data?.message || '创建角色时发生错误',
         life: 5000,
      })
   },
})

// 更新角色 mutation
const updateRoleMutation = useMutation({
   mutationFn: ({ id, data }: { id: string; data: RoleFormData }) => updateRole(id, data),
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['roles'] })
      roleDialog.value = false
      toast.add({
         severity: 'success',
         summary: '更新成功',
         detail: '角色信息已更新',
         life: 3000,
      })
   },
   onError: (error: any) => {
      toast.add({
         severity: 'error',
         summary: '更新失败',
         detail: error.response?.data?.message || '更新角色时发生错误',
         life: 5000,
      })
   },
})

// 删除角色 mutation
const deleteRoleMutation = useMutation({
   mutationFn: (id: string) => deleteRole(id),
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['roles'] })
      deleteDialog.value = false
      isDeleting.value = false
      toast.add({
         severity: 'success',
         summary: '删除成功',
         detail: '角色已成功删除',
         life: 3000,
      })
   },
   onError: (error: any) => {
      isDeleting.value = false
      toast.add({
         severity: 'error',
         summary: '删除失败',
         detail: error.response?.data?.message || '删除角色时发生错误',
         life: 5000,
      })
   },
})

// 保存权限 mutation
const savePermissionsMutation = useMutation({
   mutationFn: ({ roleId, codes }: { roleId: string; codes: string[] }) =>
      setRolePermissions(roleId, codes),
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['roles'] })
      permissionDialog.value = false
      isSavingPermissions.value = false
      toast.add({
         severity: 'success',
         summary: '保存成功',
         detail: '角色权限已更新',
         life: 3000,
      })
   },
   onError: (error: any) => {
      isSavingPermissions.value = false
      toast.add({
         severity: 'error',
         summary: '保存失败',
         detail: error.response?.data?.message || '保存权限时发生错误',
         life: 5000,
      })
   },
})

// 打开新建角色对话框
const openNewRoleDialog = () => {
   isEditing.value = false
   selectedRole.value = null
   roleDialog.value = true
}

// 编辑角色
const editRole = (role: Role) => {
   isEditing.value = true
   selectedRole.value = role
   roleDialog.value = true
}

// 打开权限配置对话框
const openPermissionDialog = async (role: Role) => {
   selectedRole.value = role
   // 获取角色当前的权限
   try {
      const permissions = await getRolePermissions(role.id)
      currentRolePermissionCodes.value = permissions.map(p => p.code)
   } catch (error) {
      currentRolePermissionCodes.value = []
   }
   permissionDialog.value = true
}

// 打开删除对话框
const openDeleteDialog = (role: Role) => {
   selectedRole.value = role
   deleteDialog.value = true
}

// 保存角色
const saveRole = (data: RoleFormData) => {
   if (isEditing.value && selectedRole.value) {
      updateRoleMutation.mutate({ id: selectedRole.value.id, data })
   } else {
      createRoleMutation.mutate(data)
   }
}

// 保存权限
const savePermissions = (groups: PermissionGroup[]) => {
   if (selectedRole.value) {
      isSavingPermissions.value = true
      const codes = extractCheckedPermissionCodes(groups)
      savePermissionsMutation.mutate({
         roleId: selectedRole.value.id,
         codes,
      })
   }
}

// 确认删除角色
const confirmDeleteRole = () => {
   if (selectedRole.value) {
      isDeleting.value = true
      deleteRoleMutation.mutate(selectedRole.value.id)
   }
}

// 排序角色：系统角色优先
const sortedRoles = computed(() => {
   if (!roles.value) return []
   return [...roles.value].sort((a, b) => {
      if (a.is_system && !b.is_system) return -1
      if (!a.is_system && b.is_system) return 1
      return 0
   })
})
</script>

<template>
   <div class="flex flex-col gap-6">
      <!-- Page Header -->
      <PageHeader title="角色权限" subtitle="管理系统角色和权限配置">
         <template #actions>
            <Button label="新建角色" icon="pi pi-plus" @click="openNewRoleDialog" />
         </template>
      </PageHeader>

      <!-- Error State -->
      <div v-if="rolesError" class="flex flex-col items-center justify-center py-16 text-center">
         <i class="pi pi-exclamation-triangle text-5xl text-amber-500 mb-4"></i>
         <p class="text-lg text-surface-600 dark:text-surface-400 mb-4">加载角色数据失败</p>
         <Button
            label="重试"
            icon="pi pi-refresh"
            severity="secondary"
            outlined
            @click="() => queryClient.invalidateQueries({ queryKey: ['roles'] })"
         />
      </div>

      <!-- Loading State -->
      <div
         v-else-if="isLoadingRoles"
         class="grid grid-cols-1 md:grid-cols-[repeat(auto-fill,minmax(320px,1fr))] gap-5"
      >
         <div
            v-for="i in 6"
            :key="i"
            class="h-52 bg-surface-100 dark:bg-surface-800 rounded-xl animate-pulse"
         />
      </div>

      <!-- Empty State -->
      <div
         v-else-if="!sortedRoles.length"
         class="flex flex-col items-center justify-center py-16 text-center"
      >
         <i class="pi pi-shield text-5xl text-surface-300 dark:text-surface-600 mb-4"></i>
         <p class="text-lg text-surface-600 dark:text-surface-400 mb-4">暂无角色</p>
         <Button label="创建第一个角色" icon="pi pi-plus" @click="openNewRoleDialog" />
      </div>

      <!-- Roles Grid -->
      <div v-else class="grid grid-cols-1 md:grid-cols-[repeat(auto-fill,minmax(320px,1fr))] gap-5">
         <RoleCard
            v-for="role in sortedRoles"
            :key="role.id"
            :role="role"
            @edit="editRole"
            @delete="openDeleteDialog"
            @configPermissions="openPermissionDialog"
         />
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
         :isLoading="createRoleMutation.isPending.value || updateRoleMutation.isPending.value"
         @save="saveRole"
      />

      <!-- Permission Dialog -->
      <PermissionDialog
         v-model:visible="permissionDialog"
         :role="selectedRole"
         :permissionGroups="permissionGroups"
         :isLoading="isSavingPermissions || isLoadingPermissions"
         @save="savePermissions"
      />

      <!-- Delete Role Dialog -->
      <DeleteRoleDialog
         v-model:visible="deleteDialog"
         :role="selectedRole"
         :isLoading="isDeleting"
         @confirm="confirmDeleteRole"
      />
   </div>
</template>
