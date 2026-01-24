<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import MultiSelect from 'primevue/multiselect'
import type { User, CreateUserFormData, UpdateUserFormData, UserStatus } from '@/types'
import type { Role } from '@/types/role'

const props = defineProps<{
   visible: boolean
   isEditing: boolean
   user?: User | null
   roles: Role[]
   loading?: boolean
}>()

const emit = defineEmits<{
   (e: 'update:visible', value: boolean): void
   (e: 'create', data: CreateUserFormData): void
   (e: 'update', id: string, data: UpdateUserFormData): void
}>()

// 状态选项
const statusOptions = [
   { label: '正常', value: 'ACTIVE' },
   { label: '已锁定', value: 'LOCKED' },
   { label: '已禁用', value: 'BANNED' },
]

// 表单数据
const formData = ref<{
   student_id: string
   name: string
   email: string
   password: string
   phone: string
   display_name: string
   status: UserStatus
   role_ids: string[]
}>({
   student_id: '',
   name: '',
   email: '',
   password: '',
   phone: '',
   display_name: '',
   status: 'ACTIVE',
   role_ids: [],
})

// 监听 visible 和 user 变化来重置表单
watch(
   () => [props.visible, props.user],
   () => {
      if (props.visible) {
         if (props.isEditing && props.user) {
            formData.value = {
               student_id: props.user.student_id,
               name: props.user.name,
               email: props.user.email,
               password: '',
               phone: props.user.phone || '',
               display_name: props.user.display_name || '',
               status: props.user.status,
               role_ids: props.user.roles?.map(r => r.id) || [],
            }
         } else {
            formData.value = {
               student_id: '',
               name: '',
               email: '',
               password: '',
               phone: '',
               display_name: '',
               status: 'ACTIVE',
               role_ids: [],
            }
         }
      }
   },
   { immediate: true }
)

const dialogVisible = computed({
   get: () => props.visible,
   set: value => emit('update:visible', value),
})

// 角色选项
const roleOptions = computed(() =>
   props.roles.map(role => ({
      label: role.name,
      value: role.id,
   }))
)

// 表单验证
const isFormValid = computed(() => {
   if (props.isEditing) {
      return formData.value.name.trim() !== '' && formData.value.email.trim() !== ''
   }
   return (
      formData.value.student_id.trim() !== '' &&
      formData.value.name.trim() !== '' &&
      formData.value.email.trim() !== ''
   )
})

const handleSave = () => {
   if (!isFormValid.value) return

   if (props.isEditing && props.user) {
      const updateData: UpdateUserFormData = {
         name: formData.value.name,
         email: formData.value.email,
         phone: formData.value.phone || undefined,
         display_name: formData.value.display_name || undefined,
         status: formData.value.status,
      }
      emit('update', props.user.id, updateData)
   } else {
      const createData: CreateUserFormData = {
         student_id: formData.value.student_id,
         name: formData.value.name,
         email: formData.value.email,
         password: formData.value.password || undefined,
         phone: formData.value.phone || undefined,
         display_name: formData.value.display_name || undefined,
         role_ids: formData.value.role_ids.length > 0 ? formData.value.role_ids : undefined,
      }
      emit('create', createData)
   }
}
</script>

<template>
   <Dialog
      v-model:visible="dialogVisible"
      :header="isEditing ? '编辑用户' : '新建用户'"
      modal
      :style="{ width: '32rem' }"
   >
      <div class="flex flex-col gap-5 py-2">
         <!-- 学号（创建时必填，编辑时只读） -->
         <div class="flex flex-col gap-2">
            <label
               for="student_id"
               class="text-sm font-medium text-surface-700 dark:text-surface-300"
            >
               学号 <span v-if="!isEditing" class="text-red-500">*</span>
            </label>
            <InputText
               id="student_id"
               v-model="formData.student_id"
               placeholder="请输入学号"
               :disabled="isEditing"
               class="w-full"
            />
         </div>

         <!-- 姓名 -->
         <div class="flex flex-col gap-2">
            <label for="name" class="text-sm font-medium text-surface-700 dark:text-surface-300">
               姓名 <span class="text-red-500">*</span>
            </label>
            <InputText id="name" v-model="formData.name" placeholder="请输入姓名" class="w-full" />
         </div>

         <!-- 显示名称 -->
         <div class="flex flex-col gap-2">
            <label
               for="display_name"
               class="text-sm font-medium text-surface-700 dark:text-surface-300"
            >
               显示名称
            </label>
            <InputText
               id="display_name"
               v-model="formData.display_name"
               placeholder="可选，用于显示的昵称"
               class="w-full"
            />
         </div>

         <!-- 邮箱 -->
         <div class="flex flex-col gap-2">
            <label for="email" class="text-sm font-medium text-surface-700 dark:text-surface-300">
               邮箱 <span class="text-red-500">*</span>
            </label>
            <InputText
               id="email"
               v-model="formData.email"
               type="email"
               placeholder="请输入邮箱"
               class="w-full"
            />
         </div>

         <!-- 手机号 -->
         <div class="flex flex-col gap-2">
            <label for="phone" class="text-sm font-medium text-surface-700 dark:text-surface-300">
               手机号
            </label>
            <InputText id="phone" v-model="formData.phone" placeholder="可选" class="w-full" />
         </div>

         <!-- 密码（仅创建时显示） -->
         <div v-if="!isEditing" class="flex flex-col gap-2">
            <label
               for="password"
               class="text-sm font-medium text-surface-700 dark:text-surface-300"
            >
               初始密码
            </label>
            <InputText
               id="password"
               v-model="formData.password"
               type="password"
               placeholder="留空则自动生成"
               class="w-full"
            />
            <small class="text-surface-500">如不填写，系统将自动生成随机密码</small>
         </div>

         <!-- 角色（创建时可选，编辑时单独管理） -->
         <div v-if="!isEditing" class="flex flex-col gap-2">
            <label for="roles" class="text-sm font-medium text-surface-700 dark:text-surface-300">
               分配角色
            </label>
            <MultiSelect
               id="roles"
               v-model="formData.role_ids"
               :options="roleOptions"
               optionLabel="label"
               optionValue="value"
               placeholder="选择角色（可多选）"
               display="chip"
               class="w-full"
            />
         </div>

         <!-- 状态（编辑时可修改） -->
         <div v-if="isEditing" class="flex flex-col gap-2">
            <label for="status" class="text-sm font-medium text-surface-700 dark:text-surface-300">
               账号状态
            </label>
            <Select
               id="status"
               v-model="formData.status"
               :options="statusOptions"
               optionLabel="label"
               optionValue="value"
               placeholder="选择状态"
               class="w-full"
            />
         </div>
      </div>

      <template #footer>
         <div class="flex justify-end gap-3">
            <Button label="取消" severity="secondary" outlined @click="dialogVisible = false" />
            <Button
               :label="isEditing ? '保存' : '创建'"
               :disabled="!isFormValid"
               :loading="loading"
               @click="handleSave"
            />
         </div>
      </template>
   </Dialog>
</template>
