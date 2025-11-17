<template>
  <v-container fill-height fluid class="d-flex align-center justify-center">
    <v-card width="450" elevation="10" class="pa-4">
      <h2 class="text-h5 text-center mb-4">Вход в систему</h2>

      <v-form @submit.prevent="handleSubmit">
        <div v-if="step === 1">
          <v-text-field
            v-model="email"
            label="Введите ваш E-mail"
            variant="filled"
            type="email"
            prepend-inner-icon="mdi-email"
            placeholder="example@mail.com"
            bg-color="#EEE0FF"
            class="mb-3"
            :rules="[v => !!v || 'Email обязателен']"
          ></v-text-field>

          <v-btn
            type="submit"
            color="#542F99"
            block
            size="large"
            class="mt-4 text-white"
            :loading="loading"
          >
            Получить код
          </v-btn>
        </div>

        <div v-else>
          <div class="text-center mb-4 text-grey">
            Код отправлен на <b>{{ email }}</b>
          </div>

          <v-otp-input
            v-model="otpCode"
            length="6"
            variant="outlined"
            class="mb-4 justify-center"
          ></v-otp-input>

          <v-btn
            type="submit"
            color="#542F99"
            block
            size="large"
            class="mt-4 text-white"
            :loading="loading"
          >
            Войти
          </v-btn>

          <v-btn
            variant="text"
            block
            class="mt-2"
            @click="step = 1"
          >
            Изменить Email
          </v-btn>
        </div>
      </v-form>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();

// Состояние
const step = ref(1); // 1 - вводим email, 2 - вводим код
const loading = ref(false);
const email = ref('');
const otpCode = ref('');

// Функция-распределитель: решает, что делать при нажатии кнопки
const handleSubmit = () => {
  if (step.value === 1) {
    sendEmail();
  } else {
    verifyCode();
  }
};

// Отправка Email на сервер
const sendEmail = async () => {
  if (!email.value) return;
  loading.value = true;

  try {
    const response = await fetch('http://localhost:8000/api/v1/auth/login/', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: email.value }) // Отправляем только email
    });

    if (response.ok) {
      step.value = 2; // Переходим ко вводу кода
    } else {
      alert('Ошибка! Возможно, неверный email.');
    }
  } catch (e) {
    console.error(e);
    alert('Ошибка сети');
  } finally {
    loading.value = false;
  }
};

const verifyCode = async () => {
  if (!otpCode.value) return;
  loading.value = true;

  try {
    // Согласно доке: POST /api/v1/auth/confirm/
    const payload = {
      email: email.value,
      code: otpCode.value
    };

    const response = await fetch('http://localhost:8000/api/v1/auth/confirm/', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    });

    if (response.ok) {
      const data = await response.json();
      localStorage.setItem('token', data.token);

      alert('Успешный вход!');
      router.push('/home'); // Перенаправляем на главную
    } else {
      alert('Неверный код');
    }
  } catch (e) {
    console.error(e);
    alert('Ошибка сети');
  } finally {
    loading.value = false;
  }
};
</script>
