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
      <div class="dialog-content">
         <div class="form-field">
            <label for="appName">应用名称</label>
            <InputText
               id="appName"
               v-model="appForm.name"
               placeholder="例如：My Application"
               class="w-full" />
         </div>

         <div class="form-field">
            <label for="appDesc">应用描述</label>
            <Textarea
               id="appDesc"
               v-model="appForm.description"
               placeholder="描述应用用途..."
               rows="2"
               class="w-full" />
         </div>

         <div class="form-field">
            <label for="redirectUris">重定向 URI</label>
            <Textarea
               id="redirectUris"
               v-model="appForm.redirectUris"
               placeholder="每行一个 URI"
               rows="3"
               class="w-full" />
            <small class="field-hint">授权完成后的回调地址</small>
         </div>

         <div class="form-field">
            <label>授权范围</label>
            <div class="scope-selector">
               <button
                  v-for="scope in availableScopes"
                  :key="scope.value"
                  type="button"
                  class="scope-btn"
                  :class="{ active: appForm.scopes.includes(scope.value) }"
                  @click="toggleScope(scope.value)">
                  {{ scope.label }}
               </button>
            </div>
         </div>

         <div class="form-field switch-field">
            <label for="trusted">可信应用</label>
            <div class="switch-wrapper">
               <InputSwitch id="trusted" v-model="appForm.trusted" />
               <span class="switch-label">{{
                  appForm.trusted ? '跳过用户授权确认' : '需要用户确认授权'
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
            <Button :label="isEditing ? '保存' : '创建'" @click="saveApp" />
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

.scope-selector {
   display: flex;
   flex-wrap: wrap;
   gap: 0.5rem;
}

.scope-btn {
   padding: 0.375rem 0.75rem;
   border: 1px solid var(--p-surface-200);
   border-radius: 9999px;
   background: transparent;
   color: var(--p-surface-600);
   font-size: 0.8125rem;
   cursor: pointer;
   transition: all 0.2s ease;
}

.scope-btn:hover {
   border-color: var(--p-orange-300);
   color: var(--p-orange-600);
}

.scope-btn.active {
   background: var(--p-orange-500);
   border-color: var(--p-orange-500);
   color: white;
}

:global(.app-dark) .scope-btn {
   border-color: var(--p-surface-700);
   color: var(--p-surface-400);
}

:global(.app-dark) .scope-btn:hover {
   border-color: var(--p-orange-400);
   color: var(--p-orange-400);
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
   font-size: 0.8125rem;
   color: var(--p-surface-500);
}

.dialog-footer {
   display: flex;
   justify-content: flex-end;
   gap: 0.75rem;
}
</style>
