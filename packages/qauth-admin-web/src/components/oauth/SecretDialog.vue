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
      <div class="flex flex-col gap-5">
         <Message severity="warn" :closable="false">
            请妥善保存此密钥，关闭对话框后将无法再次查看
         </Message>

         <div class="flex flex-col gap-2">
            <label
               class="text-sm font-medium text-surface-700 dark:text-surface-300">
               新的 Client Secret
            </label>
            <div
               class="flex items-center gap-2 py-3 px-4 bg-surface-50 dark:bg-surface-800 rounded-lg border border-surface-200 dark:border-surface-700">
               <code
                  class="flex-1 text-sm font-mono text-surface-900 dark:text-surface-100 break-all">
                  {{ secret }}
               </code>
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
         <div class="flex justify-end gap-3">
            <Button label="我已保存" @click="emit('update:visible', false)" />
         </div>
      </template>
   </Dialog>
</template>
