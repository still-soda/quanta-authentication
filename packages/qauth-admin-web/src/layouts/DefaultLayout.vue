<script setup lang="ts">
import { computed, onMounted } from 'vue';
import { useSidebarStore } from '@/stores/sidebar';
import { useThemeStore } from '@/stores/theme';
import AppSidebar from '@/components/AppSidebar.vue';
import AppTopbar from '@/components/AppTopbar.vue';
import Toast from 'primevue/toast';

const sidebarStore = useSidebarStore();
const themeStore = useThemeStore();

onMounted(() => {
   themeStore.initTheme();
});

const mainStyle = computed(() => ({
   marginLeft: sidebarStore.isCollapsed ? '5rem' : '17rem',
}));
</script>

<template>
   <div
      class="min-h-screen bg-surface-50 dark:bg-surface-950">
      <AppSidebar />
      <AppTopbar />

      <main
         class="pt-16 min-h-screen transition-[margin-left] duration-300 ease-[cubic-bezier(0.4,0,0.2,1)] lg:ml-[var(--sidebar-margin)]"
         :style="{ '--sidebar-margin': mainStyle.marginLeft }">
         <div class="p-6 max-w-[1600px] mx-auto lg:p-6 max-lg:p-4">
            <router-view v-slot="{ Component }">
               <Transition name="page" mode="out-in">
                  <component :is="Component" />
               </Transition>
            </router-view>
         </div>
      </main>

      <Toast position="top-right" />
   </div>
</template>

<style scoped>
/* Page transition - must remain as CSS for Vue transitions */
.page-enter-active,
.page-leave-active {
   transition: all 0.2s ease;
}

.page-enter-from {
   opacity: 0;
   transform: translateY(10px);
}

.page-leave-to {
   opacity: 0;
   transform: translateY(-10px);
}

/* Mobile override */
@media (max-width: 1024px) {
   main {
      margin-left: 0 !important;
   }
}
</style>
