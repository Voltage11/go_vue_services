import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authAPI } from '@/api/auth'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const isLoading = ref(false)
  const error = ref(null)

  const isAuthenticated = computed(() => !!user.value)

  const register = async (userData) => {
    isLoading.value = true
    error.value = null
    try {
      const response = await authAPI.register(userData)
      await login({ email: userData.email, password: userData.password })
      return response
    } catch (err) {
      error.value = err.response?.data?.message || 'Ошибка регистрации'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const login = async (credentials) => {
    isLoading.value = true
    error.value = null
    try {
      const response = await authAPI.login(credentials)
      // JWT токен автоматически сохраняется в cookies
      user.value = { email: credentials.email } // Можно расширить данными из ответа
      router.push('/')
      return response
    } catch (err) {
      error.value = err.response?.data?.message || 'Ошибка авторизации'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const logout = async () => {
    try {
      await authAPI.logout()
      user.value = null
      router.push('/login')
    } catch (err) {
      console.error('Ошибка выхода:', err)
    }
  }

  const clearError = () => {
    error.value = null
  }

  return {
    user,
    isLoading,
    error,
    isAuthenticated,
    register,
    login,
    logout,
    clearError
  }
})
