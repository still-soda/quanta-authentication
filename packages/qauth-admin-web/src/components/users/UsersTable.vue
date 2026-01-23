<script setup lang="ts">
import { ref } from 'vue';
import Card from 'primevue/card';
import Button from 'primevue/button';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import InputText from 'primevue/inputtext';
import Tag from 'primevue/tag';
import Avatar from 'primevue/avatar';
import Menu from 'primevue/menu';
import type { User } from './UserCell.vue';

const props = defineProps<{
   users: User[];
   selectedUsers: User[];
}>();

const emit = defineEmits<{
   (e: 'update:selectedUsers', value: User[]): void;
   (e: 'edit', user: User): void;
   (e: 'delete', user: User): void;
   (e: 'resetPassword', user: User): void;
   (e: 'disable', user: User): void;
}>();

const filters = ref({
   global: { value: null, matchMode: 'contains' },
});

const actionMenu = ref();
const currentUser = ref<User | null>(null);

const actionMenuItems = ref([
   {
      label: '编辑',
      icon: 'pi pi-pencil',
      command: () => currentUser.value && emit('edit', currentUser.value),
   },
   {
      label: '重置密码',
      icon: 'pi pi-refresh',
      command: () =>
         currentUser.value && emit('resetPassword', currentUser.value),
   },
   { separator: true },
   {
      label: '禁用',
      icon: 'pi pi-ban',
      command: () => currentUser.value && emit('disable', currentUser.value),
   },
   {
      label: '删除',
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => currentUser.value && emit('delete', currentUser.value),
   },
]);

const getStatusSeverity = (status: string) => {
   const map: Record<
      string,
      'success' | 'warn' | 'danger' | 'info' | 'secondary'
   > = {
      active: 'success',
      inactive: 'secondary',
      locked: 'danger',
      pending: 'warn',
   };
   return map[status] || 'info';
};

const getStatusLabel = (status: string) => {
   const map: Record<string, string> = {
      active: '正常',
      inactive: '未激活',
      locked: '已锁定',
      pending: '待审核',
   };
   return map[status] || status;
};

const getRoleSeverity = (role: string) => {
   const map: Record<string, 'danger' | 'warn' | 'info'> = {
      管理员: 'danger',
      开发者: 'warn',
      普通用户: 'info',
   };
   return map[role] || 'info';
};

const openActionMenu = (event: Event, user: User) => {
   currentUser.value = user;
   actionMenu.value.toggle(event);
};

const onSelectionChange = (selection: User[]) => {
   emit('update:selectedUsers', selection);
};
</script>

<template>
   <Card class="users-card">
      <template #content>
         <!-- Toolbar -->
         <div class="table-toolbar">
            <div class="toolbar-left">
               <div class="search-box">
                  <i class="pi pi-search"></i>
                  <InputText
                     v-model="filters['global'].value"
                     placeholder="搜索用户..."
                     class="search-input" />
               </div>
            </div>
            <div class="toolbar-right">
               <Button
                  v-if="selectedUsers.length > 0"
                  :label="`已选择 ${selectedUsers.length} 项`"
                  icon="pi pi-trash"
                  severity="danger"
                  outlined />
            </div>
         </div>

         <!-- Data Table -->
         <DataTable
            :selection="selectedUsers"
            @update:selection="onSelectionChange"
            v-model:filters="filters"
            :value="users"
            :rows="10"
            :paginator="true"
            :rowsPerPageOptions="[5, 10, 20, 50]"
            dataKey="id"
            filterDisplay="row"
            :globalFilterFields="['name', 'email', 'role']"
            class="users-table"
            stripedRows>
            <Column selectionMode="multiple" headerStyle="width: 3rem"></Column>

            <Column
               field="name"
               header="用户"
               sortable
               style="min-width: 14rem">
               <template #body="{ data }">
                  <div class="user-cell">
                     <Avatar
                        :image="data.avatar"
                        shape="circle"
                        size="normal" />
                     <div class="user-info">
                        <span class="user-name">{{ data.name }}</span>
                        <span class="user-email">{{ data.email }}</span>
                     </div>
                  </div>
               </template>
            </Column>

            <Column field="role" header="角色" sortable style="min-width: 8rem">
               <template #body="{ data }">
                  <Tag :severity="getRoleSeverity(data.role)" rounded>
                     {{ data.role }}
                  </Tag>
               </template>
            </Column>

            <Column
               field="status"
               header="状态"
               sortable
               style="min-width: 8rem">
               <template #body="{ data }">
                  <Tag :severity="getStatusSeverity(data.status)">
                     {{ getStatusLabel(data.status) }}
                  </Tag>
               </template>
            </Column>

            <Column
               field="lastLogin"
               header="最后登录"
               sortable
               style="min-width: 10rem">
               <template #body="{ data }">
                  <span class="text-muted">{{ data.lastLogin }}</span>
               </template>
            </Column>

            <Column
               field="createdAt"
               header="创建时间"
               sortable
               style="min-width: 8rem">
               <template #body="{ data }">
                  <span class="text-muted">{{ data.createdAt }}</span>
               </template>
            </Column>

            <Column header="操作" style="min-width: 6rem">
               <template #body="{ data }">
                  <Button
                     icon="pi pi-ellipsis-v"
                     text
                     rounded
                     severity="secondary"
                     @click="openActionMenu($event, data)" />
               </template>
            </Column>

            <template #empty>
               <div class="empty-state">
                  <i class="pi pi-users"></i>
                  <p>暂无用户数据</p>
               </div>
            </template>
         </DataTable>
      </template>
   </Card>

   <!-- Action Menu -->
   <Menu ref="actionMenu" :model="actionMenuItems" popup />
</template>

<style scoped>
.users-card {
   border-radius: 16px;
   border: 1px solid var(--p-surface-100);
}

:global(.app-dark) .users-card {
   border-color: var(--p-surface-800);
}

.table-toolbar {
   display: flex;
   justify-content: space-between;
   align-items: center;
   margin-bottom: 1rem;
   flex-wrap: wrap;
   gap: 1rem;
}

.search-box {
   position: relative;
   display: flex;
   align-items: center;
}

.search-box i {
   position: absolute;
   left: 0.875rem;
   color: var(--p-surface-400);
   font-size: 0.875rem;
}

.search-input {
   padding-left: 2.5rem;
   min-width: 18rem;
   height: 2.5rem;
   border-radius: 10px;
}

.users-table {
   font-size: 0.875rem;
}

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

.text-muted {
   color: var(--p-surface-500);
   font-size: 0.8125rem;
}

.empty-state {
   display: flex;
   flex-direction: column;
   align-items: center;
   justify-content: center;
   padding: 3rem;
   color: var(--p-surface-400);
}

.empty-state i {
   font-size: 3rem;
   margin-bottom: 1rem;
}

@media (max-width: 768px) {
   .search-input {
      min-width: 100%;
   }
}
</style>
