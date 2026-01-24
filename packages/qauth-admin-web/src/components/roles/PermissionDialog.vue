<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import Checkbox from 'primevue/checkbox'
import Accordion from 'primevue/accordion'
import AccordionPanel from 'primevue/accordionpanel'
import AccordionHeader from 'primevue/accordionheader'
import AccordionContent from 'primevue/accordioncontent'
import InputText from 'primevue/inputtext'
import type { Role, PermissionGroup } from '@/types'

const props = defineProps<{
   visible: boolean
   role: Role | null
   permissionGroups: PermissionGroup[]
   isLoading?: boolean
}>()

const emit = defineEmits<{
   (e: 'update:visible', value: boolean): void
   (e: 'save', groups: PermissionGroup[]): void
}>()

// 本地权限组副本，用于编辑
const localGroups = ref<PermissionGroup[]>([])

// 搜索关键词
const searchKeyword = ref('')

// 监听 permissionGroups 变化，创建本地副本
watch(
   () => props.permissionGroups,
   newGroups => {
      localGroups.value = JSON.parse(JSON.stringify(newGroups))
   },
   { immediate: true, deep: true }
)

// 监听对话框打开，重置本地副本
watch(
   () => props.visible,
   newVal => {
      if (newVal) {
         localGroups.value = JSON.parse(JSON.stringify(props.permissionGroups))
         searchKeyword.value = ''
      }
   }
)

const dialogVisible = computed({
   get: () => props.visible,
   set: value => emit('update:visible', value),
})

// 过滤后的权限组
const filteredGroups = computed(() => {
   if (!searchKeyword.value.trim()) {
      return localGroups.value
   }

   const keyword = searchKeyword.value.toLowerCase()
   return localGroups.value
      .map(group => ({
         ...group,
         permissions: group.permissions.filter(
            p => p.name.toLowerCase().includes(keyword) || p.code.toLowerCase().includes(keyword)
         ),
      }))
      .filter(group => group.permissions.length > 0)
})

const toggleGroupPermissions = (group: PermissionGroup, checked: boolean) => {
   // 找到本地组并更新
   const localGroup = localGroups.value.find(g => g.resource === group.resource)
   if (localGroup) {
      localGroup.permissions.forEach(p => (p.checked = checked))
   }
}

const isGroupChecked = (group: PermissionGroup) => {
   const localGroup = localGroups.value.find(g => g.resource === group.resource)
   if (!localGroup) return false
   return localGroup.permissions.every(p => p.checked)
}

const isGroupIndeterminate = (group: PermissionGroup) => {
   const localGroup = localGroups.value.find(g => g.resource === group.resource)
   if (!localGroup) return false
   const checkedCount = localGroup.permissions.filter(p => p.checked).length
   return checkedCount > 0 && checkedCount < localGroup.permissions.length
}

// 统计选中的权限数量
const selectedCount = computed(() => {
   return localGroups.value.reduce((acc, group) => {
      return acc + group.permissions.filter(p => p.checked).length
   }, 0)
})

// 总权限数量
const totalCount = computed(() => {
   return localGroups.value.reduce((acc, group) => {
      return acc + group.permissions.length
   }, 0)
})

// 全选/取消全选
const selectAll = () => {
   const allChecked = selectedCount.value === totalCount.value
   localGroups.value.forEach(group => {
      group.permissions.forEach(p => (p.checked = !allChecked))
   })
}

const savePermissions = () => {
   emit('save', localGroups.value)
}

// 更新单个权限的选中状态
const updatePermissionChecked = (groupResource: string, permId: string, checked: boolean) => {
   const group = localGroups.value.find(g => g.resource === groupResource)
   if (group) {
      const perm = group.permissions.find(p => p.id === permId)
      if (perm) {
         perm.checked = checked
      }
   }
}

// 获取单个权限的选中状态
const getPermissionChecked = (groupResource: string, permId: string) => {
   const group = localGroups.value.find(g => g.resource === groupResource)
   return group?.permissions.find(p => p.id === permId)?.checked || false
}
</script>

<template>
   <Dialog
      v-model:visible="dialogVisible"
      :header="`配置权限 - ${role?.name || ''}`"
      modal
      :closable="!isLoading"
      :closeOnEscape="!isLoading"
      :style="{ width: '40rem' }"
   >
      <!-- 搜索和统计 -->
      <div class="flex items-center justify-between gap-4 mb-4">
         <div class="flex-1">
            <InputText
               v-model="searchKeyword"
               placeholder="搜索权限..."
               class="w-full"
               :disabled="isLoading"
            >
               <template #prefix>
                  <i class="pi pi-search text-surface-400"></i>
               </template>
            </InputText>
         </div>
         <div class="flex items-center gap-3">
            <span class="text-sm text-surface-500 dark:text-surface-400">
               已选 {{ selectedCount }} / {{ totalCount }}
            </span>
            <Button
               :label="selectedCount === totalCount ? '取消全选' : '全选'"
               severity="secondary"
               text
               size="small"
               :disabled="isLoading"
               @click="selectAll"
            />
         </div>
      </div>

      <!-- 权限列表 -->
      <div class="max-h-[55vh] overflow-y-auto pr-1">
         <div v-if="filteredGroups.length === 0" class="py-8 text-center">
            <i class="pi pi-search text-3xl text-surface-300 dark:text-surface-600 mb-3"></i>
            <p class="text-surface-500 dark:text-surface-400">
               {{ searchKeyword ? '未找到匹配的权限' : '暂无可配置的权限' }}
            </p>
         </div>

         <Accordion v-else multiple :value="filteredGroups.map(g => g.resource)">
            <AccordionPanel
               v-for="group in filteredGroups"
               :key="group.resource"
               :value="group.resource"
            >
               <AccordionHeader>
                  <div class="flex justify-between items-center w-full pr-2">
                     <div class="flex items-center gap-3">
                        <i :class="group.icon" class="text-base text-primary-500"></i>
                        <span class="font-medium">{{ group.name }}</span>
                        <span class="text-xs text-surface-400 dark:text-surface-500">
                           ({{
                              localGroups
                                 .find(g => g.resource === group.resource)
                                 ?.permissions.filter(p => p.checked).length || 0
                           }}/{{ group.permissions.length }})
                        </span>
                     </div>
                     <Checkbox
                        :modelValue="isGroupChecked(group)"
                        :indeterminate="isGroupIndeterminate(group)"
                        :disabled="isLoading"
                        @update:modelValue="toggleGroupPermissions(group, $event as boolean)"
                        @click.stop
                     />
                  </div>
               </AccordionHeader>
               <AccordionContent>
                  <div class="grid grid-cols-2 gap-3 py-2">
                     <div
                        v-for="perm in group.permissions"
                        :key="perm.id"
                        class="flex items-start gap-3 p-2 rounded-lg hover:bg-surface-50 dark:hover:bg-surface-800 transition-colors"
                     >
                        <Checkbox
                           :modelValue="getPermissionChecked(group.resource, perm.id)"
                           @update:modelValue="
                              updatePermissionChecked(group.resource, perm.id, $event as boolean)
                           "
                           :inputId="perm.id"
                           :binary="true"
                           :disabled="isLoading"
                        />
                        <label :for="perm.id" class="flex flex-col gap-0.5 text-sm cursor-pointer">
                           <span class="text-surface-700 dark:text-surface-300">
                              {{ perm.name }}
                           </span>
                           <code
                              class="text-[0.6875rem] font-mono text-surface-400 dark:text-surface-500"
                           >
                              {{ perm.code }}
                           </code>
                        </label>
                     </div>
                  </div>
               </AccordionContent>
            </AccordionPanel>
         </Accordion>
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
            <Button label="保存权限" :loading="isLoading" @click="savePermissions" />
         </div>
      </template>
   </Dialog>
</template>
