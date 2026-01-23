<script setup lang="ts">
import { computed } from 'vue';
import Dialog from 'primevue/dialog';
import Button from 'primevue/button';
import Message from 'primevue/message';

defineProps<{
   visible: boolean;
   secret: string;
}>();

const emit = defineEmits<{
   (e: 'update:visible', value: boolean): void;
}>();

const copyToClipboard = (text: string) => {
   navigator.clipboard.writeText(text);
};
</script>

<template>
   <Dialog
      :visible="visible"
      @update:visible="emit('update:visible', $event)"
      header="重置客户端密钥"
      modal
      :style="{ width: '28rem' }">
      <div class="secret-dialog-content">
         <Message severity="warn" :closable="false">
            请妥善保存此密钥，关闭对话框后将无法再次查看
         </Message>

         <div class="secret-display">
            <label>新的 Client Secret</label>
            <div class="secret-value">
               <code>{{ secret }}</code>
               <Button
                  icon="pi pi-copy"
                  text
                  rounded
                  severity="secondary"
                  @click="copyToClipboard(secret)"
                  v-tooltip.top="'复制'" />
            </div>
         </div>
      </div>

      <template #footer>
         <div class="dialog-footer">
            <Button label="我已保存" @click="emit('update:visible', false)" />
         </div>
      </template>
   </Dialog>
</template>

<style scoped>
.secret-dialog-content {
   display: flex;
   flex-direction: column;
   gap: 1.25rem;
}

.secret-display {
   display: flex;
   flex-direction: column;
   gap: 0.5rem;
}

.secret-display label {
   font-size: 0.875rem;
   font-weight: 500;
   color: var(--p-surface-700);
}

:global(.app-dark) .secret-display label {
   color: var(--p-surface-300);
}

.secret-value {
   display: flex;
   align-items: center;
   gap: 0.5rem;
   padding: 0.75rem 1rem;
   background: var(--p-surface-50);
   border-radius: 8px;
   border: 1px solid var(--p-surface-200);
}

:global(.app-dark) .secret-value {
   background: var(--p-surface-800);
   border-color: var(--p-surface-700);
}

.secret-value code {
   flex: 1;
   font-size: 0.875rem;
   font-family: monospace;
   color: var(--p-surface-900);
   word-break: break-all;
}

:global(.app-dark) .secret-value code {
   color: var(--p-surface-100);
}

.dialog-footer {
   display: flex;
   justify-content: flex-end;
   gap: 0.75rem;
}
</style>
