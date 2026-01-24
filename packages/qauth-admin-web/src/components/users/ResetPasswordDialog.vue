<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Message from 'primevue/message'
import type { User } from '@/types'

const props = defineProps<{
   visible: boolean
   user: User | null
   loading?: boolean
   generatedPassword?: string
}>()

const emit = defineEmits<{
   (e: 'update:visible', value: boolean): void
   (e: 'reset', userId: string, newPassword?: string): void
}>()

const newPassword = ref('')
const useRandomPassword = ref(true)

// 监听 visible 变化来重置表单
watch(
   () => props.visible,
   visible => {
      if (visible) {
         newPassword.value = ''
         useRandomPassword.value = true
      }
   }
)

const dialogVisible = computed({
   get: () => props.visible,
   set: value => emit('update:visible', value),
})

const handleReset = () => {
   if (props.user) {
      emit('reset', props.user.id, useRandomPassword.value ? undefined : newPassword.value)
   }
}

const isValid = computed(() => {
   if (useRandomPassword.value) return true
   return newPassword.value.length >= 6
})
</script>

<template>
   <Dialog v-model:visible="dialogVisible" header="重置用户密码" modal :style="{ width: '28rem' }">
      <div v-if="user" class="flex flex-col gap-4">
         <!-- 用户信息 -->
         <div class="p-4 bg-surface-50 dark:bg-surface-800 rounded-lg">
            <div class="flex flex-col">
               <span class="font-semibold text-surface-900 dark:text-surface-100">
                  {{ user.display_name || user.name }}
               </span>
               <span class="text-sm text-surface-500">{{ user.email }}</span>
            </div>
         </div>

         <!-- 显示生成的密码 -->
         <Message v-if="generatedPassword" severity="success" :closable="false">
            <div class="flex flex-col gap-2">
               <span>密码已重置成功！新密码为：</span>
               <code class="bg-surface-100 dark:bg-surface-700 px-2 py-1 rounded font-mono text-lg">
                  {{ generatedPassword }}
               </code>
               <span class="text-xs">请妥善保存此密码，关闭对话框后将无法再次查看</span>
            </div>
         </Message>

         <!-- 密码设置选项 -->
         <div v-else class="flex flex-col gap-4">
            <div class="flex items-center gap-3">
               <input
                  type="radio"
                  id="random"
                  :checked="useRandomPassword"
                  @change="useRandomPassword = true"
               />
               <label for="random" class="text-sm text-surface-700 dark:text-surface-300">
                  自动生成随机密码
               </label>
            </div>

            <div class="flex items-center gap-3">
               <input
                  type="radio"
                  id="custom"
                  :checked="!useRandomPassword"
                  @change="useRandomPassword = false"
               />
               <label for="custom" class="text-sm text-surface-700 dark:text-surface-300">
                  手动设置密码
               </label>
            </div>

            <div v-if="!useRandomPassword" class="flex flex-col gap-2">
               <label
                  for="new_password"
                  class="text-sm font-medium text-surface-700 dark:text-surface-300"
               >
                  新密码
               </label>
               <InputText
                  id="new_password"
                  v-model="newPassword"
                  type="password"
                  placeholder="请输入新密码（至少6位）"
                  class="w-full"
               />
            </div>

            <Message severity="warn" :closable="false">
               重置密码后，用户需要使用新密码登录
            </Message>
         </div>
      </div>

      <template #footer>
         <div class="flex justify-end gap-3">
            <Button
               :label="generatedPassword ? '关闭' : '取消'"
               severity="secondary"
               outlined
               @click="dialogVisible = false"
            />
            <Button
               v-if="!generatedPassword"
               label="重置密码"
               severity="warn"
               :disabled="!isValid"
               :loading="loading"
               @click="handleReset"
            />
         </div>
      </template>
   </Dialog>
</template>
