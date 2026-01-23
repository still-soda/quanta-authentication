<script setup lang="ts">
import { ref, computed } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import Button from 'primevue/button';
import PageHeader from '@/components/shared/PageHeader.vue';
import SimpleStatCard from '@/components/shared/SimpleStatCard.vue';
import SearchBox from '@/components/shared/SearchBox.vue';
import OAuthAppCard from '@/components/oauth/OAuthAppCard.vue';
import OAuthAppDialog from '@/components/oauth/OAuthAppDialog.vue';
import SecretDialog from '@/components/oauth/SecretDialog.vue';
import type { OAuthApp, OAuthAppFormData, SimpleStatData } from '@/types';
import {
   getOAuthApps,
   createOAuthApp,
   updateOAuthApp,
   regenerateClientSecret,
} from '@/apis/oauth';

const queryClient = useQueryClient();

// 使用 TanStack Query 获取 OAuth 应用数据
const { data: apps, isLoading } = useQuery({
   queryKey: ['oauth-apps'],
   queryFn: getOAuthApps,
});

// 创建应用 mutation
const createAppMutation = useMutation({
   mutationFn: createOAuthApp,
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['oauth-apps'] });
      appDialog.value = false;
   },
});

// 更新应用 mutation
const updateAppMutation = useMutation({
   mutationFn: ({ id, data }: { id: number; data: Partial<OAuthAppFormData> }) =>
      updateOAuthApp(id, data),
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['oauth-apps'] });
      appDialog.value = false;
   },
});

// 重新生成密钥 mutation
const regenerateSecretMutation = useMutation({
   mutationFn: regenerateClientSecret,
   onSuccess: (data) => {
      newSecret.value = data.secret;
      secretDialog.value = true;
   },
});

const selectedApp = ref<OAuthApp | null>(null);
const appDialog = ref(false);
const secretDialog = ref(false);
const isEditing = ref(false);
const searchQuery = ref('');
const newSecret = ref('');

const filteredApps = computed(() => {
   const appList = apps.value || [];
   if (!searchQuery.value) return appList;
   const query = searchQuery.value.toLowerCase();
   return appList.filter(
      (app) =>
         app.name.toLowerCase().includes(query) ||
         app.clientId.toLowerCase().includes(query),
   );
});

const stats = computed<SimpleStatData[]>(() => {
   const appList = apps.value || [];
   return [
      {
         title: '总应用',
         value: appList.length,
         icon: 'pi pi-th-large',
         color: 'blue',
      },
      {
         title: '生产环境',
         value: appList.filter((a) => a.status === 'active').length,
         icon: 'pi pi-check-circle',
         color: 'green',
      },
      {
         title: '开发中',
         value: appList.filter((a) => a.status === 'development').length,
         icon: 'pi pi-code',
         color: 'orange',
      },
      {
         title: '已弃用',
         value: appList.filter((a) => a.status === 'deprecated').length,
         icon: 'pi pi-history',
         color: 'gray',
      },
   ];
});

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
   regenerateSecretMutation.mutate(app.id);
};

const saveApp = (data: OAuthAppFormData) => {
   if (isEditing.value && selectedApp.value) {
      updateAppMutation.mutate({ id: selectedApp.value.id, data });
   } else {
      createAppMutation.mutate(data);
   }
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

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4">
         <template v-if="isLoading">
            <div
               v-for="i in 4"
               :key="i"
               class="h-20 bg-surface-100 dark:bg-surface-800 rounded-xl animate-pulse" />
         </template>
         <template v-else>
            <SimpleStatCard v-for="stat in stats" :key="stat.title" :stat="stat" />
         </template>
      </div>

      <!-- Search -->
      <SearchBox v-model="searchQuery" />

      <!-- Apps Grid -->
      <div
         v-if="isLoading"
         class="grid grid-cols-1 md:grid-cols-[repeat(auto-fill,minmax(340px,1fr))] gap-5">
         <div
            v-for="i in 5"
            :key="i"
            class="h-64 bg-surface-100 dark:bg-surface-800 rounded-xl animate-pulse" />
      </div>
      <div
         v-else
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
