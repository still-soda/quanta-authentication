<script setup lang="ts">
import { ref, computed } from 'vue';
import Button from 'primevue/button';
import PageHeader from '@/components/shared/PageHeader.vue';
import MiniStats, {
   type MiniStatItem,
} from '@/components/shared/MiniStats.vue';
import SearchBox from '@/components/shared/SearchBox.vue';
import OAuthAppCard, {
   type OAuthApp,
} from '@/components/oauth/OAuthAppCard.vue';
import OAuthAppDialog, {
   type OAuthAppFormData,
} from '@/components/oauth/OAuthAppDialog.vue';
import SecretDialog from '@/components/oauth/SecretDialog.vue';

// OAuth 应用数据
const apps = ref<OAuthApp[]>([
   {
      id: 1,
      name: 'Web Dashboard',
      clientId: 'web_dashboard_prod',
      description: '主要的 Web 管理后台应用',
      icon: 'pi pi-desktop',
      iconBg: 'linear-gradient(135deg, #f97316 0%, #ea580c 100%)',
      redirectUris: ['https://dashboard.example.com/callback'],
      scopes: ['openid', 'profile', 'email', 'admin'],
      grantTypes: ['authorization_code', 'refresh_token'],
      status: 'active',
      trusted: true,
      createdAt: '2024-01-15',
      lastUsed: '2026-01-23',
      requestCount: 125840,
   },
   {
      id: 2,
      name: 'Mobile App',
      clientId: 'mobile_app_v2',
      description: 'iOS 和 Android 移动应用',
      icon: 'pi pi-mobile',
      iconBg: 'linear-gradient(135deg, #3b82f6 0%, #2563eb 100%)',
      redirectUris: ['myapp://callback', 'https://mobile.example.com/callback'],
      scopes: ['openid', 'profile', 'email', 'offline_access'],
      grantTypes: ['authorization_code', 'refresh_token'],
      status: 'active',
      trusted: true,
      createdAt: '2024-03-20',
      lastUsed: '2026-01-23',
      requestCount: 89562,
   },
   {
      id: 3,
      name: 'Partner API',
      clientId: 'partner_api_client',
      description: '第三方合作伙伴 API 访问',
      icon: 'pi pi-link',
      iconBg: 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
      redirectUris: ['https://partner.example.com/oauth/callback'],
      scopes: ['openid', 'profile', 'read:users'],
      grantTypes: ['client_credentials'],
      status: 'active',
      trusted: false,
      createdAt: '2024-06-10',
      lastUsed: '2026-01-22',
      requestCount: 34521,
   },
   {
      id: 4,
      name: 'Internal Tools',
      clientId: 'internal_tools_dev',
      description: '内部开发工具和脚本',
      icon: 'pi pi-wrench',
      iconBg: 'linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%)',
      redirectUris: ['http://localhost:3000/callback'],
      scopes: ['openid', 'profile', 'admin'],
      grantTypes: ['authorization_code', 'client_credentials'],
      status: 'development',
      trusted: true,
      createdAt: '2024-09-05',
      lastUsed: '2026-01-20',
      requestCount: 8421,
   },
   {
      id: 5,
      name: 'Legacy System',
      clientId: 'legacy_system_v1',
      description: '旧版系统兼容接口',
      icon: 'pi pi-history',
      iconBg: 'linear-gradient(135deg, #6b7280 0%, #4b5563 100%)',
      redirectUris: ['https://legacy.example.com/auth'],
      scopes: ['openid', 'profile'],
      grantTypes: ['authorization_code'],
      status: 'deprecated',
      trusted: true,
      createdAt: '2023-05-01',
      lastUsed: '2026-01-15',
      requestCount: 2150,
   },
]);

const selectedApp = ref<OAuthApp | null>(null);
const appDialog = ref(false);
const secretDialog = ref(false);
const isEditing = ref(false);
const searchQuery = ref('');
const newSecret = ref('sk_live_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx');

const filteredApps = computed(() => {
   if (!searchQuery.value) return apps.value;
   const query = searchQuery.value.toLowerCase();
   return apps.value.filter(
      (app) =>
         app.name.toLowerCase().includes(query) ||
         app.clientId.toLowerCase().includes(query),
   );
});

const stats = computed<MiniStatItem[]>(() => [
   { label: '总应用', value: apps.value.length },
   {
      label: '生产环境',
      value: apps.value.filter((a) => a.status === 'active').length,
      colorClass: 'text-green-500',
   },
   {
      label: '开发中',
      value: apps.value.filter((a) => a.status === 'development').length,
      colorClass: 'text-yellow-500',
   },
   {
      label: '已弃用',
      value: apps.value.filter((a) => a.status === 'deprecated').length,
      colorClass: 'text-gray-400',
   },
]);

const openNewAppDialog = () => {
   isEditing.value = false;
   selectedApp.value = null;
   appDialog.value = true;
};

const editApp = (app: OAuthApp) => {
   isEditing.value = true;
   selectedApp.value = app;
   appDialog.value = true;
};

const viewApp = (app: OAuthApp) => {
   selectedApp.value = app;
};

const regenerateSecret = (app: OAuthApp) => {
   selectedApp.value = app;
   newSecret.value = 'sk_live_' + Math.random().toString(36).substring(2, 34);
   secretDialog.value = true;
};

const saveApp = (data: OAuthAppFormData) => {
   console.log('Saving app:', data);
   appDialog.value = false;
};
</script>

<template>
   <div class="flex flex-col gap-6">
      <!-- Page Header -->
      <PageHeader title="OAuth 应用" subtitle="管理 OAuth 2.0 客户端应用和授权">
         <template #actions>
            <Button
               label="新建应用"
               icon="pi pi-plus"
               @click="openNewAppDialog" />
         </template>
      </PageHeader>

      <!-- Stats -->
      <MiniStats :items="stats" />

      <!-- Search -->
      <SearchBox v-model="searchQuery" />

      <!-- Apps Grid -->
      <div
         class="grid grid-cols-1 md:grid-cols-[repeat(auto-fill,minmax(340px,1fr))] gap-5">
         <OAuthAppCard
            v-for="app in filteredApps"
            :key="app.id"
            :app="app"
            @view="viewApp"
            @edit="editApp"
            @regenerateSecret="regenerateSecret" />

         <!-- Empty State -->
         <div
            v-if="filteredApps.length === 0"
            class="col-span-full flex flex-col items-center justify-center p-16 text-surface-400">
            <i class="pi pi-key text-5xl mb-4"></i>
            <p>未找到匹配的应用</p>
         </div>
      </div>

      <!-- App Dialog -->
      <OAuthAppDialog
         v-model:visible="appDialog"
         :isEditing="isEditing"
         :initialData="
            selectedApp
               ? {
                    name: selectedApp.name,
                    description: selectedApp.description,
                    redirectUris: selectedApp.redirectUris.join('\n'),
                    scopes: [...selectedApp.scopes],
                    trusted: selectedApp.trusted,
                 }
               : undefined
         "
         @save="saveApp" />

      <!-- Secret Dialog -->
      <SecretDialog v-model:visible="secretDialog" :secret="newSecret" />
   </div>
</template>
