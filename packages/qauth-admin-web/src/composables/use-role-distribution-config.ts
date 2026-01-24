/**
 * 角色分布图配置持久化管理
 */
import { ref, watch, computed } from 'vue'

// 本地存储 key
const STORAGE_KEY = 'qauth-role-distribution-config'

// 系统内置角色 codes
export const BUILTIN_ROLE_CODES = ['system_super_admin', 'system_admin', 'system_user'] as const

export type BuiltinRoleCode = (typeof BUILTIN_ROLE_CODES)[number]

export interface RoleDistributionConfig {
   selectedRoleCodes: string[]
}

// 默认配置：只显示系统内置角色
const defaultConfig: RoleDistributionConfig = {
   selectedRoleCodes: [...BUILTIN_ROLE_CODES],
}

/**
 * 从本地存储加载配置
 */
function loadConfig(): RoleDistributionConfig {
   try {
      const stored = localStorage.getItem(STORAGE_KEY)
      if (stored) {
         const parsed = JSON.parse(stored) as RoleDistributionConfig
         // 确保至少包含系统内置角色
         const selectedRoleCodes = parsed.selectedRoleCodes || []
         return { selectedRoleCodes }
      }
   } catch {
      console.warn('Failed to load role distribution config from localStorage')
   }
   return { ...defaultConfig }
}

/**
 * 保存配置到本地存储
 */
function saveConfig(config: RoleDistributionConfig): void {
   try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(config))
   } catch {
      console.warn('Failed to save role distribution config to localStorage')
   }
}

/**
 * 角色分布配置 composable
 */
export function useRoleDistributionConfig() {
   const config = ref<RoleDistributionConfig>(loadConfig())

   // 监听配置变化，自动保存
   watch(
      config,
      newConfig => {
         saveConfig(newConfig)
      },
      { deep: true }
   )

   /**
    * 获取选中的角色 codes
    */
   const selectedRoleCodes = computed({
      get: () => config.value.selectedRoleCodes,
      set: (value: string[]) => {
         config.value.selectedRoleCodes = value
      },
   })

   /**
    * 切换角色选中状态
    */
   function toggleRole(roleCode: string): void {
      const index = config.value.selectedRoleCodes.indexOf(roleCode)
      if (index === -1) {
         config.value.selectedRoleCodes.push(roleCode)
      } else {
         config.value.selectedRoleCodes.splice(index, 1)
      }
   }

   /**
    * 设置选中的角色
    */
   function setSelectedRoles(roleCodes: string[]): void {
      config.value.selectedRoleCodes = [...roleCodes]
   }

   /**
    * 重置为默认配置（只显示系统内置角色）
    */
   function resetToDefault(): void {
      config.value = { ...defaultConfig, selectedRoleCodes: [...defaultConfig.selectedRoleCodes] }
   }

   /**
    * 检查角色是否被选中
    */
   function isRoleSelected(roleCode: string): boolean {
      return config.value.selectedRoleCodes.includes(roleCode)
   }

   /**
    * 检查是否是系统内置角色
    */
   function isBuiltinRole(roleCode: string): boolean {
      return BUILTIN_ROLE_CODES.includes(roleCode as BuiltinRoleCode)
   }

   return {
      selectedRoleCodes,
      toggleRole,
      setSelectedRoles,
      resetToDefault,
      isRoleSelected,
      isBuiltinRole,
   }
}
