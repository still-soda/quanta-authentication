<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Button from 'primevue/button'
import { useToast } from 'primevue/usetoast'

const router = useRouter()
const toast = useToast()

// Form state
const email = ref('')
const reason = ref('')
const isLoading = ref(false)
const isSubmitted = ref(false)

// Email validation
const isEmailValid = computed(() => {
   const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
   return emailRegex.test(email.value)
})

const isFormValid = computed(() => {
   return isEmailValid.value && reason.value.trim().length >= 10
})

// Handle submit
const handleSubmit = async () => {
   if (!isFormValid.value) {
      toast.add({
         severity: 'warn',
         summary: '表单验证',
         detail: '请填写有效的邮箱地址和详细的申请原因（至少10个字符）',
         life: 3000,
      })
      return
   }

   isLoading.value = true

   try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1500))

      isSubmitted.value = true
      toast.add({
         severity: 'success',
         summary: '申请已提交',
         detail: '密码重置申请已发送，管理员将会审核您的申请',
         life: 4000,
      })
   } catch {
      toast.add({
         severity: 'error',
         summary: '提交失败',
         detail: '申请提交失败，请稍后重试',
         life: 3000,
      })
   } finally {
      isLoading.value = false
   }
}

// Navigation
const goToLogin = () => router.push('/auth/login')
const submitAnother = () => {
   isSubmitted.value = false
   email.value = ''
   reason.value = ''
}
</script>

<template>
   <div class="forgot-password-view">
      <!-- Success state -->
      <template v-if="isSubmitted">
         <div class="text-center">
            <!-- Success icon -->
            <div class="mb-6">
               <div
                  class="inline-flex items-center justify-center w-20 h-20 rounded-full bg-linear-to-br from-green-400/20 to-green-600/20 border border-green-500/30"
               >
                  <i class="pi pi-check-circle text-4xl text-green-400 animate-bounce-in"></i>
               </div>
            </div>

            <h2 class="text-xl font-semibold text-white mb-3">申请已提交</h2>
            <p class="text-sm text-white/50 mb-6 leading-relaxed">
               我们已收到您的密码重置申请。<br />
               管理员将在 <span class="text-primary-400">24小时内</span> 审核您的请求，<br />
               审核结果将通过邮件通知您。
            </p>

            <div class="p-4 rounded-xl bg-white/3 border border-white/10 mb-6">
               <p class="text-xs text-white/40 mb-1">申请邮箱</p>
               <p class="text-sm text-white font-medium">{{ email }}</p>
            </div>

            <div class="space-y-3">
               <Button
                  label="返回登录"
                  class="auth-button w-full cursor-pointer"
                  severity="contrast"
                  @click="goToLogin"
               />
               <button
                  type="button"
                  class="text-sm text-white/50 hover:text-white/70 transition-colors"
                  @click="submitAnother"
               >
                  提交另一个申请
               </button>
            </div>
         </div>
      </template>

      <!-- Form state -->
      <template v-else>
         <!-- Header -->
         <div class="mb-8">
            <button
               type="button"
               class="inline-flex cursor-pointer items-center gap-2 text-sm text-white/50 hover:text-white/80 transition-colors mb-4"
               @click="goToLogin"
            >
               <i class="pi pi-arrow-left text-xs"></i>
               返回登录
            </button>
            <h2 class="text-xl font-semibold text-white mb-2">忘记密码</h2>
            <p class="text-sm text-white/50 leading-relaxed">
               请填写以下信息，我们将审核您的身份后<br />
               帮助您重置密码。
            </p>
         </div>

         <!-- Form -->
         <form @submit.prevent="handleSubmit" class="space-y-5">
            <!-- Email field -->
            <div class="form-field">
               <label class="form-label">
                  <i class="pi pi-envelope mr-2 opacity-60"></i>
                  联系邮箱
                  <span class="text-red-400 ml-1">*</span>
               </label>
               <div class="input-wrapper">
                  <InputText
                     v-model="email"
                     placeholder="请输入您的注册邮箱"
                     class="auth-input"
                     :disabled="isLoading"
                     :invalid="email.length > 0 && !isEmailValid"
                     autocomplete="email"
                  />
               </div>
               <p v-if="email.length > 0 && !isEmailValid" class="text-xs text-red-400 mt-1">
                  请输入有效的邮箱地址
               </p>
            </div>

            <!-- Reason field -->
            <div class="form-field">
               <label class="form-label">
                  <i class="pi pi-file-edit mr-2 opacity-60"></i>
                  申请原因
                  <span class="text-red-400 ml-1">*</span>
               </label>
               <div class="input-wrapper">
                  <Textarea
                     v-model="reason"
                     placeholder="请简要说明忘记密码的原因，以及您的身份验证信息（如注册时间、常用功能等）"
                     class="auth-textarea"
                     :disabled="isLoading"
                     rows="4"
                     autoResize
                  />
               </div>
               <p class="text-xs text-white/30 mt-1">{{ reason.length }}/10 最少字符</p>
            </div>

            <!-- Info box -->
            <div class="p-4 rounded-xl bg-primary-500/10 border border-primary-500/20">
               <div class="flex gap-3">
                  <i class="pi pi-info-circle text-primary-400 mt-0.5 shrink-0"></i>
                  <div class="text-xs text-white/60 leading-relaxed">
                     <p class="font-medium text-white/80 mb-1">温馨提示</p>
                     <p>
                        为确保账户安全，密码重置需要管理员人工审核。请提供尽可能详细的身份验证信息以加快审核速度。
                     </p>
                  </div>
               </div>
            </div>

            <!-- Submit button -->
            <Button
               type="submit"
               :label="isLoading ? '提交中...' : '提交申请'"
               :loading="isLoading"
               :disabled="!isFormValid || isLoading"
               class="auth-button w-full"
               severity="contrast"
            />
         </form>
      </template>
   </div>
</template>

<style scoped>
.forgot-password-view {
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

:deep(.auth-input.p-invalid) {
   border-color: var(--p-red-500);
}

/* Textarea styling */
:deep(.auth-textarea) {
   width: 100%;
   padding: 0.875rem 1rem;
   background: rgba(255, 255, 255, 0.03);
   border: 1px solid rgba(255, 255, 255, 0.1);
   border-radius: 0.75rem;
   color: white;
   font-size: 0.9375rem;
   resize: none;
   transition: all 0.2s ease;
   font-family: inherit;
}

:deep(.auth-textarea::placeholder) {
   color: rgba(255, 255, 255, 0.3);
}

:deep(.auth-textarea:hover) {
   border-color: rgba(255, 255, 255, 0.2);
   background: rgba(255, 255, 255, 0.05);
}

:deep(.auth-textarea:focus) {
   outline: none;
   border-color: var(--p-orange-500);
   background: rgba(255, 255, 255, 0.05);
   box-shadow: 0 0 0 3px rgba(251, 146, 60, 0.15);
}

:deep(.auth-textarea:disabled) {
   opacity: 0.5;
   cursor: not-allowed;
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
   transform: translateY(-1px);
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

/* Success animation */
@keyframes bounce-in {
   0% {
      transform: scale(0);
      opacity: 0;
   }
   50% {
      transform: scale(1.2);
   }
   100% {
      transform: scale(1);
      opacity: 1;
   }
}

.animate-bounce-in {
   animation: bounce-in 0.5s ease-out forwards;
}
</style>
