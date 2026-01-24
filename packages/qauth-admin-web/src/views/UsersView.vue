<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import Button from 'primevue/button'
import Toast from 'primevue/toast'
import ConfirmDialog from 'primevue/confirmdialog'
import PageHeader from '@/components/shared/PageHeader.vue'
import SimpleStatCard from '@/components/shared/SimpleStatCard.vue'
import UsersTable from '@/components/users/UsersTable.vue'
import UserDialog from '@/components/users/UserDialog.vue'
import UserRolesDialog from '@/components/users/UserRolesDialog.vue'
import ResetPasswordDialog from '@/components/users/ResetPasswordDialog.vue'
import type {
   SimpleStatData,
   User,
   CreateUserFormData,
   UpdateUserFormData,
   ListUsersParams,
   UserStatus,
} from '@/types'
import {
   getUsers,
   createUser,
   updateUser,
   deleteUser,
   resetUserPassword,
   getUserStatusCounts,
   setUserRoles,
} from '@/apis/users'
import { getRoles } from '@/apis/roles'

const queryClient = useQueryClient()
const toast = useToast()
const confirm = useConfirm()

// 分页和筛选参数
const page = ref(1)
const pageSize = ref(10)
const searchKeyword = ref('')
const statusFilter = ref<UserStatus | '' | undefined>(undefined)

// 查询参数
const queryParams = computed<ListUsersParams>(() => ({
   page: page.value,
   page_size: pageSize.value,
   search: searchKeyword.value || undefined,
   status: statusFilter.value,
}))

// 获取用户数据
const {
   data: usersResult,
   isLoading: isLoadingUsers,
} = useQuery({
   queryKey: ['users', queryParams],
   queryFn: () => getUsers(queryParams.value),
})

// 获取用户状态统计
const { data: statusCounts, isLoading: isLoadingCounts } = useQuery({
   queryKey: ['userStatusCounts'],
   queryFn: getUserStatusCounts,
})

// 获取角色列表
const { data: roles } = useQuery({
   queryKey: ['roles'],
   queryFn: getRoles,
})

// 创建用户 mutation
const createUserMutation = useMutation({
   mutationFn: (data: CreateUserFormData) => createUser(data),
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
      queryClient.invalidateQueries({ queryKey: ['userStatusCounts'] })
      userDialog.value = false
      toast.add({
         severity: 'success',
         summary: '成功',
         detail: '用户创建成功',
         life: 3000,
      })
   },
   onError: (error: any) => {
      toast.add({
         severity: 'error',
         summary: '错误',
         detail: error.message || '创建用户失败',
         life: 5000,
      })
   },
})

// 更新用户 mutation
const updateUserMutation = useMutation({
   mutationFn: ({ id, data }: { id: string; data: UpdateUserFormData }) => updateUser(id, data),
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
      queryClient.invalidateQueries({ queryKey: ['userStatusCounts'] })
      userDialog.value = false
      toast.add({
         severity: 'success',
         summary: '成功',
         detail: '用户更新成功',
         life: 3000,
      })
   },
   onError: (error: any) => {
      toast.add({
         severity: 'error',
         summary: '错误',
         detail: error.message || '更新用户失败',
         life: 5000,
      })
   },
})

// 删除用户 mutation
const deleteUserMutation = useMutation({
   mutationFn: deleteUser,
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
      queryClient.invalidateQueries({ queryKey: ['userStatusCounts'] })
      toast.add({
         severity: 'success',
         summary: '成功',
         detail: '用户已删除',
         life: 3000,
      })
   },
   onError: (error: any) => {
      toast.add({
         severity: 'error',
         summary: '错误',
         detail: error.message || '删除用户失败',
         life: 5000,
      })
   },
})

// 重置密码 mutation
const resetPasswordMutation = useMutation({
   mutationFn: ({ userId, password }: { userId: string; password?: string }) =>
      resetUserPassword(userId, password),
   onSuccess: data => {
      resetPasswordDialog.value = false
      toast.add({
         severity: 'success',
         summary: '成功',
         detail: `密码已重置为: ${data.new_password}`,
         life: 10000,
      })
   },
   onError: (error: any) => {
      toast.add({
         severity: 'error',
         summary: '错误',
         detail: error.message || '重置密码失败',
         life: 5000,
      })
   },
})

// 设置用户角色 mutation
const setUserRolesMutation = useMutation({
   mutationFn: ({ userId, roleIds }: { userId: string; roleIds: string[] }) =>
      setUserRoles(userId, roleIds),
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
      userRolesDialog.value = false
      toast.add({
         severity: 'success',
         summary: '成功',
         detail: '用户角色已更新',
         life: 3000,
      })
   },
   onError: (error: any) => {
      toast.add({
         severity: 'error',
         summary: '错误',
         detail: error.message || '更新用户角色失败',
         life: 5000,
      })
   },
})

// 启用用户 mutation
const enableUserMutation = useMutation({
   mutationFn: (userId: string) => updateUser(userId, { status: 'ACTIVE' }),
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
      queryClient.invalidateQueries({ queryKey: ['userStatusCounts'] })
      toast.add({
         severity: 'success',
         summary: '成功',
         detail: '用户已启用',
         life: 3000,
      })
   },
   onError: (error: any) => {
      toast.add({
         severity: 'error',
         summary: '错误',
         detail: error.message || '启用用户失败',
         life: 5000,
      })
   },
})

const selectedUsers = ref<User[]>([])
const userDialog = ref(false)
const userRolesDialog = ref(false)
const resetPasswordDialog = ref(false)
const isEditing = ref(false)
const currentUser = ref<User | null>(null)

// 统计数据
const stats = computed<SimpleStatData[]>(() => {
   const counts = statusCounts.value || { ACTIVE: 0, LOCKED: 0, BANNED: 0 }
   const total = (counts.ACTIVE || 0) + (counts.LOCKED || 0) + (counts.BANNED || 0)
   return [
      {
         title: '总用户',
         value: total,
         icon: 'pi pi-users',
         color: 'blue',
      },
      {
         title: '活跃用户',
         value: counts.ACTIVE || 0,
         icon: 'pi pi-check-circle',
         color: 'green',
      },
      {
         title: '已禁用',
         value: counts.BANNED || 0,
         icon: 'pi pi-ban',
         color: 'gray',
      },
      {
         title: '已锁定',
         value: counts.LOCKED || 0,
         icon: 'pi pi-lock',
         color: 'red',
      },
   ]
})

// 用户列表
const users = computed(() => usersResult.value?.users || [])
const totalRecords = computed(() => usersResult.value?.total || 0)

const openNewUserDialog = () => {
   isEditing.value = false
   currentUser.value = null
   userDialog.value = true
}

const editUser = (user: User) => {
   isEditing.value = true
   currentUser.value = user
   userDialog.value = true
}

const handleDelete = (user: User) => {
   confirm.require({
      message: `确定要删除用户 "${user.name}" 吗？此操作不可恢复。`,
      header: '确认删除',
      icon: 'pi pi-exclamation-triangle',
      rejectClass: 'p-button-secondary p-button-outlined',
      rejectLabel: '取消',
      acceptLabel: '删除',
      acceptClass: 'p-button-danger',
      accept: () => {
         deleteUserMutation.mutate(user.id)
      },
   })
}

const handleResetPassword = (user: User) => {
   currentUser.value = user
   resetPasswordDialog.value = true
}

const handleResetPasswordConfirm = (password?: string) => {
   if (currentUser.value) {
      resetPasswordMutation.mutate({ userId: currentUser.value.id, password })
   }
}

const handleDisable = (user: User) => {
   confirm.require({
      message: `确定要禁用用户 "${user.name}" 吗？`,
      header: '确认禁用',
      icon: 'pi pi-exclamation-triangle',
      rejectClass: 'p-button-secondary p-button-outlined',
      rejectLabel: '取消',
      acceptLabel: '禁用',
      acceptClass: 'p-button-warning',
      accept: () => {
         updateUserMutation.mutate({ id: user.id, data: { status: 'BANNED' } })
      },
   })
}

const handleEnable = (user: User) => {
   confirm.require({
      message: `确定要启用用户 "${user.name}" 吗？`,
      header: '确认启用',
      icon: 'pi pi-check-circle',
      rejectClass: 'p-button-secondary p-button-outlined',
      rejectLabel: '取消',
      acceptLabel: '启用',
      acceptClass: 'p-button-success',
      accept: () => {
         enableUserMutation.mutate(user.id)
      },
   })
}

const handleManageRoles = (user: User) => {
   currentUser.value = user
   userRolesDialog.value = true
}

const handleSaveRoles = (userId: string, roleIds: string[]) => {
   setUserRolesMutation.mutate({ userId, roleIds })
}

const handleSearch = (keyword: string) => {
   searchKeyword.value = keyword
   page.value = 1
}

const handlePageChange = (event: { page: number; rows: number }) => {
   page.value = event.page + 1
   pageSize.value = event.rows
}

// 监听筛选变化，重置页码
watch([statusFilter], () => {
   page.value = 1
})
</script>

<template>
   <div class="flex flex-col gap-6">
      <Toast />
      <ConfirmDialog />

      <!-- Page Header -->
      <PageHeader title="用户管理" subtitle="管理系统用户账号和权限">
         <template #actions>
            <Button label="新建用户" icon="pi pi-plus" @click="openNewUserDialog" />
         </template>
      </PageHeader>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4">
         <template v-if="isLoadingCounts">
            <div
               v-for="i in 4"
               :key="i"
               class="h-20 bg-surface-100 dark:bg-surface-800 rounded-xl animate-pulse"
            />
         </template>
         <template v-else>
            <SimpleStatCard v-for="stat in stats" :key="stat.title" :stat="stat" />
         </template>
      </div>

      <!-- Users Table -->
      <div
         v-if="isLoadingUsers"
         class="h-96 bg-surface-100 dark:bg-surface-800 rounded-xl animate-pulse"
      />
      <UsersTable
         v-else
         :users="users"
         :loading="isLoadingUsers"
         :total-records="totalRecords"
         :rows="pageSize"
         v-model:selectedUsers="selectedUsers"
         @edit="editUser"
         @delete="handleDelete"
         @resetPassword="handleResetPassword"
         @disable="handleDisable"
         @enable="handleEnable"
         @manageRoles="handleManageRoles"
         @search="handleSearch"
         @page="handlePageChange"
      />

      <!-- User Dialog -->
      <UserDialog
         v-model:visible="userDialog"
         :isEditing="isEditing"
         :user="currentUser"
         :roles="roles || []"
         :loading="createUserMutation.isPending.value || updateUserMutation.isPending.value"
         @create="(data) => createUserMutation.mutate(data)"
         @update="(id, data) => updateUserMutation.mutate({ id, data })"
      />

      <!-- User Roles Dialog -->
      <UserRolesDialog
         v-model:visible="userRolesDialog"
         :user="currentUser"
         :roles="roles || []"
         :loading="setUserRolesMutation.isPending.value"
         @save="handleSaveRoles"
      />

      <!-- Reset Password Dialog -->
      <ResetPasswordDialog
         v-model:visible="resetPasswordDialog"
         :user="currentUser"
         :loading="resetPasswordMutation.isPending.value"
         @confirm="handleResetPasswordConfirm"
      />
   </div>
</template>
