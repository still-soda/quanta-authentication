<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useThemeStore } from '@/stores/theme'
import { APP_NAME } from '@/config'

const themeStore = useThemeStore()

onMounted(() => {
   themeStore.initTheme()
})

// 动态粒子效果
const particles = ref(
   Array.from({ length: 50 }, (_, i) => ({
      id: i,
      x: Math.random() * 100,
      y: Math.random() * 100,
      size: Math.random() * 4 + 1,
      duration: Math.random() * 20 + 10,
      delay: Math.random() * 5,
   }))
)
</script>

<template>
   <div
      class="auth-layout min-h-screen flex items-center justify-center relative overflow-hidden bg-[#0a0a0f] dark:bg-[#0a0a0f]"
   >
      <!-- Animated background -->
      <div class="absolute inset-0 overflow-hidden">
         <!-- Gradient orbs -->
         <div
            class="absolute -top-40 -left-40 w-96 h-96 bg-linear-to-br from-primary-500/30 via-primary-600/20 to-transparent rounded-full blur-3xl animate-orb-1"
         />
         <div
            class="absolute top-1/2 -right-32 w-80 h-80 bg-linear-to-br from-purple-500/25 via-purple-600/15 to-transparent rounded-full blur-3xl animate-orb-2"
         />
         <div
            class="absolute -bottom-20 left-1/3 w-72 h-72 bg-linear-to-br from-cyan-500/20 via-cyan-600/10 to-transparent rounded-full blur-3xl animate-orb-3"
         />

         <!-- Grid pattern -->
         <div
            class="absolute inset-0 opacity-[0.03]"
            style="
               background-image:
                  linear-gradient(rgba(255, 255, 255, 0.1) 1px, transparent 1px),
                  linear-gradient(90deg, rgba(255, 255, 255, 0.1) 1px, transparent 1px);
               background-size: 60px 60px;
            "
         />

         <!-- Floating particles -->
         <div
            v-for="particle in particles"
            :key="particle.id"
            class="absolute rounded-full bg-white/10"
            :style="{
               left: `${particle.x}%`,
               top: `${particle.y}%`,
               width: `${particle.size}px`,
               height: `${particle.size}px`,
               animation: `float ${particle.duration}s ease-in-out infinite`,
               animationDelay: `${particle.delay}s`,
            }"
         />

         <!-- Noise texture overlay -->
         <div class="absolute inset-0 opacity-[0.015] noise-overlay" />
      </div>

      <!-- Main content -->
      <div class="relative z-10 w-full max-w-md mx-auto px-4">
         <!-- Logo and branding -->
         <div class="text-center mb-8">
            <div class="relative w-full flex justify-center h-28">
               <img src="/logo.png" alt="Logo" class="w-32 absolute top-0" />
            </div>
            <h1 class="text-2xl font-bold tracking-tight text-white mb-1 font-display">
               {{ APP_NAME }}
            </h1>
            <p class="text-sm text-white/50">统一身份认证平台</p>
         </div>

         <!-- Card container -->
         <div
            class="auth-card relative backdrop-blur-xl bg-white/3 border border-white/8 rounded-3xl p-8 shadow-2xl shadow-black/20"
         >
            <!-- Card glow effect -->
            <div
               class="absolute -inset-px bg-linear-to-b from-white/10 via-transparent to-transparent rounded-3xl pointer-events-none"
            />
            <router-view v-slot="{ Component }">
               <Transition name="auth-page" mode="out-in">
                  <component :is="Component" />
               </Transition>
            </router-view>
         </div>

         <!-- Footer -->
         <div class="mt-8 text-center">
            <p class="text-xs text-white/30">
               © {{ new Date().getFullYear() }} {{ APP_NAME }}. 保留所有权利。
            </p>
         </div>
      </div>
   </div>
</template>

<style>
/* Floating particle animation */
@keyframes float {
   0%,
   100% {
      transform: translateY(0) translateX(0);
      opacity: 0.3;
   }
   25% {
      transform: translateY(-20px) translateX(10px);
      opacity: 0.6;
   }
   50% {
      transform: translateY(-10px) translateX(-5px);
      opacity: 0.4;
   }
   75% {
      transform: translateY(-30px) translateX(15px);
      opacity: 0.5;
   }
}
</style>

<style scoped>
.auth-layout {
   font-family:
      'Inter',
      -apple-system,
      BlinkMacSystemFont,
      'Segoe UI',
      Roboto,
      sans-serif;
}

.font-display {
   font-family:
      'Inter',
      -apple-system,
      BlinkMacSystemFont,
      sans-serif;
   letter-spacing: -0.025em;
}

/* Orb animations */
@keyframes orb-1 {
   0%,
   100% {
      transform: translate(0, 0) scale(1);
   }
   33% {
      transform: translate(30px, 50px) scale(1.1);
   }
   66% {
      transform: translate(-20px, 30px) scale(0.95);
   }
}

@keyframes orb-2 {
   0%,
   100% {
      transform: translate(0, 0) scale(1);
   }
   33% {
      transform: translate(-40px, -30px) scale(0.9);
   }
   66% {
      transform: translate(20px, -50px) scale(1.05);
   }
}

@keyframes orb-3 {
   0%,
   100% {
      transform: translate(0, 0) scale(1);
   }
   33% {
      transform: translate(50px, -20px) scale(1.15);
   }
   66% {
      transform: translate(-30px, 40px) scale(0.9);
   }
}

.animate-orb-1 {
   animation: orb-1 25s ease-in-out infinite;
}

.animate-orb-2 {
   animation: orb-2 30s ease-in-out infinite;
}

.animate-orb-3 {
   animation: orb-3 20s ease-in-out infinite;
}

/* Noise overlay */
.noise-overlay {
   background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)'/%3E%3C/svg%3E");
}

/* Auth card subtle animation */
.auth-card {
   transition:
      transform 0.3s ease,
      box-shadow 0.3s ease;
}

/* Page transitions */
.auth-page-enter-active,
.auth-page-leave-active {
   transition: all 0.25s ease;
}

.auth-page-enter-from {
   opacity: 0;
   transform: translateX(20px);
}

.auth-page-leave-to {
   opacity: 0;
   transform: translateX(-20px);
}
</style>
