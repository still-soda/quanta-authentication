/**
 * API 模块入口
 * 所有 API 方法返回 mock 数据（100ms 延迟）
 */

// 模拟延迟函数
export const delay = (ms: number = 100) => new Promise(resolve => setTimeout(resolve, ms));

// 模拟 API 响应
export async function mockResponse<T>(data: T, delayMs: number = 100): Promise<T> {
  await delay(delayMs);
  return data;
}

// 导出所有 API 模块
export * from './dashboard';
export * from './users';
export * from './roles';
export * from './oauth';
export * from './audit';
export * from './notifications';
export * from './organizations';
export * from './profile';
export * from './settings';
