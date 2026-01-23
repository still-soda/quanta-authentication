<script setup lang="ts">
import { ref, reactive } from 'vue';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import InputNumber from 'primevue/inputnumber';
import Textarea from 'primevue/textarea';
import ToggleSwitch from 'primevue/toggleswitch';
import Select from 'primevue/select';
import Divider from 'primevue/divider';
import Message from 'primevue/message';
import Tabs from 'primevue/tabs';
import TabList from 'primevue/tablist';
import Tab from 'primevue/tab';
import TabPanels from 'primevue/tabpanels';
import TabPanel from 'primevue/tabpanel';
import PageHeader from '@/components/shared/PageHeader.vue';
import { SETTING_GROUPS, LANGUAGE_OPTIONS, TIMEZONE_OPTIONS, ENCRYPTION_OPTIONS, STORAGE_TYPE_OPTIONS, DEFAULT_SETTINGS } from '@/config';

const settingGroups = SETTING_GROUPS;

const activeTab = ref('general');

// 基本设置
const generalSettings = reactive({
   siteName: 'Quanta 认证中心',
   siteDescription: '企业级统一身份认证平台',
   adminEmail: 'admin@example.com',
   defaultLanguage: 'zh-CN',
   timezone: 'Asia/Shanghai',
   maintenanceMode: false,
   allowRegistration: true,
   requireEmailVerification: true,
});

// 安全设置
const securitySettings = reactive({
   passwordMinLength: 8,
   passwordRequireUppercase: true,
   passwordRequireLowercase: true,
   passwordRequireNumbers: true,
   passwordRequireSpecial: false,
   maxLoginAttempts: 5,
   lockoutDuration: 30,
   sessionTimeout: 60,
   enableMfa: false,
   enforceHttps: true,
   allowedIpRanges: '',
   csrfProtection: true,
});

// OAuth 配置
const oauthSettings = reactive({
   accessTokenLifetime: 3600,
   refreshTokenLifetime: 604800,
   authCodeLifetime: 600,
   allowImplicitGrant: false,
   allowPasswordGrant: false,
   allowClientCredentials: true,
   requirePkce: true,
   allowedScopes: 'openid profile email',
   jwksRotationDays: 90,
});

// 邮件设置
const emailSettings = reactive({
   smtpHost: 'smtp.example.com',
   smtpPort: 587,
   smtpUser: 'noreply@example.com',
   smtpPassword: '********',
   smtpEncryption: 'tls',
   fromAddress: 'noreply@example.com',
   fromName: 'Quanta Auth',
   enableEmailNotifications: true,
});

// 存储设置
const storageSettings = reactive({
   storageType: 'local',
   localPath: '/var/data/uploads',
   s3Bucket: '',
   s3Region: '',
   s3AccessKey: '',
   s3SecretKey: '',
   maxFileSize: 10,
   allowedFileTypes: 'jpg,png,gif,pdf,doc,docx',
});

const languageOptions = LANGUAGE_OPTIONS;
const timezoneOptions = TIMEZONE_OPTIONS;
const encryptionOptions = ENCRYPTION_OPTIONS;
const storageTypeOptions = STORAGE_TYPE_OPTIONS;

// 保存状态
const isSaving = ref(false);
const saveSuccess = ref(false);

const saveSettings = async () => {
   isSaving.value = true;
   // 模拟保存
   await new Promise((resolve) => setTimeout(resolve, 1000));
   isSaving.value = false;
   saveSuccess.value = true;
   setTimeout(() => {
      saveSuccess.value = false;
   }, 3000);
};

const resetSettings = () => {
   // Reset logic here
   console.log('Reset settings');
};
</script>

<template>
   <div class="flex flex-col gap-6">
      <!-- Page Header -->
      <PageHeader title="系统设置" subtitle="配置和管理系统各项参数">
         <template #actions>
            <Button
               label="重置"
               icon="pi pi-refresh"
               severity="secondary"
               outlined
               @click="resetSettings" />
            <Button
               label="保存设置"
               icon="pi pi-check"
               :loading="isSaving"
               @click="saveSettings" />
         </template>
      </PageHeader>

      <!-- Success Message -->
      <Transition name="slide-fade">
         <Message v-if="saveSuccess" severity="success" :closable="false">
            设置已成功保存！
         </Message>
      </Transition>

      <!-- Settings Tabs -->
      <div
         class="bg-surface-0 dark:bg-surface-900 border border-surface-200 dark:border-surface-800 rounded-xl overflow-hidden">
         <Tabs v-model:value="activeTab">
            <TabList
               :pt="{
                  root: {
                     class: 'bg-surface-50 dark:bg-surface-800 border-b border-surface-200 dark:border-surface-700 px-4',
                  },
               }">
               <Tab
                  v-for="group in settingGroups"
                  :key="group.id"
                  :value="group.id"
                  :pt="{
                     root: {
                        class: 'flex items-center gap-2 px-4 py-3.5 text-sm font-medium',
                     },
                  }">
                  <i :class="group.icon"></i>
                  <span class="max-sm:hidden">{{ group.label }}</span>
               </Tab>
            </TabList>

            <TabPanels :pt="{ root: { class: 'p-6' } }">
               <!-- 基本设置 -->
               <TabPanel value="general">
                  <div class="max-w-2xl">
                     <h3
                        class="text-lg font-semibold text-surface-900 dark:text-surface-100 mb-1">
                        基本设置
                     </h3>
                     <p class="text-sm text-surface-500 mb-6">
                        配置系统的基本信息和默认参数
                     </p>

                     <div class="flex flex-col gap-6">
                        <div class="grid gap-6 sm:grid-cols-2">
                           <div class="flex flex-col gap-2">
                              <label
                                 class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                 站点名称
                              </label>
                              <InputText
                                 v-model="generalSettings.siteName"
                                 class="w-full" />
                           </div>
                           <div class="flex flex-col gap-2">
                              <label
                                 class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                 管理员邮箱
                              </label>
                              <InputText
                                 v-model="generalSettings.adminEmail"
                                 type="email"
                                 class="w-full" />
                           </div>
                        </div>

                        <div class="flex flex-col gap-2">
                           <label
                              class="text-sm font-medium text-surface-700 dark:text-surface-300">
                              站点描述
                           </label>
                           <Textarea
                              v-model="generalSettings.siteDescription"
                              rows="2"
                              class="w-full" />
                        </div>

                        <div class="grid gap-6 sm:grid-cols-2">
                           <div class="flex flex-col gap-2">
                              <label
                                 class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                 默认语言
                              </label>
                              <Select
                                 v-model="generalSettings.defaultLanguage"
                                 :options="languageOptions"
                                 optionLabel="label"
                                 optionValue="value"
                                 class="w-full" />
                           </div>
                           <div class="flex flex-col gap-2">
                              <label
                                 class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                 时区
                              </label>
                              <Select
                                 v-model="generalSettings.timezone"
                                 :options="timezoneOptions"
                                 optionLabel="label"
                                 optionValue="value"
                                 class="w-full" />
                           </div>
                        </div>

                        <Divider />

                        <div class="flex flex-col gap-4">
                           <div
                              class="flex items-center justify-between p-4 bg-surface-50 dark:bg-surface-800 rounded-lg">
                              <div>
                                 <div
                                    class="font-medium text-surface-900 dark:text-surface-100">
                                    允许注册
                                 </div>
                                 <div class="text-sm text-surface-500">
                                    允许新用户注册账号
                                 </div>
                              </div>
                              <ToggleSwitch
                                 v-model="generalSettings.allowRegistration" />
                           </div>

                           <div
                              class="flex items-center justify-between p-4 bg-surface-50 dark:bg-surface-800 rounded-lg">
                              <div>
                                 <div
                                    class="font-medium text-surface-900 dark:text-surface-100">
                                    邮箱验证
                                 </div>
                                 <div class="text-sm text-surface-500">
                                    新注册用户需要验证邮箱
                                 </div>
                              </div>
                              <ToggleSwitch
                                 v-model="
                                    generalSettings.requireEmailVerification
                                 " />
                           </div>

                           <div
                              class="flex items-center justify-between p-4 bg-red-50 dark:bg-red-900/20 rounded-lg border border-red-200 dark:border-red-800">
                              <div>
                                 <div
                                    class="font-medium text-red-700 dark:text-red-400">
                                    维护模式
                                 </div>
                                 <div
                                    class="text-sm text-red-600 dark:text-red-500">
                                    开启后仅管理员可访问系统
                                 </div>
                              </div>
                              <ToggleSwitch
                                 v-model="generalSettings.maintenanceMode" />
                           </div>
                        </div>
                     </div>
                  </div>
               </TabPanel>

               <!-- 安全设置 -->
               <TabPanel value="security">
                  <div class="max-w-2xl">
                     <h3
                        class="text-lg font-semibold text-surface-900 dark:text-surface-100 mb-1">
                        安全设置
                     </h3>
                     <p class="text-sm text-surface-500 mb-6">
                        配置密码策略和登录安全参数
                     </p>

                     <div class="flex flex-col gap-6">
                        <!-- 密码策略 -->
                        <div
                           class="p-4 bg-surface-50 dark:bg-surface-800 rounded-xl">
                           <h4
                              class="text-sm font-semibold text-surface-900 dark:text-surface-100 mb-4 flex items-center gap-2">
                              <i class="pi pi-lock text-primary-500"></i>
                              密码策略
                           </h4>
                           <div class="grid gap-4 sm:grid-cols-2">
                              <div class="flex flex-col gap-2">
                                 <label
                                    class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                    最小长度
                                 </label>
                                 <InputNumber
                                    v-model="securitySettings.passwordMinLength"
                                    :min="6"
                                    :max="32"
                                    class="w-full" />
                              </div>
                           </div>
                           <div class="grid gap-3 mt-4">
                              <label
                                 class="flex items-center gap-3 cursor-pointer">
                                 <ToggleSwitch
                                    v-model="
                                       securitySettings.passwordRequireUppercase
                                    " />
                                 <span
                                    class="text-sm text-surface-700 dark:text-surface-300"
                                    >要求大写字母</span
                                 >
                              </label>
                              <label
                                 class="flex items-center gap-3 cursor-pointer">
                                 <ToggleSwitch
                                    v-model="
                                       securitySettings.passwordRequireLowercase
                                    " />
                                 <span
                                    class="text-sm text-surface-700 dark:text-surface-300"
                                    >要求小写字母</span
                                 >
                              </label>
                              <label
                                 class="flex items-center gap-3 cursor-pointer">
                                 <ToggleSwitch
                                    v-model="
                                       securitySettings.passwordRequireNumbers
                                    " />
                                 <span
                                    class="text-sm text-surface-700 dark:text-surface-300"
                                    >要求数字</span
                                 >
                              </label>
                              <label
                                 class="flex items-center gap-3 cursor-pointer">
                                 <ToggleSwitch
                                    v-model="
                                       securitySettings.passwordRequireSpecial
                                    " />
                                 <span
                                    class="text-sm text-surface-700 dark:text-surface-300"
                                    >要求特殊字符</span
                                 >
                              </label>
                           </div>
                        </div>

                        <!-- 登录安全 -->
                        <div
                           class="p-4 bg-surface-50 dark:bg-surface-800 rounded-xl">
                           <h4
                              class="text-sm font-semibold text-surface-900 dark:text-surface-100 mb-4 flex items-center gap-2">
                              <i class="pi pi-sign-in text-primary-500"></i>
                              登录安全
                           </h4>
                           <div class="grid gap-4 sm:grid-cols-3">
                              <div class="flex flex-col gap-2">
                                 <label
                                    class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                    最大尝试次数
                                 </label>
                                 <InputNumber
                                    v-model="securitySettings.maxLoginAttempts"
                                    :min="3"
                                    :max="10"
                                    class="w-full" />
                              </div>
                              <div class="flex flex-col gap-2">
                                 <label
                                    class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                    锁定时长 (分钟)
                                 </label>
                                 <InputNumber
                                    v-model="securitySettings.lockoutDuration"
                                    :min="5"
                                    :max="120"
                                    class="w-full" />
                              </div>
                              <div class="flex flex-col gap-2">
                                 <label
                                    class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                    会话超时 (分钟)
                                 </label>
                                 <InputNumber
                                    v-model="securitySettings.sessionTimeout"
                                    :min="15"
                                    :max="480"
                                    class="w-full" />
                              </div>
                           </div>
                        </div>

                        <!-- 其他安全选项 -->
                        <div class="flex flex-col gap-4">
                           <div
                              class="flex items-center justify-between p-4 bg-surface-50 dark:bg-surface-800 rounded-lg">
                              <div>
                                 <div
                                    class="font-medium text-surface-900 dark:text-surface-100">
                                    双因素认证
                                 </div>
                                 <div class="text-sm text-surface-500">
                                    启用 TOTP 双因素认证
                                 </div>
                              </div>
                              <ToggleSwitch v-model="securitySettings.enableMfa" />
                           </div>
                           <div
                              class="flex items-center justify-between p-4 bg-surface-50 dark:bg-surface-800 rounded-lg">
                              <div>
                                 <div
                                    class="font-medium text-surface-900 dark:text-surface-100">
                                    强制 HTTPS
                                 </div>
                                 <div class="text-sm text-surface-500">
                                    所有请求必须使用 HTTPS
                                 </div>
                              </div>
                              <ToggleSwitch
                                 v-model="securitySettings.enforceHttps" />
                           </div>
                           <div
                              class="flex items-center justify-between p-4 bg-surface-50 dark:bg-surface-800 rounded-lg">
                              <div>
                                 <div
                                    class="font-medium text-surface-900 dark:text-surface-100">
                                    CSRF 保护
                                 </div>
                                 <div class="text-sm text-surface-500">
                                    启用跨站请求伪造保护
                                 </div>
                              </div>
                              <ToggleSwitch
                                 v-model="securitySettings.csrfProtection" />
                           </div>
                        </div>
                     </div>
                  </div>
               </TabPanel>

               <!-- OAuth 配置 -->
               <TabPanel value="oauth">
                  <div class="max-w-2xl">
                     <h3
                        class="text-lg font-semibold text-surface-900 dark:text-surface-100 mb-1">
                        OAuth 配置
                     </h3>
                     <p class="text-sm text-surface-500 mb-6">
                        配置 OAuth 2.0 / OIDC 协议参数
                     </p>

                     <div class="flex flex-col gap-6">
                        <!-- Token 有效期 -->
                        <div
                           class="p-4 bg-surface-50 dark:bg-surface-800 rounded-xl">
                           <h4
                              class="text-sm font-semibold text-surface-900 dark:text-surface-100 mb-4 flex items-center gap-2">
                              <i class="pi pi-clock text-primary-500"></i>
                              Token 有效期
                           </h4>
                           <div class="grid gap-4 sm:grid-cols-3">
                              <div class="flex flex-col gap-2">
                                 <label
                                    class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                    Access Token (秒)
                                 </label>
                                 <InputNumber
                                    v-model="oauthSettings.accessTokenLifetime"
                                    :min="300"
                                    :max="86400"
                                    class="w-full" />
                              </div>
                              <div class="flex flex-col gap-2">
                                 <label
                                    class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                    Refresh Token (秒)
                                 </label>
                                 <InputNumber
                                    v-model="oauthSettings.refreshTokenLifetime"
                                    :min="3600"
                                    :max="2592000"
                                    class="w-full" />
                              </div>
                              <div class="flex flex-col gap-2">
                                 <label
                                    class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                    Auth Code (秒)
                                 </label>
                                 <InputNumber
                                    v-model="oauthSettings.authCodeLifetime"
                                    :min="60"
                                    :max="1800"
                                    class="w-full" />
                              </div>
                           </div>
                        </div>

                        <!-- 授权方式 -->
                        <div class="flex flex-col gap-4">
                           <div
                              class="flex items-center justify-between p-4 bg-surface-50 dark:bg-surface-800 rounded-lg">
                              <div>
                                 <div
                                    class="font-medium text-surface-900 dark:text-surface-100">
                                    Implicit Grant
                                 </div>
                                 <div class="text-sm text-surface-500">
                                    允许简化授权流程（不推荐）
                                 </div>
                              </div>
                              <ToggleSwitch
                                 v-model="oauthSettings.allowImplicitGrant" />
                           </div>
                           <div
                              class="flex items-center justify-between p-4 bg-surface-50 dark:bg-surface-800 rounded-lg">
                              <div>
                                 <div
                                    class="font-medium text-surface-900 dark:text-surface-100">
                                    Password Grant
                                 </div>
                                 <div class="text-sm text-surface-500">
                                    允许密码授权（不推荐）
                                 </div>
                              </div>
                              <ToggleSwitch
                                 v-model="oauthSettings.allowPasswordGrant" />
                           </div>
                           <div
                              class="flex items-center justify-between p-4 bg-surface-50 dark:bg-surface-800 rounded-lg">
                              <div>
                                 <div
                                    class="font-medium text-surface-900 dark:text-surface-100">
                                    Client Credentials
                                 </div>
                                 <div class="text-sm text-surface-500">
                                    允许客户端凭据授权
                                 </div>
                              </div>
                              <ToggleSwitch
                                 v-model="oauthSettings.allowClientCredentials" />
                           </div>
                           <div
                              class="flex items-center justify-between p-4 bg-green-50 dark:bg-green-900/20 rounded-lg border border-green-200 dark:border-green-800">
                              <div>
                                 <div
                                    class="font-medium text-green-700 dark:text-green-400">
                                    强制 PKCE
                                 </div>
                                 <div
                                    class="text-sm text-green-600 dark:text-green-500">
                                    要求公共客户端使用 PKCE
                                 </div>
                              </div>
                              <ToggleSwitch v-model="oauthSettings.requirePkce" />
                           </div>
                        </div>

                        <div class="flex flex-col gap-2">
                           <label
                              class="text-sm font-medium text-surface-700 dark:text-surface-300">
                              允许的 Scopes
                           </label>
                           <Textarea
                              v-model="oauthSettings.allowedScopes"
                              rows="2"
                              class="w-full"
                              placeholder="openid profile email" />
                           <span class="text-xs text-surface-500">
                              多个 scope 用空格分隔
                           </span>
                        </div>
                     </div>
                  </div>
               </TabPanel>

               <!-- 邮件服务 -->
               <TabPanel value="email">
                  <div class="max-w-2xl">
                     <h3
                        class="text-lg font-semibold text-surface-900 dark:text-surface-100 mb-1">
                        邮件服务
                     </h3>
                     <p class="text-sm text-surface-500 mb-6">
                        配置 SMTP 邮件发送服务
                     </p>

                     <div class="flex flex-col gap-6">
                        <div class="grid gap-6 sm:grid-cols-2">
                           <div class="flex flex-col gap-2">
                              <label
                                 class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                 SMTP 服务器
                              </label>
                              <InputText
                                 v-model="emailSettings.smtpHost"
                                 class="w-full" />
                           </div>
                           <div class="flex flex-col gap-2">
                              <label
                                 class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                 端口
                              </label>
                              <InputNumber
                                 v-model="emailSettings.smtpPort"
                                 :min="1"
                                 :max="65535"
                                 class="w-full" />
                           </div>
                        </div>

                        <div class="grid gap-6 sm:grid-cols-2">
                           <div class="flex flex-col gap-2">
                              <label
                                 class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                 用户名
                              </label>
                              <InputText
                                 v-model="emailSettings.smtpUser"
                                 class="w-full" />
                           </div>
                           <div class="flex flex-col gap-2">
                              <label
                                 class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                 密码
                              </label>
                              <InputText
                                 v-model="emailSettings.smtpPassword"
                                 type="password"
                                 class="w-full" />
                           </div>
                        </div>

                        <div class="flex flex-col gap-2">
                           <label
                              class="text-sm font-medium text-surface-700 dark:text-surface-300">
                              加密方式
                           </label>
                           <Select
                              v-model="emailSettings.smtpEncryption"
                              :options="encryptionOptions"
                              optionLabel="label"
                              optionValue="value"
                              class="w-full sm:w-48" />
                        </div>

                        <Divider />

                        <div class="grid gap-6 sm:grid-cols-2">
                           <div class="flex flex-col gap-2">
                              <label
                                 class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                 发件人地址
                              </label>
                              <InputText
                                 v-model="emailSettings.fromAddress"
                                 type="email"
                                 class="w-full" />
                           </div>
                           <div class="flex flex-col gap-2">
                              <label
                                 class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                 发件人名称
                              </label>
                              <InputText
                                 v-model="emailSettings.fromName"
                                 class="w-full" />
                           </div>
                        </div>

                        <div
                           class="flex items-center justify-between p-4 bg-surface-50 dark:bg-surface-800 rounded-lg">
                           <div>
                              <div
                                 class="font-medium text-surface-900 dark:text-surface-100">
                                 启用邮件通知
                              </div>
                              <div class="text-sm text-surface-500">
                                 发送系统通知和验证邮件
                              </div>
                           </div>
                           <ToggleSwitch
                              v-model="emailSettings.enableEmailNotifications" />
                        </div>

                        <Button
                           label="测试邮件发送"
                           icon="pi pi-send"
                           severity="secondary"
                           outlined
                           class="w-fit" />
                     </div>
                  </div>
               </TabPanel>

               <!-- 存储配置 -->
               <TabPanel value="storage">
                  <div class="max-w-2xl">
                     <h3
                        class="text-lg font-semibold text-surface-900 dark:text-surface-100 mb-1">
                        存储配置
                     </h3>
                     <p class="text-sm text-surface-500 mb-6">
                        配置文件存储服务和限制
                     </p>

                     <div class="flex flex-col gap-6">
                        <div class="flex flex-col gap-2">
                           <label
                              class="text-sm font-medium text-surface-700 dark:text-surface-300">
                              存储类型
                           </label>
                           <Select
                              v-model="storageSettings.storageType"
                              :options="storageTypeOptions"
                              optionLabel="label"
                              optionValue="value"
                              class="w-full sm:w-60" />
                        </div>

                        <!-- 本地存储配置 -->
                        <div
                           v-if="storageSettings.storageType === 'local'"
                           class="flex flex-col gap-2">
                           <label
                              class="text-sm font-medium text-surface-700 dark:text-surface-300">
                              本地路径
                           </label>
                           <InputText
                              v-model="storageSettings.localPath"
                              class="w-full" />
                        </div>

                        <!-- S3 配置 -->
                        <div
                           v-if="storageSettings.storageType === 's3'"
                           class="flex flex-col gap-4 p-4 bg-surface-50 dark:bg-surface-800 rounded-xl">
                           <h4
                              class="text-sm font-semibold text-surface-900 dark:text-surface-100 flex items-center gap-2">
                              <i class="pi pi-cloud text-primary-500"></i>
                              Amazon S3 配置
                           </h4>
                           <div class="grid gap-4 sm:grid-cols-2">
                              <div class="flex flex-col gap-2">
                                 <label
                                    class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                    Bucket
                                 </label>
                                 <InputText
                                    v-model="storageSettings.s3Bucket"
                                    class="w-full" />
                              </div>
                              <div class="flex flex-col gap-2">
                                 <label
                                    class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                    Region
                                 </label>
                                 <InputText
                                    v-model="storageSettings.s3Region"
                                    class="w-full" />
                              </div>
                              <div class="flex flex-col gap-2">
                                 <label
                                    class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                    Access Key
                                 </label>
                                 <InputText
                                    v-model="storageSettings.s3AccessKey"
                                    class="w-full" />
                              </div>
                              <div class="flex flex-col gap-2">
                                 <label
                                    class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                    Secret Key
                                 </label>
                                 <InputText
                                    v-model="storageSettings.s3SecretKey"
                                    type="password"
                                    class="w-full" />
                              </div>
                           </div>
                        </div>

                        <Divider />

                        <div class="grid gap-6 sm:grid-cols-2">
                           <div class="flex flex-col gap-2">
                              <label
                                 class="text-sm font-medium text-surface-700 dark:text-surface-300">
                                 最大文件大小 (MB)
                              </label>
                              <InputNumber
                                 v-model="storageSettings.maxFileSize"
                                 :min="1"
                                 :max="100"
                                 class="w-full" />
                           </div>
                        </div>

                        <div class="flex flex-col gap-2">
                           <label
                              class="text-sm font-medium text-surface-700 dark:text-surface-300">
                              允许的文件类型
                           </label>
                           <InputText
                              v-model="storageSettings.allowedFileTypes"
                              class="w-full"
                              placeholder="jpg,png,gif,pdf" />
                           <span class="text-xs text-surface-500">
                              多个类型用逗号分隔
                           </span>
                        </div>
                     </div>
                  </div>
               </TabPanel>
            </TabPanels>
         </Tabs>
      </div>
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
