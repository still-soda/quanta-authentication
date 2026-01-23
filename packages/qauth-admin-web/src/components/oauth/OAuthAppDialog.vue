<script setup lang="ts">
import { ref, computed } from 'vue';
import Dialog from 'primevue/dialog';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Textarea from 'primevue/textarea';
import InputSwitch from 'primevue/inputswitch';

export interface OAuthAppFormData {
   name: string;
   description: string;
   redirectUris: string;
   scopes: string[];
   trusted: boolean;
}

const props = defineProps<{
   visible: boolean;
   isEditing: boolean;
   initialData?: Partial<OAuthAppFormData>;
}>();

const emit = defineEmits<{
   (e: 'update:visible', value: boolean): void;
   (e: 'save', data: OAuthAppFormData): void;
}>();

const appForm = ref<OAuthAppFormData>({
   name: props.initialData?.name || '',
   description: props.initialData?.description || '',
   redirectUris: props.initialData?.redirectUris || '',
   scopes: props.initialData?.scopes || ['openid', 'profile'],
   trusted: props.initialData?.trusted || false,
});

const availableScopes = [
   { label: 'OpenID', value: 'openid' },
   { label: 'Profile', value: 'profile' },
   { label: 'Email', value: 'email' },
   { label: 'Admin', value: 'admin' },
   { label: 'Read Users', value: 'read:users' },
   { label: 'Write Users', value: 'write:users' },
   { label: 'Offline Access', value: 'offline_access' },
];

const dialogVisible = computed({
   get: () => props.visible,
   set: (value) => emit('update:visible', value),
});

const toggleScope = (scope: string) => {
   const index = appForm.value.scopes.indexOf(scope);
   if (index > -1) {
      appForm.value.scopes.splice(index, 1);
   } else {
      appForm.value.scopes.push(scope);
   }
};

const saveApp = () => {
   emit('save', appForm.value);
};

const resetForm = (data?: Partial<OAuthAppFormData>) => {
   appForm.value = {
      name: data?.name || '',
      description: data?.description || '',
      redirectUris: data?.redirectUris || '',
      scopes: data?.scopes || ['openid', 'profile'],
      trusted: data?.trusted || false,
   };
};

defineExpose({ resetForm });
</script>

<template>
   <Dialog
      v-model:visible="dialogVisible"
      :header="isEditing ? '编辑应用' : '新建应用'"
      modal
      :style="{ width: '32rem' }">
      <div class="flex flex-col gap-5 py-2">
         <div class="flex flex-col gap-2">
            <label
               for="appName"
               class="text-sm font-medium text-surface-700 dark:text-surface-300">
               应用名称
            </label>
            <InputText
               id="appName"
               v-model="appForm.name"
               placeholder="例如：My Application"
               class="w-full" />
         </div>

         <div class="flex flex-col gap-2">
            <label
               for="appDesc"
               class="text-sm font-medium text-surface-700 dark:text-surface-300">
               应用描述
            </label>
            <Textarea
               id="appDesc"
               v-model="appForm.description"
               placeholder="描述应用用途..."
               rows="2"
               class="w-full" />
         </div>

         <div class="flex flex-col gap-2">
            <label
               for="redirectUris"
               class="text-sm font-medium text-surface-700 dark:text-surface-300">
               重定向 URI
            </label>
            <Textarea
               id="redirectUris"
               v-model="appForm.redirectUris"
               placeholder="每行一个 URI"
               rows="3"
               class="w-full" />
            <small class="text-xs text-surface-400">
               授权完成后的回调地址
            </small>
         </div>

         <div class="flex flex-col gap-2">
            <label
               class="text-sm font-medium text-surface-700 dark:text-surface-300">
               授权范围
            </label>
            <div class="flex flex-wrap gap-2">
               <button
                  v-for="scope in availableScopes"
                  :key="scope.value"
                  type="button"
                  class="py-1.5 px-3 border border-surface-200 dark:border-surface-700 rounded-full bg-transparent text-surface-600 dark:text-surface-400 text-[0.8125rem] cursor-pointer transition-all duration-200 ease hover:border-primary-300 hover:text-primary-600 dark:hover:border-primary-400 dark:hover:text-primary-400"
                  :class="{
                     'bg-primary-500! border-primary-500! text-white!':
                        appForm.scopes.includes(scope.value),
                  }"
                  @click="toggleScope(scope.value)">
                  {{ scope.label }}
               </button>
            </div>
         </div>

         <div class="flex flex-row justify-between items-center gap-2">
            <label
               for="trusted"
               class="text-sm font-medium text-surface-700 dark:text-surface-300">
               可信应用
            </label>
            <div class="flex items-center gap-3">
               <InputSwitch id="trusted" v-model="appForm.trusted" />
               <span class="text-[0.8125rem] text-surface-500">
                  {{ appForm.trusted ? '跳过用户授权确认' : '需要用户确认授权' }}
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
            <Button :label="isEditing ? '保存' : '创建'" @click="saveApp" />
         </div>
      </template>
   </Dialog>
</template>
