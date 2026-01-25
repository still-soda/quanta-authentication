<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import InputSwitch from 'primevue/inputswitch'
import Select from 'primevue/select'
import FileUpload from 'primevue/fileupload'
import type { OAuthAppFormData, OAuthClientStatus, OAuthOption } from '@/types'
import { OAUTH_STATUS_OPTIONS } from '@/config'
import { getOAuthClientOptions } from '@/apis/oauth'
import { httpClient } from '@/apis'

const props = defineProps<{
   visible: boolean
   isEditing: boolean
   initialData?: Partial<OAuthAppFormData>
}>()

const emit = defineEmits<{
   (e: 'update:visible', value: boolean): void
   (e: 'save', data: OAuthAppFormData): void
}>()

// 从后端获取可选配置
const { data: clientOptions, isLoading: optionsLoading } = useQuery({
   queryKey: ['oauth-client-options'],
   queryFn: getOAuthClientOptions,
   staleTime: 5 * 60 * 1000, // 5分钟缓存
})

// 授权范围和授权类型选项
const scopeOptions = computed<OAuthOption[]>(() => clientOptions.value?.scopes || [])
const grantTypeOptions = computed<OAuthOption[]>(() => clientOptions.value?.grant_types || [])

const defaultFormData: OAuthAppFormData = {
   name: '',
   description: '',
   domain: '',
   redirect_uris: [],
   scopes: ['openid', 'profile'],
   grant_types: ['authorization_code', 'refresh_token'],
   status: 'development',
   trusted: false,
   logo: '',
}

const appForm = ref<OAuthAppFormData>({ ...defaultFormData })
const redirectUrisText = ref('')
const logoPreview = ref<string>('')
const isUploading = ref(false)

const dialogVisible = computed({
   get: () => props.visible,
   set: value => emit('update:visible', value),
})

// 当 dialog 打开时重置表单
watch(
   () => props.visible,
   visible => {
      if (visible) {
         if (props.initialData) {
            appForm.value = {
               ...defaultFormData,
               ...props.initialData,
            }
            redirectUrisText.value = props.initialData.redirect_uris?.join('\n') || ''
            logoPreview.value = props.initialData.logo || ''
         } else {
            appForm.value = { ...defaultFormData }
            redirectUrisText.value = ''
            logoPreview.value = ''
         }
      }
   }
)

const toggleScope = (scope: string) => {
   const index = appForm.value.scopes.indexOf(scope)
   if (index > -1) {
      appForm.value.scopes.splice(index, 1)
   } else {
      appForm.value.scopes.push(scope)
   }
}

const toggleGrantType = (grantType: string) => {
   const index = appForm.value.grant_types.indexOf(grantType)
   if (index > -1) {
      appForm.value.grant_types.splice(index, 1)
   } else {
      appForm.value.grant_types.push(grantType)
   }
}

const isFormValid = computed(() => {
   return (
      appForm.value.name.trim() !== '' &&
      appForm.value.domain.trim() !== '' &&
      appForm.value.scopes.length > 0 &&
      appForm.value.grant_types.length > 0
   )
})

// 处理 Logo 上传
const handleLogoUpload = async (event: { files: File[] }) => {
   const file = event.files[0]
   if (!file) return

   isUploading.value = true
   try {
      const formData = new FormData()
      formData.append('file', file)

      const response = await httpClient.post('/_/v1/resources/upload', formData, {
         headers: { 'Content-Type': 'multipart/form-data' },
      })

      const fileName = response.data.data.fileName
      appForm.value.logo = fileName
      logoPreview.value = URL.createObjectURL(file)
   } catch (error) {
      console.error('Upload failed:', error)
   } finally {
      isUploading.value = false
   }
}

// 移除 Logo
const removeLogo = () => {
   appForm.value.logo = ''
   logoPreview.value = ''
}

const saveApp = () => {
   // 将文本框中的 URI 转换为数组
   const redirectUris = redirectUrisText.value
      .split('\n')
      .map(uri => uri.trim())
      .filter(uri => uri !== '')

   emit('save', {
      ...appForm.value,
      redirect_uris: redirectUris,
   })
}

const resetForm = (data?: Partial<OAuthAppFormData>) => {
   appForm.value = {
      ...defaultFormData,
      ...data,
   }
   redirectUrisText.value = data?.redirect_uris?.join('\n') || ''
   logoPreview.value = data?.logo || ''
}

defineExpose({ resetForm })
</script>

<template>
   <Dialog
      v-model:visible="dialogVisible"
      :header="isEditing ? '编辑应用' : '新建应用'"
      modal
      :style="{ width: '36rem' }"
      :breakpoints="{ '640px': '90vw' }"
   >
      <div class="flex flex-col gap-5 py-2">
         <!-- Logo 上传 -->
         <div class="flex flex-col gap-2">
            <label class="text-sm font-medium text-surface-700 dark:text-surface-300">
               应用 Logo
            </label>
            <div class="flex items-center gap-4">
               <div
                  class="w-16 h-16 flex items-center justify-center rounded-xl border-2 border-dashed border-surface-200 dark:border-surface-700 overflow-hidden bg-surface-50 dark:bg-surface-800"
               >
                  <img
                     v-if="logoPreview"
                     :src="logoPreview"
                     alt="Logo preview"
                     class="w-full h-full object-cover"
                  />
                  <i v-else class="pi pi-image text-2xl text-surface-400"></i>
               </div>
               <div class="flex flex-col gap-2">
                  <FileUpload
                     mode="basic"
                     accept="image/*"
                     :maxFileSize="2000000"
                     :auto="true"
                     chooseLabel="上传 Logo"
                     :disabled="isUploading"
                     @select="handleLogoUpload"
                     class="p-button-sm p-button-outlined"
                  />
                  <Button
                     v-if="logoPreview"
                     label="移除"
                     icon="pi pi-times"
                     severity="secondary"
                     text
                     size="small"
                     @click="removeLogo"
                  />
               </div>
            </div>
            <small class="text-xs text-surface-400">
               可选，支持 JPG、PNG 格式，最大 2MB。不上传时将显示默认图标
            </small>
         </div>

         <!-- 应用名称 -->
         <div class="flex flex-col gap-2">
            <label for="appName" class="text-sm font-medium text-surface-700 dark:text-surface-300">
               应用名称 <span class="text-red-500">*</span>
            </label>
            <InputText
               id="appName"
               v-model="appForm.name"
               placeholder="例如：My Application"
               class="w-full"
            />
         </div>

         <!-- 应用描述 -->
         <div class="flex flex-col gap-2">
            <label for="appDesc" class="text-sm font-medium text-surface-700 dark:text-surface-300">
               应用描述
            </label>
            <Textarea
               id="appDesc"
               v-model="appForm.description"
               placeholder="描述应用用途..."
               rows="2"
               class="w-full"
            />
         </div>

         <!-- 域名 -->
         <div class="flex flex-col gap-2">
            <label for="domain" class="text-sm font-medium text-surface-700 dark:text-surface-300">
               域名 <span class="text-red-500">*</span>
            </label>
            <InputText
               id="domain"
               v-model="appForm.domain"
               placeholder="例如：https://example.com"
               class="w-full"
            />
            <small class="text-xs text-surface-400"> 应用的主域名 </small>
         </div>

         <!-- 重定向 URI -->
         <div class="flex flex-col gap-2">
            <label
               for="redirectUris"
               class="text-sm font-medium text-surface-700 dark:text-surface-300"
            >
               重定向 URI
            </label>
            <Textarea
               id="redirectUris"
               v-model="redirectUrisText"
               placeholder="每行一个 URI&#10;https://example.com/callback&#10;myapp://callback"
               rows="3"
               class="w-full font-mono text-sm"
            />
            <small class="text-xs text-surface-400"> 授权完成后的回调地址，每行一个 </small>
         </div>

         <!-- 授权类型 -->
         <div class="flex flex-col gap-2">
            <label class="text-sm font-medium text-surface-700 dark:text-surface-300">
               授权类型 <span class="text-red-500">*</span>
            </label>
            <div v-if="optionsLoading" class="text-surface-400 text-sm">加载中...</div>
            <div v-else class="flex flex-wrap gap-2">
               <button
                  v-for="grantType in grantTypeOptions"
                  :key="grantType.value"
                  type="button"
                  class="py-1.5 px-3 border border-surface-200 dark:border-surface-700 rounded-full bg-transparent text-surface-600 dark:text-surface-400 text-[0.8125rem] cursor-pointer transition-all duration-200 ease hover:border-primary-300 hover:text-primary-600 dark:hover:border-primary-400 dark:hover:text-primary-400"
                  :class="{
                     'bg-primary-500! border-primary-500! text-white!':
                        appForm.grant_types.includes(grantType.value),
                  }"
                  @click="toggleGrantType(grantType.value)"
               >
                  {{ grantType.label }}
               </button>
            </div>
         </div>

         <!-- 授权范围 -->
         <div class="flex flex-col gap-2">
            <label class="text-sm font-medium text-surface-700 dark:text-surface-300">
               授权范围 <span class="text-red-500">*</span>
            </label>
            <div v-if="optionsLoading" class="text-surface-400 text-sm">加载中...</div>
            <div v-else class="flex flex-wrap gap-2">
               <button
                  v-for="scope in scopeOptions"
                  :key="scope.value"
                  type="button"
                  class="py-1.5 px-3 border border-surface-200 dark:border-surface-700 rounded-full bg-transparent text-surface-600 dark:text-surface-400 text-[0.8125rem] cursor-pointer transition-all duration-200 ease hover:border-primary-300 hover:text-primary-600 dark:hover:border-primary-400 dark:hover:text-primary-400"
                  :class="{
                     'bg-primary-500! border-primary-500! text-white!': appForm.scopes.includes(
                        scope.value
                     ),
                  }"
                  @click="toggleScope(scope.value)"
               >
                  {{ scope.label }}
               </button>
            </div>
         </div>

         <!-- 状态 -->
         <div class="flex flex-col gap-2">
            <label for="status" class="text-sm font-medium text-surface-700 dark:text-surface-300">
               状态
            </label>
            <Select
               id="status"
               v-model="appForm.status"
               :options="OAUTH_STATUS_OPTIONS"
               optionLabel="label"
               optionValue="value"
               placeholder="选择状态"
               class="w-full"
            />
         </div>

         <!-- 可信应用 -->
         <div class="flex flex-row justify-between items-center gap-2">
            <label for="trusted" class="text-sm font-medium text-surface-700 dark:text-surface-300">
               可信应用
            </label>
            <div class="flex items-center gap-3">
               <InputSwitch id="trusted" v-model="appForm.trusted" />
               <span class="text-[0.8125rem] text-surface-500">
                  {{ appForm.trusted ? '跳过用户授权确认' : '需要用户确认授权' }}
               </span>
            </div>
         </div>
      </div>

      <template #footer>
         <div class="flex justify-end gap-3">
            <Button label="取消" severity="secondary" outlined @click="dialogVisible = false" />
            <Button
               :label="isEditing ? '保存' : '创建'"
               :disabled="!isFormValid"
               @click="saveApp"
            />
         </div>
      </template>
   </Dialog>
</template>
