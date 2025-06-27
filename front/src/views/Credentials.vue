<script setup>
import { Form, Field } from 'vee-validate';
import * as Yup from 'yup';
import { ref } from 'vue';
import { fetchWrapper } from '@/helpers/fetch-wrapper.js';

const credentialsType = 2;
const baseUrl = `${import.meta.env.VITE_API_URL}/v1`;

const schema = Yup.object().shape({
  login: Yup.string()
    .required('Обязательно к заполнению')
    .min(1, 'Максимум 1 символ')
    .max(150, 'Максимум 20 символов'),
  password: Yup.string()
    .required('Обязательно к заполнению'),
  meta: Yup.string()
    .max(1000, 'Максимум 1000 символов'),
});

const errorSubmitForm = ref(false);
const successSubmitForm = ref(false);
const isSubmitting = ref(false);
const isFormSubmitted = ref(false);
const showPassword = ref(false);

async function onSubmit(values, { resetForm }) {
  const { login, password, meta } = values;

  const requestData = {
    credentials: {
      login: login.trim(),
      password: password,
    },
    meta: {
      content: meta
    },
    type: credentialsType
  };

  isFormSubmitted.value = true;
  isFormSubmitted.value = true;

  try {
    const response = await fetchWrapper.post(`${baseUrl}/data/save`, requestData);

    if (response && response.error) {
      errorSubmitForm.value = response.error || 'Ошибка сервера';
      successSubmitForm.value = null;
    } else {
      resetForm();
      errorSubmitForm.value = null;
      successSubmitForm.value = 'Данные успешно отправлены!'

      setTimeout(() => {
        successSubmitForm.value = null;
      }, 5000);
    }
  } catch (error) {
    successSubmitForm.value = null;
    errorSubmitForm.value = 'Ошибка отправки:' + error
    errorSubmitForm.value = error.message || 'Ошибка сети. Пожалуйста, попробуйте позже.';
  } finally {
    isSubmitting.value = false;
    isFormSubmitted.value = false;
  }
}
</script>

<template>
  <div class="flex justify-center flex-col w-100 m-auto">
    <div v-if="errorSubmitForm" class="p-4 mb-4 text-sm text-white rounded-lg bg-red-800" role="alert" style="margin-bottom: 15px;">
      <span class="font-medium">{{ errorSubmitForm }}</span>
    </div>

    <div v-if="successSubmitForm" class="p-4 mb-4 text-sm text-white rounded-lg bg-green-800" role="alert" style="margin-bottom: 15px;">
      <span class="font-medium">{{ successSubmitForm }}</span>
    </div>

    <div>
      <Form
        class="max-w-sm mx-auto"
        @submit="onSubmit"
        :validation-schema="schema"
        v-slot="{ errors }"
        autocomplete="off"
      >
        <!-- Login Field -->
        <div class="mb-4">
          <label for="login"
                 class="block text-sm font-medium text-gray-900 dark:text-white mb-2">Логин</label>
          <Field
            name="login"
            type="text"
            id="login"
            autocomplete="off"
            class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            :class="{ 'border-red-500': errors.login }"
          />
          <div v-if="errors.login" class="text-red-500 text-xs mt-1">
            {{ errors.login }}
          </div>
        </div>

        <!-- Password Field with toggle -->
        <div style="margin-top: 25px;">
          <label for="password"
                 class="block text-sm font-medium text-gray-900 dark:text-white mb-2">Пароль</label>
          <div class="relative">
            <Field
              name="password"
              :type="showPassword ? 'text' : 'password'"
              id="password"
              autocomplete="off"
              class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
              :class="{ 'border-red-500': errors.password }"
            />
            <button
              type="button"
              class="absolute inset-y-0 right-0 pr-3 flex items-center text-sm leading-5"
              @click="showPassword = !showPassword"
            >
              {{ showPassword ? '👁️' : '👁️‍🗨️' }}
            </button>
          </div>
          <div v-if="errors.password" class="text-red-500 text-xs mt-1">
            {{ errors.password }}
          </div>
        </div>

        <!-- Meta Field -->
        <div class="mb-6" style="margin-top: 25px;">
          <label for="meta"
                 class="block text-sm font-medium text-gray-900 dark:text-white mb-2">Дополнительная информация</label>
          <div class="relative">
            <Field
              name="meta"
              as="textarea"
              rows="5"
              id="meta"
              autocomplete="off"
              class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
              placeholder="Примечания, описание и т.д."
            />
          </div>
          <div v-if="errors.meta" class="text-red-500 text-xs mt-1">
            {{ errors.meta }}
          </div>
        </div>

        <!-- Submit Button -->
        <div class="text-center" style="margin-top: 25px;">
          <button
            class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800 disabled:opacity-50"
            type="submit"
            :disabled="isSubmitting || isFormSubmitted"
          >
            <span v-if="isSubmitting" class="inline-block animate-spin mr-2">↻</span>
            <span>&nbsp;{{ isSubmitting ? 'Сохранение...' : 'Сохранить' }}</span>
          </button>
        </div>

        <div v-if="errors.apiError" class="mt-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
          Errors: {{ errors.apiError }}
        </div>
      </Form>
    </div>
  </div>
</template>

<style scoped>

</style>
