<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'

const props = defineProps<{
   visible: boolean
   title?: string
   itemName: string
   itemType?: string
}>()

const emit = defineEmits<{
   (e: 'update:visible', value: boolean): void
   (e: 'confirm'): void
}>()

const confirmInput = ref('')

const dialogVisible = computed({
   get: () => props.visible,
   set: value => emit('update:visible', value),
})

const isConfirmValid = computed(() => {
   return confirmInput.value === props.itemName
})

watch(
   () => props.visible,
   visible => {
      if (visible) {
         confirmInput.value = ''
      }
   }
)

const handleConfirm = () => {
   if (isConfirmValid.value) {
      emit('confirm')
      dialogVisible.value = false
   }
}
</script>

<template>
   <Dialog
      v-model:visible="dialogVisible"
      :header="title || '删除确认'"
      modal
      :style="{ width: '28rem' }"
      :breakpoints="{ '640px': '90vw' }"
   >
      <div class="flex flex-col gap-4">
         <div class="flex items-center gap-3 p-3 bg-red-50 dark:bg-red-900/20 rounded-lg">
            <i class="pi pi-exclamation-triangle text-2xl text-red-500"></i>
            <div class="flex flex-col gap-1">
               <span class="text-sm font-medium text-surface-700 dark:text-surface-300">
                  此操作不可撤销
               </span>
               <span class="text-xs text-surface-500">
                  删除后，{{ itemType || '该项目' }}的所有数据将被永久删除。
               </span>
            </div>
         </div>

         <div class="flex flex-col gap-2">
            <label class="text-sm text-surface-600 dark:text-surface-400">
               请输入
               <code
                  class="px-1.5 py-0.5 bg-surface-100 dark:bg-surface-700 rounded text-primary-600 dark:text-primary-400 font-mono text-sm"
               >
                  {{ itemName }}
               </code>
               以确认删除
            </label>
            <InputText
               v-model="confirmInput"
               :placeholder="`输入 ${itemName}`"
               class="w-full font-mono"
               @keyup.enter="handleConfirm"
            />
         </div>
      </div>

      <template #footer>
         <div class="flex justify-end gap-3">
            <Button label="取消" severity="secondary" outlined @click="dialogVisible = false" />
            <Button
               label="确认删除"
               severity="danger"
               :disabled="!isConfirmValid"
               @click="handleConfirm"
            />
         </div>
      </template>
   </Dialog>
</template>
