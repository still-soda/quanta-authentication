<script setup lang="ts">
import { ref } from 'vue'
import Card from 'primevue/card'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputText from 'primevue/inputtext'
import Tag from 'primevue/tag'
import Avatar from 'primevue/avatar'
import Menu from 'primevue/menu'
import type { User, UserStatus } from '@/types'
import { UserStatusNames, UserStatusSeverity } from '@/types'

const props = defineProps<{
   users: User[]
   selectedUsers: User[]
   loading?: boolean
   totalRecords?: number
   rows?: number
   sortField?: string
   sortOrder?: 1 | -1 | 0
}>()

const emit = defineEmits<{
   (e: 'update:selectedUsers', value: User[]): void
   (e: 'edit', user: User): void
   (e: 'delete', user: User): void
   (e: 'resetPassword', user: User): void
   (e: 'disable', user: User): void
   (e: 'enable', user: User): void
   (e: 'manageRoles', user: User): void
   (e: 'search', value: string): void
   (e: 'page', event: any): void
   (
      e: 'sort',
      event: { sortField: string | undefined; sortOrder: 1 | -1 | 0 | null | undefined }
   ): void
}>()

const filters = ref({
   global: { value: null, matchMode: 'contains' },
})

const actionMenu = ref()
const currentUser = ref<User | null>(null)

const actionMenuItems = ref([
   {
      label: '编辑',
      icon: 'pi pi-pencil',
      command: () => currentUser.value && emit('edit', currentUser.value),
   },
   {
      label: '管理角色',
      icon: 'pi pi-shield',
      command: () => currentUser.value && emit('manageRoles', currentUser.value),
   },
   {
      label: '重置密码',
      icon: 'pi pi-refresh',
      command: () => currentUser.value && emit('resetPassword', currentUser.value),
   },
   { separator: true },
   {
      label: '禁用',
      icon: 'pi pi-ban',
      visible: () => currentUser.value?.status === 'ACTIVE',
      command: () => currentUser.value && emit('disable', currentUser.value),
   },
   {
      label: '启用',
      icon: 'pi pi-check',
      visible: () => currentUser.value?.status !== 'ACTIVE',
      command: () => currentUser.value && emit('enable', currentUser.value),
   },
   {
      label: '删除',
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => currentUser.value && emit('delete', currentUser.value),
   },
])

const getStatusSeverity = (status: UserStatus) => {
   return UserStatusSeverity[status] || 'secondary'
}

const getStatusLabel = (status: UserStatus) => {
   return UserStatusNames[status] || status
}

const getRoleSeverity = (roleName: string) => {
   // 系统管理员角色显示为红色
   if (roleName.includes('超级管理员') || roleName.includes('super')) {
      return 'danger'
   }
   if (roleName.includes('管理员') || roleName.includes('admin')) {
      return 'warn'
   }
   return 'info'
}

const openActionMenu = (event: Event, user: User) => {
   currentUser.value = user
   actionMenu.value.toggle(event)
}

const onSelectionChange = (selection: User[]) => {
   emit('update:selectedUsers', selection)
}

const formatDateTime = (dateStr?: string) => {
   if (!dateStr) return '-'
   const date = new Date(dateStr)
   return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
   })
}

const formatDate = (dateStr?: string) => {
   if (!dateStr) return '-'
   const date = new Date(dateStr)
   return date.toLocaleDateString('zh-CN')
}

// 生成头像 URL
const getAvatarUrl = (user: User) => {
   if (user.avatar_id) {
      return `/uploads/${user.avatar_id}`
   }
   return `https://api.dicebear.com/7.x/avataaars/svg?seed=${user.id}`
}

const handleSearch = () => {
   emit('search', filters.value.global.value || '')
}
</script>

<template>
   <Card class="rounded-2xl border border-surface-100 dark:border-surface-800">
      <template #content>
         <!-- Toolbar -->
         <div class="flex justify-between items-center mb-4 flex-wrap gap-4">
            <div class="relative flex items-center">
               <i class="pi pi-search absolute left-3.5 text-surface-400 text-sm"></i>
               <InputText
                  v-model="filters['global'].value"
                  placeholder="搜索用户（姓名、邮箱、学号）..."
                  class="pl-10 min-w-80 h-10 rounded-2xl max-md:min-w-full"
                  @keyup.enter="handleSearch"
               />
            </div>
            <div>
               <Button
                  v-if="selectedUsers.length > 0"
                  :label="`已选择 ${selectedUsers.length} 项`"
                  icon="pi pi-trash"
                  severity="danger"
                  outlined
               />
            </div>
         </div>

         <!-- Data Table -->
         <DataTable
            :selection="selectedUsers"
            @update:selection="onSelectionChange"
            v-model:filters="filters"
            :value="users"
            :rows="rows || 10"
            :totalRecords="totalRecords"
            :paginator="true"
            :lazy="true"
            :sortField="sortField"
            :sortOrder="sortOrder"
            removableSort
            @page="event => emit('page', event)"
            @sort="
               event =>
                  emit('sort', {
                     sortField: typeof event.sortField === 'string' ? event.sortField : undefined,
                     sortOrder: event.sortOrder,
                  })
            "
            :rowsPerPageOptions="[10, 20, 50]"
            paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown CurrentPageReport"
            currentPageReportTemplate="显示 {first} 到 {last} 条，共 {totalRecords} 条"
            dataKey="id"
            :globalFilterFields="['name', 'email', 'student_id']"
            class="text-sm"
            :loading="loading"
            stripedRows
         >
            <Column selectionMode="multiple" headerStyle="width: 3rem"></Column>

            <Column field="name" header="用户" sortable style="min-width: 14rem">
               <template #body="{ data }">
                  <div class="flex items-center gap-3.5">
                     <Avatar :image="getAvatarUrl(data)" shape="circle" size="normal" />
                     <div class="flex flex-col gap-0.5">
                        <span class="font-semibold text-surface-900 dark:text-surface-100">
                           {{ data.display_name || data.name }}
                        </span>
                        <span class="text-[0.8125rem] text-surface-500">
                           {{ data.email }}
                        </span>
                     </div>
                  </div>
               </template>
            </Column>

            <Column field="student_id" header="学号" sortable style="min-width: 8rem">
               <template #body="{ data }">
                  <span class="font-mono text-surface-600 dark:text-surface-400">
                     {{ data.student_id }}
                  </span>
               </template>
            </Column>

            <Column field="roles" header="角色" style="min-width: 10rem">
               <template #body="{ data }">
                  <div class="flex flex-wrap gap-1">
                     <Tag
                        v-for="roleName in data.role_names?.slice(0, 2)"
                        :key="roleName"
                        :severity="getRoleSeverity(roleName)"
                        rounded
                     >
                        {{ roleName }}
                     </Tag>
                     <Tag v-if="data.role_names?.length > 2" severity="secondary" rounded>
                        +{{ data.role_names.length - 2 }}
                     </Tag>
                     <span
                        v-if="!data.role_names?.length"
                        class="text-surface-400 text-[0.8125rem]"
                     >
                        未分配
                     </span>
                  </div>
               </template>
            </Column>

            <Column field="status" header="状态" sortable style="min-width: 6rem">
               <template #body="{ data }">
                  <Tag :severity="getStatusSeverity(data.status)">
                     {{ getStatusLabel(data.status) }}
                  </Tag>
               </template>
            </Column>

            <Column field="last_login_at" header="最后登录" sortable style="min-width: 10rem">
               <template #body="{ data }">
                  <span class="text-surface-500 text-[0.8125rem]">
                     {{ formatDateTime(data.last_login_at) }}
                  </span>
               </template>
            </Column>

            <Column field="created_at" header="创建时间" sortable style="min-width: 8rem">
               <template #body="{ data }">
                  <span class="text-surface-500 text-[0.8125rem]">
                     {{ formatDate(data.created_at) }}
                  </span>
               </template>
            </Column>

            <Column header="操作" style="min-width: 6rem">
               <template #body="{ data }">
                  <Button
                     icon="pi pi-ellipsis-v"
                     text
                     rounded
                     severity="secondary"
                     @click="openActionMenu($event, data)"
                  />
               </template>
            </Column>

            <template #empty>
               <div class="flex flex-col items-center justify-center p-12 text-surface-400">
                  <i class="pi pi-users text-5xl mb-4"></i>
                  <p>暂无用户数据</p>
               </div>
            </template>
         </DataTable>
      </template>
   </Card>

   <!-- Action Menu -->
   <Menu ref="actionMenu" :model="actionMenuItems" popup />
</template>
