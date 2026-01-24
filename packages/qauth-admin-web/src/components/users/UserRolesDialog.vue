<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import Checkbox from 'primevue/checkbox'
import Tag from 'primevue/tag'
import type { User } from '@/types'
import type { Role } from '@/types/role'

const props = defineProps<{
   visible: boolean
   user: User | null
   roles: Role[]
   loading?: boolean
}>()

const emit = defineEmits<{
   (e: 'update:visible', value: boolean): void
   (e: 'save', userId: string, roleIds: string[]): void
}>()

// 选中的角色 ID
const selectedRoleIds = ref<string[]>([])

// 监听 visible 和 user 变化来重置选中状态
watch(
   () => [props.visible, props.user],
   () => {
      if (props.visible && props.user) {
         selectedRoleIds.value = props.user.roles?.map(r => r.id) || []
      }
   },
   { immediate: true }
)

const dialogVisible = computed({
   get: () => props.visible,
   set: value => emit('update:visible', value),
})

// 判断角色是否是系统角色
const isSystemRole = (role: Role) => role.is_system

// 切换角色选中状态
const toggleRole = (roleId: string) => {
   const index = selectedRoleIds.value.indexOf(roleId)
   if (index === -1) {
      selectedRoleIds.value.push(roleId)
   } else {
      selectedRoleIds.value.splice(index, 1)
   }
}

// 检查是否有变化
const hasChanges = computed(() => {
   const currentIds = props.user?.roles?.map(r => r.id) || []
   if (currentIds.length !== selectedRoleIds.value.length) return true
   return !currentIds.every(id => selectedRoleIds.value.includes(id))
})

const handleSave = () => {
   if (props.user) {
      emit('save', props.user.id, selectedRoleIds.value)
   }
}
</script>

<template>
   <Dialog v-model:visible="dialogVisible" header="管理用户角色" modal :style="{ width: '36rem' }">
      <div v-if="user" class="flex flex-col gap-4">
         <!-- 用户信息 -->
         <div class="p-4 bg-surface-50 dark:bg-surface-800 rounded-lg">
            <div class="flex items-center gap-3">
               <div class="flex flex-col">
                  <span class="font-semibold text-surface-900 dark:text-surface-100">
                     {{ user.display_name || user.name }}
                  </span>
                  <span class="text-sm text-surface-500">{{ user.email }}</span>
               </div>
            </div>
         </div>

         <!-- 角色列表 -->
         <div class="flex flex-col gap-2">
            <label class="text-sm font-medium text-surface-700 dark:text-surface-300">
               选择角色
            </label>
            <div
               class="border border-surface-200 dark:border-surface-700 rounded-lg divide-y divide-surface-200 dark:divide-surface-700"
            >
               <div
                  v-for="role in roles"
                  :key="role.id"
                  class="flex items-center justify-between p-3 hover:bg-surface-50 dark:hover:bg-surface-800 transition-colors cursor-pointer"
                  @click="toggleRole(role.id)"
               >
                  <div class="flex items-center gap-3">
                     <Checkbox
                        :modelValue="selectedRoleIds.includes(role.id)"
                        :binary="true"
                        @click.stop
                     />
                     <div class="flex flex-col">
                        <div class="flex items-center gap-2">
                           <span class="font-medium text-surface-900 dark:text-surface-100">
                              {{ role.name }}
                           </span>
                           <Tag v-if="isSystemRole(role)" severity="secondary" class="text-xs">
                              系统
                           </Tag>
                        </div>
                        <span class="text-sm text-surface-500">
                           {{ role.description || role.code }}
                        </span>
                     </div>
                  </div>
                  <div class="text-sm text-surface-400">{{ role.user_count }} 个用户</div>
               </div>
               <div v-if="roles.length === 0" class="p-8 text-center text-surface-400">
                  <i class="pi pi-shield text-3xl mb-2"></i>
                  <p>暂无可用角色</p>
               </div>
            </div>
         </div>

         <!-- 已选择的角色 -->
         <div v-if="selectedRoleIds.length > 0" class="flex flex-col gap-2">
            <label class="text-sm font-medium text-surface-700 dark:text-surface-300">
               已选择 {{ selectedRoleIds.length }} 个角色
            </label>
            <div class="flex flex-wrap gap-2">
               <Tag v-for="roleId in selectedRoleIds" :key="roleId" severity="info" rounded>
                  {{ roles.find(r => r.id === roleId)?.name }}
               </Tag>
            </div>
         </div>
      </div>

      <template #footer>
         <div class="flex justify-end gap-3">
            <Button label="取消" severity="secondary" outlined @click="dialogVisible = false" />
            <Button label="保存" :disabled="!hasChanges" :loading="loading" @click="handleSave" />
         </div>
      </template>
   </Dialog>
</template>
