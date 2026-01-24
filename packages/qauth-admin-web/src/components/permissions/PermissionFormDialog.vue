<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Select from 'primevue/select'
import type { PermissionFormData, PermissionAction } from '@/types'

const props = defineProps<{
   visible: boolean
   isEditing: boolean
   initialData?: Partial<PermissionFormData>
   isLoading?: boolean
}>()

const emit = defineEmits<{
   (e: 'update:visible', value: boolean): void
   (e: 'save', data: PermissionFormData): void
}>()

// 资源选项
const resourceOptions = [
   { label: 'OAuth 应用', value: 'oauth_clients' },
   { label: '系统管理', value: 'system' },
   { label: '角色管理', value: 'roles' },
   { label: '审计日志', value: 'audit' },
   { label: '用户管理', value: 'users' },
   { label: '权限管理', value: 'permissions' },
]

// 操作类型选项
const actionOptions = [
   { label: '创建 (Create)', value: 1 },
   { label: '查看 (Read)', value: 2 },
   { label: '更新 (Update)', value: 3 },
   { label: '删除 (Delete)', value: 4 },
]

const permissionForm = ref<PermissionFormData>({
   resource: '',
   action: 2 as PermissionAction,
   code: '',
   description: '',
})

const dialogVisible = computed({
   get: () => props.visible,
   set: value => emit('update:visible', value),
})

// 监听 visible 变化，重置表单
watch(
   () => props.visible,
   newVal => {
      if (newVal) {
         permissionForm.value = {
            resource: props.initialData?.resource || '',
            action: props.initialData?.action || (2 as PermissionAction),
            code: props.initialData?.code || '',
            description: props.initialData?.description || '',
         }
      }
   }
)

// 监听 initialData 变化
watch(
   () => props.initialData,
   newData => {
      if (props.visible && newData) {
         permissionForm.value = {
            resource: newData.resource || '',
            action: newData.action || (2 as PermissionAction),
            code: newData.code || '',
            description: newData.description || '',
         }
      }
   },
   { deep: true }
)

// 表单验证
const isFormValid = computed(() => {
   return (
      permissionForm.value.resource.trim() !== '' &&
      permissionForm.value.action !== undefined &&
      permissionForm.value.code.trim() !== ''
   )
})

// 自动生成 code
const generateCode = () => {
   if (!props.isEditing && permissionForm.value.resource && permissionForm.value.action) {
      const actionMap: Record<number, string> = {
         1: 'create',
         2: 'view',
         3: 'update',
         4: 'delete',
      }
      const action = actionMap[permissionForm.value.action] || 'view'
      permissionForm.value.code = `${permissionForm.value.resource}_${action}`
   }
}

// 监听资源和操作变化
watch(
   () => [permissionForm.value.resource, permissionForm.value.action],
   () => {
      if (!props.isEditing) {
         generateCode()
      }
   }
)

const savePermission = () => {
   if (isFormValid.value) {
      emit('save', permissionForm.value)
   }
}
</script>

<template>
   <Dialog
      v-model:visible="dialogVisible"
      :header="isEditing ? '编辑权限' : '新建权限'"
      modal
      :closable="!isLoading"
      :closeOnEscape="!isLoading"
      :style="{ width: '32rem' }"
   >
      <div class="flex flex-col gap-5 py-2">
         <div class="flex flex-col gap-2">
            <label
               for="permResource"
               class="text-sm font-medium text-surface-700 dark:text-surface-300"
            >
               资源 <span class="text-red-500">*</span>
            </label>
            <Select
               id="permResource"
               v-model="permissionForm.resource"
               :options="resourceOptions"
               optionLabel="label"
               optionValue="value"
               placeholder="选择资源"
               class="w-full"
               :disabled="isLoading"
            />
         </div>

         <div class="flex flex-col gap-2">
            <label
               for="permAction"
               class="text-sm font-medium text-surface-700 dark:text-surface-300"
            >
               操作类型 <span class="text-red-500">*</span>
            </label>
            <Select
               id="permAction"
               v-model="permissionForm.action"
               :options="actionOptions"
               optionLabel="label"
               optionValue="value"
               placeholder="选择操作类型"
               class="w-full"
               :disabled="isLoading"
            />
         </div>

         <div class="flex flex-col gap-2">
            <label
               for="permCode"
               class="text-sm font-medium text-surface-700 dark:text-surface-300"
            >
               权限代码 <span class="text-red-500">*</span>
            </label>
            <InputText
               id="permCode"
               v-model="permissionForm.code"
               placeholder="例如：oauth_client_create"
               class="w-full"
               :disabled="isLoading || isEditing"
            />
            <small class="text-xs text-surface-400"> 唯一标识符，创建后不可修改 </small>
         </div>

         <div class="flex flex-col gap-2">
            <label
               for="permDesc"
               class="text-sm font-medium text-surface-700 dark:text-surface-300"
            >
               描述
            </label>
            <Textarea
               id="permDesc"
               v-model="permissionForm.description"
               placeholder="权限用途描述..."
               rows="3"
               class="w-full"
               :disabled="isLoading"
            />
         </div>
      </div>

      <template #footer>
         <div class="flex justify-end gap-3">
            <Button
               label="取消"
               severity="secondary"
               outlined
               :disabled="isLoading"
               @click="dialogVisible = false"
            />
            <Button
               :label="isEditing ? '保存' : '创建'"
               :loading="isLoading"
               :disabled="!isFormValid"
               @click="savePermission"
            />
         </div>
      </template>
   </Dialog>
</template>
