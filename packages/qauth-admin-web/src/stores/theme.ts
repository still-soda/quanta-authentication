import { defineStore } from 'pinia';
import { ref, watch } from 'vue';

export const useThemeStore = defineStore('theme', () => {
   const isDark = ref(false);

   const initTheme = () => {
      const savedTheme = localStorage.getItem('theme');
      if (savedTheme) {
         isDark.value = savedTheme === 'dark';
      } else {
         isDark.value = window.matchMedia(
            '(prefers-color-scheme: dark)',
         ).matches;
      }
      applyTheme();
   };

   const applyTheme = () => {
      if (isDark.value) {
         document.documentElement.classList.add('app-dark');
      } else {
         document.documentElement.classList.remove('app-dark');
      }
   };

   const toggleTheme = () => {
      isDark.value = !isDark.value;
      localStorage.setItem('theme', isDark.value ? 'dark' : 'light');
      applyTheme();
   };

   watch(isDark, () => {
      applyTheme();
   });

   return {
      isDark,
      initTheme,
      toggleTheme,
   };
});
