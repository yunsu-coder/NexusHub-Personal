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
      path: '/code',
      name: 'Code',
      component: () => import('../views/CodeEditor.vue'),
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
      path: '/calculator',
      name: 'Calculator',
      component: () => import('../views/Calculator.vue'),
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
    }
  ]
})

// 路由守卫 - 可选的身份验证
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')

  // 如果已登录且访问欢迎页，重定向到仪表盘
  if (to.path === '/' && token) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router

