<script setup lang="ts">
import { ref, computed, reactive } from 'vue'
import { useRouter } from 'vue-router'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import { useToast } from 'primevue/usetoast'
import { submitRegister } from '@/apis/auth'

const router = useRouter()
const toast = useToast()

// Form state - 适配服务端注册接口
const formData = reactive({
   studentId: '',
   name: '',
   email: '',
   password: '',
   confirmPassword: '',
})

const isLoading = ref(false)
const isSubmitted = ref(false)
const submittedData = reactive({
   studentId: '',
   name: '',
   email: '',
})

// Validation
const isStudentIdValid = computed(() => {
   // 学号不为空且长度合理
   return formData.studentId.trim().length >= 5 && formData.studentId.trim().length <= 20
})

const isEmailValid = computed(() => {
   const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
   return emailRegex.test(formData.email)
})

const isNameValid = computed(() => {
   return formData.name.trim().length >= 2
})

const isPasswordValid = computed(() => {
   // 密码至少6位
   return formData.password.length >= 6
})

const isConfirmPasswordValid = computed(() => {
   return formData.confirmPassword === formData.password && formData.confirmPassword.length > 0
})

const isFormValid = computed(() => {
   return (
      isStudentIdValid.value &&
      isEmailValid.value &&
      isNameValid.value &&
      isPasswordValid.value &&
      isConfirmPasswordValid.value
   )
})

// Handle submit
const handleSubmit = async () => {
   if (!isFormValid.value) {
      toast.add({
         severity: 'warn',
         summary: '表单验证',
         detail: '请完整填写所有必填信息',
         life: 3000,
      })
      return
   }

   isLoading.value = true

   try {
      // 调用注册 API
      const response = await submitRegister({
         student_id: formData.studentId.trim(),
         password: formData.password,
         email: formData.email.trim(),
         name: formData.name.trim(),
      })

      if (response.success) {
         // 保存提交的数据用于显示
         submittedData.studentId = formData.studentId
         submittedData.name = formData.name
         submittedData.email = formData.email

         isSubmitted.value = true
         toast.add({
            severity: 'success',
            summary: '注册成功',
            detail: response.message || '您的账户已创建成功，请前往登录',
            life: 4000,
         })
      } else {
         throw new Error(response.message || '注册失败')
      }
   } catch (error: unknown) {
      const err = error as { message?: string }
      toast.add({
         severity: 'error',
         summary: '注册失败',
         detail: err.message || '注册失败，请稍后重试',
         life: 3000,
      })
   } finally {
      isLoading.value = false
   }
}

// Navigation
const goToLogin = () => router.push('/auth/login')
</script>

<template>
   <div class="register-view">
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

            <h2 class="text-xl font-semibold text-white mb-3">注册成功</h2>
            <p class="text-sm text-white/50 mb-6 leading-relaxed">
               您的账户已成功创建！<br />
               请使用学号和密码登录系统。
            </p>

            <!-- Registration summary -->
            <div class="p-4 rounded-xl bg-white/3 border border-white/10 mb-6 text-left">
               <p class="text-xs text-white/40 mb-3 uppercase tracking-wider">账户信息</p>
               <div class="space-y-2">
                  <div class="flex justify-between items-center">
                     <span class="text-sm text-white/50">学号</span>
                     <span class="text-sm text-white font-medium">{{
                        submittedData.studentId
                     }}</span>
                  </div>
                  <div class="flex justify-between items-center">
                     <span class="text-sm text-white/50">姓名</span>
                     <span class="text-sm text-white font-medium">{{ submittedData.name }}</span>
                  </div>
                  <div class="flex justify-between items-center">
                     <span class="text-sm text-white/50">邮箱</span>
                     <span class="text-sm text-white font-medium">{{ submittedData.email }}</span>
                  </div>
               </div>
            </div>

            <Button
               label="前往登录"
               class="auth-button w-full cursor-pointer"
               severity="contrast"
               @click="goToLogin"
            />
         </div>
      </template>

      <!-- Form state -->
      <template v-else>
         <!-- Header -->
         <div class="mb-6">
            <button
               type="button"
               class="inline-flex cursor-pointer items-center gap-2 text-sm text-white/50 hover:text-white/80 transition-colors mb-4"
               @click="goToLogin"
            >
               <i class="pi pi-arrow-left text-xs"></i>
               返回登录
            </button>
            <h2 class="text-xl font-semibold text-white mb-2">注册账户</h2>
            <p class="text-sm text-white/50">填写以下信息创建您的账户</p>
         </div>

         <!-- Form -->
         <form @submit.prevent="handleSubmit" class="space-y-5">
            <!-- Student ID field -->
            <div class="form-field">
               <label class="form-label">
                  <i class="pi pi-id-card mr-2 opacity-60"></i>
                  学号
                  <span class="text-red-400 ml-1">*</span>
               </label>
               <div class="input-wrapper">
                  <InputText
                     v-model="formData.studentId"
                     placeholder="请输入您的学号"
                     class="auth-input"
                     :disabled="isLoading"
                     :invalid="formData.studentId.length > 0 && !isStudentIdValid"
                     autocomplete="username"
                  />
               </div>
               <p
                  v-if="formData.studentId.length > 0 && !isStudentIdValid"
                  class="text-xs text-red-400 mt-1"
               >
                  学号需为5-20位
               </p>
            </div>

            <!-- Name field -->
            <div class="form-field">
               <label class="form-label">
                  <i class="pi pi-user mr-2 opacity-60"></i>
                  姓名
                  <span class="text-red-400 ml-1">*</span>
               </label>
               <div class="input-wrapper">
                  <InputText
                     v-model="formData.name"
                     placeholder="请输入您的真实姓名"
                     class="auth-input"
                     :disabled="isLoading"
                     autocomplete="name"
                  />
               </div>
            </div>

            <!-- Email field -->
            <div class="form-field">
               <label class="form-label">
                  <i class="pi pi-envelope mr-2 opacity-60"></i>
                  邮箱
                  <span class="text-red-400 ml-1">*</span>
               </label>
               <div class="input-wrapper">
                  <InputText
                     v-model="formData.email"
                     placeholder="用于接收系统通知"
                     class="auth-input"
                     :disabled="isLoading"
                     :invalid="formData.email.length > 0 && !isEmailValid"
                     autocomplete="email"
                  />
               </div>
               <p
                  v-if="formData.email.length > 0 && !isEmailValid"
                  class="text-xs text-red-400 mt-1"
               >
                  请输入有效的邮箱地址
               </p>
            </div>

            <!-- Password field -->
            <div class="form-field">
               <label class="form-label">
                  <i class="pi pi-lock mr-2 opacity-60"></i>
                  密码
                  <span class="text-red-400 ml-1">*</span>
               </label>
               <div class="input-wrapper">
                  <InputText
                     v-model="formData.password"
                     type="password"
                     placeholder="请设置密码（至少6位）"
                     class="auth-input"
                     :disabled="isLoading"
                     :invalid="formData.password.length > 0 && !isPasswordValid"
                     autocomplete="new-password"
                  />
               </div>
               <p
                  v-if="formData.password.length > 0 && !isPasswordValid"
                  class="text-xs text-red-400 mt-1"
               >
                  密码至少6位
               </p>
            </div>

            <!-- Confirm Password field -->
            <div class="form-field">
               <label class="form-label">
                  <i class="pi pi-lock mr-2 opacity-60"></i>
                  确认密码
                  <span class="text-red-400 ml-1">*</span>
               </label>
               <div class="input-wrapper">
                  <InputText
                     v-model="formData.confirmPassword"
                     type="password"
                     placeholder="请再次输入密码"
                     class="auth-input"
                     :disabled="isLoading"
                     :invalid="formData.confirmPassword.length > 0 && !isConfirmPasswordValid"
                     autocomplete="new-password"
                  />
               </div>
               <p
                  v-if="formData.confirmPassword.length > 0 && !isConfirmPasswordValid"
                  class="text-xs text-red-400 mt-1"
               >
                  两次输入的密码不一致
               </p>
            </div>

            <!-- Info box -->
            <div class="p-4 rounded-xl bg-primary-500/10 border border-primary-500/20">
               <div class="flex gap-3">
                  <i class="pi pi-info-circle text-primary-400 mt-0.5 shrink-0"></i>
                  <div class="text-xs text-white/60 leading-relaxed">
                     <p class="font-medium text-white/80 mb-1">注册说明</p>
                     <p>
                        请确保您提供的学号和邮箱真实有效。注册成功后，您可以立即使用学号和密码登录系统。
                     </p>
                  </div>
               </div>
            </div>

            <!-- Submit button -->
            <Button
               type="submit"
               :label="isLoading ? '注册中...' : '立即注册'"
               :loading="isLoading"
               :disabled="!isFormValid || isLoading"
               class="auth-button w-full"
               severity="contrast"
            />
         </form>

         <!-- Login link -->
         <div class="mt-6 text-center">
            <p class="text-sm text-white/40">
               已有账户？
               <button
                  type="button"
                  class="text-primary-400 cursor-pointer hover:text-primary-300 transition-colors font-medium ml-1"
                  @click="goToLogin"
               >
                  立即登录
               </button>
            </p>
         </div>
      </template>
   </div>
</template>

<style scoped>
.register-view {
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

/* Select styling */
:deep(.auth-select.p-select) {
   width: 100%;
   background: rgba(255, 255, 255, 0.03);
   border: 1px solid rgba(255, 255, 255, 0.1);
   border-radius: 0.75rem;
   transition: all 0.2s ease;
}

:deep(.auth-select.p-select .p-select-label) {
   padding: 0.875rem 1rem;
   color: white;
   font-size: 0.9375rem;
}

:deep(.auth-select.p-select .p-select-label.p-placeholder) {
   color: rgba(255, 255, 255, 0.3);
}

:deep(.auth-select.p-select .p-select-dropdown) {
   color: white;
   padding-right: 0.5rem;
}

:deep(.auth-select.p-select:not(.p-disabled):hover) {
   border-color: rgba(255, 255, 255, 0.2);
   background: rgba(255, 255, 255, 0.05);
}

:deep(.auth-select.p-select:not(.p-disabled).p-focus) {
   outline: none;
   border-color: var(--p-orange-500);
   background: rgba(255, 255, 255, 0.05);
   box-shadow: 0 0 0 3px rgba(251, 146, 60, 0.15);
}

:deep(.auth-select.p-select.p-disabled) {
   opacity: 0.5;
   cursor: not-allowed;
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

/* Secondary button */
:deep(.auth-button-secondary) {
   padding: 0.875rem 1.5rem;
   font-size: 0.9375rem;
   font-weight: 600;
   border-radius: 0.75rem;
   background: rgba(255, 255, 255, 0.05);
   border: 1px solid rgba(255, 255, 255, 0.15);
   color: white;
   transition: all 0.2s ease;
}

:deep(.auth-button-secondary:hover) {
   background: rgba(255, 255, 255, 0.1);
   border-color: rgba(255, 255, 255, 0.25);
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
