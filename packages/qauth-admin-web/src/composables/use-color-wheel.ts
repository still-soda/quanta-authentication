/**
 * 色轮颜色生成工具
 *
 * 使用 HSL 色彩空间进行颜色偏转，
 * 选定美观的饱和度和亮度范围，通过色相偏转产生足够多的不重复颜色。
 */

// 预设的美观基础色调（HSL 色相值，0-360）
// 这些色调经过精心挑选，视觉上更加和谐
const BEAUTIFUL_BASE_HUES = [
   210, // 蓝色
   160, // 青绿色
   260, // 紫色
   340, // 粉红色
   35, // 橙色
   120, // 绿色
   180, // 蓝绿色
   290, // 紫罗兰
   20, // 珊瑚色
   90, // 黄绿色
] as const

// 美观的饱和度范围（适中的饱和度看起来更舒适）
const SATURATION_RANGE = { min: 55, max: 70 }

// 亮度范围（避免太亮或太暗）
const LIGHTNESS_LIGHT = { min: 55, max: 65 }
const LIGHTNESS_DARK = { min: 50, max: 60 }

/**
 * HSL 转 Hex 颜色
 */
function hslToHex(h: number, s: number, l: number): string {
   s = s / 100
   l = l / 100
   const a = s * Math.min(l, 1 - l)
   const f = (n: number) => {
      const k = (n + h / 30) % 12
      const color = l - a * Math.max(Math.min(k - 3, 9 - k, 1), -1)
      return Math.round(255 * color)
         .toString(16)
         .padStart(2, '0')
   }
   return `#${f(0)}${f(8)}${f(4)}`
}

/**
 * 基于索引生成唯一的色相值
 * 使用黄金角度偏转确保颜色分布均匀
 */
function generateHue(index: number, totalColors: number): number {
   // 如果颜色数量在预设范围内，使用预设色调
   if (index < BEAUTIFUL_BASE_HUES.length && totalColors <= BEAUTIFUL_BASE_HUES.length) {
      return BEAUTIFUL_BASE_HUES[index] ?? 0
   }

   // 对于更多颜色，使用黄金角度（约 137.5°）进行色相偏转
   // 黄金角度能确保颜色在色轮上均匀分布且不重复
   const goldenAngle = 137.508
   const baseHue = BEAUTIFUL_BASE_HUES[index % BEAUTIFUL_BASE_HUES.length] ?? 0
   const offset = Math.floor(index / BEAUTIFUL_BASE_HUES.length) * goldenAngle
   return (baseHue + offset) % 360
}

/**
 * 生成美观的颜色数组
 *
 * @param count - 需要的颜色数量
 * @param isDarkMode - 是否为深色模式
 * @returns 颜色数组（Hex 格式）
 */
export function generateBeautifulColors(count: number, isDarkMode = false): string[] {
   const colors: string[] = []
   const lightnessRange = isDarkMode ? LIGHTNESS_DARK : LIGHTNESS_LIGHT

   for (let i = 0; i < count; i++) {
      const hue = generateHue(i, count)

      // 稍微变化饱和度和亮度，增加层次感
      const saturationOffset = ((i % 3) - 1) * 5
      const lightnessOffset = (i % 2) * 5

      const saturation = Math.min(
         SATURATION_RANGE.max,
         Math.max(
            SATURATION_RANGE.min,
            (SATURATION_RANGE.min + SATURATION_RANGE.max) / 2 + saturationOffset
         )
      )

      const lightness = Math.min(
         lightnessRange.max,
         Math.max(
            lightnessRange.min,
            (lightnessRange.min + lightnessRange.max) / 2 + lightnessOffset
         )
      )

      colors.push(hslToHex(hue, saturation, lightness))
   }

   return colors
}

/**
 * 为特定角色代码生成稳定的颜色
 * 使用角色代码的哈希值确保同一角色始终获得相同颜色
 *
 * @param roleCode - 角色代码
 * @param isDarkMode - 是否为深色模式
 * @returns Hex 颜色值
 */
export function getColorForRole(roleCode: string, isDarkMode = false): string {
   // 计算角色代码的简单哈希
   let hash = 0
   for (let i = 0; i < roleCode.length; i++) {
      const char = roleCode.charCodeAt(i)
      hash = (hash << 5) - hash + char
      hash = hash & hash // 转换为32位整数
   }

   // 使用哈希值计算色相
   const hue = Math.abs(hash) % 360

   const lightnessRange = isDarkMode ? LIGHTNESS_DARK : LIGHTNESS_LIGHT
   const saturation = (SATURATION_RANGE.min + SATURATION_RANGE.max) / 2
   const lightness = (lightnessRange.min + lightnessRange.max) / 2

   return hslToHex(hue, saturation, lightness)
}

/**
 * 为多个角色生成颜色映射
 *
 * @param roleCodes - 角色代码数组
 * @param isDarkMode - 是否为深色模式
 * @returns 角色代码到颜色的映射
 */
export function generateRoleColorMap(
   roleCodes: string[],
   isDarkMode = false
): Record<string, string> {
   const colors = generateBeautifulColors(roleCodes.length, isDarkMode)
   const colorMap: Record<string, string> = {}

   roleCodes.forEach((code, index) => {
      colorMap[code] = colors[index] ?? '#6b7280'
   })

   return colorMap
}
