import { createApp } from 'vue';
import './style.css';
import 'primeicons/primeicons.css';
import App from './App.vue';
import PrimeVue from 'primevue/config';
import Aura from 'primevue/resources/themes/aura/theme.css';
import { definePreset } from '@primeuix/themes';

const getColorLevels = (color: string) => {
   const colorLevels: Record<number, string> = {};
   [50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 950].forEach((level) => {
      colorLevels[level] = `{${color}.${level}}`;
   });
   return colorLevels;
};

const AppPreset = definePreset(Aura, {
   semantic: {
      primary: getColorLevels('orange'),
      surface: {
         0: '#ffffff',
         ...getColorLevels('soho'),
      },
      highlight: {
         background: '{primary.100}',
         color: '{primary.700}',
      },
   },
});

createApp(App)
   .use(PrimeVue, {
      ripple: true,
      theme: {
         preset: AppPreset,
         options: {
            darkModeSelector: '.app-dark',
            cssLayer: {
               name: 'primevue',
               order: 'theme, base, primevue',
            },
         },
      },
   })
   .mount('#app');
