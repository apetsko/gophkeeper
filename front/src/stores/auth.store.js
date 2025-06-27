import { defineStore } from 'pinia';
import router from '@/router/index.js'

import { fetchWrapper } from '@/helpers/fetch-wrapper.js';

const baseUrl = `${import.meta.env.VITE_API_URL}/v1`;

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('user')),
    returnUrl: null
  }),
  actions: {
    async login(username, password) {
      const user = await fetchWrapper.post(`${baseUrl}/login`, { username, password });

      this.user = user;

      localStorage.setItem('user', JSON.stringify(user));

      router.push(this.returnUrl || '/');
    },
    async register(username, password) {
      try {
        const user = await fetchWrapper.post(`${baseUrl}/signup`, { username, password });

        this.user = user;

        localStorage.setItem('user', JSON.stringify(user));

        router.push(this.returnUrl || '/');
      } catch (error) {
        throw new Error(error.response?.data?.message || 'Registration failed');
      }
    },
    logout() {
      this.user = null;
      localStorage.removeItem('user');
      router.push('/login');
    }
  }
});
