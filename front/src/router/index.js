import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth.store.js'
import HomeView from '../views/HomeView.vue'
import Credentials from '../views/Credentials.vue'
import BankCard from '../views/BankCard.vue'
import BinaryData from '../views/BinaryData.vue'
import LoginPage from '../views/LoginPage.vue'
import SignupPage from '../views/SignupPage.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: LoginPage,
    },
    {
      path: '/signup',
      name: 'Signup',
      component: SignupPage,
    },
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { requiresAuth: true },
    },
    {
      path: '/credentials',
      name: 'credentials',
      component: Credentials,
      meta: { requiresAuth: true },
    },
    {
      path: '/bank-card',
      name: 'bank-card',
      component: BankCard,
      meta: { requiresAuth: true },
    },
    {
      path: '/binary-data',
      name: 'binary-data',
      component: BinaryData,
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach((to, from, next) => {
  const auth = useAuthStore();

  if (to.meta.requiresAuth && !auth.user) {
    auth.returnUrl = to.fullPath;
    next('/login')
  } else {
    next()
  }
})

export default router
