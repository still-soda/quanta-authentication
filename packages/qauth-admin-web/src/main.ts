import { createApp } from 'vue';
import { createPinia } from 'pinia';
import './style.css';
import 'primeicons/primeicons.css';
import App from './App.vue';
import router from './router';
import PrimeVue from 'primevue/config';
import Aura from '@primeuix/themes/aura';
import { definePreset } from '@primeuix/themes';
import ToastService from 'primevue/toastservice';
import ConfirmationService from 'primevue/confirmationservice';
import Tooltip from 'primevue/tooltip';
import Ripple from 'primevue/ripple';
import { DARK_MODE_SELECTOR } from './config';

const AppPreset = definePreset(Aura, {
   primitive: {
      borderRadius: {
         none: '0',
         xs: '1px',
         sm: '2px',
         md: '4px',
         lg: '6px',
         xl: '8px',
      },
   },
   semantic: {
      primary: {
         50: '{orange.50}',
         100: '{orange.100}',
         200: '{orange.200}',
         300: '{orange.300}',
         400: '{orange.400}',
         500: '{orange.500}',
         600: '{orange.600}',
         700: '{orange.700}',
         800: '{orange.800}',
         900: '{orange.900}',
         950: '{orange.950}',
      },
      colorScheme: {
         light: {
            surface: {
               0: '#ffffff',
               50: '{slate.50}',
               100: '{slate.100}',
               200: '{slate.200}',
               300: '{slate.300}',
               400: '{slate.400}',
               500: '{slate.500}',
               600: '{slate.600}',
               700: '{slate.700}',
               800: '{slate.800}',
               900: '{slate.900}',
               950: '{slate.950}',
            },
            primary: {
               color: '{orange.500}',
               inverseColor: '#ffffff',
               hoverColor: '{orange.600}',
               activeColor: '{orange.700}',
            },
            highlight: {
               background: '{orange.50}',
               focusBackground: '{orange.100}',
               color: '{orange.700}',
               focusColor: '{orange.800}',
            },
         },
         dark: {
            surface: {
               0: '#ffffff',
               50: '{zinc.50}',
               100: '{zinc.100}',
               200: '{zinc.200}',
               300: '{zinc.300}',
               400: '{zinc.400}',
               500: '{zinc.500}',
               600: '{zinc.600}',
               700: '{zinc.700}',
               800: '{zinc.800}',
               900: '{zinc.900}',
               950: '{zinc.950}',
            },
            primary: {
               color: '{orange.400}',
               inverseColor: '{zinc.950}',
               hoverColor: '{orange.300}',
               activeColor: '{orange.200}',
            },
            highlight: {
               background: 'rgba(251, 146, 60, 0.16)',
               focusBackground: 'rgba(251, 146, 60, 0.24)',
               color: 'rgba(255,255,255,.87)',
               focusColor: 'rgba(255,255,255,.87)',
            },
         },
      },
   },
});

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.use(router);
app.use(PrimeVue, {
   ripple: true,
   theme: {
      preset: AppPreset,
      options: {
         darkModeSelector: DARK_MODE_SELECTOR,
         cssLayer: { name: 'primevue', order: 'theme, base, primevue' },
      },
   },
});
app.use(ToastService);
app.use(ConfirmationService);
app.directive('tooltip', Tooltip);
app.directive('ripple', Ripple);

app.mount('#app');
