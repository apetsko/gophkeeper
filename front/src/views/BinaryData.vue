<script setup>
import {Form, Field} from 'vee-validate';
import * as Yup from 'yup';
import {ref} from 'vue';
import {fetchWrapper} from '@/helpers/fetch-wrapper.js';

const binaryDataType = 3;
const baseUrl = `${import.meta.env.VITE_API_URL}/v1`;

const schema = Yup.object().shape({
  file: Yup.mixed()
    .required('Файл обязателен для загрузки')
    .test('fileSize', 'Файл слишком большой', value => {
      if (!value) return true;
      return value.size <= 10 * 1024 * 1024; // 10MB
    }),
  meta: Yup.string()
    .max(1000, 'Максимум 1000 символов'),
});

const errorSubmitForm = ref(false);
const successSubmitForm = ref(false);
const isSubmitting = ref(false);
const isFormSubmitted = ref(false);
const fileName = ref('');

// Функция для чтения файла как base64
function readFileAsBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result.split(',')[1]); // Удаляем префикс data:...
    reader.onerror = error => reject(error);
    reader.readAsDataURL(file);
  });
}

async function onSubmit(values, {resetForm}) {
  const {file, meta} = values;

  isSubmitting.value = true;
  isFormSubmitted.value = true;

  try {
    const fileBase64 = await readFileAsBase64(file);

    const payload = {
      binary_data: {
        name: file.name,
        type: file.type,
        size: file.size,
        data: fileBase64,
      },
      meta: {
        content: meta
      },
      type: binaryDataType
    };

    const response = await fetchWrapper.post(`${baseUrl}/data/save`, payload);

    if (response?.error) {
      errorSubmitForm.value = response.error;
      successSubmitForm.value = null;
    } else {
      resetForm();
      fileName.value = '';
      successSubmitForm.value = 'Файл успешно загружен!';
      setTimeout(() => successSubmitForm.value = null, 5000);
    }
  } catch (error) {
    successSubmitForm.value = null;
    errorSubmitForm.value = error.message || 'Ошибка загрузки файла';
  } finally {
    isSubmitting.value = false;
    isFormSubmitted.value = false;
  }
}

// Обработчик изменения файла
function handleFileChange(event, setFieldValue) {
  const file = event.target.files[0];
  if (file) {
    setFieldValue('file', file);
    fileName.value = file.name;
  }
}
</script>

<template>
  <div class="flex justify-center flex-col w-100 m-auto">
    <div v-if="errorSubmitForm" class="p-4 mb-4 text-sm text-white rounded-lg bg-red-800"
         role="alert" style="margin-bottom: 15px;">
      <span class="font-medium">{{ errorSubmitForm }}</span>
    </div>

    <div v-if="successSubmitForm" class="p-4 mb-4 text-sm text-white rounded-lg bg-green-800"
         role="alert" style="margin-bottom: 15px;">
      <span class="font-medium">{{ successSubmitForm }}</span>
    </div>

    <div>
      <Form
        class="max-w-sm mx-auto"
        @submit="onSubmit"
        :validation-schema="schema"
        v-slot="{ errors, setFieldValue }"
        autocomplete="off"
      >
        <!-- File Upload Field -->
        <div class="mb-4">
          <label for="file" class="block text-sm font-medium text-gray-900 dark:text-white mb-2">
            Выберите файл
          </label>
          <div class="flex items-center">
            <label class="cursor-pointer">
              <span
                class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">
                Выбрать файл
              </span>
              <input
                type="file"
                id="file"
                name="file"
                class="hidden"
                @change="handleFileChange($event, setFieldValue)"
              />
            </label>
            <span class="ml-3 text-sm text-gray-700 dark:text-gray-300" v-if="fileName">
              {{ fileName }}
            </span>
          </div>
          <div v-if="errors.file" class="text-red-500 text-xs mt-1">
            {{ errors.file }}
          </div>
        </div>

        <!-- Meta Field -->
        <div class="mb-6" style="margin-top: 25px;">
          <label for="meta"
                 class="block text-sm font-medium text-gray-900 dark:text-white mb-2">Дополнительная
            информация</label>
          <div class="relative">
            <Field
              name="meta"
              as="textarea"
              rows="5"
              id="meta"
              autocomplete="off"
              class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
              placeholder="Описание файла, комментарии и т.д."
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
            <span>&nbsp;{{ isSubmitting ? 'Загрузка...' : 'Загрузить файл' }}</span>
          </button>
        </div>
      </Form>
    </div>
  </div>
</template>

<style scoped>
</style>
