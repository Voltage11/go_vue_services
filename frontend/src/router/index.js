import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue'),
      meta: { requiresAuth: true }
    },
    // {
    //   path: '/login',
    //   name: 'login',
    //   component: () => import('@/views/LoginView.vue'),
    //   meta: { requiresGuest: true }
    // },
    // {
    //   path: '/register',
    //   name: 'register',
    //   component: () => import('@/views/RegisterView.vue'),
    //   meta: { requiresGuest: true }
    // }
  ]
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.meta.requiresGuest && authStore.isAuthenticated) {
    next('/')
  } else {
    next()
  }
})

export default router
