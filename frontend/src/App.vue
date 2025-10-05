<script setup>
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import HeaderWithMenu from './components/HeaderWithMenu.vue'

const authStore = useAuthStore()

// При загрузке приложения можно проверить авторизацию
onMounted(() => {
  if (authStore.isAuthenticated) {
    // Здесь можно выполнить действия при авторизации
  }
  // Здесь можно добавить проверку авторизации при загрузке
})
</script>


<template>
  <div id="app">
    <HeaderWithMenu />
    <div class="container-fluid">
      qqqq
      <router-view />
    </div>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

#app {
  min-height: 100vh;
  width: 100%;
}
</style>
