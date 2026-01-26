<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useQuery, useMutation } from '@tanstack/vue-query'
import { useToast } from 'primevue/usetoast'
import Tabs from 'primevue/tabs'
import TabList from 'primevue/tablist'
import Tab from 'primevue/tab'
import TabPanels from 'primevue/tabpanels'
import TabPanel from 'primevue/tabpanel'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Select from 'primevue/select'
import Checkbox from 'primevue/checkbox'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import ConfirmDialog from 'primevue/confirmdialog'
import { useConfirm } from 'primevue/useconfirm'
import type {
   OAuthApp,
   AppGroupAdmin,
   AppGroupPermission,
   AppGroupRole,
   AppGroupAdminType,
   AppGroupPermissionFormData,
   AppGroupRoleFormData,
   AppGroupUserInfo,
   User,
} from '@/types'
import {
   AppGroupAdminTypeLabels,
   AppGroupAdminTypeDescriptions,
   PermissionActionLabels,
   PermissionActionOptions,
} from '@/types'
import {
   getAppGroupAdmins,
   addAppGroupAdmin,
   removeAppGroupAdmin,
   getAppGroupPermissions,
   createAppGroupPermission,
   updateAppGroupPermission,
   deleteAppGroupPermission,
   getAppGroupRoles,
   createAppGroupRole,
   updateAppGroupRole,
   deleteAppGroupRole,
   getAppGroupRolePermissions,
   setAppGroupRolePermissions,
   getAppGroupRoleUsers,
   assignAppGroupRoleToUser,
   revokeAppGroupRoleFromUser,
} from '@/apis/app-group'
import { getUsers } from '@/apis/users'

const props = defineProps<{
   visible: boolean
   client: OAuthApp | null
}>()

const emit = defineEmits<{
   (e: 'update:visible', value: boolean): void
}>()

const toast = useToast()
const confirm = useConfirm()

const dialogVisible = computed({
   get: () => props.visible,
   set: value => emit('update:visible', value),
})

const clientId = computed(() => props.client?.id || '')

// ======================== 管理员相关 ========================
const { data: admins, refetch: refetchAdmins } = useQuery({
   queryKey: ['app-group-admins', clientId],
   queryFn: () => getAppGroupAdmins(clientId.value),
   enabled: computed(() => !!clientId.value && props.visible),
})

const addAdminDialog = ref(false)
const newAdminUserId = ref('')
const newAdminType = ref<AppGroupAdminType>('admin')

const adminTypeOptions = [
   { label: AppGroupAdminTypeLabels.admin, value: 'admin' as AppGroupAdminType },
   { label: AppGroupAdminTypeLabels.role_manager, value: 'role_manager' as AppGroupAdminType },
   {
      label: AppGroupAdminTypeLabels.permission_manager,
      value: 'permission_manager' as AppGroupAdminType,
   },
]

const addAdminMutation = useMutation({
   mutationFn: () =>
      addAppGroupAdmin(clientId.value, {
         user_id: newAdminUserId.value,
         admin_type: newAdminType.value,
      }),
   onSuccess: () => {
      refetchAdmins()
      addAdminDialog.value = false
      newAdminUserId.value = ''
      toast.add({ severity: 'success', summary: '添加成功', detail: '管理员已添加', life: 3000 })
   },
   onError: (error: Error) => {
      toast.add({ severity: 'error', summary: '添加失败', detail: error.message, life: 5000 })
   },
})

const removeAdminMutation = useMutation({
   mutationFn: ({ userId, adminType }: { userId: string; adminType: AppGroupAdminType }) =>
      removeAppGroupAdmin(clientId.value, userId, adminType),
   onSuccess: () => {
      refetchAdmins()
      toast.add({ severity: 'success', summary: '移除成功', detail: '管理员已移除', life: 3000 })
   },
   onError: (error: Error) => {
      toast.add({ severity: 'error', summary: '移除失败', detail: error.message, life: 5000 })
   },
})

const confirmRemoveAdmin = (admin: AppGroupAdmin) => {
   if (admin.admin_type === 'owner') {
      toast.add({
         severity: 'warn',
         summary: '无法移除',
         detail: '应用所有者不能被移除',
         life: 3000,
      })
      return
   }
   confirm.require({
      message: `确定要移除 "${admin.user_name}" 的 ${AppGroupAdminTypeLabels[admin.admin_type]} 权限吗？`,
      header: '移除管理员',
      icon: 'pi pi-exclamation-triangle',
      rejectLabel: '取消',
      acceptLabel: '确认移除',
      acceptClass: 'p-button-danger',
      accept: () =>
         removeAdminMutation.mutate({ userId: admin.user_id, adminType: admin.admin_type }),
   })
}

// ======================== 权限相关 ========================
const { data: permissions, refetch: refetchPermissions } = useQuery({
   queryKey: ['app-group-permissions', clientId],
   queryFn: () => getAppGroupPermissions(clientId.value),
   enabled: computed(() => !!clientId.value && props.visible),
})

const permissionDialog = ref(false)
const isEditingPermission = ref(false)
const selectedPermission = ref<AppGroupPermission | null>(null)
const permissionForm = ref<AppGroupPermissionFormData>({
   resource: '',
   action: 1,
   code: '',
   name: '',
   description: '',
})

const openNewPermissionDialog = () => {
   isEditingPermission.value = false
   selectedPermission.value = null
   permissionForm.value = { resource: '', action: 1, code: '', name: '', description: '' }
   permissionDialog.value = true
}

const editPermission = (perm: AppGroupPermission) => {
   isEditingPermission.value = true
   selectedPermission.value = perm
   permissionForm.value = {
      resource: perm.resource,
      action: perm.action,
      code: perm.code,
      name: perm.name,
      description: perm.description,
   }
   permissionDialog.value = true
}

const createPermissionMutation = useMutation({
   mutationFn: (data: AppGroupPermissionFormData) => createAppGroupPermission(clientId.value, data),
   onSuccess: () => {
      refetchPermissions()
      permissionDialog.value = false
      toast.add({ severity: 'success', summary: '创建成功', detail: '权限已创建', life: 3000 })
   },
   onError: (error: Error) => {
      toast.add({ severity: 'error', summary: '创建失败', detail: error.message, life: 5000 })
   },
})

const updatePermissionMutation = useMutation({
   mutationFn: ({
      id,
      data,
   }: {
      id: string
      data: Pick<AppGroupPermissionFormData, 'name' | 'description'>
   }) => updateAppGroupPermission(clientId.value, id, data),
   onSuccess: () => {
      refetchPermissions()
      permissionDialog.value = false
      toast.add({ severity: 'success', summary: '更新成功', detail: '权限已更新', life: 3000 })
   },
   onError: (error: Error) => {
      toast.add({ severity: 'error', summary: '更新失败', detail: error.message, life: 5000 })
   },
})

const deletePermissionMutation = useMutation({
   mutationFn: (id: string) => deleteAppGroupPermission(clientId.value, id),
   onSuccess: () => {
      refetchPermissions()
      toast.add({ severity: 'success', summary: '删除成功', detail: '权限已删除', life: 3000 })
   },
   onError: (error: Error) => {
      toast.add({ severity: 'error', summary: '删除失败', detail: error.message, life: 5000 })
   },
})

const savePermission = () => {
   if (isEditingPermission.value && selectedPermission.value) {
      updatePermissionMutation.mutate({
         id: selectedPermission.value.id,
         data: { name: permissionForm.value.name, description: permissionForm.value.description },
      })
   } else {
      createPermissionMutation.mutate(permissionForm.value)
   }
}

const confirmDeletePermission = (perm: AppGroupPermission) => {
   confirm.require({
      message: `确定要删除权限 "${perm.name}" 吗？此操作不可撤销。`,
      header: '删除权限',
      icon: 'pi pi-exclamation-triangle',
      rejectLabel: '取消',
      acceptLabel: '确认删除',
      acceptClass: 'p-button-danger',
      accept: () => deletePermissionMutation.mutate(perm.id),
   })
}

// ======================== 角色相关 ========================
const { data: roles, refetch: refetchRoles } = useQuery({
   queryKey: ['app-group-roles', clientId],
   queryFn: () => getAppGroupRoles(clientId.value),
   enabled: computed(() => !!clientId.value && props.visible),
})

const roleDialog = ref(false)
const isEditingRole = ref(false)
const selectedRole = ref<AppGroupRole | null>(null)
const roleForm = ref<AppGroupRoleFormData>({
   code: '',
   name: '',
   description: '',
   is_default: false,
})

const openNewRoleDialog = () => {
   isEditingRole.value = false
   selectedRole.value = null
   roleForm.value = { code: '', name: '', description: '', is_default: false }
   roleDialog.value = true
}

const editRole = (role: AppGroupRole) => {
   isEditingRole.value = true
   selectedRole.value = role
   roleForm.value = {
      code: role.code,
      name: role.name,
      description: role.description,
      is_default: role.is_default,
   }
   roleDialog.value = true
}

const createRoleMutation = useMutation({
   mutationFn: (data: AppGroupRoleFormData) => createAppGroupRole(clientId.value, data),
   onSuccess: () => {
      refetchRoles()
      roleDialog.value = false
      toast.add({ severity: 'success', summary: '创建成功', detail: '角色已创建', life: 3000 })
   },
   onError: (error: Error) => {
      toast.add({ severity: 'error', summary: '创建失败', detail: error.message, life: 5000 })
   },
})

const updateRoleMutation = useMutation({
   mutationFn: ({
      id,
      data,
   }: {
      id: string
      data: Pick<AppGroupRoleFormData, 'name' | 'description' | 'is_default'>
   }) => updateAppGroupRole(clientId.value, id, data),
   onSuccess: () => {
      refetchRoles()
      roleDialog.value = false
      toast.add({ severity: 'success', summary: '更新成功', detail: '角色已更新', life: 3000 })
   },
   onError: (error: Error) => {
      toast.add({ severity: 'error', summary: '更新失败', detail: error.message, life: 5000 })
   },
})

const deleteRoleMutation = useMutation({
   mutationFn: (id: string) => deleteAppGroupRole(clientId.value, id),
   onSuccess: () => {
      refetchRoles()
      toast.add({ severity: 'success', summary: '删除成功', detail: '角色已删除', life: 3000 })
   },
   onError: (error: Error) => {
      toast.add({ severity: 'error', summary: '删除失败', detail: error.message, life: 5000 })
   },
})

const saveRole = () => {
   if (isEditingRole.value && selectedRole.value) {
      updateRoleMutation.mutate({
         id: selectedRole.value.id,
         data: {
            name: roleForm.value.name,
            description: roleForm.value.description,
            is_default: roleForm.value.is_default,
         },
      })
   } else {
      createRoleMutation.mutate(roleForm.value)
   }
}

const confirmDeleteRole = (role: AppGroupRole) => {
   confirm.require({
      message: `确定要删除角色 "${role.name}" 吗？此操作不可撤销。`,
      header: '删除角色',
      icon: 'pi pi-exclamation-triangle',
      rejectLabel: '取消',
      acceptLabel: '确认删除',
      acceptClass: 'p-button-danger',
      accept: () => deleteRoleMutation.mutate(role.id),
   })
}

// ======================== 角色权限分配 ========================
const rolePermissionsDialog = ref(false)
const selectedRoleForPermissions = ref<AppGroupRole | null>(null)
const rolePermissionIds = ref<string[]>([])

const { data: rolePermissions, refetch: refetchRolePermissions } = useQuery({
   queryKey: ['app-group-role-permissions', selectedRoleForPermissions],
   queryFn: () => getAppGroupRolePermissions(clientId.value, selectedRoleForPermissions.value!.id),
   enabled: computed(() => !!selectedRoleForPermissions.value),
})

watch(rolePermissions, perms => {
   rolePermissionIds.value = perms?.map(p => p.id) || []
})

const openRolePermissionsDialog = (role: AppGroupRole) => {
   selectedRoleForPermissions.value = role
   rolePermissionsDialog.value = true
   refetchRolePermissions()
}

const setRolePermissionsMutation = useMutation({
   mutationFn: () =>
      setAppGroupRolePermissions(clientId.value, selectedRoleForPermissions.value!.id, {
         permission_ids: rolePermissionIds.value,
      }),
   onSuccess: () => {
      refetchRoles()
      rolePermissionsDialog.value = false
      toast.add({ severity: 'success', summary: '保存成功', detail: '角色权限已更新', life: 3000 })
   },
   onError: (error: Error) => {
      toast.add({ severity: 'error', summary: '保存失败', detail: error.message, life: 5000 })
   },
})

// ======================== 角色用户分配 ========================
const roleUsersDialog = ref(false)
const selectedRoleForUsers = ref<AppGroupRole | null>(null)
const newRoleUserId = ref('')

const { data: roleUsers, refetch: refetchRoleUsers } = useQuery({
   queryKey: ['app-group-role-users', selectedRoleForUsers],
   queryFn: () => getAppGroupRoleUsers(clientId.value, selectedRoleForUsers.value!.id),
   enabled: computed(() => !!selectedRoleForUsers.value),
})

const openRoleUsersDialog = (role: AppGroupRole) => {
   selectedRoleForUsers.value = role
   roleUsersDialog.value = true
   refetchRoleUsers()
}

const assignUserMutation = useMutation({
   mutationFn: () =>
      assignAppGroupRoleToUser(clientId.value, selectedRoleForUsers.value!.id, {
         user_id: newRoleUserId.value,
      }),
   onSuccess: () => {
      refetchRoleUsers()
      refetchRoles()
      newRoleUserId.value = ''
      toast.add({ severity: 'success', summary: '分配成功', detail: '用户已分配角色', life: 3000 })
   },
   onError: (error: Error) => {
      toast.add({ severity: 'error', summary: '分配失败', detail: error.message, life: 5000 })
   },
})

const revokeUserMutation = useMutation({
   mutationFn: (userId: string) =>
      revokeAppGroupRoleFromUser(clientId.value, selectedRoleForUsers.value!.id, userId),
   onSuccess: () => {
      refetchRoleUsers()
      refetchRoles()
      toast.add({ severity: 'success', summary: '移除成功', detail: '用户角色已移除', life: 3000 })
   },
   onError: (error: Error) => {
      toast.add({ severity: 'error', summary: '移除失败', detail: error.message, life: 5000 })
   },
})

const confirmRevokeUser = (user: AppGroupUserInfo) => {
   confirm.require({
      message: `确定要从 "${user.name}" 移除此角色吗？`,
      header: '移除用户角色',
      icon: 'pi pi-exclamation-triangle',
      rejectLabel: '取消',
      acceptLabel: '确认移除',
      acceptClass: 'p-button-danger',
      accept: () => revokeUserMutation.mutate(user.id),
   })
}

// 获取用户列表用于选择
const { data: allUsers } = useQuery({
   queryKey: ['users-for-select'],
   queryFn: () => getUsers({ page: 1, page_size: 1000 }),
   enabled: computed(() => props.visible),
})

const userOptions = computed(() => {
   return (
      allUsers.value?.users?.map((u: User) => ({
         label: `${u.name} (${u.student_id})`,
         value: u.id,
      })) || []
   )
})

const getAdminTypeSeverity = (type: AppGroupAdminType) => {
   switch (type) {
      case 'owner':
         return 'danger'
      case 'admin':
         return 'warn'
      case 'role_manager':
         return 'info'
      case 'permission_manager':
         return 'secondary'
      default:
         return 'secondary'
   }
}
</script>

<template>
   <Dialog
      v-model:visible="dialogVisible"
      :header="`${client?.name || ''} - 应用组权限管理`"
      modal
      :style="{ width: '64rem' }"
      :breakpoints="{ '960px': '95vw' }"
      class="app-group-dialog"
   >
      <Tabs value="admins">
         <TabList>
            <Tab value="admins">
               <i class="pi pi-user-edit mr-2"></i>
               管理员
            </Tab>
            <Tab value="permissions">
               <i class="pi pi-shield mr-2"></i>
               权限
            </Tab>
            <Tab value="roles">
               <i class="pi pi-users mr-2"></i>
               角色
            </Tab>
         </TabList>

         <TabPanels>
            <!-- 管理员标签页 -->
            <TabPanel value="admins">
               <div class="flex flex-col gap-4 pt-4">
                  <div class="flex justify-between items-center">
                     <p class="text-sm text-surface-500 m-0">
                        管理应用组的管理员，他们可以管理权限和角色
                     </p>
                     <Button
                        label="添加管理员"
                        icon="pi pi-plus"
                        size="small"
                        @click="addAdminDialog = true"
                     />
                  </div>

                  <DataTable
                     :value="admins || []"
                     stripedRows
                     size="small"
                     class="rounded-lg overflow-hidden"
                  >
                     <Column field="user_name" header="用户"></Column>
                     <Column field="user_email" header="邮箱"></Column>
                     <Column field="admin_type" header="类型">
                        <template #body="{ data }">
                           <Tag
                              :value="AppGroupAdminTypeLabels[data.admin_type as AppGroupAdminType]"
                              :severity="getAdminTypeSeverity(data.admin_type)"
                           />
                        </template>
                     </Column>
                     <Column field="granted_at" header="授权时间">
                        <template #body="{ data }">
                           {{ new Date(data.granted_at).toLocaleString('zh-CN') }}
                        </template>
                     </Column>
                     <Column header="操作" :style="{ width: '100px' }">
                        <template #body="{ data }">
                           <Button
                              v-if="data.admin_type !== 'owner'"
                              icon="pi pi-trash"
                              severity="danger"
                              text
                              size="small"
                              @click="confirmRemoveAdmin(data)"
                           />
                           <span v-else class="text-xs text-surface-400">不可移除</span>
                        </template>
                     </Column>
                  </DataTable>
               </div>
            </TabPanel>

            <!-- 权限标签页 -->
            <TabPanel value="permissions">
               <div class="flex flex-col gap-4 pt-4">
                  <div class="flex justify-between items-center">
                     <p class="text-sm text-surface-500 m-0">定义应用组的自定义权限</p>
                     <Button
                        label="新建权限"
                        icon="pi pi-plus"
                        size="small"
                        @click="openNewPermissionDialog"
                     />
                  </div>

                  <DataTable
                     :value="permissions || []"
                     stripedRows
                     size="small"
                     class="rounded-lg overflow-hidden"
                  >
                     <Column field="name" header="名称"></Column>
                     <Column field="code" header="代码">
                        <template #body="{ data }">
                           <code
                              class="text-xs bg-surface-100 dark:bg-surface-700 px-2 py-1 rounded"
                              >{{ data.code }}</code
                           >
                        </template>
                     </Column>
                     <Column field="resource" header="资源"></Column>
                     <Column field="action" header="操作">
                        <template #body="{ data }">
                           {{ PermissionActionLabels[data.action] }}
                        </template>
                     </Column>
                     <Column field="description" header="描述">
                        <template #body="{ data }">
                           <span class="text-sm text-surface-500">{{
                              data.description || '-'
                           }}</span>
                        </template>
                     </Column>
                     <Column header="操作" :style="{ width: '120px' }">
                        <template #body="{ data }">
                           <div class="flex gap-1">
                              <Button
                                 icon="pi pi-pencil"
                                 text
                                 size="small"
                                 @click="editPermission(data)"
                              />
                              <Button
                                 icon="pi pi-trash"
                                 severity="danger"
                                 text
                                 size="small"
                                 @click="confirmDeletePermission(data)"
                              />
                           </div>
                        </template>
                     </Column>
                  </DataTable>

                  <div v-if="!permissions?.length" class="text-center py-8 text-surface-400">
                     <i class="pi pi-shield text-3xl mb-2"></i>
                     <p>暂无自定义权限</p>
                  </div>
               </div>
            </TabPanel>

            <!-- 角色标签页 -->
            <TabPanel value="roles">
               <div class="flex flex-col gap-4 pt-4">
                  <div class="flex justify-between items-center">
                     <p class="text-sm text-surface-500 m-0">定义应用组的角色并分配权限</p>
                     <Button
                        label="新建角色"
                        icon="pi pi-plus"
                        size="small"
                        @click="openNewRoleDialog"
                     />
                  </div>

                  <DataTable
                     :value="roles || []"
                     stripedRows
                     size="small"
                     class="rounded-lg overflow-hidden"
                  >
                     <Column field="name" header="名称">
                        <template #body="{ data }">
                           <div class="flex items-center gap-2">
                              {{ data.name }}
                              <Tag v-if="data.is_default" value="默认" severity="info" />
                           </div>
                        </template>
                     </Column>
                     <Column field="code" header="代码">
                        <template #body="{ data }">
                           <code
                              class="text-xs bg-surface-100 dark:bg-surface-700 px-2 py-1 rounded"
                              >{{ data.code }}</code
                           >
                        </template>
                     </Column>
                     <Column field="permission_count" header="权限数">
                        <template #body="{ data }">
                           <Button
                              :label="String(data.permission_count)"
                              link
                              size="small"
                              @click="openRolePermissionsDialog(data)"
                           />
                        </template>
                     </Column>
                     <Column field="user_count" header="用户数">
                        <template #body="{ data }">
                           <Button
                              :label="String(data.user_count)"
                              link
                              size="small"
                              @click="openRoleUsersDialog(data)"
                           />
                        </template>
                     </Column>
                     <Column header="操作" :style="{ width: '180px' }">
                        <template #body="{ data }">
                           <div class="flex gap-1">
                              <Button
                                 icon="pi pi-shield"
                                 v-tooltip="'管理权限'"
                                 text
                                 size="small"
                                 @click="openRolePermissionsDialog(data)"
                              />
                              <Button
                                 icon="pi pi-users"
                                 v-tooltip="'管理用户'"
                                 text
                                 size="small"
                                 @click="openRoleUsersDialog(data)"
                              />
                              <Button
                                 icon="pi pi-pencil"
                                 text
                                 size="small"
                                 @click="editRole(data)"
                              />
                              <Button
                                 icon="pi pi-trash"
                                 severity="danger"
                                 text
                                 size="small"
                                 @click="confirmDeleteRole(data)"
                              />
                           </div>
                        </template>
                     </Column>
                  </DataTable>

                  <div v-if="!roles?.length" class="text-center py-8 text-surface-400">
                     <i class="pi pi-users text-3xl mb-2"></i>
                     <p>暂无自定义角色</p>
                  </div>
               </div>
            </TabPanel>
         </TabPanels>
      </Tabs>

      <template #footer>
         <Button label="关闭" severity="secondary" @click="dialogVisible = false" />
      </template>
   </Dialog>

   <!-- 添加管理员对话框 -->
   <Dialog v-model:visible="addAdminDialog" header="添加管理员" modal :style="{ width: '28rem' }">
      <div class="flex flex-col gap-4">
         <div class="flex flex-col gap-2">
            <label class="text-sm font-medium">选择用户</label>
            <Select
               v-model="newAdminUserId"
               :options="userOptions"
               optionLabel="label"
               optionValue="value"
               placeholder="请选择用户"
               filter
               class="w-full"
            />
         </div>
         <div class="flex flex-col gap-2">
            <label class="text-sm font-medium">管理员类型</label>
            <Select
               v-model="newAdminType"
               :options="adminTypeOptions"
               optionLabel="label"
               optionValue="value"
               class="w-full"
            />
            <p class="text-xs text-surface-400 m-0">
               {{ AppGroupAdminTypeDescriptions[newAdminType] }}
            </p>
         </div>
      </div>
      <template #footer>
         <Button label="取消" severity="secondary" @click="addAdminDialog = false" />
         <Button
            label="添加"
            :loading="addAdminMutation.isPending.value"
            @click="addAdminMutation.mutate()"
         />
      </template>
   </Dialog>

   <!-- 权限编辑对话框 -->
   <Dialog
      v-model:visible="permissionDialog"
      :header="isEditingPermission ? '编辑权限' : '新建权限'"
      modal
      :style="{ width: '28rem' }"
   >
      <div class="flex flex-col gap-4">
         <div class="flex flex-col gap-2">
            <label class="text-sm font-medium">权限名称 <span class="text-red-500">*</span></label>
            <InputText v-model="permissionForm.name" placeholder="例如：查看报表" class="w-full" />
         </div>
         <div v-if="!isEditingPermission" class="flex flex-col gap-2">
            <label class="text-sm font-medium">权限代码 <span class="text-red-500">*</span></label>
            <InputText
               v-model="permissionForm.code"
               placeholder="例如：view_report"
               class="w-full"
            />
            <p class="text-xs text-surface-400 m-0">系统会自动添加应用前缀</p>
         </div>
         <div v-if="!isEditingPermission" class="flex flex-col gap-2">
            <label class="text-sm font-medium">资源 <span class="text-red-500">*</span></label>
            <InputText
               v-model="permissionForm.resource"
               placeholder="例如：reports"
               class="w-full"
            />
         </div>
         <div v-if="!isEditingPermission" class="flex flex-col gap-2">
            <label class="text-sm font-medium">操作类型 <span class="text-red-500">*</span></label>
            <Select
               v-model="permissionForm.action"
               :options="PermissionActionOptions"
               optionLabel="label"
               optionValue="value"
               class="w-full"
            />
         </div>
         <div class="flex flex-col gap-2">
            <label class="text-sm font-medium">描述</label>
            <Textarea
               v-model="permissionForm.description"
               placeholder="权限的详细说明..."
               rows="3"
               class="w-full"
            />
         </div>
      </div>
      <template #footer>
         <Button label="取消" severity="secondary" @click="permissionDialog = false" />
         <Button
            :label="isEditingPermission ? '保存' : '创建'"
            :loading="
               createPermissionMutation.isPending.value || updatePermissionMutation.isPending.value
            "
            @click="savePermission"
         />
      </template>
   </Dialog>

   <!-- 角色编辑对话框 -->
   <Dialog
      v-model:visible="roleDialog"
      :header="isEditingRole ? '编辑角色' : '新建角色'"
      modal
      :style="{ width: '28rem' }"
   >
      <div class="flex flex-col gap-4">
         <div class="flex flex-col gap-2">
            <label class="text-sm font-medium">角色名称 <span class="text-red-500">*</span></label>
            <InputText v-model="roleForm.name" placeholder="例如：报表管理员" class="w-full" />
         </div>
         <div v-if="!isEditingRole" class="flex flex-col gap-2">
            <label class="text-sm font-medium">角色代码 <span class="text-red-500">*</span></label>
            <InputText v-model="roleForm.code" placeholder="例如：report_admin" class="w-full" />
            <p class="text-xs text-surface-400 m-0">系统会自动添加应用前缀</p>
         </div>
         <div class="flex flex-col gap-2">
            <label class="text-sm font-medium">描述</label>
            <Textarea
               v-model="roleForm.description"
               placeholder="角色的详细说明..."
               rows="3"
               class="w-full"
            />
         </div>
         <div class="flex items-center gap-2">
            <Checkbox v-model="roleForm.is_default" :binary="true" inputId="isDefault" />
            <label for="isDefault" class="text-sm">设为默认角色（新用户自动分配）</label>
         </div>
      </div>
      <template #footer>
         <Button label="取消" severity="secondary" @click="roleDialog = false" />
         <Button
            :label="isEditingRole ? '保存' : '创建'"
            :loading="createRoleMutation.isPending.value || updateRoleMutation.isPending.value"
            @click="saveRole"
         />
      </template>
   </Dialog>

   <!-- 角色权限分配对话框 -->
   <Dialog
      v-model:visible="rolePermissionsDialog"
      :header="`${selectedRoleForPermissions?.name || ''} - 权限分配`"
      modal
      :style="{ width: '36rem' }"
   >
      <div class="flex flex-col gap-4">
         <p class="text-sm text-surface-500 m-0">为此角色选择权限</p>

         <div
            v-if="permissions?.length"
            class="border border-surface-200 dark:border-surface-700 rounded-lg p-4 max-h-80 overflow-y-auto"
         >
            <div
               v-for="perm in permissions"
               :key="perm.id"
               class="flex items-center gap-3 py-2 border-b border-surface-100 dark:border-surface-800 last:border-b-0"
            >
               <Checkbox v-model="rolePermissionIds" :value="perm.id" :inputId="perm.id" />
               <label :for="perm.id" class="flex-1 cursor-pointer">
                  <div class="font-medium">{{ perm.name }}</div>
                  <div class="text-xs text-surface-400">{{ perm.code }}</div>
               </label>
            </div>
         </div>

         <div v-else class="text-center py-4 text-surface-400">暂无可用权限，请先创建权限</div>
      </div>
      <template #footer>
         <Button label="取消" severity="secondary" @click="rolePermissionsDialog = false" />
         <Button
            label="保存"
            :loading="setRolePermissionsMutation.isPending.value"
            @click="setRolePermissionsMutation.mutate()"
         />
      </template>
   </Dialog>

   <!-- 角色用户管理对话框 -->
   <Dialog
      v-model:visible="roleUsersDialog"
      :header="`${selectedRoleForUsers?.name || ''} - 用户管理`"
      modal
      :style="{ width: '36rem' }"
   >
      <div class="flex flex-col gap-4">
         <div class="flex gap-2">
            <Select
               v-model="newRoleUserId"
               :options="userOptions"
               optionLabel="label"
               optionValue="value"
               placeholder="选择要添加的用户"
               filter
               class="flex-1"
            />
            <Button
               label="添加"
               :loading="assignUserMutation.isPending.value"
               :disabled="!newRoleUserId"
               @click="assignUserMutation.mutate()"
            />
         </div>

         <DataTable
            :value="roleUsers || []"
            stripedRows
            size="small"
            class="rounded-lg overflow-hidden"
         >
            <Column field="name" header="姓名"></Column>
            <Column field="student_id" header="学号"></Column>
            <Column field="email" header="邮箱"></Column>
            <Column header="操作" :style="{ width: '80px' }">
               <template #body="{ data }">
                  <Button
                     icon="pi pi-trash"
                     severity="danger"
                     text
                     size="small"
                     @click="confirmRevokeUser(data)"
                  />
               </template>
            </Column>
         </DataTable>

         <div v-if="!roleUsers?.length" class="text-center py-4 text-surface-400">
            此角色暂无用户
         </div>
      </div>
      <template #footer>
         <Button label="关闭" severity="secondary" @click="roleUsersDialog = false" />
      </template>
   </Dialog>

   <ConfirmDialog />
</template>

<style scoped>
.app-group-dialog :deep(.p-tabview-panels) {
   padding: 1rem 0;
}
</style>
