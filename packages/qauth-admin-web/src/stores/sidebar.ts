import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useSidebarStore = defineStore('sidebar', () => {
   const isCollapsed = ref(false);
   const isMobileOpen = ref(false);

   const toggleCollapsed = () => {
      isCollapsed.value = !isCollapsed.value;
   };

   const toggleMobile = () => {
      isMobileOpen.value = !isMobileOpen.value;
   };

   const closeMobile = () => {
      isMobileOpen.value = false;
   };

   return {
      isCollapsed,
      isMobileOpen,
      toggleCollapsed,
      toggleMobile,
      closeMobile,
   };
});
