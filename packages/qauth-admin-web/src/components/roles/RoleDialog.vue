<script setup lang="ts">
import { ref, computed } from 'vue';
import Dialog from 'primevue/dialog';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Textarea from 'primevue/textarea';

export interface RoleFormData {
   name: string;
   code: string;
   description: string;
}

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
      <div class="dialog-content">
         <div class="form-field">
            <label for="roleName">角色名称</label>
            <InputText
               id="roleName"
               v-model="roleForm.name"
               placeholder="例如：内容编辑"
               class="w-full" />
         </div>

         <div class="form-field">
            <label for="roleCode">角色标识</label>
            <InputText
               id="roleCode"
               v-model="roleForm.code"
               placeholder="例如：content_editor"
               class="w-full" />
            <small class="field-hint">唯一标识符，仅支持小写字母和下划线</small>
         </div>

         <div class="form-field">
            <label for="roleDesc">描述</label>
            <Textarea
               id="roleDesc"
               v-model="roleForm.description"
               placeholder="角色职责描述..."
               rows="3"
               class="w-full" />
         </div>
      </div>

      <template #footer>
         <div class="dialog-footer">
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

<style scoped>
.dialog-content {
   display: flex;
   flex-direction: column;
   gap: 1.25rem;
   padding: 0.5rem 0;
}

.form-field {
   display: flex;
   flex-direction: column;
   gap: 0.5rem;
}

.form-field label {
   font-size: 0.875rem;
   font-weight: 500;
   color: var(--p-surface-700);
}

:global(.app-dark) .form-field label {
   color: var(--p-surface-300);
}

.field-hint {
   font-size: 0.75rem;
   color: var(--p-surface-400);
}

.dialog-footer {
   display: flex;
   justify-content: flex-end;
   gap: 0.75rem;
}
</style>
