import { defineStore } from 'pinia';
import { ref, watch } from 'vue';
import { STORAGE_KEYS, DARK_MODE_SELECTOR } from '@/config';

export const useThemeStore = defineStore('theme', () => {
   const isDark = ref(false);
   const darkClass = DARK_MODE_SELECTOR.replace('.', '');

   const applyTheme = () => {
      document.documentElement.classList.toggle(darkClass, isDark.value);
   };

   const initTheme = () => {
      const saved = localStorage.getItem(STORAGE_KEYS.THEME);
      isDark.value = saved ? saved === 'dark' : window.matchMedia('(prefers-color-scheme: dark)').matches;
      applyTheme();
   };

   const toggleTheme = () => {
      isDark.value = !isDark.value;
      localStorage.setItem(STORAGE_KEYS.THEME, isDark.value ? 'dark' : 'light');
      applyTheme();
   };

   watch(isDark, applyTheme);

   return { isDark, initTheme, toggleTheme };
});
