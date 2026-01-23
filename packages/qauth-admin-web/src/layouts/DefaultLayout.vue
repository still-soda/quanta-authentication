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
   <div class="app-layout">
      <AppSidebar />
      <AppTopbar />

      <main class="main-content" :style="mainStyle">
         <div class="content-wrapper">
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
.app-layout {
   min-height: 100vh;
   background: var(--p-surface-50);
}

:global(.app-dark) .app-layout {
   background: var(--p-surface-950);
}

.main-content {
   padding-top: 4rem;
   min-height: 100vh;
   transition: margin-left 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.content-wrapper {
   padding: 1.5rem;
   max-width: 1600px;
   margin: 0 auto;
}

/* Page transition */
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

/* Mobile */
@media (max-width: 1024px) {
   .main-content {
      margin-left: 0 !important;
   }

   .content-wrapper {
      padding: 1rem;
   }
}
</style>
