import { createRouter, createWebHistory } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Welcome',
      component: () => import('../views/Welcome.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: () => import('../views/Dashboard.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/files',
      name: 'Files',
      component: () => import('../views/FileManager.vue'),
      meta: { requiresAuth: false }
    },

    {
      path: '/notes',
      name: 'Notes',
      component: () => import('../views/Notes.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/data',
      name: 'DataAnalysis',
      component: () => import('../views/DataAnalysis.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/todos',
      name: 'TodoList',
      component: () => import('../views/TodoList.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/chat',
      name: 'Chat',
      component: () => import('../views/AIChat.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/settings',
      name: 'Settings',
      component: () => import('../views/Settings.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/collection',
      name: 'Collection',
      component: () => import('../views/Collection.vue'),
      meta: { requiresAuth: false }
    }
  ]
})

// 公共模式：不进行鉴权重定向
// 直接放行所有路由
router.beforeEach((to, from, next) => {
  next()
})

export default router
