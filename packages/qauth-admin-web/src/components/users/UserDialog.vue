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
      :style="{ width: '28rem' }">
      <div class="flex flex-col gap-5 py-2">
         <div class="flex flex-col gap-2">
            <label
               for="name"
               class="text-sm font-medium text-surface-700 dark:text-surface-300">
               姓名
            </label>
            <InputText
               id="name"
               v-model="userForm.name"
               placeholder="请输入姓名"
               class="w-full" />
         </div>

         <div class="flex flex-col gap-2">
            <label
               for="email"
               class="text-sm font-medium text-surface-700 dark:text-surface-300">
               邮箱
            </label>
            <InputText
               id="email"
               v-model="userForm.email"
               placeholder="请输入邮箱"
               class="w-full" />
         </div>

         <div class="flex flex-col gap-2">
            <label
               for="role"
               class="text-sm font-medium text-surface-700 dark:text-surface-300">
               角色
            </label>
            <Select
               id="role"
               v-model="userForm.role"
               :options="roleOptions"
               optionLabel="label"
               optionValue="value"
               placeholder="选择角色"
               class="w-full" />
         </div>

         <div class="flex flex-row justify-between items-center gap-2">
            <label
               for="status"
               class="text-sm font-medium text-surface-700 dark:text-surface-300">
               账号状态
            </label>
            <div class="flex items-center gap-3">
               <InputSwitch id="status" v-model="userForm.status" />
               <span class="text-sm text-surface-600">
                  {{ userForm.status ? '启用' : '禁用' }}
               </span>
            </div>
         </div>
      </div>

      <template #footer>
         <div class="flex justify-end gap-3">
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
