<script setup>
// import { isAuthenticated, user } from '@/stores/auth';
import MenuMain from './MenuMain.vue';

</script>

<template>
  <div class="header">
    <div class="row">
      <div class="col text-start left-menu">
        <MenuMain />
      </div>

      <div class="col text-end">
        <!-- Если пользователь авторизован -->
         <template v-if="isAuthenticated">
          <!-- <router-link to="/profile">{{ user }}</router-link>
          <router-link to="/logout">Выйти</router-link> -->

         </template>
         <template v-else>
          <!-- <router-link to="/login">Войти</router-link>
          <router-link to="/register">Регистрация</router-link> -->
         </template>


      </div>
    </div>
  </div>
</template>


<style scoped>
  .header {
    width: 100%;
    height: 60px;
    background: var(--bg-color-dark);
  }
  .left-menu{
    display: flex;
    align-items: center;
    gap: 20px;
    padding-top: 5px;
    padding-bottom: 5px;
    margin-left: 30px;
  }
</style>
