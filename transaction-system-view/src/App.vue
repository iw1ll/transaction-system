<template>
    <div class="app-container">
        <header class="app-header">
            <h1>Wallet</h1>
            <nav>
                <router-link to="/">Главная</router-link>
            </nav>
        </header>

        <main class="main-content">
            <MainPage />
        </main>

        <footer class="app-footer">
            <p>&copy; 2025 Wallet</p>
        </footer>
    </div>
</template>

<script setup lang="ts">
import MainPage from '@/views/MainPage.vue';
</script>

<style>
/* Глобальные стили */
:root {
  --primary-color: #42b883;
  --secondary-color: #35495e;
  --error-color: #ff4444;
  --background-color: #f8f9fa;
}

* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: 'Arial', sans-serif;
  line-height: 1.6;
  background-color: var(--background-color);
}

/* Стили layout */
.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  background-color: var(--secondary-color);
  color: white;
  padding: 1rem;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.app-header h1 {
  font-size: 1.5rem;
  margin-bottom: 0.5rem;
}

.main-content {
  flex: 1;
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

.app-footer {
  background-color: var(--secondary-color);
  color: white;
  text-align: center;
  padding: 1rem;
  margin-top: auto;
}

/* Общие стили компонентов */
.button {
  padding: 0.5rem 1rem;
  background-color: var(--primary-color);
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: opacity 0.2s;
}

.button:hover {
  opacity: 0.9;
}

.input-field {
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  width: 100%;
}

.error-message {
  color: var(--error-color);
  padding: 0.5rem;
  margin-top: 0.5rem;
}
</style>