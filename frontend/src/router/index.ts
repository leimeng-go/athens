import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import AdminLayout from '@/layouts/AdminLayout.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    component: AdminLayout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '仪表盘', icon: 'Odometer' }
      },
      {
        path: 'dashboard/system',
        name: 'SystemStatus',
        component: () => import('../views/dashboard/SystemStatus.vue'),
        meta: { title: '系统状态', icon: 'Monitor' }
      },
      {
        path: 'dashboard/activities',
        name: 'Activities',
        component: () => import('../views/dashboard/Activities.vue'),
        meta: { title: '最近活动', icon: 'Bell' }
      },
      {
        path: 'repositories',
        name: 'Repositories',
        component: () => import('@/views/repositories/RepositoryList.vue'),
        meta: { title: '仓库管理', icon: 'Files' }
      },
      {
        path: 'upload',
        name: 'Upload',
        component: () => import('@/views/upload/UploadModule.vue'),
        meta: { title: '上传模块', icon: 'Upload' }
      },
      {
        path: 'download',
        name: 'Download',
        component: () => import('@/views/download/DownloadModule.vue'),
        meta: { title: '下载模块', icon: 'Download' }
      },
      {
        path: 'download/modules/:path',
        name: 'ModuleDetail',
        component: () => import('@/views/download/ModuleDetail.vue'),
        meta: { title: '模块详情', icon: 'InfoFilled' },
        props: true
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('@/views/Settings.vue'),
        meta: { title: '系统设置', icon: 'Setting' }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue')
  }
]

const router = createRouter({
  history: createWebHistory('/admin'),
  routes
})

// 全局前置守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = `${to.meta.title || 'Athens Admin'} - Athens Go Module Proxy`
  next()
})

export default router