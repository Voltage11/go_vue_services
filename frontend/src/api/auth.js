import axios from 'axios'

const API_BASE_URL = 'http://localhost:8080/api'

const api = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true // для работы с cookies
})

export const authAPI = {
  async register(userData) {
    const response = await api.post('/auth/register', userData)
    return response.data
  },

  async login(credentials) {
    const response = await api.post('/auth/login', credentials)
    return response.data
  },

  async logout() {
    const response = await api.post('/auth/logout')
    return response.data
  }
}

export default api
