import { createRouter, createWebHistory } from 'vue-router';
import DefaultLayout from '@/layouts/DefaultLayout.vue';
import { APP_NAME } from '@/config';

const router = createRouter({
   history: createWebHistory(import.meta.env.BASE_URL),
   routes: [
      {
         path: '/',
         component: DefaultLayout,
         children: [
            {
               path: '',
               name: 'dashboard',
               component: () => import('@/views/DashboardView.vue'),
               meta: { title: '仪表盘' },
            },
            {
               path: 'users',
               name: 'users',
               component: () => import('@/views/UsersView.vue'),
               meta: { title: '用户管理' },
            },
            {
               path: 'roles',
               name: 'roles',
               component: () => import('@/views/RolesView.vue'),
               meta: { title: '角色权限' },
            },
            {
               path: 'oauth',
               name: 'oauth',
               component: () => import('@/views/OAuthView.vue'),
               meta: { title: 'OAuth 应用' },
            },
            {
               path: 'organizations',
               name: 'organizations',
               component: () => import('@/views/OrganizationsView.vue'),
               meta: { title: '组织架构' },
            },
            {
               path: 'audit',
               name: 'audit',
               component: () => import('@/views/AuditView.vue'),
               meta: { title: '审计日志' },
            },
            {
               path: 'settings',
               name: 'settings',
               component: () => import('@/views/SettingsView.vue'),
               meta: { title: '系统设置' },
            },
            {
               path: 'profile',
               name: 'profile',
               component: () => import('@/views/ProfileView.vue'),
               meta: { title: '个人资料' },
            },
            {
               path: 'notifications',
               name: 'notifications',
               component: () => import('@/views/NotificationsView.vue'),
               meta: { title: '全部通知' },
            },
         ],
      },
   ],
});

router.beforeEach((to, _from, next) => {
   const title = to.meta?.title as string;
   document.title = title ? `${title} - ${APP_NAME}` : APP_NAME;
   next();
});

export default router;
