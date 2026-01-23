<script setup lang="ts">
import { computed } from 'vue';
import Dialog from 'primevue/dialog';
import Button from 'primevue/button';
import Checkbox from 'primevue/checkbox';
import Accordion from 'primevue/accordion';
import AccordionPanel from 'primevue/accordionpanel';
import AccordionHeader from 'primevue/accordionheader';
import AccordionContent from 'primevue/accordioncontent';
import type { Role, PermissionGroup } from '@/types';

const props = defineProps<{
   visible: boolean;
   role: Role | null;
   permissionGroups: PermissionGroup[];
}>();

const emit = defineEmits<{
   (e: 'update:visible', value: boolean): void;
   (e: 'save', groups: PermissionGroup[]): void;
}>();

const dialogVisible = computed({
   get: () => props.visible,
   set: (value) => emit('update:visible', value),
});

const toggleGroupPermissions = (group: PermissionGroup, checked: boolean) => {
   group.permissions.forEach((p) => (p.checked = checked));
};

const isGroupChecked = (group: PermissionGroup) => {
   return group.permissions.every((p) => p.checked);
};

const isGroupIndeterminate = (group: PermissionGroup) => {
   const checkedCount = group.permissions.filter((p) => p.checked).length;
   return checkedCount > 0 && checkedCount < group.permissions.length;
};

const savePermissions = () => {
   emit('save', props.permissionGroups);
};
</script>

<template>
   <Dialog
      v-model:visible="dialogVisible"
      :header="`配置权限 - ${role?.name}`"
      modal
      :style="{ width: '36rem' }">
      <div class="max-h-[60vh] overflow-y-auto">
         <Accordion multiple>
            <AccordionPanel
               v-for="group in permissionGroups"
               :key="group.name"
               :value="group.name">
               <AccordionHeader>
                  <div
                     class="flex justify-between items-center w-full pr-2">
                     <div class="flex items-center gap-3">
                        <i :class="group.icon" class="text-base text-primary-500"></i>
                        <span>{{ group.name }}</span>
                     </div>
                     <Checkbox
                        :modelValue="isGroupChecked(group)"
                        :indeterminate="isGroupIndeterminate(group)"
                        @update:modelValue="
                           toggleGroupPermissions(group, $event as boolean)
                        "
                        @click.stop />
                  </div>
               </AccordionHeader>
               <AccordionContent>
                  <div class="flex flex-col gap-3 py-2">
                     <div
                        v-for="perm in group.permissions"
                        :key="perm.id"
                        class="flex items-center gap-3">
                        <Checkbox
                           v-model="perm.checked"
                           :inputId="perm.id"
                           :binary="true" />
                        <label
                           :for="perm.id"
                           class="flex items-center gap-3 text-sm text-surface-700 dark:text-surface-300 cursor-pointer">
                           {{ perm.name }}
                           <code
                              class="py-0.5 px-1.5 bg-surface-100 dark:bg-surface-800 rounded text-xs font-mono text-surface-500 dark:text-surface-400">
                              {{ perm.id }}
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
               @click="dialogVisible = false" />
            <Button label="保存权限" @click="savePermissions" />
         </div>
      </template>
   </Dialog>
</template>
