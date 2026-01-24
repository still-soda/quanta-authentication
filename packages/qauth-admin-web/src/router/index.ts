import { createRouter, createWebHistory } from 'vue-router'
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import AuthLayout from '@/layouts/AuthLayout.vue'
import { APP_NAME } from '@/config'

// Token 检查函数
const isAuthenticated = (): boolean => {
   const token = localStorage.getItem('access_token')
   return !!token
}

const router = createRouter({
   history: createWebHistory(import.meta.env.BASE_URL),
   routes: [
      // Auth routes (public)
      {
         path: '/auth',
         component: AuthLayout,
         meta: { requiresAuth: false },
         children: [
            {
               path: '',
               redirect: '/auth/login',
            },
            {
               path: 'login',
               name: 'login',
               component: () => import('@/views/LoginView.vue'),
               meta: { title: '登录', requiresAuth: false },
            },
            {
               path: 'register',
               name: 'register',
               component: () => import('@/views/RegisterView.vue'),
               meta: { title: '注册', requiresAuth: false },
            },
            {
               path: 'forgot-password',
               name: 'forgot-password',
               component: () => import('@/views/ForgotPasswordView.vue'),
               meta: { title: '忘记密码', requiresAuth: false },
            },
         ],
      },
      // Main app routes (protected)
      {
         path: '/',
         component: DefaultLayout,
         meta: { requiresAuth: true },
         children: [
            {
               path: '',
               name: 'dashboard',
               component: () => import('@/views/DashboardView.vue'),
               meta: { title: '仪表盘', requiresAuth: true },
            },
            {
               path: 'users',
               name: 'users',
               component: () => import('@/views/UsersView.vue'),
               meta: { title: '用户管理', requiresAuth: true },
            },
            {
               path: 'roles',
               name: 'roles',
               component: () => import('@/views/RolesView.vue'),
               meta: { title: '角色权限', requiresAuth: true },
            },
            {
               path: 'oauth',
               name: 'oauth',
               component: () => import('@/views/OAuthView.vue'),
               meta: { title: 'OAuth 应用', requiresAuth: true },
            },
            {
               path: 'organizations',
               name: 'organizations',
               component: () => import('@/views/OrganizationsView.vue'),
               meta: { title: '组织架构', requiresAuth: true },
            },
            {
               path: 'audit',
               name: 'audit',
               component: () => import('@/views/AuditView.vue'),
               meta: { title: '审计日志', requiresAuth: true },
            },
            {
               path: 'settings',
               name: 'settings',
               component: () => import('@/views/SettingsView.vue'),
               meta: { title: '系统设置', requiresAuth: true },
            },
            {
               path: 'profile',
               name: 'profile',
               component: () => import('@/views/ProfileView.vue'),
               meta: { title: '个人资料', requiresAuth: true },
            },
            {
               path: 'notifications',
               name: 'notifications',
               component: () => import('@/views/NotificationsView.vue'),
               meta: { title: '全部通知', requiresAuth: true },
            },
         ],
      },
   ],
})

// 路由守卫
router.beforeEach((to, _from, next) => {
   // 设置页面标题
   const title = to.meta?.title as string
   document.title = title ? `${title} - ${APP_NAME}` : APP_NAME

   // 认证检查
   const requiresAuth = to.meta?.requiresAuth !== false

   if (requiresAuth && !isAuthenticated()) {
      // 需要认证但未登录，跳转到登录页
      next({
         path: '/auth/login',
         query: { redirect: to.fullPath },
      })
   } else if (!requiresAuth && isAuthenticated() && to.path.startsWith('/auth')) {
      // 已登录但访问登录相关页面，跳转到首页
      next({ path: '/' })
   } else {
      next()
   }
})

export default router
