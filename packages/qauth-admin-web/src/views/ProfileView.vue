<script setup lang="ts">
import { ref, reactive, computed } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Textarea from 'primevue/textarea';
import Avatar from 'primevue/avatar';
import FileUpload from 'primevue/fileupload';
import Tag from 'primevue/tag';
import Divider from 'primevue/divider';
import Dialog from 'primevue/dialog';
import Password from 'primevue/password';
import Message from 'primevue/message';
import ToggleSwitch from 'primevue/toggleswitch';
import PageHeader from '@/components/shared/PageHeader.vue';
import {
   getProfile,
   updateProfile,
   getUserRoles,
   getLoginHistory,
   getSecuritySettings,
   updateSecuritySettings,
   changePassword,
   uploadAvatar,
} from '@/apis/profile';

const queryClient = useQueryClient();

// 使用 TanStack Query 获取数据
const { data: profile, isLoading: isLoadingProfile } = useQuery({
   queryKey: ['profile'],
   queryFn: getProfile,
});

const { data: userRoles, isLoading: isLoadingRoles } = useQuery({
   queryKey: ['profile', 'roles'],
   queryFn: getUserRoles,
});

const { data: loginHistory, isLoading: isLoadingHistory } = useQuery({
   queryKey: ['profile', 'loginHistory'],
   queryFn: getLoginHistory,
});

const { data: securitySettingsData, isLoading: isLoadingSettings } = useQuery({
   queryKey: ['profile', 'security'],
   queryFn: getSecuritySettings,
});

// 更新资料 mutation
const updateProfileMutation = useMutation({
   mutationFn: updateProfile,
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profile'] });
      isEditing.value = false;
      saveSuccess.value = true;
      setTimeout(() => {
         saveSuccess.value = false;
      }, 3000);
   },
});

// 更新安全设置 mutation
const updateSecurityMutation = useMutation({
   mutationFn: updateSecuritySettings,
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profile', 'security'] });
   },
});

// 修改密码 mutation
const changePasswordMutation = useMutation({
   mutationFn: changePassword,
   onSuccess: () => {
      passwordDialog.value = false;
      passwordForm.currentPassword = '';
      passwordForm.newPassword = '';
      passwordForm.confirmPassword = '';
   },
});

// 安全设置响应式绑定
const securitySettings = computed({
   get: () => securitySettingsData.value || { mfaEnabled: false, emailNotifications: true, loginAlerts: true },
   set: (val) => updateSecurityMutation.mutate(val),
});

// 编辑状态
const isEditing = ref(false);
const editForm = reactive({
   displayName: '',
   phone: '',
   bio: '',
});

// 密码修改对话框
const passwordDialog = ref(false);
const passwordForm = reactive({
   currentPassword: '',
   newPassword: '',
   confirmPassword: '',
});

// 头像上传对话框
const avatarDialog = ref(false);

// 保存状态
const saveSuccess = ref(false);

// 开始编辑
const startEditing = () => {
   if (profile.value) {
      editForm.displayName = profile.value.displayName;
      editForm.phone = profile.value.phone;
      editForm.bio = profile.value.bio;
   }
   isEditing.value = true;
};

// 取消编辑
const cancelEditing = () => {
   isEditing.value = false;
};

// 保存资料
const saveProfile = () => {
   updateProfileMutation.mutate(editForm);
};

// 修改密码
const handleChangePassword = () => {
   if (passwordForm.newPassword !== passwordForm.confirmPassword) {
      return;
   }
   changePasswordMutation.mutate({
      currentPassword: passwordForm.currentPassword,
      newPassword: passwordForm.newPassword,
   });
};

// 头像上传
const onAvatarUpload = (event: any) => {
   console.log('Avatar uploaded:', event);
   avatarDialog.value = false;
};

// 获取状态标签
const getStatusSeverity = (status: string) => {
   switch (status) {
      case 'ACTIVE':
         return 'success';
      case 'LOCKED':
         return 'warn';
      case 'BANNED':
         return 'danger';
      default:
         return 'secondary';
   }
};

const getStatusLabel = (status: string) => {
   switch (status) {
      case 'ACTIVE':
         return '正常';
      case 'LOCKED':
         return '已锁定';
      case 'BANNED':
         return '已禁用';
      default:
         return status;
   }
};

// 密码强度
const passwordStrength = computed(() => {
   const pwd = passwordForm.newPassword;
   if (!pwd) return 0;
   let strength = 0;
   if (pwd.length >= 8) strength++;
   if (/[a-z]/.test(pwd)) strength++;
   if (/[A-Z]/.test(pwd)) strength++;
   if (/[0-9]/.test(pwd)) strength++;
   if (/[^a-zA-Z0-9]/.test(pwd)) strength++;
   return strength;
});

const passwordStrengthLabel = computed(() => {
   const labels = ['极弱', '弱', '一般', '强', '很强'];
   return labels[Math.min(passwordStrength.value, 4)] || '';
});

const passwordStrengthColor = computed(() => {
   const colors = [
      'bg-red-500',
      'bg-orange-500',
      'bg-yellow-500',
      'bg-blue-500',
      'bg-green-500',
   ];
   return colors[Math.min(passwordStrength.value, 4)] || '';
});

// 切换安全设置
const toggleMfa = (value: boolean) => {
   updateSecurityMutation.mutate({ mfaEnabled: value });
};

const toggleEmailNotifications = (value: boolean) => {
   updateSecurityMutation.mutate({ emailNotifications: value });
};

const toggleLoginAlerts = (value: boolean) => {
   updateSecurityMutation.mutate({ loginAlerts: value });
};
</script>

<template>
   <div class="flex flex-col gap-6">
      <!-- Page Header -->
      <PageHeader title="个人资料" subtitle="管理您的账号信息和安全设置">
         <template #actions>
            <template v-if="isEditing">
               <Button
                  label="取消"
                  severity="secondary"
                  outlined
                  @click="cancelEditing" />
               <Button
                  label="保存"
                  icon="pi pi-check"
                  :loading="updateProfileMutation.isPending.value"
                  @click="saveProfile" />
            </template>
            <template v-else>
               <Button
                  label="编辑资料"
                  icon="pi pi-pencil"
                  @click="startEditing" />
            </template>
         </template>
      </PageHeader>

      <!-- Success Message -->
      <Transition name="slide-fade">
         <Message v-if="saveSuccess" severity="success" :closable="false">
            个人资料已成功更新！
         </Message>
      </Transition>

      <!-- Loading State -->
      <div
         v-if="isLoadingProfile"
         class="flex items-center justify-center py-24">
         <i class="pi pi-spin pi-spinner text-3xl text-surface-400"></i>
      </div>

      <div v-else-if="profile" class="grid grid-cols-1 lg:grid-cols-[1fr_360px] gap-6">
         <!-- Main Profile Section -->
         <div class="flex flex-col gap-6">
            <!-- Profile Card -->
            <div
               class="bg-surface-0 dark:bg-surface-900 border border-surface-200 dark:border-surface-800 rounded-xl overflow-hidden">
               <!-- Header with gradient -->
               <div
                  class="h-28 bg-linear-to-br from-primary-400 via-primary-500 to-primary-600 relative">
                  <div
                     class="absolute inset-0 bg-[radial-gradient(circle_at_30%_50%,rgba(255,255,255,0.1),transparent)]"></div>
               </div>

               <!-- Avatar and basic info -->
               <div class="px-6 pb-6">
                  <div class="flex items-end gap-5 -mt-12 relative z-10">
                     <div class="relative group">
                        <div
                           class="w-24 h-24 rounded-2xl overflow-hidden bg-surface-100 dark:bg-surface-700 border-4 border-surface-0 dark:border-surface-900 shadow-xl">
                           <img
                              :src="profile.avatar"
                              :alt="profile.name"
                              class="w-full h-full object-cover" />
                        </div>
                        <button
                           class="absolute inset-0 flex items-center justify-center bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity rounded-2xl cursor-pointer border-none"
                           @click="avatarDialog = true">
                           <i class="pi pi-camera text-white text-xl"></i>
                        </button>
                     </div>
                     <div class="flex-1 pb-1">
                        <div class="flex items-center gap-3 mb-1">
                           <h2
                              class="text-xl font-bold text-surface-900 dark:text-surface-100 m-0">
                              {{ profile.name }}
                           </h2>
                           <Tag
                              :severity="getStatusSeverity(profile.status)"
                              :value="getStatusLabel(profile.status)"
                              :pt="{ root: { class: 'text-xs' } }" />
                        </div>
                        <p class="text-sm text-surface-500 m-0">
                           {{ profile.email }}
                           <i
                              v-if="profile.emailVerified"
                              class="pi pi-verified text-green-500 ml-1"
                              v-tooltip="'邮箱已验证'"></i>
                        </p>
                     </div>
                  </div>

                  <Divider />

                  <!-- Profile Form -->
                  <div class="grid gap-5 sm:grid-cols-2">
                     <div class="flex flex-col gap-2">
                        <label
                           class="text-sm font-medium text-surface-500">
                           学号 / 工号
                        </label>
                        <InputText
                           :modelValue="profile.studentId"
                           disabled
                           class="w-full" />
                     </div>
                     <div class="flex flex-col gap-2">
                        <label
                           class="text-sm font-medium text-surface-500">
                           邮箱
                        </label>
                        <InputText
                           :modelValue="profile.email"
                           disabled
                           class="w-full" />
                     </div>
                     <div class="flex flex-col gap-2">
                        <label
                           class="text-sm font-medium text-surface-700 dark:text-surface-300">
                           显示名称
                        </label>
                        <InputText
                           v-if="isEditing"
                           v-model="editForm.displayName"
                           class="w-full" />
                        <InputText
                           v-else
                           :modelValue="profile.displayName"
                           disabled
                           class="w-full" />
                     </div>
                     <div class="flex flex-col gap-2">
                        <label
                           class="text-sm font-medium text-surface-700 dark:text-surface-300">
                           手机号
                        </label>
                        <InputText
                           v-if="isEditing"
                           v-model="editForm.phone"
                           class="w-full" />
                        <InputText
                           v-else
                           :modelValue="profile.phone"
                           disabled
                           class="w-full" />
                     </div>
                     <div class="flex flex-col gap-2 sm:col-span-2">
                        <label
                           class="text-sm font-medium text-surface-700 dark:text-surface-300">
                           个人简介
                        </label>
                        <Textarea
                           v-if="isEditing"
                           v-model="editForm.bio"
                           rows="3"
                           class="w-full" />
                        <Textarea
                           v-else
                           :modelValue="profile.bio"
                           disabled
                           rows="3"
                           class="w-full" />
                     </div>
                  </div>
               </div>
            </div>

            <!-- Login History -->
            <div
               class="bg-surface-0 dark:bg-surface-900 border border-surface-200 dark:border-surface-800 rounded-xl p-6">
               <h3
                  class="text-base font-semibold text-surface-900 dark:text-surface-100 m-0 mb-4 flex items-center gap-2">
                  <i class="pi pi-history text-primary-500"></i>
                  登录历史
               </h3>
               <div
                  v-if="isLoadingHistory"
                  class="flex items-center justify-center py-8">
                  <i class="pi pi-spin pi-spinner text-xl text-surface-400"></i>
               </div>
               <div v-else class="flex flex-col gap-3">
                  <div
                     v-for="record in loginHistory"
                     :key="record.id"
                     class="flex items-center gap-4 p-3 bg-surface-50 dark:bg-surface-800 rounded-lg">
                     <div
                        class="w-10 h-10 rounded-lg flex items-center justify-center shrink-0"
                        :class="
                           record.status === 'success'
                              ? 'bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400'
                              : 'bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400'
                        ">
                        <i
                           :class="
                              record.status === 'success'
                                 ? 'pi pi-check'
                                 : 'pi pi-times'
                           "></i>
                     </div>
                     <div class="flex-1 min-w-0">
                        <div
                           class="text-sm font-medium text-surface-900 dark:text-surface-100">
                           {{ record.device }}
                        </div>
                        <div class="text-xs text-surface-500">
                           {{ record.ip }} · {{ record.location }}
                        </div>
                     </div>
                     <div class="text-xs text-surface-500 tabular-nums shrink-0">
                        {{ record.time }}
                     </div>
                  </div>
               </div>
            </div>
         </div>

         <!-- Sidebar -->
         <div class="flex flex-col gap-6">
            <!-- Account Info -->
            <div
               class="bg-surface-0 dark:bg-surface-900 border border-surface-200 dark:border-surface-800 rounded-xl p-6">
               <h3
                  class="text-base font-semibold text-surface-900 dark:text-surface-100 m-0 mb-4 flex items-center gap-2">
                  <i class="pi pi-info-circle text-primary-500"></i>
                  账号信息
               </h3>
               <div class="flex flex-col gap-4">
                  <div
                     class="flex items-center justify-between py-2 border-b border-surface-100 dark:border-surface-800">
                     <span class="text-sm text-surface-500">注册时间</span>
                     <span
                        class="text-sm font-medium text-surface-900 dark:text-surface-100">
                        {{ profile.createdAt.split(' ')[0] }}
                     </span>
                  </div>
                  <div
                     class="flex items-center justify-between py-2 border-b border-surface-100 dark:border-surface-800">
                     <span class="text-sm text-surface-500">最后登录</span>
                     <span
                        class="text-sm font-medium text-surface-900 dark:text-surface-100">
                        {{ profile.lastLogin }}
                     </span>
                  </div>
                  <div class="flex items-center justify-between py-2">
                     <span class="text-sm text-surface-500">账号 ID</span>
                     <span
                        class="text-xs font-mono text-surface-400 dark:text-surface-500">
                        {{ profile.id }}
                     </span>
                  </div>
               </div>
            </div>

            <!-- Roles -->
            <div
               class="bg-surface-0 dark:bg-surface-900 border border-surface-200 dark:border-surface-800 rounded-xl p-6">
               <h3
                  class="text-base font-semibold text-surface-900 dark:text-surface-100 m-0 mb-4 flex items-center gap-2">
                  <i class="pi pi-shield text-primary-500"></i>
                  角色权限
               </h3>
               <div
                  v-if="isLoadingRoles"
                  class="flex items-center justify-center py-4">
                  <i class="pi pi-spin pi-spinner text-xl text-surface-400"></i>
               </div>
               <div v-else class="flex flex-wrap gap-2">
                  <Tag
                     v-for="role in userRoles"
                     :key="role.code"
                     :severity="role.isSystem ? 'warn' : 'secondary'"
                     :value="role.name"
                     :pt="{ root: { class: 'text-xs px-2.5 py-1' } }" />
               </div>
            </div>

            <!-- Security Settings -->
            <div
               class="bg-surface-0 dark:bg-surface-900 border border-surface-200 dark:border-surface-800 rounded-xl p-6">
               <h3
                  class="text-base font-semibold text-surface-900 dark:text-surface-100 m-0 mb-4 flex items-center gap-2">
                  <i class="pi pi-lock text-primary-500"></i>
                  安全设置
               </h3>
               <div class="flex flex-col gap-4">
                  <Button
                     label="修改密码"
                     icon="pi pi-key"
                     severity="secondary"
                     outlined
                     class="w-full"
                     @click="passwordDialog = true" />

                  <Divider class="my-2" />

                  <div
                     v-if="isLoadingSettings"
                     class="flex items-center justify-center py-4">
                     <i class="pi pi-spin pi-spinner text-xl text-surface-400"></i>
                  </div>
                  <template v-else-if="securitySettings">
                     <div
                        class="flex items-center justify-between p-3 bg-surface-50 dark:bg-surface-800 rounded-lg">
                        <div>
                           <div
                              class="text-sm font-medium text-surface-900 dark:text-surface-100">
                              双因素认证
                           </div>
                           <div class="text-xs text-surface-500">
                              {{ securitySettings.mfaEnabled ? '已启用' : '未启用' }}
                           </div>
                        </div>
                        <ToggleSwitch
                           :modelValue="securitySettings.mfaEnabled"
                           @update:modelValue="toggleMfa" />
                     </div>

                     <div
                        class="flex items-center justify-between p-3 bg-surface-50 dark:bg-surface-800 rounded-lg">
                        <div>
                           <div
                              class="text-sm font-medium text-surface-900 dark:text-surface-100">
                              邮件通知
                           </div>
                           <div class="text-xs text-surface-500">
                              接收重要通知邮件
                           </div>
                        </div>
                        <ToggleSwitch
                           :modelValue="securitySettings.emailNotifications"
                           @update:modelValue="toggleEmailNotifications" />
                     </div>

                     <div
                        class="flex items-center justify-between p-3 bg-surface-50 dark:bg-surface-800 rounded-lg">
                        <div>
                           <div
                              class="text-sm font-medium text-surface-900 dark:text-surface-100">
                              登录提醒
                           </div>
                           <div class="text-xs text-surface-500">
                              新设备登录时发送提醒
                           </div>
                        </div>
                        <ToggleSwitch
                           :modelValue="securitySettings.loginAlerts"
                           @update:modelValue="toggleLoginAlerts" />
                     </div>
                  </template>
               </div>
            </div>
         </div>
      </div>

      <!-- Password Dialog -->
      <Dialog
         v-model:visible="passwordDialog"
         header="修改密码"
         modal
         :style="{ width: '420px' }"
         :pt="{
            root: { class: 'border-none shadow-2xl rounded-2xl' },
            header: {
               class: 'border-b border-surface-200 dark:border-surface-700 px-6 py-4',
            },
            content: { class: 'p-6' },
         }">
         <div class="flex flex-col gap-5">
            <div class="flex flex-col gap-2">
               <label
                  class="text-sm font-medium text-surface-700 dark:text-surface-300">
                  当前密码
               </label>
               <Password
                  v-model="passwordForm.currentPassword"
                  :feedback="false"
                  toggleMask
                  class="w-full"
                  inputClass="w-full" />
            </div>
            <div class="flex flex-col gap-2">
               <label
                  class="text-sm font-medium text-surface-700 dark:text-surface-300">
                  新密码
               </label>
               <Password
                  v-model="passwordForm.newPassword"
                  :feedback="false"
                  toggleMask
                  class="w-full"
                  inputClass="w-full" />
               <!-- Password strength indicator -->
               <div v-if="passwordForm.newPassword" class="flex items-center gap-2">
                  <div class="flex-1 h-1.5 bg-surface-200 dark:bg-surface-700 rounded-full overflow-hidden">
                     <div
                        :class="passwordStrengthColor"
                        class="h-full transition-all duration-300"
                        :style="{
                           width: `${(passwordStrength / 5) * 100}%`,
                        }"></div>
                  </div>
                  <span class="text-xs text-surface-500 w-12">{{
                     passwordStrengthLabel
                  }}</span>
               </div>
            </div>
            <div class="flex flex-col gap-2">
               <label
                  class="text-sm font-medium text-surface-700 dark:text-surface-300">
                  确认新密码
               </label>
               <Password
                  v-model="passwordForm.confirmPassword"
                  :feedback="false"
                  toggleMask
                  class="w-full"
                  inputClass="w-full" />
               <small
                  v-if="
                     passwordForm.confirmPassword &&
                     passwordForm.newPassword !== passwordForm.confirmPassword
                  "
                  class="text-red-500 text-xs">
                  两次输入的密码不一致
               </small>
            </div>
         </div>
         <template #footer>
            <div class="flex justify-end gap-3">
               <Button
                  label="取消"
                  severity="secondary"
                  outlined
                  @click="passwordDialog = false" />
               <Button
                  label="确认修改"
                  :disabled="
                     !passwordForm.currentPassword ||
                     !passwordForm.newPassword ||
                     passwordForm.newPassword !== passwordForm.confirmPassword
                  "
                  :loading="changePasswordMutation.isPending.value"
                  @click="handleChangePassword" />
            </div>
         </template>
      </Dialog>

      <!-- Avatar Upload Dialog -->
      <Dialog
         v-model:visible="avatarDialog"
         header="更换头像"
         modal
         :style="{ width: '480px' }"
         :pt="{
            root: { class: 'border-none shadow-2xl rounded-2xl' },
            header: {
               class: 'border-b border-surface-200 dark:border-surface-700 px-6 py-4',
            },
            content: { class: 'p-6' },
         }">
         <div class="flex flex-col items-center gap-6">
            <div
               class="w-32 h-32 rounded-2xl overflow-hidden bg-surface-100 dark:bg-surface-700">
               <img
                  v-if="profile"
                  :src="profile.avatar"
                  :alt="profile.name"
                  class="w-full h-full object-cover" />
            </div>
            <FileUpload
               mode="basic"
               accept="image/*"
               :maxFileSize="2000000"
               chooseLabel="选择图片"
               :auto="true"
               customUpload
               @uploader="onAvatarUpload"
               :pt="{
                  root: { class: 'w-full' },
                  chooseButton: { class: 'w-full' },
               }" />
            <p class="text-xs text-surface-500 text-center m-0">
               支持 JPG、PNG、GIF 格式，最大 2MB
            </p>
         </div>
      </Dialog>
   </div>
</template>

<style scoped>
.slide-fade-enter-active,
.slide-fade-leave-active {
   transition: all 0.3s ease;
}

.slide-fade-enter-from,
.slide-fade-leave-to {
   opacity: 0;
   transform: translateY(-10px);
}
</style>
