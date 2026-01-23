<script setup lang="ts">
import { ref, computed } from 'vue';
import Dialog from 'primevue/dialog';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import InputSwitch from 'primevue/inputswitch';
import Select from 'primevue/select';

export interface UserFormData {
   name: string;
   email: string;
   role: string;
   status: boolean;
}

const props = defineProps<{
   visible: boolean;
   isEditing: boolean;
   initialData?: Partial<UserFormData>;
}>();

const emit = defineEmits<{
   (e: 'update:visible', value: boolean): void;
   (e: 'save', data: UserFormData): void;
}>();

const userForm = ref<UserFormData>({
   name: props.initialData?.name || '',
   email: props.initialData?.email || '',
   role: props.initialData?.role || '',
   status: props.initialData?.status ?? true,
});

const roleOptions = ref([
   { label: '管理员', value: '管理员' },
   { label: '开发者', value: '开发者' },
   { label: '普通用户', value: '普通用户' },
]);

const dialogVisible = computed({
   get: () => props.visible,
   set: (value) => emit('update:visible', value),
});

const saveUser = () => {
   emit('save', userForm.value);
};

// 监听 initialData 变化重置表单
const resetForm = (data?: Partial<UserFormData>) => {
   userForm.value = {
      name: data?.name || '',
      email: data?.email || '',
      role: data?.role || '',
      status: data?.status ?? true,
   };
};

defineExpose({ resetForm });
</script>

<template>
   <Dialog
      v-model:visible="dialogVisible"
      :header="isEditing ? '编辑用户' : '新建用户'"
      modal
      :style="{ width: '28rem' }"
      class="user-dialog">
      <div class="dialog-content">
         <div class="form-field">
            <label for="name">姓名</label>
            <InputText
               id="name"
               v-model="userForm.name"
               placeholder="请输入姓名"
               class="w-full" />
         </div>

         <div class="form-field">
            <label for="email">邮箱</label>
            <InputText
               id="email"
               v-model="userForm.email"
               placeholder="请输入邮箱"
               class="w-full" />
         </div>

         <div class="form-field">
            <label for="role">角色</label>
            <Select
               id="role"
               v-model="userForm.role"
               :options="roleOptions"
               optionLabel="label"
               optionValue="value"
               placeholder="选择角色"
               class="w-full" />
         </div>

         <div class="form-field switch-field">
            <label for="status">账号状态</label>
            <div class="switch-wrapper">
               <InputSwitch id="status" v-model="userForm.status" />
               <span class="switch-label">{{
                  userForm.status ? '启用' : '禁用'
               }}</span>
            </div>
         </div>
      </div>

      <template #footer>
         <div class="dialog-footer">
            <Button
               label="取消"
               severity="secondary"
               outlined
               @click="dialogVisible = false" />
            <Button :label="isEditing ? '保存' : '创建'" @click="saveUser" />
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

.switch-field {
   flex-direction: row;
   justify-content: space-between;
   align-items: center;
}

.switch-wrapper {
   display: flex;
   align-items: center;
   gap: 0.75rem;
}

.switch-label {
   font-size: 0.875rem;
   color: var(--p-surface-600);
}

.dialog-footer {
   display: flex;
   justify-content: flex-end;
   gap: 0.75rem;
}
</style>
