<script setup>
import { Form, Field } from 'vee-validate';
import * as Yup from 'yup';
import { useAuthStore } from '@/stores/auth.store.js';

const schema = Yup.object().shape({
  username: Yup.string()
    .required('Username is required')
    .min(3, 'Username must be at least 3 characters'),
  password: Yup.string()
    .required('Password is required')
    .min(8, 'Password must be at least 8 characters'),
  confirmPassword: Yup.string()
    .required('Please confirm your password')
    .oneOf([Yup.ref('password')], 'Passwords must match')
});

function onSubmit(values, { setErrors }) {
  const authStore = useAuthStore();
  const { username, password } = values;

  return authStore.register(username, password)
    .catch(error => setErrors({ apiError: error }));
}
</script>

<template>
  <div class="flex justify-center flex-col w-100 m-auto">
    <div style="margin-bottom: 25px">
      <h2 class="text-3xl text-center">Регистрация</h2>
    </div>
    <div>
      <Form
        class="max-w-sm mx-auto"
        @submit="onSubmit"
        :validation-schema="schema"
        v-slot="{ errors, isSubmitting }"
      >
        <div>
          <label
            for="username"
            class="block text-sm font-medium text-gray-900 dark:text-white"
            style="margin-bottom: 10px;"
          >
            Username
          </label>
          <Field
            name="username"
            type="text"
            id="username"
            class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            required
            :class="{ 'is-invalid': errors.username }"
          />
          <div class="invalid-feedback">{{ errors.username }}</div>
        </div>

        <div style="margin-top: 15px;">
          <label
            for="password"
            class="block text-sm font-medium text-gray-900 dark:text-white"
            style="margin-bottom: 10px;"
          >
            Password
          </label>
          <Field
            name="password"
            type="password"
            id="password"
            class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            required
            :class="{ 'is-invalid': errors.password }"
          />
          <div class="invalid-feedback">{{ errors.password }}</div>
        </div>

        <div style="margin-top: 15px;">
          <label
            for="confirmPassword"
            class="block text-sm font-medium text-gray-900 dark:text-white"
            style="margin-bottom: 10px;"
          >
            Confirm Password
          </label>
          <Field
            name="confirmPassword"
            type="password"
            id="confirmPassword"
            class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            required
            :class="{ 'is-invalid': errors.confirmPassword }"
          />
          <div class="invalid-feedback">{{ errors.confirmPassword }}</div>
        </div>

        <div class="form-group text-center" style="margin-top: 25px;">
          <button
            class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
            type="submit"
            :disabled="isSubmitting"
          >
            <span v-show="isSubmitting" class="spinner-border spinner-border-sm mr-1"></span>
            Зарегистрироваться
          </button>
        </div>

        <div v-if="errors.apiError" class="alert alert-danger mt-3 text-red-500 text-center">
          {{ errors.apiError }}
        </div>

        <div class="mt-4 text-center text-sm text-gray-600 dark:text-gray-400">
          Уже есть аккаунт?
          <router-link
            to="/login"
            class="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
          >
            Войдите здесь
          </router-link>
        </div>
      </Form>
    </div>
  </div>
</template>

<style scoped>
.invalid-feedback {
  color: #dc3545;
  font-size: 0.875rem;
  margin-top: 0.25rem;
}

.is-invalid {
  border-color: #dc3545 !important;
}

.alert-danger {
  color: #721c24;
  background-color: #f8d7da;
  border-color: #f5c6cb;
  padding: 0.75rem 1.25rem;
  margin-top: 1rem;
  border-radius: 0.25rem;
}
</style>
