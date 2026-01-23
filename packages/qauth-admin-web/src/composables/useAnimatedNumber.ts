import { ref, watch, onMounted, onUnmounted, type Ref, computed } from 'vue';

export interface AnimatedNumberOptions {
   /** 动画持续时间（毫秒），默认 800ms */
   duration?: number;
   /** 缓动函数，默认 easeOutExpo */
   easing?: (t: number) => number;
   /** 小数位数，默认 0 */
   decimals?: number;
   /** 数字格式化函数 */
   formatter?: (value: number) => string;
   /** 是否在挂载时自动开始动画，默认 true */
   autoStart?: boolean;
   /** 延迟开始时间（毫秒），默认 0 */
   delay?: number;
   /** 是否根据目标值的位数进行前导零填充，默认 false */
   padStart?: boolean;
   /** 是否使用千分位分隔符，默认 true */
   useGrouping?: boolean;
}

// 内置缓动函数
export const easings = {
   /** 线性 */
   linear: (t: number) => t,
   /** 缓出指数 - 快速开始，平滑结束 */
   easeOutExpo: (t: number) => (t === 1 ? 1 : 1 - Math.pow(2, -10 * t)),
   /** 缓出立方 */
   easeOutCubic: (t: number) => 1 - Math.pow(1 - t, 3),
   /** 缓出四次方 */
   easeOutQuart: (t: number) => 1 - Math.pow(1 - t, 4),
   /** 弹性效果 */
   easeOutElastic: (t: number) => {
      const c4 = (2 * Math.PI) / 3;
      return t === 0
         ? 0
         : t === 1
           ? 1
           : Math.pow(2, -10 * t) * Math.sin((t * 10 - 0.75) * c4) + 1;
   },
   /** 回弹效果 */
   easeOutBack: (t: number) => {
      const c1 = 1.70158;
      const c3 = c1 + 1;
      return 1 + c3 * Math.pow(t - 1, 3) + c1 * Math.pow(t - 1, 2);
   },
};

/**
 * 数字滚动动画组合式函数
 *
 * @example
 * ```vue
 * <script setup>
 * import { useAnimatedNumber } from '@/composables/useAnimatedNumber';
 *
 * const targetValue = ref(1234);
 * const { displayValue, formattedValue } = useAnimatedNumber(targetValue, {
 *   duration: 1000,
 *   decimals: 0,
 * });
 * </script>
 *
 * <template>
 *   <span>{{ formattedValue }}</span>
 * </template>
 * ```
 */
export function useAnimatedNumber(
   target: Ref<number> | number,
   options: AnimatedNumberOptions = {},
) {
   const {
      duration = 800,
      easing = easings.easeOutExpo,
      decimals = 0,
      formatter,
      autoStart = true,
      delay = 0,
      padStart = false,
      useGrouping = true,
   } = options;

   // 当前显示的数值
   const displayValue = ref(0);
   // 动画是否正在进行
   const isAnimating = ref(false);
   // 目标值的整数部分位数（用于 padStart）
   const targetDigits = ref(1);

   // 动画状态
   let animationId: number | null = null;
   let startTime: number | null = null;
   let startValue = 0;
   let endValue = 0;
   let delayTimeoutId: ReturnType<typeof setTimeout> | null = null;

   // 获取目标值
   const getTargetValue = () => {
      return typeof target === 'number' ? target : target.value;
   };

   // 计算数字的整数部分位数
   const getDigitCount = (num: number) => {
      const absInt = Math.floor(Math.abs(num));
      if (absInt === 0) return 1;
      return Math.floor(Math.log10(absInt)) + 1;
   };

   // 停止动画
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

   // 动画帧函数
   const animate = (currentTime: number) => {
      if (startTime === null) {
         startTime = currentTime;
      }

      const elapsed = currentTime - startTime;
      const progress = Math.min(elapsed / duration, 1);
      const easedProgress = easing(progress);

      // 计算当前值
      displayValue.value = startValue + (endValue - startValue) * easedProgress;

      if (progress < 1) {
         animationId = requestAnimationFrame(animate);
      } else {
         // 确保最终值精确
         displayValue.value = endValue;
         isAnimating.value = false;
         animationId = null;
         startTime = null;
      }
   };

   // 开始动画
   const start = (newTarget?: number) => {
      stop();

      startValue = displayValue.value;
      endValue = newTarget !== undefined ? newTarget : getTargetValue();

      // 更新目标位数
      targetDigits.value = getDigitCount(endValue);

      // 如果起始值和结束值相同，直接设置
      if (startValue === endValue) {
         return;
      }

      isAnimating.value = true;

      if (delay > 0) {
         delayTimeoutId = setTimeout(() => {
            animationId = requestAnimationFrame(animate);
         }, delay);
      } else {
         animationId = requestAnimationFrame(animate);
      }
   };

   // 立即设置值（无动画）
   const set = (value: number) => {
      stop();
      displayValue.value = value;
   };

   // 重置到 0
   const reset = () => {
      stop();
      displayValue.value = 0;
   };

   // 从 0 重新开始动画
   const restart = () => {
      reset();
      start();
   };

   // 格式化后的值
   const formattedValue = computed(() => {
      const value = Number(displayValue.value.toFixed(decimals));
      if (formatter) {
         return formatter(value);
      }

      // 如果需要 padStart，先生成不带分组的数字字符串
      if (padStart && targetDigits.value > 1) {
         const intPart = Math.floor(Math.abs(value));
         const decimalPart =
            decimals > 0 ? (value % 1).toFixed(decimals).slice(1) : '';
         const paddedInt = String(intPart).padStart(targetDigits.value, '0');
         const sign = value < 0 ? '-' : '';

         if (useGrouping) {
            // 对填充后的整数部分添加千分位分隔符
            const formattedInt = paddedInt.replace(
               /\B(?=(\d{3})+(?!\d))/g,
               ',',
            );
            return `${sign}${formattedInt}${decimalPart}`;
         }
         return `${sign}${paddedInt}${decimalPart}`;
      }

      // 默认使用千分位分隔符
      return value.toLocaleString('zh-CN', {
         minimumFractionDigits: decimals,
         maximumFractionDigits: decimals,
         useGrouping,
      });
   });

   // 监听目标值变化
   if (typeof target !== 'number') {
      watch(
         target,
         (newValue) => {
            start(newValue);
         },
         { immediate: false },
      );
   }

   // 挂载时自动开始
   onMounted(() => {
      if (autoStart) {
         start();
      }
   });

   // 卸载时清理
   onUnmounted(() => {
      stop();
   });

   return {
      /** 当前显示的数值（原始数值） */
      displayValue,
      /** 格式化后的显示值 */
      formattedValue,
      /** 是否正在动画中 */
      isAnimating,
      /** 开始动画 */
      start,
      /** 停止动画 */
      stop,
      /** 立即设置值 */
      set,
      /** 重置到 0 */
      reset,
      /** 从 0 重新开始动画 */
      restart,
   };
}
