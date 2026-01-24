<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import InputText from 'primevue/inputtext'
import Checkbox from 'primevue/checkbox'
import Button from 'primevue/button'
import { useToast } from 'primevue/usetoast'
import { login } from '@/apis/auth'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const toast = useToast()
const authStore = useAuthStore()

// Form state
const studentId = ref('')
const password = ref('')
const rememberMe = ref(false)
const isLoading = ref(false)
const showPassword = ref(false)

// Validation
const isFormValid = computed(() => {
   return studentId.value.trim().length > 0 && password.value.length >= 6
})

// Handle login
const handleLogin = async () => {
   if (!isFormValid.value) {
      toast.add({
         severity: 'warn',
         summary: '表单验证',
         detail: '请填写正确的学号和密码',
         life: 3000,
      })
      return
   }

   isLoading.value = true

   try {
      // 调用登录 API
      const response = await login({
         student_id: studentId.value.trim(),
         password: password.value,
      })

      // 保存认证信息到 store
      authStore.setAuth(
         {
            accessToken: response.access_token,
            refreshToken: response.refresh_token,
         },
         response.user
      )

      toast.add({
         severity: 'success',
         summary: '登录成功',
         detail: `欢迎回来，${response.user.name || response.user.student_id}！`,
         life: 2000,
      })

      setTimeout(() => {
         router.replace('/')
      }, 500)
   } catch (error: unknown) {
      const err = error as { response?: { data?: { message?: string } }; message?: string }
      const message = err.response?.data?.message || err.message || '用户名或密码错误，请重试'

      toast.add({
         severity: 'error',
         summary: '登录失败',
         detail: message,
         life: 3000,
      })
      isLoading.value = false
   }
}

// Navigation
const goToRegister = () => router.push('/auth/register')
const goToForgotPassword = () => router.push('/auth/forgot-password')
</script>

<template>
   <div class="login-view">
      <!-- Header -->
      <div class="mb-8">
         <h2 class="text-xl font-semibold text-white mb-2">欢迎回来</h2>
         <p class="text-sm text-white/50">请登录您的账户以继续</p>
      </div>

      <!-- Login Form -->
      <form @submit.prevent="handleLogin" class="space-y-5">
         <!-- Student ID field -->
         <div class="form-field">
            <label class="form-label">
               <i class="pi pi-id-card mr-2 opacity-60"></i>
               学号
            </label>
            <div class="input-wrapper">
               <InputText
                  v-model="studentId"
                  placeholder="请输入学号"
                  class="auth-input"
                  :disabled="isLoading"
                  autocomplete="username"
               />
            </div>
         </div>

         <!-- Password field -->
         <div class="form-field">
            <label class="form-label">
               <i class="pi pi-lock mr-2 opacity-60"></i>
               密码
            </label>
            <div class="input-wrapper relative">
               <InputText
                  v-if="!showPassword"
                  v-model="password"
                  type="password"
                  placeholder="请输入密码"
                  class="auth-input pr-12"
                  :disabled="isLoading"
                  autocomplete="current-password"
               />
               <InputText
                  v-else
                  v-model="password"
                  type="text"
                  placeholder="请输入密码"
                  class="auth-input pr-12"
                  :disabled="isLoading"
                  autocomplete="current-password"
               />
               <button
                  type="button"
                  class="cursor-pointer absolute right-4 top-1/2 -translate-y-1/2 text-white/40 hover:text-white/70 transition-colors"
                  @click="showPassword = !showPassword"
               >
                  <i :class="showPassword ? 'pi pi-eye-slash' : 'pi pi-eye'"></i>
               </button>
            </div>
         </div>

         <!-- Remember me & Forgot password -->
         <div class="flex items-center justify-between">
            <label class="flex items-center gap-2 cursor-pointer select-none group">
               <Checkbox v-model="rememberMe" :binary="true" class="auth-checkbox" />
               <span class="text-sm text-white/50 group-hover:text-white/70 transition-colors">
                  记住我
               </span>
            </label>
            <button
               type="button"
               class="cursor-pointer text-sm text-primary-400 hover:text-primary-300 transition-colors font-medium"
               @click="goToForgotPassword"
            >
               忘记密码？
            </button>
         </div>

         <!-- Login button -->
         <Button
            type="submit"
            :label="isLoading ? '登录中...' : '登 录'"
            :loading="isLoading"
            :disabled="!isFormValid || isLoading"
            class="auth-button w-full"
            severity="contrast"
         />
      </form>

      <!-- Register link -->
      <div class="mt-8 text-center">
         <p class="text-sm text-white/40">
            还没有账户？
            <button
               type="button"
               class="cursor-pointer text-primary-400 hover:text-primary-300 transition-colors font-medium ml-1"
               @click="goToRegister"
            >
               申请注册
            </button>
         </p>
      </div>
   </div>
</template>

<style scoped>
.login-view {
   position: relative;
}

.form-field {
   display: flex;
   flex-direction: column;
   gap: 0.5rem;
}

.form-label {
   display: flex;
   align-items: center;
   font-size: 0.8125rem;
   font-weight: 500;
   color: rgba(255, 255, 255, 0.7);
}

/* Input styling */
:deep(.auth-input) {
   width: 100%;
   padding: 0.875rem 1rem;
   background: rgba(255, 255, 255, 0.03);
   border: 1px solid rgba(255, 255, 255, 0.1);
   border-radius: 0.75rem;
   color: white;
   font-size: 0.9375rem;
   transition: all 0.2s ease;
}

:deep(.auth-input::placeholder) {
   color: rgba(255, 255, 255, 0.3);
}

:deep(.auth-input:hover) {
   border-color: rgba(255, 255, 255, 0.2);
   background: rgba(255, 255, 255, 0.05);
}

:deep(.auth-input:focus) {
   outline: none;
   border-color: var(--p-orange-500);
   background: rgba(255, 255, 255, 0.05);
   box-shadow: 0 0 0 3px rgba(251, 146, 60, 0.15);
}

:deep(.auth-input:disabled) {
   opacity: 0.5;
   cursor: not-allowed;
}

/* Checkbox styling */
:deep(.auth-checkbox) {
   width: 1.125rem;
   height: 1.125rem;
}

:deep(.auth-checkbox .p-checkbox-box) {
   background: rgba(255, 255, 255, 0.03);
   border: 1px solid rgba(255, 255, 255, 0.2);
   border-radius: 0.375rem;
}

:deep(.auth-checkbox.p-checkbox-checked .p-checkbox-box) {
   background: var(--p-orange-500);
   border-color: var(--p-orange-500);
}

/* Button styling */
:deep(.auth-button) {
   padding: 0.875rem 1.5rem;
   font-size: 0.9375rem;
   font-weight: 600;
   border-radius: 0.75rem;
   background: linear-gradient(135deg, var(--p-orange-500) 0%, var(--p-orange-600) 100%);
   border: none;
   color: white;
   box-shadow: 0 4px 16px -4px rgba(251, 146, 60, 0.4);
   transition: all 0.2s ease;
}

:deep(.auth-button:hover:not(:disabled)) {
   box-shadow: 0 6px 20px -4px rgba(251, 146, 60, 0.5);
}

:deep(.auth-button:active:not(:disabled)) {
   transform: translateY(0);
}

:deep(.auth-button:disabled) {
   opacity: 0.6;
   cursor: not-allowed;
   transform: none;
}

/* OAuth button */
.oauth-button {
   display: flex;
   align-items: center;
   justify-content: center;
   gap: 0.75rem;
   padding: 0.75rem 1rem;
   background: rgba(255, 255, 255, 0.03);
   border: 1px solid rgba(255, 255, 255, 0.1);
   border-radius: 0.75rem;
   color: rgba(255, 255, 255, 0.8);
   font-size: 0.875rem;
   font-weight: 500;
   cursor: pointer;
   transition: all 0.2s ease;
}

.oauth-button:hover:not(:disabled) {
   background: rgba(255, 255, 255, 0.06);
   border-color: rgba(255, 255, 255, 0.2);
   color: white;
}

.oauth-button:disabled {
   opacity: 0.5;
   cursor: not-allowed;
}
</style>
