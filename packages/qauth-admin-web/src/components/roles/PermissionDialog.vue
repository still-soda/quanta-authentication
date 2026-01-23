<script setup lang="ts">
import { ref, computed } from 'vue';
import Dialog from 'primevue/dialog';
import Button from 'primevue/button';
import Checkbox from 'primevue/checkbox';
import Accordion from 'primevue/accordion';
import AccordionPanel from 'primevue/accordionpanel';
import AccordionHeader from 'primevue/accordionheader';
import AccordionContent from 'primevue/accordioncontent';
import type { Role } from './RoleCard.vue';

export interface Permission {
   id: string;
   name: string;
   checked: boolean;
}

export interface PermissionGroup {
   name: string;
   icon: string;
   permissions: Permission[];
}

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
      :style="{ width: '36rem' }"
      class="permission-dialog">
      <div class="permission-content">
         <Accordion multiple>
            <AccordionPanel
               v-for="group in permissionGroups"
               :key="group.name"
               :value="group.name">
               <AccordionHeader>
                  <div class="permission-group-header">
                     <div class="group-info">
                        <i :class="group.icon"></i>
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
                  <div class="permission-list">
                     <div
                        v-for="perm in group.permissions"
                        :key="perm.id"
                        class="permission-item">
                        <Checkbox
                           v-model="perm.checked"
                           :inputId="perm.id"
                           :binary="true" />
                        <label :for="perm.id" class="permission-label">
                           {{ perm.name }}
                           <code class="permission-code">{{ perm.id }}</code>
                        </label>
                     </div>
                  </div>
               </AccordionContent>
            </AccordionPanel>
         </Accordion>
      </div>

      <template #footer>
         <div class="dialog-footer">
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

<style scoped>
.permission-content {
   max-height: 60vh;
   overflow-y: auto;
}

.permission-group-header {
   display: flex;
   justify-content: space-between;
   align-items: center;
   width: 100%;
   padding-right: 0.5rem;
}

.group-info {
   display: flex;
   align-items: center;
   gap: 0.75rem;
}

.group-info i {
   font-size: 1rem;
   color: var(--p-orange-500);
}

.permission-list {
   display: flex;
   flex-direction: column;
   gap: 0.75rem;
   padding: 0.5rem 0;
}

.permission-item {
   display: flex;
   align-items: center;
   gap: 0.75rem;
}

.permission-label {
   display: flex;
   align-items: center;
   gap: 0.75rem;
   font-size: 0.875rem;
   color: var(--p-surface-700);
   cursor: pointer;
}

:global(.app-dark) .permission-label {
   color: var(--p-surface-300);
}

.permission-code {
   padding: 0.125rem 0.375rem;
   background: var(--p-surface-100);
   border-radius: 4px;
   font-size: 0.75rem;
   font-family: monospace;
   color: var(--p-surface-500);
}

:global(.app-dark) .permission-code {
   background: var(--p-surface-800);
   color: var(--p-surface-400);
}

.dialog-footer {
   display: flex;
   justify-content: flex-end;
   gap: 0.75rem;
}
</style>
