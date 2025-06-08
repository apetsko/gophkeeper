<script setup>
import { storeToRefs } from 'pinia'
import { RouterLink, RouterView } from 'vue-router'
import { useAuthStore } from '@/stores/auth.store.js'

const authStore = useAuthStore()
const { user } = storeToRefs(authStore)
</script>

<template>
  <header>
    <div class="wrapper">
      <nav v-show="user" class="navbar navbar-expand navbar-dark bg-dark">
        <div class="navbar-nav flex justify-between">
          <div>
            <RouterLink to="/">Главная</RouterLink>
            <RouterLink to="/credentials">Логин/пароль</RouterLink>
            <RouterLink to="/bank-card">Банковская карта</RouterLink>
            <RouterLink to="/binary-data">Файл</RouterLink>
          </div>
          <div>
            <a @click="authStore.logout()" class="nav-item nav-link">Выйти</a>
          </div>
        </div>
      </nav>
      <div class="container pt-4 pb-4">
        <RouterView />
      </div>
    </div>
  </header>
</template>

<style scoped>
header {
  line-height: 1.5;
  width: 100%;
}

nav {
  width: 100%;
  font-size: 12px;
  text-align: center;
  margin-top: 2rem;
}

nav a.router-link-exact-active {
  color: var(--color-text);
}

nav a.router-link-exact-active:hover {
  background-color: transparent;
}

nav a {
  display: inline-block;
  padding: 0 1rem;
  border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
  border: 0;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
    width: 100%;
  }

  nav {
    margin-left: -1rem;
    font-size: 1rem;

    padding: 1rem 0;
    margin-top: 1rem;
  }
}
</style>
