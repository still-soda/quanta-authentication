<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import PageHeader from '@/components/shared/PageHeader.vue'
import SimpleStatCard from '@/components/shared/SimpleStatCard.vue'
import SearchBox from '@/components/shared/SearchBox.vue'
import OAuthAppCard from '@/components/oauth/OAuthAppCard.vue'
import OAuthAppDialog from '@/components/oauth/OAuthAppDialog.vue'
import SecretDialog from '@/components/oauth/SecretDialog.vue'
import DeleteConfirmDialog from '@/components/shared/DeleteConfirmDialog.vue'
import AppGroupPermissionDialog from '@/components/oauth/AppGroupPermissionDialog.vue'
import ConfirmDialog from 'primevue/confirmdialog'
import { useConfirm } from 'primevue/useconfirm'
import type { OAuthApp, OAuthAppFormData, SimpleStatData, ListOAuthAppsParams } from '@/types'
import {
   getOAuthApps,
   createOAuthApp,
   updateOAuthApp,
   deleteOAuthApp,
   regenerateClientSecret,
} from '@/apis/oauth'

const queryClient = useQueryClient()
const toast = useToast()
const confirm = useConfirm()

// 查询参数
const queryParams = ref<ListOAuthAppsParams>({
   page: 1,
   page_size: 50,
   search: '',
   status: undefined,
})

// 获取 OAuth 应用数据
const {
   data: appsData,
   isLoading,
   refetch,
} = useQuery({
   queryKey: ['oauth-apps', queryParams],
   queryFn: () => getOAuthApps(queryParams.value),
})

const apps = computed(() => appsData.value?.items || [])
const total = computed(() => appsData.value?.total || 0)

// 创建应用 mutation
const createAppMutation = useMutation({
   mutationFn: createOAuthApp,
   onSuccess: data => {
      queryClient.invalidateQueries({ queryKey: ['oauth-apps'] })
      appDialog.value = false
      // 显示新创建的密钥
      newSecret.value = data.client_secret
      secretDialog.value = true
      toast.add({
         severity: 'success',
         summary: '创建成功',
         detail: `OAuth 应用 "${data.client.name}" 已创建`,
         life: 3000,
      })
   },
   onError: (error: Error) => {
      toast.add({
         severity: 'error',
         summary: '创建失败',
         detail: error.message || '创建 OAuth 应用时发生错误',
         life: 5000,
      })
   },
})

// 更新应用 mutation
const updateAppMutation = useMutation({
   mutationFn: ({ id, data }: { id: string; data: Partial<OAuthAppFormData> }) =>
      updateOAuthApp(id, data),
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['oauth-apps'] })
      appDialog.value = false
      toast.add({
         severity: 'success',
         summary: '更新成功',
         detail: 'OAuth 应用信息已更新',
         life: 3000,
      })
   },
   onError: (error: Error) => {
      toast.add({
         severity: 'error',
         summary: '更新失败',
         detail: error.message || '更新 OAuth 应用时发生错误',
         life: 5000,
      })
   },
})

// 删除应用 mutation
const deleteAppMutation = useMutation({
   mutationFn: deleteOAuthApp,
   onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['oauth-apps'] })
      toast.add({
         severity: 'success',
         summary: '删除成功',
         detail: 'OAuth 应用已删除',
         life: 3000,
      })
   },
   onError: (error: Error) => {
      toast.add({
         severity: 'error',
         summary: '删除失败',
         detail: error.message || '删除 OAuth 应用时发生错误',
         life: 5000,
      })
   },
})

// 重新生成密钥 mutation
const regenerateSecretMutation = useMutation({
   mutationFn: regenerateClientSecret,
   onSuccess: data => {
      newSecret.value = data.secret
      secretDialog.value = true
      toast.add({
         severity: 'success',
         summary: '密钥已重新生成',
         detail: '请妥善保存新的客户端密钥',
         life: 3000,
      })
   },
   onError: (error: Error) => {
      toast.add({
         severity: 'error',
         summary: '操作失败',
         detail: error.message || '重新生成密钥时发生错误',
         life: 5000,
      })
   },
})

const selectedApp = ref<OAuthApp | null>(null)
const appDialog = ref(false)
const secretDialog = ref(false)
const viewDialog = ref(false)
const deleteDialog = ref(false)
const appGroupDialog = ref(false)
const appToDelete = ref<OAuthApp | null>(null)
const isEditing = ref(false)
const searchQuery = ref('')
const newSecret = ref('')

// 搜索防抖
let searchTimeout: ReturnType<typeof setTimeout>
watch(searchQuery, value => {
   clearTimeout(searchTimeout)
   searchTimeout = setTimeout(() => {
      queryParams.value.search = value
      queryParams.value.page = 1
   }, 300)
})

const filteredApps = computed(() => {
   return apps.value
})

const stats = computed<SimpleStatData[]>(() => {
   const appList = apps.value || []
   return [
      {
         title: '总应用',
         value: total.value,
         icon: 'pi pi-th-large',
         color: 'blue',
      },
      {
         title: '生产环境',
         value: appList.filter(a => a.status === 'active').length,
         icon: 'pi pi-check-circle',
         color: 'green',
      },
      {
         title: '开发中',
         value: appList.filter(a => a.status === 'development').length,
         icon: 'pi pi-code',
         color: 'orange',
      },
      {
         title: '已弃用',
         value: appList.filter(a => a.status === 'deprecated').length,
         icon: 'pi pi-history',
         color: 'gray',
      },
   ]
})

const openNewAppDialog = () => {
   isEditing.value = false
   selectedApp.value = null
   appDialog.value = true
}

const editApp = (app: OAuthApp) => {
   isEditing.value = true
   selectedApp.value = app
   appDialog.value = true
}

const viewApp = (app: OAuthApp) => {
   selectedApp.value = app
   viewDialog.value = true
}

const confirmDeleteApp = (app: OAuthApp) => {
   appToDelete.value = app
   deleteDialog.value = true
}

const handleDeleteConfirm = () => {
   if (appToDelete.value) {
      deleteAppMutation.mutate(appToDelete.value.id)
      appToDelete.value = null
   }
}

const regenerateSecret = (app: OAuthApp) => {
   confirm.require({
      message: `确定要重新生成 "${app.name}" 的客户端密钥吗？旧密钥将立即失效。`,
      header: '重新生成密钥',
      icon: 'pi pi-exclamation-triangle',
      rejectLabel: '取消',
      acceptLabel: '确认',
      acceptClass: 'p-button-warning',
      accept: () => {
         selectedApp.value = app
         regenerateSecretMutation.mutate(app.id)
      },
   })
}

const manageAppGroupPermissions = (app: OAuthApp) => {
   selectedApp.value = app
   appGroupDialog.value = true
}

const saveApp = (data: OAuthAppFormData) => {
   if (isEditing.value && selectedApp.value) {
      updateAppMutation.mutate({ id: selectedApp.value.id, data })
   } else {
      createAppMutation.mutate(data)
   }
}
</script>

<template>
   <div class="flex flex-col gap-6">
      <!-- Page Header -->
      <PageHeader title="OAuth 应用" subtitle="管理 OAuth 2.0 客户端应用和授权">
         <template #actions>
            <Button label="新建应用" icon="pi pi-plus" @click="openNewAppDialog" />
         </template>
      </PageHeader>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4">
         <template v-if="isLoading">
            <div
               v-for="i in 4"
               :key="i"
               class="h-20 bg-surface-100 dark:bg-surface-800 rounded-xl animate-pulse"
            />
         </template>
         <template v-else>
            <SimpleStatCard v-for="stat in stats" :key="stat.title" :stat="stat" />
         </template>
      </div>

      <!-- Search -->
      <SearchBox v-model="searchQuery" placeholder="搜索应用名称或 Client ID..." />

      <!-- Apps Grid -->
      <div
         v-if="isLoading"
         class="grid grid-cols-1 md:grid-cols-[repeat(auto-fill,minmax(340px,1fr))] gap-5"
      >
         <div
            v-for="i in 5"
            :key="i"
            class="h-64 bg-surface-100 dark:bg-surface-800 rounded-xl animate-pulse"
         />
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-[repeat(auto-fill,minmax(340px,1fr))] gap-5">
         <OAuthAppCard
            v-for="app in filteredApps"
            :key="app.id"
            :app="app"
            @view="viewApp"
            @edit="editApp"
            @delete="confirmDeleteApp"
            @regenerateSecret="regenerateSecret"
            @managePermissions="manageAppGroupPermissions"
         />

         <!-- Empty State -->
         <div
            v-if="filteredApps.length === 0"
            class="col-span-full flex flex-col items-center justify-center p-16 text-surface-400"
         >
            <i class="pi pi-key text-5xl mb-4"></i>
            <p>{{ searchQuery ? '未找到匹配的应用' : '暂无 OAuth 应用，点击"新建应用"创建' }}</p>
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
                    domain: selectedApp.domain,
                    redirect_uris: selectedApp.redirect_uris,
                    scopes: [...selectedApp.scopes],
                    grant_types: [...selectedApp.grant_types],
                    status: selectedApp.status,
                    trusted: selectedApp.trusted,
                    logo: selectedApp.logo,
                    icon: selectedApp.icon,
                    icon_bg: selectedApp.icon_bg,
                 }
               : undefined
         "
         @save="saveApp"
      />

      <!-- View App Dialog -->
      <Dialog
         v-model:visible="viewDialog"
         :header="selectedApp?.name || '应用详情'"
         modal
         :style="{ width: '32rem' }"
         :breakpoints="{ '640px': '90vw' }"
      >
         <div v-if="selectedApp" class="flex flex-col gap-4">
            <div
               class="flex items-center gap-4 pb-4 border-b border-surface-100 dark:border-surface-800"
            >
               <div
                  class="w-16 h-16 flex items-center justify-center rounded-xl text-white text-2xl shadow-lg overflow-hidden"
                  :style="{ background: selectedApp.logo ? 'transparent' : selectedApp.icon_bg }"
               >
                  <img
                     v-if="selectedApp.logo"
                     :src="`/api/uploads/${selectedApp.logo}`"
                     alt="App Logo"
                     class="w-full h-full object-cover"
                  />
                  <i v-else :class="selectedApp.icon"></i>
               </div>
               <div class="flex flex-col">
                  <h3 class="text-lg font-semibold text-surface-900 dark:text-surface-100 m-0">
                     {{ selectedApp.name }}
                  </h3>
                  <p class="text-sm text-surface-500 m-0">
                     {{ selectedApp.description || '暂无描述' }}
                  </p>
               </div>
            </div>

            <div class="grid grid-cols-2 gap-4 text-sm">
               <div class="flex flex-col gap-1">
                  <label class="text-xs text-surface-400 uppercase tracking-wider">Client ID</label>
                  <code
                     class="text-surface-700 dark:text-surface-300 font-mono text-xs break-all"
                     >{{ selectedApp.client_id }}</code
                  >
               </div>
               <div class="flex flex-col gap-1">
                  <label class="text-xs text-surface-400 uppercase tracking-wider">状态</label>
                  <span class="text-surface-700 dark:text-surface-300">{{
                     selectedApp.status === 'active'
                        ? '生产环境'
                        : selectedApp.status === 'development'
                          ? '开发中'
                          : '已弃用'
                  }}</span>
               </div>
               <div class="flex flex-col gap-1">
                  <label class="text-xs text-surface-400 uppercase tracking-wider">域名</label>
                  <span class="text-surface-700 dark:text-surface-300">{{
                     selectedApp.domain
                  }}</span>
               </div>
               <div class="flex flex-col gap-1">
                  <label class="text-xs text-surface-400 uppercase tracking-wider">可信应用</label>
                  <span class="text-surface-700 dark:text-surface-300">{{
                     selectedApp.trusted ? '是' : '否'
                  }}</span>
               </div>
               <div class="flex flex-col gap-1 col-span-2">
                  <label class="text-xs text-surface-400 uppercase tracking-wider"
                     >重定向 URI</label
                  >
                  <div class="flex flex-col gap-1">
                     <code
                        v-for="uri in selectedApp.redirect_uris"
                        :key="uri"
                        class="text-surface-700 dark:text-surface-300 font-mono text-xs break-all"
                        >{{ uri }}</code
                     >
                     <span v-if="!selectedApp.redirect_uris?.length" class="text-surface-400"
                        >未配置</span
                     >
                  </div>
               </div>
               <div class="flex flex-col gap-1">
                  <label class="text-xs text-surface-400 uppercase tracking-wider">授权范围</label>
                  <div class="flex flex-wrap gap-1">
                     <span
                        v-for="scope in selectedApp.scopes"
                        :key="scope"
                        class="px-2 py-0.5 bg-surface-100 dark:bg-surface-700 rounded text-xs"
                        >{{ scope }}</span
                     >
                  </div>
               </div>
               <div class="flex flex-col gap-1">
                  <label class="text-xs text-surface-400 uppercase tracking-wider">授权类型</label>
                  <div class="flex flex-wrap gap-1">
                     <span
                        v-for="gt in selectedApp.grant_types"
                        :key="gt"
                        class="px-2 py-0.5 bg-surface-100 dark:bg-surface-700 rounded text-xs"
                        >{{ gt }}</span
                     >
                  </div>
               </div>
               <div class="flex flex-col gap-1">
                  <label class="text-xs text-surface-400 uppercase tracking-wider">请求次数</label>
                  <span class="text-surface-700 dark:text-surface-300">{{
                     selectedApp.request_count
                  }}</span>
               </div>
               <div class="flex flex-col gap-1">
                  <label class="text-xs text-surface-400 uppercase tracking-wider">最后使用</label>
                  <span class="text-surface-700 dark:text-surface-300">{{
                     selectedApp.last_used_at
                        ? new Date(selectedApp.last_used_at).toLocaleString('zh-CN')
                        : '从未使用'
                  }}</span>
               </div>
            </div>
         </div>
         <template #footer>
            <Button label="关闭" severity="secondary" @click="viewDialog = false" />
         </template>
      </Dialog>

      <!-- Secret Dialog -->
      <SecretDialog v-model:visible="secretDialog" :secret="newSecret" />

      <!-- Delete Confirm Dialog -->
      <DeleteConfirmDialog
         v-model:visible="deleteDialog"
         :itemName="appToDelete?.name || ''"
         itemType="OAuth 应用"
         title="删除 OAuth 应用"
         @confirm="handleDeleteConfirm"
      />

      <!-- App Group Permission Dialog -->
      <AppGroupPermissionDialog v-model:visible="appGroupDialog" :client="selectedApp" />

      <!-- Confirm Dialog -->
      <ConfirmDialog />
   </div>
</template>
