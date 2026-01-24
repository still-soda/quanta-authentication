<script setup lang="ts">
import { computed } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import type { Role } from '@/types'

const props = defineProps<{
   visible: boolean
   role: Role | null
   isLoading?: boolean
}>()

const emit = defineEmits<{
   (e: 'update:visible', value: boolean): void
   (e: 'confirm'): void
}>()

const dialogVisible = computed({
   get: () => props.visible,
   set: value => emit('update:visible', value),
})

const handleConfirm = () => {
   emit('confirm')
}
</script>

<template>
   <Dialog
      v-model:visible="dialogVisible"
      header="确认删除"
      modal
      :closable="!isLoading"
      :closeOnEscape="!isLoading"
      :style="{ width: '26rem' }"
   >
      <div class="flex flex-col gap-4 py-2">
         <div class="flex items-start gap-4">
            <div
               class="w-12 h-12 flex items-center justify-center rounded-full bg-red-100 dark:bg-red-900/30 shrink-0"
            >
               <i class="pi pi-exclamation-triangle text-xl text-red-500"></i>
            </div>
            <div>
               <p class="text-surface-700 dark:text-surface-300 m-0 mb-2">
                  确定要删除角色
                  <strong class="text-surface-900 dark:text-surface-100">
                     {{ role?.name }}
                  </strong>
                  吗？
               </p>
               <p class="text-sm text-surface-500 dark:text-surface-400 m-0">
                  此操作将同时移除该角色下所有用户的角色关联，且不可恢复。
               </p>
            </div>
         </div>

         <!-- 角色信息摘要 -->
         <div
            v-if="role"
            class="flex items-center gap-4 p-3 bg-surface-50 dark:bg-surface-800 rounded-lg"
         >
            <div class="flex items-center gap-2">
               <i class="pi pi-users text-sm text-surface-400"></i>
               <span class="text-sm text-surface-600 dark:text-surface-400">
                  {{ role.user_count }} 个用户
               </span>
            </div>
            <div class="flex items-center gap-2">
               <i class="pi pi-key text-sm text-surface-400"></i>
               <span class="text-sm text-surface-600 dark:text-surface-400">
                  {{ role.permission_count }} 个权限
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
               :disabled="isLoading"
               @click="dialogVisible = false"
            />
            <Button
               label="确认删除"
               severity="danger"
               :loading="isLoading"
               @click="handleConfirm"
            />
         </div>
      </template>
   </Dialog>
</template>
