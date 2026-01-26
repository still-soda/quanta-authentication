<script setup lang="ts">
import { ref, computed } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputText from 'primevue/inputtext'
import Tag from 'primevue/tag'
import PageHeader from '@/components/shared/PageHeader.vue'
import PermissionFormDialog from '@/components/permissions/PermissionFormDialog.vue'
import DeletePermissionDialog from '@/components/permissions/DeletePermissionDialog.vue'
import type { Permission, PermissionFormData, PermissionAction } from '@/types'
import { getPermissions, createPermission, updatePermission, deletePermission } from '@/apis/roles'

const queryClient = useQueryClient()
const toast = useToast()

// 获取权限数据（获取所有用于管理）
const {
   data: permissions,
   isLoading,
   error,
} = useQuery({
   queryKey: ['permissions', 'all'],
   queryFn: () => getPermissions({ all: true }),
})

// 状态管理
const selectedPermission = ref<Permission | null>(null)
const formDialog = ref(false)
const deleteDialog = ref(false)
const isEditing = ref(false)
const isDeleting = ref(false)
const searchKeyword = ref('')

// 资源名称映射
const resourceNames: Record<string, string> = {
   oauth_clients: 'OAuth 应用',
   system: '系统管理',
   roles: '角色管理',
   audit: '审计日志',
   users: '用户管理',
   permissions: '权限管理',
}

// 操作类型映射
const actionNames: Record<number, string> = {
   1: '创建',
   2: '查看',
   3: '更新',
   4: '删除',
}

const actionSeverities: Record<number, string> = {
   1: 'success',
   2: 'info',
   3: 'warn',
   4: 'danger',
}

// 过滤后的权限列表
const filteredPermissions = computed(() => {
   if (!permissions.value) return []
   if (!searchKeyword.value.trim()) return permissions.value

   const keyword = searchKeyword.value.toLowerCase()
   return permissions.value.filter(
      p =>
         p.code.toLowerCase().includes(keyword) ||
         p.description.toLowerCase().includes(keyword) ||
         p.resource.toLowerCase().includes(keyword)
   )
})

// 创建权限 mutation
const createMutation = useMutation({
   mutationFn: (data: PermissionFormData) =>
      createPermission({
         resource: data.resource,
         action: data.action,
         code: data.code,
         description: data.description,
      }),
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['permissions'] })
      formDialog.value = false
      toast.add({
         severity: 'success',
         summary: '创建成功',
         detail: '权限已成功创建',
         life: 3000,
      })
   },
   onError: (error: any) => {
      toast.add({
         severity: 'error',
         summary: '创建失败',
         detail: error.response?.data?.message || '创建权限时发生错误',
         life: 5000,
      })
   },
})

// 更新权限 mutation
const updateMutation = useMutation({
   mutationFn: ({ id, data }: { id: string; data: PermissionFormData }) =>
      updatePermission(id, {
         resource: data.resource,
         action: data.action,
         code: data.code,
         description: data.description,
      }),
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['permissions'] })
      formDialog.value = false
      toast.add({
         severity: 'success',
         summary: '更新成功',
         detail: '权限信息已更新',
         life: 3000,
      })
   },
   onError: (error: any) => {
      toast.add({
         severity: 'error',
         summary: '更新失败',
         detail: error.response?.data?.message || '更新权限时发生错误',
         life: 5000,
      })
   },
})

// 删除权限 mutation
const deleteMutation = useMutation({
   mutationFn: (id: string) => deletePermission(id),
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['permissions'] })
      deleteDialog.value = false
      isDeleting.value = false
      toast.add({
         severity: 'success',
         summary: '删除成功',
         detail: '权限已成功删除',
         life: 3000,
      })
   },
   onError: (error: any) => {
      isDeleting.value = false
      toast.add({
         severity: 'error',
         summary: '删除失败',
         detail: error.response?.data?.message || '删除权限时发生错误',
         life: 5000,
      })
   },
})

// 打开新建对话框
const openCreateDialog = () => {
   isEditing.value = false
   selectedPermission.value = null
   formDialog.value = true
}

// 打开编辑对话框
const openEditDialog = (permission: Permission) => {
   isEditing.value = true
   selectedPermission.value = permission
   formDialog.value = true
}

// 打开删除对话框
const openDeleteDialog = (permission: Permission) => {
   selectedPermission.value = permission
   deleteDialog.value = true
}

// 保存权限
const savePermission = (data: PermissionFormData) => {
   if (isEditing.value && selectedPermission.value) {
      updateMutation.mutate({ id: selectedPermission.value.id, data })
   } else {
      createMutation.mutate(data)
   }
}

// 确认删除
const confirmDelete = () => {
   if (selectedPermission.value) {
      isDeleting.value = true
      deleteMutation.mutate(selectedPermission.value.id)
   }
}

// 格式化日期
const formatDate = (dateStr: string) => {
   if (!dateStr) return '-'
   const date = new Date(dateStr)
   return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
   })
}
</script>

<template>
   <div class="flex flex-col gap-6">
      <!-- Page Header -->
      <PageHeader title="权限管理" subtitle="管理系统中的所有权限定义">
         <template #actions>
            <Button label="新建权限" icon="pi pi-plus" @click="openCreateDialog" />
         </template>
      </PageHeader>

      <!-- Search -->
      <div class="flex items-center gap-4">
         <div class="flex-1 max-w-md">
            <InputText
               v-model="searchKeyword"
               placeholder="搜索权限代码、描述或资源..."
               class="w-full"
            />
         </div>
         <span class="text-sm text-surface-500 dark:text-surface-400">
            共 {{ filteredPermissions.length }} 条权限
         </span>
      </div>

      <!-- Error State -->
      <div v-if="error" class="flex flex-col items-center justify-center py-16 text-center">
         <i class="pi pi-exclamation-triangle text-5xl text-amber-500 mb-4"></i>
         <p class="text-lg text-surface-600 dark:text-surface-400 mb-4">加载权限数据失败</p>
         <Button
            label="重试"
            icon="pi pi-refresh"
            severity="secondary"
            outlined
            @click="() => queryClient.invalidateQueries({ queryKey: ['permissions'] })"
         />
      </div>

      <!-- Data Table -->
      <DataTable
         v-else
         :value="filteredPermissions"
         :loading="isLoading"
         paginator
         :rows="15"
         :rowsPerPageOptions="[10, 15, 25, 50]"
         stripedRows
         removableSort
         class="rounded-xl border border-surface-100 dark:border-surface-800 overflow-hidden"
      >
         <template #empty>
            <div class="flex flex-col items-center justify-center py-12">
               <i class="pi pi-lock text-4xl text-surface-300 dark:text-surface-600 mb-3"></i>
               <p class="text-surface-500 dark:text-surface-400">
                  {{ searchKeyword ? '未找到匹配的权限' : '暂无权限数据' }}
               </p>
            </div>
         </template>

         <Column field="code" header="权限代码" sortable style="min-width: 200px">
            <template #body="{ data }">
               <code
                  class="text-sm font-mono text-primary-600 dark:text-primary-400 bg-primary-50 dark:bg-primary-900/20 px-2 py-0.5 rounded"
               >
                  {{ data.code }}
               </code>
            </template>
         </Column>

         <Column field="description" header="描述" style="min-width: 200px">
            <template #body="{ data }">
               <span class="text-surface-700 dark:text-surface-300">
                  {{ data.description || '-' }}
               </span>
            </template>
         </Column>

         <Column field="resource" header="资源" sortable style="min-width: 120px">
            <template #body="{ data }">
               <span class="text-surface-600 dark:text-surface-400">
                  {{ resourceNames[data.resource] || data.resource }}
               </span>
            </template>
         </Column>

         <Column field="action" header="操作类型" sortable style="min-width: 100px">
            <template #body="{ data }">
               <Tag :severity="actionSeverities[data.action] || 'secondary'" class="text-xs">
                  {{ actionNames[data.action] || data.action }}
               </Tag>
            </template>
         </Column>

         <Column field="created_at" header="创建时间" sortable style="min-width: 120px">
            <template #body="{ data }">
               <span class="text-sm text-surface-500 dark:text-surface-400">
                  {{ formatDate(data.created_at) }}
               </span>
            </template>
         </Column>

         <Column header="操作" style="width: 100px" frozen alignFrozen="right">
            <template #body="{ data }">
               <div class="flex items-center gap-1">
                  <Button
                     icon="pi pi-pencil"
                     severity="secondary"
                     text
                     size="small"
                     v-tooltip.top="'编辑'"
                     @click="openEditDialog(data)"
                  />
                  <Button
                     icon="pi pi-trash"
                     severity="danger"
                     text
                     size="small"
                     v-tooltip.top="'删除'"
                     @click="openDeleteDialog(data)"
                  />
               </div>
            </template>
         </Column>
      </DataTable>

      <!-- Permission Form Dialog -->
      <PermissionFormDialog
         v-model:visible="formDialog"
         :isEditing="isEditing"
         :initialData="
            selectedPermission
               ? {
                    resource: selectedPermission.resource,
                    action: selectedPermission.action,
                    code: selectedPermission.code,
                    description: selectedPermission.description,
                 }
               : undefined
         "
         :isLoading="createMutation.isPending.value || updateMutation.isPending.value"
         @save="savePermission"
      />

      <!-- Delete Permission Dialog -->
      <DeletePermissionDialog
         v-model:visible="deleteDialog"
         :permission="selectedPermission"
         :isLoading="isDeleting"
         @confirm="confirmDelete"
      />
   </div>
</template>
