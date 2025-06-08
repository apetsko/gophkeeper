<script setup>
import {ref, onMounted} from 'vue';
import {fetchWrapper} from '@/helpers/fetch-wrapper.js';

const baseUrl = `${import.meta.env.VITE_API_URL}/v1`;

// Данные таблицы и пагинация
const records = ref([]);
const totalCount = ref(0);
const currentPage = ref(1);
const itemsPerPage = ref(10);
const isLoading = ref(false);
const error = ref(null);

// Модалка и текущие детали
const isModalOpen = ref(false);
const modalData = ref(null);
const modalError = ref(null);
const modalLoading = ref(false);

// Получение данных с сервера
const fetchRecords = async () => {
  try {
    isLoading.value = true;
    error.value = null;

    const response = await fetchWrapper.get(
      `${baseUrl}/data/list?page=${currentPage.value}&limit=${itemsPerPage.value}`
    );

    records.value = response.records;
    totalCount.value = response.count;
  } catch (err) {
    error.value = err.message || 'Ошибка при загрузке данных';
    console.error('Ошибка:', err);
  } finally {
    isLoading.value = false;
  }
};

// Обработчик изменения страницы
const handlePageChange = (page) => {
  currentPage.value = page;
  fetchRecords();
};

const fetchDataDetails = async (id) => {
  try {
    modalLoading.value = true;
    modalError.value = null;

    const response = await fetchWrapper.get(`${baseUrl}/data/view?id=${id}`);

    if (response.type === 'DATA_TYPE_BINARY_DATA') {
      downloadBinaryFile(response.binaryData);
    } else {
      modalData.value = response;
      isModalOpen.value = true;
    }
  } catch (err) {
    modalError.value = err.message || 'Ошибка при загрузке данных';
  } finally {
    modalLoading.value = false;
  }
};

function downloadBinaryFile(binaryData) {
  const byteCharacters = atob(binaryData.data);
  const byteNumbers = new Array(byteCharacters.length);

  for (let i = 0; i < byteCharacters.length; i++) {
    byteNumbers[i] = byteCharacters.charCodeAt(i);
  }

  const byteArray = new Uint8Array(byteNumbers);

  const blob = new Blob([byteArray], { type: binaryData.type });

  const link = document.createElement('a');
  link.href = URL.createObjectURL(blob);
  link.download = binaryData.name || 'file';

  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);

  URL.revokeObjectURL(link.href);
}

// Функция удаления записи
const deleteRecord = async (id) => {
  if (!confirm('Вы уверены, что хотите удалить эту запись?')) {
    return;
  }

  try {
    await fetchWrapper.delete(`${baseUrl}/data/delete?id=${id}`);

    // Удаляем запись из массива без перезагрузки
    records.value = records.value.filter(record => record.id !== id);
    totalCount.value -= 1;

    // Если на текущей странице не осталось записей и это не первая страница,
    // переходим на предыдущую страницу
    if (records.value.length === 0 && currentPage.value > 1) {
      currentPage.value -= 1;
      fetchRecords();
    }
  } catch (err) {
    error.value = err.message || 'Ошибка при удалении записи';
    console.error('Ошибка:', err);
  }
};

// Обработчик действий
const handleAction = async (itemId, itemType) => {
  await fetchDataDetails(itemId);
};

// Загрузка данных при монтировании
onMounted(() => {
  fetchRecords();
});

const closeModal = () => {
  isModalOpen.value = false;
  modalData.value = null;
  modalError.value = null;
};
</script>

<template>
  <main>
    <div v-if="isLoading" class="text-center py-4">Загрузка данных...</div>

    <div v-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
      {{ error }}
    </div>

    <div v-if="!isLoading && !error" class="relative overflow-x-auto">
      <table class="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
        <thead
          class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
        <tr>
          <th class="px-6 py-3">#</th>
          <th class="px-6 py-3">Тип</th>
          <th class="px-6 py-3">Meta</th>
          <th class="px-6 py-3">Создано</th>
          <th class="px-6 py-3">Действия</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="item in records" :key="item.id" class="dark:bg-gray-800 dark:border-gray-700">
          <th class="px-6 py-4 font-medium whitespace-nowrap">{{ item.id }}</th>
          <td class="px-6 py-4">{{ item.type }}</td>
          <td class="px-6 py-4">{{ item.meta?.content }}</td>
          <td class="px-6 py-4">{{ item.createdAt }}</td>
          <td class="px-6 py-4 space-x-2">
            <button
              @click="handleAction(item.id, item.type)"
              class="text-blue-600 hover:underline"
            >
              <template v-if="item.type !== 'binary_data'">Просмотр</template>
              <template v-else>Скачать</template>
            </button>
            <br>
            <button
              @click="deleteRecord(item.id)"
              class="text-red-600 hover:underline"
            >
              Удалить
            </button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Модалка -->
    <div v-if="isModalOpen"
         class="fixed inset-0 flex items-center justify-center z-50">
      <div class="rounded-lg p-6 max-w-lg w-full relative bg-white text-gray-600">
        <button
          @click="closeModal"
          class="absolute top-2 right-2 text-gray-600 hover:text-gray-900 font-bold"
          aria-label="Закрыть"
        >✖</button>

        <div v-if="modalLoading">Загрузка...</div>
        <div v-if="modalError" class="text-red-600">{{ modalError }}</div>

        <div v-if="modalData && !modalLoading">
          <h2 class="text-xl font-bold mb-4">Детали данных</h2>
          <div v-if="modalData.type === 'DATA_TYPE_BANK_CARD'">
            <p><strong>Номер карты:</strong> {{ modalData.bankCard.cardNumber }}</p>
            <p><strong>Срок действия:</strong> {{ modalData.bankCard.expiryDate }}</p>
            <p><strong>CVV:</strong> {{ modalData.bankCard.cvv }}</p>
            <p><strong>Владелец карты:</strong> {{ modalData.bankCard.cardholder }}</p>
          </div>
          <div v-else-if="modalData.type === 'DATA_TYPE_CREDENTIALS'">
            <p><strong>Логин:</strong> {{ modalData.credentials.login }}</p>
            <p><strong>Пароль:</strong> {{ modalData.credentials.password }}</p>
          </div>
          <div v-else>
            <p>Тип данных не поддерживается для просмотра</p>
          </div>
        </div>
      </div>
    </div>
  </main>
</template>
