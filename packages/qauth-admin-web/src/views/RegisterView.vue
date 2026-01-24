<script setup lang="ts">
import { ref, computed, reactive } from 'vue'
import { useRouter } from 'vue-router'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Select from 'primevue/select'
import Button from 'primevue/button'
import { useToast } from 'primevue/usetoast'

const router = useRouter()
const toast = useToast()

// Form state
const formData = reactive({
   username: '',
   email: '',
   fullName: '',
   department: '',
   phone: '',
   purpose: '',
})

const isLoading = ref(false)
const isSubmitted = ref(false)
const currentStep = ref(1)

// Department options
const departments = [
   { label: '技术部', value: 'tech' },
   { label: '产品部', value: 'product' },
   { label: '运营部', value: 'operations' },
   { label: '市场部', value: 'marketing' },
   { label: '人力资源', value: 'hr' },
   { label: '财务部', value: 'finance' },
   { label: '其他', value: 'other' },
]

// Validation
const isEmailValid = computed(() => {
   const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
   return emailRegex.test(formData.email)
})

const isUsernameValid = computed(() => {
   const usernameRegex = /^[a-zA-Z0-9_]{3,20}$/
   return usernameRegex.test(formData.username)
})

const isPhoneValid = computed(() => {
   if (!formData.phone) return true // Optional field
   const phoneRegex = /^1[3-9]\d{9}$/
   return phoneRegex.test(formData.phone)
})

const isStep1Valid = computed(() => {
   return isUsernameValid.value && isEmailValid.value && formData.fullName.trim().length >= 2
})

const isStep2Valid = computed(() => {
   return formData.department && isPhoneValid.value && formData.purpose.trim().length >= 20
})

const isFormValid = computed(() => isStep1Valid.value && isStep2Valid.value)

// Handle step navigation
const nextStep = () => {
   if (currentStep.value === 1 && isStep1Valid.value) {
      currentStep.value = 2
   }
}

const prevStep = () => {
   if (currentStep.value === 2) {
      currentStep.value = 1
   }
}

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
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 2000))

      isSubmitted.value = true
      toast.add({
         severity: 'success',
         summary: '申请已提交',
         detail: '您的注册申请已成功提交，请等待管理员审核',
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
                  <i class="pi pi-send text-4xl text-green-400 animate-bounce-in"></i>
               </div>
            </div>

            <h2 class="text-xl font-semibold text-white mb-3">注册申请已提交</h2>
            <p class="text-sm text-white/50 mb-6 leading-relaxed">
               感谢您的申请！管理员将在 <span class="text-primary-400">1-3 个工作日</span> 内<br />
               审核您的注册请求，审核结果将通过邮件通知您。
            </p>

            <!-- Application summary -->
            <div class="p-4 rounded-xl bg-white/3 border border-white/10 mb-6 text-left">
               <p class="text-xs text-white/40 mb-3 uppercase tracking-wider">申请摘要</p>
               <div class="space-y-2">
                  <div class="flex justify-between items-center">
                     <span class="text-sm text-white/50">用户名</span>
                     <span class="text-sm text-white font-medium">{{ formData.username }}</span>
                  </div>
                  <div class="flex justify-between items-center">
                     <span class="text-sm text-white/50">邮箱</span>
                     <span class="text-sm text-white font-medium">{{ formData.email }}</span>
                  </div>
                  <div class="flex justify-between items-center">
                     <span class="text-sm text-white/50">姓名</span>
                     <span class="text-sm text-white font-medium">{{ formData.fullName }}</span>
                  </div>
               </div>
            </div>

            <Button
               label="返回登录"
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
            <h2 class="text-xl font-semibold text-white mb-2">申请注册</h2>
            <p class="text-sm text-white/50">填写以下信息申请账户，管理员将审核您的申请</p>
         </div>

         <!-- Progress steps -->
         <div class="flex items-center gap-3 mb-6">
            <div
               class="flex items-center gap-2 px-3 py-1.5 rounded-full transition-all"
               :class="
                  currentStep >= 1
                     ? 'bg-primary-500/20 text-primary-400'
                     : 'bg-white/5 text-white/40'
               "
            >
               <span
                  class="w-5 h-5 rounded-full flex items-center justify-center text-xs font-medium"
                  :class="
                     currentStep >= 1 ? 'bg-primary-500 text-white' : 'bg-white/10 text-white/50'
                  "
               >
                  {{ currentStep > 1 ? '✓' : '1' }}
               </span>
               <span class="text-xs font-medium">基本信息</span>
            </div>
            <div class="flex-1 h-px bg-white/10"></div>
            <div
               class="flex items-center gap-2 px-3 py-1.5 rounded-full transition-all"
               :class="
                  currentStep >= 2
                     ? 'bg-primary-500/20 text-primary-400'
                     : 'bg-white/5 text-white/40'
               "
            >
               <span
                  class="w-5 h-5 rounded-full flex items-center justify-center text-xs font-medium"
                  :class="
                     currentStep >= 2 ? 'bg-primary-500 text-white' : 'bg-white/10 text-white/50'
                  "
               >
                  2
               </span>
               <span class="text-xs font-medium">详细信息</span>
            </div>
         </div>

         <!-- Form -->
         <form @submit.prevent="handleSubmit" class="space-y-5">
            <!-- Step 1: Basic Info -->
            <template v-if="currentStep === 1">
               <!-- Username field -->
               <div class="form-field">
                  <label class="form-label">
                     <i class="pi pi-user mr-2 opacity-60"></i>
                     用户名
                     <span class="text-red-400 ml-1">*</span>
                  </label>
                  <div class="input-wrapper">
                     <InputText
                        v-model="formData.username"
                        placeholder="3-20位字母、数字或下划线"
                        class="auth-input"
                        :disabled="isLoading"
                        :invalid="formData.username.length > 0 && !isUsernameValid"
                        autocomplete="username"
                     />
                  </div>
                  <p
                     v-if="formData.username.length > 0 && !isUsernameValid"
                     class="text-xs text-red-400 mt-1"
                  >
                     用户名需为3-20位字母、数字或下划线
                  </p>
               </div>

               <!-- Email field -->
               <div class="form-field">
                  <label class="form-label">
                     <i class="pi pi-envelope mr-2 opacity-60"></i>
                     联系邮箱
                     <span class="text-red-400 ml-1">*</span>
                  </label>
                  <div class="input-wrapper">
                     <InputText
                        v-model="formData.email"
                        placeholder="用于接收审核结果和系统通知"
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

               <!-- Full name field -->
               <div class="form-field">
                  <label class="form-label">
                     <i class="pi pi-id-card mr-2 opacity-60"></i>
                     真实姓名
                     <span class="text-red-400 ml-1">*</span>
                  </label>
                  <div class="input-wrapper">
                     <InputText
                        v-model="formData.fullName"
                        placeholder="请输入您的真实姓名"
                        class="auth-input"
                        :disabled="isLoading"
                        autocomplete="name"
                     />
                  </div>
               </div>

               <!-- Next button -->
               <Button
                  type="button"
                  label="下一步"
                  icon="pi pi-arrow-right"
                  iconPos="right"
                  :disabled="!isStep1Valid"
                  class="auth-button w-full"
                  severity="contrast"
                  @click="nextStep"
               />
            </template>

            <!-- Step 2: Additional Info -->
            <template v-if="currentStep === 2">
               <!-- Department field -->
               <div class="form-field">
                  <label class="form-label">
                     <i class="pi pi-building mr-2 opacity-60"></i>
                     所属部门
                     <span class="text-red-400 ml-1">*</span>
                  </label>
                  <div class="input-wrapper">
                     <Select
                        v-model="formData.department"
                        :options="departments"
                        optionLabel="label"
                        optionValue="value"
                        placeholder="请选择您的部门"
                        class="auth-select w-full"
                        :disabled="isLoading"
                     />
                  </div>
               </div>

               <!-- Phone field -->
               <div class="form-field">
                  <label class="form-label">
                     <i class="pi pi-phone mr-2 opacity-60"></i>
                     联系电话
                     <span class="text-white/30 ml-1 text-xs">（选填）</span>
                  </label>
                  <div class="input-wrapper">
                     <InputText
                        v-model="formData.phone"
                        placeholder="请输入您的手机号码"
                        class="auth-input"
                        :disabled="isLoading"
                        :invalid="formData.phone.length > 0 && !isPhoneValid"
                        autocomplete="tel"
                     />
                  </div>
                  <p
                     v-if="formData.phone.length > 0 && !isPhoneValid"
                     class="text-xs text-red-400 mt-1"
                  >
                     请输入有效的手机号码
                  </p>
               </div>

               <!-- Purpose field -->
               <div class="form-field">
                  <label class="form-label">
                     <i class="pi pi-file-edit mr-2 opacity-60"></i>
                     申请目的
                     <span class="text-red-400 ml-1">*</span>
                  </label>
                  <div class="input-wrapper">
                     <Textarea
                        v-model="formData.purpose"
                        placeholder="请详细说明您申请账户的目的，包括您的工作职责、需要使用的系统功能等（至少20个字符）"
                        class="auth-textarea"
                        :disabled="isLoading"
                        rows="4"
                        autoResize
                     />
                  </div>
                  <p class="text-xs text-white/30 mt-1">
                     {{ formData.purpose.length }}/20 最少字符
                  </p>
               </div>

               <!-- Info box -->
               <div class="p-4 rounded-xl bg-primary-500/10 border border-primary-500/20">
                  <div class="flex gap-3">
                     <i class="pi pi-info-circle text-primary-400 mt-0.5 shrink-0"></i>
                     <div class="text-xs text-white/60 leading-relaxed">
                        <p class="font-medium text-white/80 mb-1">审核说明</p>
                        <p>
                           注册申请将由管理员人工审核。请确保您提供的信息真实有效，以便我们快速处理您的申请。
                        </p>
                     </div>
                  </div>
               </div>

               <!-- Buttons -->
               <div class="flex gap-3">
                  <Button
                     type="button"
                     label="上一步"
                     icon="pi pi-arrow-left"
                     severity="secondary"
                     outlined
                     class="auth-button-secondary flex-1"
                     @click="prevStep"
                  />
                  <Button
                     type="submit"
                     :label="isLoading ? '提交中...' : '提交申请'"
                     :loading="isLoading"
                     :disabled="!isStep2Valid || isLoading"
                     class="auth-button flex-1"
                     severity="contrast"
                  />
               </div>
            </template>
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
