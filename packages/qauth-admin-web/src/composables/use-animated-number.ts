import { ref, watch, onMounted, onUnmounted, type Ref, computed } from 'vue';
import type { AnimatedNumberOptions } from '@/types';
import { easings, ANIMATION_DEFAULTS } from '@/config';

export { easings };

export function useAnimatedNumber(
   target: Ref<number> | number,
   options: AnimatedNumberOptions = {},
) {
   const {
      duration = ANIMATION_DEFAULTS.DURATION,
      easing = easings.easeOutExpo,
      decimals = ANIMATION_DEFAULTS.DECIMALS,
      formatter,
      autoStart = ANIMATION_DEFAULTS.AUTO_START,
      delay = ANIMATION_DEFAULTS.DELAY,
      padStart = ANIMATION_DEFAULTS.PAD_START,
      useGrouping = ANIMATION_DEFAULTS.USE_GROUPING,
   } = options;

   const displayValue = ref(0);
   const isAnimating = ref(false);
   const targetDigits = ref(1);

   let animationId: number | null = null;
   let startTime: number | null = null;
   let startValue = 0;
   let endValue = 0;
   let delayTimeoutId: ReturnType<typeof setTimeout> | null = null;

   const getTargetValue = () => (typeof target === 'number' ? target : target.value);

   const getDigitCount = (num: number) => {
      const absInt = Math.floor(Math.abs(num));
      return absInt === 0 ? 1 : Math.floor(Math.log10(absInt)) + 1;
   };

   const stop = () => {
      if (animationId !== null) {
         cancelAnimationFrame(animationId);
         animationId = null;
      }
      if (delayTimeoutId !== null) {
         clearTimeout(delayTimeoutId);
         delayTimeoutId = null;
      }
      startTime = null;
      isAnimating.value = false;
   };

   const animate = (currentTime: number) => {
      if (startTime === null) startTime = currentTime;

      const elapsed = currentTime - startTime;
      const progress = Math.min(elapsed / duration, 1);
      displayValue.value = startValue + (endValue - startValue) * easing(progress);

      if (progress < 1) {
         animationId = requestAnimationFrame(animate);
      } else {
         displayValue.value = endValue;
         isAnimating.value = false;
         animationId = null;
         startTime = null;
      }
   };

   const start = (newTarget?: number) => {
      stop();

      startValue = displayValue.value;
      endValue = newTarget !== undefined ? newTarget : getTargetValue();
      targetDigits.value = getDigitCount(endValue);

      if (startValue === endValue) return;

      isAnimating.value = true;
      if (delay > 0) {
         delayTimeoutId = setTimeout(() => {
            animationId = requestAnimationFrame(animate);
         }, delay);
      } else {
         animationId = requestAnimationFrame(animate);
      }
   };

   const set = (value: number) => {
      stop();
      displayValue.value = value;
   };

   const reset = () => {
      stop();
      displayValue.value = 0;
   };

   const restart = () => {
      reset();
      start();
   };

   const formattedValue = computed(() => {
      const value = Number(displayValue.value.toFixed(decimals));
      if (formatter) return formatter(value);

      if (padStart && targetDigits.value > 1) {
         const intPart = Math.floor(Math.abs(value));
         const decimalPart = decimals > 0 ? (value % 1).toFixed(decimals).slice(1) : '';
         const paddedInt = String(intPart).padStart(targetDigits.value, '0');
         const sign = value < 0 ? '-' : '';
         if (useGrouping) {
            const formattedInt = paddedInt.replace(/\B(?=(\d{3})+(?!\d))/g, ',');
            return `${sign}${formattedInt}${decimalPart}`;
         }
         return `${sign}${paddedInt}${decimalPart}`;
      }

      return value.toLocaleString('zh-CN', { minimumFractionDigits: decimals, maximumFractionDigits: decimals, useGrouping });
   });

   if (typeof target !== 'number') {
      watch(target, (newValue) => start(newValue), { immediate: false });
   }

   onMounted(() => { if (autoStart) start(); });
   onUnmounted(() => stop());

   return { displayValue, formattedValue, isAnimating, start, stop, set, reset, restart };
}
