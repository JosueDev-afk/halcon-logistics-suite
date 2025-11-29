import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/track',
    },
    {
      path: '/track',
      name: 'track',
      component: () => import('@/views/PublicTracking.vue'),
      meta: { public: true },
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/Login.vue'),
      meta: { public: true },
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('@/views/Dashboard.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          redirect: '/dashboard/orders',
        },
        {
          path: 'orders',
          name: 'orders',
          component: () => import('@/views/orders/OrderList.vue'),
        },
        {
          path: 'orders/create',
          name: 'orders-create',
          component: () => import('@/views/orders/OrderCreate.vue'),
          meta: { roles: ['Sales'] },
        },
        {
          path: 'orders/:id',
          name: 'orders-detail',
          component: () => import('@/views/orders/OrderDetail.vue'),
        },
        {
          path: 'recycle-bin',
          name: 'recycle-bin',
          component: () => import('@/views/orders/RecycleBin.vue'),
          meta: { roles: ['Admin', 'Sales'] },
        },
        {
          path: 'users',
          name: 'users',
          component: () => import('@/views/users/UserList.vue'),
          meta: { roles: ['Admin'] },
        },
        {
          path: 'users/create',
          name: 'users-create',
          component: () => import('@/views/users/UserForm.vue'),
          meta: { roles: ['Admin'] },
        },
        {
          path: 'users/:id/edit',
          name: 'users-edit',
          component: () => import('@/views/users/UserForm.vue'),
          meta: { roles: ['Admin'] },
        },
      ],
    },
  ],
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth)
  const isPublic = to.matched.some((record) => record.meta.public)

  if (requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.path === '/login' && authStore.isAuthenticated) {
    next('/dashboard')
  } else if (to.meta.roles) {
    const allowedRoles = to.meta.roles as string[]
    if (authStore.userRole && allowedRoles.includes(authStore.userRole)) {
      next()
    } else {
      next('/dashboard')
    }
  } else {
    next()
  }
})

export default router
