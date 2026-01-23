<script setup lang="ts">
import { ref, computed } from 'vue';
import Dialog from 'primevue/dialog';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Textarea from 'primevue/textarea';
import type { RoleFormData } from '@/types';

const props = defineProps<{
   visible: boolean;
   isEditing: boolean;
   initialData?: Partial<RoleFormData>;
}>();

const emit = defineEmits<{
   (e: 'update:visible', value: boolean): void;
   (e: 'save', data: RoleFormData): void;
}>();

const roleForm = ref<RoleFormData>({
   name: props.initialData?.name || '',
   code: props.initialData?.code || '',
   description: props.initialData?.description || '',
});

const dialogVisible = computed({
   get: () => props.visible,
   set: (value) => emit('update:visible', value),
});

const saveRole = () => {
   emit('save', roleForm.value);
};

const resetForm = (data?: Partial<RoleFormData>) => {
   roleForm.value = {
      name: data?.name || '',
      code: data?.code || '',
      description: data?.description || '',
   };
};

defineExpose({ resetForm });
</script>

<template>
   <Dialog
      v-model:visible="dialogVisible"
      :header="isEditing ? '编辑角色' : '新建角色'"
      modal
      :style="{ width: '28rem' }">
      <div class="flex flex-col gap-5 py-2">
         <div class="flex flex-col gap-2">
            <label
               for="roleName"
               class="text-sm font-medium text-surface-700 dark:text-surface-300">
               角色名称
            </label>
            <InputText
               id="roleName"
               v-model="roleForm.name"
               placeholder="例如：内容编辑"
               class="w-full" />
         </div>

         <div class="flex flex-col gap-2">
            <label
               for="roleCode"
               class="text-sm font-medium text-surface-700 dark:text-surface-300">
               角色标识
            </label>
            <InputText
               id="roleCode"
               v-model="roleForm.code"
               placeholder="例如：content_editor"
               class="w-full" />
            <small class="text-xs text-surface-400">
               唯一标识符，仅支持小写字母和下划线
            </small>
         </div>

         <div class="flex flex-col gap-2">
            <label
               for="roleDesc"
               class="text-sm font-medium text-surface-700 dark:text-surface-300">
               描述
            </label>
            <Textarea
               id="roleDesc"
               v-model="roleForm.description"
               placeholder="角色职责描述..."
               rows="3"
               class="w-full" />
         </div>
      </div>

      <template #footer>
         <div class="flex justify-end gap-3">
            <Button
               label="取消"
               severity="secondary"
               outlined
               @click="dialogVisible = false" />
            <Button :label="isEditing ? '保存' : '创建'" @click="saveRole" />
         </div>
      </template>
   </Dialog>
</template>
