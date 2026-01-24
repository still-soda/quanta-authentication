<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import type { RoleFormData } from '@/types'

const props = defineProps<{
   visible: boolean
   isEditing: boolean
   initialData?: Partial<RoleFormData>
   isLoading?: boolean
}>()

const emit = defineEmits<{
   (e: 'update:visible', value: boolean): void
   (e: 'save', data: RoleFormData): void
}>()

const roleForm = ref<RoleFormData>({
   name: '',
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
         roleForm.value = {
            name: props.initialData?.name || '',
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
         roleForm.value = {
            name: newData.name || '',
            code: newData.code || '',
            description: newData.description || '',
         }
      }
   },
   { deep: true }
)

// 表单验证
const isFormValid = computed(() => {
   return roleForm.value.name.trim() !== '' && roleForm.value.code.trim() !== ''
})

// 自动生成 code
const generateCode = () => {
   if (!props.isEditing && roleForm.value.name && !roleForm.value.code) {
      // 简单的中文转拼音（仅作为示例，实际可能需要库支持）
      roleForm.value.code = roleForm.value.name
         .toLowerCase()
         .replace(/[\s-]+/g, '_')
         .replace(/[^\w_]/g, '')
   }
}

const saveRole = () => {
   if (isFormValid.value) {
      emit('save', roleForm.value)
   }
}
</script>

<template>
   <Dialog
      v-model:visible="dialogVisible"
      :header="isEditing ? '编辑角色' : '新建角色'"
      modal
      :closable="!isLoading"
      :closeOnEscape="!isLoading"
      :style="{ width: '28rem' }"
   >
      <div class="flex flex-col gap-5 py-2">
         <div class="flex flex-col gap-2">
            <label
               for="roleName"
               class="text-sm font-medium text-surface-700 dark:text-surface-300"
            >
               角色名称 <span class="text-red-500">*</span>
            </label>
            <InputText
               id="roleName"
               v-model="roleForm.name"
               placeholder="例如：内容编辑"
               class="w-full"
               :disabled="isLoading"
               @blur="generateCode"
            />
         </div>

         <div class="flex flex-col gap-2">
            <label
               for="roleCode"
               class="text-sm font-medium text-surface-700 dark:text-surface-300"
            >
               角色标识 <span class="text-red-500">*</span>
            </label>
            <InputText
               id="roleCode"
               v-model="roleForm.code"
               placeholder="例如：content_editor"
               class="w-full"
               :disabled="isLoading || isEditing"
            />
            <small class="text-xs text-surface-400"> 唯一标识符，创建后不可修改 </small>
         </div>

         <div class="flex flex-col gap-2">
            <label
               for="roleDesc"
               class="text-sm font-medium text-surface-700 dark:text-surface-300"
            >
               描述
            </label>
            <Textarea
               id="roleDesc"
               v-model="roleForm.description"
               placeholder="角色职责描述..."
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
               @click="saveRole"
            />
         </div>
      </template>
   </Dialog>
</template>
