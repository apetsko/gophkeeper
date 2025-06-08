<script setup>
import { Form, Field } from 'vee-validate';
import * as Yup from 'yup';
import { ref } from 'vue';
import { fetchWrapper } from '@/helpers/fetch-wrapper.js';

const bankCardType = 1;
const baseUrl = `${import.meta.env.VITE_API_URL}/v1`;

const schema = Yup.object().shape({
  cardNumber: Yup.string()
    .required('Обязательно к заполнению')
    .test('card-number', 'Номер карты должен содержать 16 цифр', value => {
      if (!value) return false;
      const digitsOnly = value.replace(/\s/g, '');
      return /^\d{16}$/.test(digitsOnly);
    }),
  cardholder: Yup.string()
    .required('Обязательно к заполнению')
    .max(100, 'Максимум 100 символов'),
  expiryDate: Yup.string()
    .required('Обязательно к заполнению')
    .matches(/^(0[1-9]|1[0-2])\/?([0-9]{2})$/, 'Формат: ММ/ГГ'),
  cvv: Yup.string()
    .required('Обязательно к заполнению')
    .matches(/^\d{3,4}$/, 'CVV должен содержать 3 или 4 цифры'),
  meta: Yup.string()
    .max(1000, 'Максимум 1000 символов'),
});

const errorSubmitForm = ref(false);
const successSubmitForm = ref(false);
const isSubmitting = ref(false);
const isFormSubmitted = ref(false);
const showCvv = ref(false);

function formatCardNumber(value) {
  if (!value) return '';
  // Удаляем все не-цифры
  const digitsOnly = value.replace(/\D/g, '');
  // Форматируем по 4 цифры с пробелами

  console.log(digitsOnly.replace(/(\d{4})(?=\d)/g, '$1 ').substring(0, 19))

  return digitsOnly.replace(/(\d{4})(?=\d)/g, '$1 ').substring(0, 19);
}

function formatExpiryDate(value) {
  return value
    .replace(/[^0-9]/g, '')
    .replace(/^([2-9])$/g, '0$1')
    .replace(/^(1{1})([3-9]{1})$/g, '0$1/$2')
    .replace(/^0{1,}/g, '0')
    .replace(/^([0-1]{1}[0-9]{1})([0-9]{1,2}).*/g, '$1/$2');
}

async function onSubmit(values, { resetForm }) {
  const { cardNumber, cardholder, expiryDate, cvv, meta } = values;

  const requestData = {
    bank_card: {
      card_number: cardNumber.replace(/\s/g, ''),
      cardholder: cardholder.trim(),
      expiry_date: expiryDate,
      cvv: cvv,
    },
    meta: {
      content: meta
    },
    type: bankCardType
  };

  isSubmitting.value = true;
  isFormSubmitted.value = true;

  try {
    const response = await fetchWrapper.post(`${baseUrl}/data/save`, requestData);

    if (response && response.error) {
      errorSubmitForm.value = response.error || 'Ошибка сервера';
      successSubmitForm.value = null;
    } else {
      resetForm();
      errorSubmitForm.value = null;
      successSubmitForm.value = 'Данные карты успешно сохранены!'

      setTimeout(() => {
        successSubmitForm.value = null;
      }, 5000);
    }
  } catch (error) {
    successSubmitForm.value = null;
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
        v-slot="{ errors, setFieldValue }"
        autocomplete="on"
      >
        <!-- Card Number Field -->
        <div class="mb-4">
          <label for="cardNumber" class="block text-sm font-medium text-gray-900 dark:text-white mb-2">
            Номер карты
          </label>
          <Field
            name="cardNumber"
            type="text"
            id="cardNumber"
            v-slot="{ field }"
            class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            :class="{ 'border-red-500': errors.cardNumber }"
            @input="e => setFieldValue('cardNumber', formatCardNumber(e.target.value))"
            placeholder="0000 0000 0000 0000"
          />
          <div v-if="errors.cardNumber" class="text-red-500 text-xs mt-1">
            {{ errors.cardNumber }}
          </div>
        </div>

        <!-- Card Holder Field -->
        <div class="mb-4">
          <label for="cardholder" class="block text-sm font-medium text-gray-900 dark:text-white mb-2">
            Имя владельца
          </label>
          <Field
            name="cardholder"
            type="text"
            id="cardholder"
            class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            :class="{ 'border-red-500': errors.cardholder }"
            placeholder="IVAN IVANOV"
          />
          <div v-if="errors.cardholder" class="text-red-500 text-xs mt-1">
            {{ errors.cardholder }}
          </div>
        </div>

        <!-- Expiry Date and CVV -->
        <div class="grid grid-cols-2 gap-4 mb-4">
          <div>
            <label for="expiryDate" class="block text-sm font-medium text-gray-900 dark:text-white mb-2">
              Срок действия
            </label>
            <Field
              name="expiryDate"
              type="text"
              id="expiryDate"
              v-slot="{ field }"
              class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
              :class="{ 'border-red-500': errors.expiryDate }"
              @input="e => setFieldValue('expiryDate', formatExpiryDate(e.target.value))"
              placeholder="ММ/ГГ"
            />
            <div v-if="errors.expiryDate" class="text-red-500 text-xs mt-1">
              {{ errors.expiryDate }}
            </div>
          </div>

          <div>
            <label for="cvv" class="block text-sm font-medium text-gray-900 dark:text-white mb-2">
              CVV/CVC
            </label>
            <div class="relative">
              <Field
                name="cvv"
                :type="showCvv ? 'text' : 'password'"
                id="cvv"
                class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                :class="{ 'border-red-500': errors.cvv }"
                placeholder="123"
                maxlength="4"
              />
              <button
                type="button"
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-sm leading-5"
                @click="showCvv = !showCvv"
              >
                {{ showCvv ? '👁️' : '👁️‍🗨️' }}
              </button>
            </div>
            <div v-if="errors.cvv" class="text-red-500 text-xs mt-1">
              {{ errors.cvv }}
            </div>
          </div>
        </div>

        <!-- Meta Field -->
        <div class="mb-6">
          <label for="meta" class="block text-sm font-medium text-gray-900 dark:text-white mb-2">
            Дополнительная информация
          </label>
          <div class="relative">
            <Field
              name="meta"
              as="textarea"
              rows="3"
              id="meta"
              class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
              placeholder="Примечания, описание карты и т.д."
            />
          </div>
          <div v-if="errors.meta" class="text-red-500 text-xs mt-1">
            {{ errors.meta }}
          </div>
        </div>

        <div class="text-center">
          <button
            class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800 disabled:opacity-50"
            type="submit"
            :disabled="isSubmitting || isFormSubmitted"
          >
            <span v-if="isSubmitting" class="inline-block animate-spin mr-2">↻</span>
            <span>&nbsp;{{ isSubmitting ? 'Сохранение...' : 'Сохранить карту' }}</span>
          </button>
        </div>
      </Form>
    </div>
  </div>
</template>

<style scoped>

</style>
