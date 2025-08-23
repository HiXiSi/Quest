import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/utils/api'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(null)

  const isLoggedIn = computed(() => !!token.value && !!user.value)

  // 设置token
  const setToken = (newToken) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  // 设置用户信息
  const setUser = (userData) => {
    user.value = userData
  }

  // 登录
  const login = async (credentials) => {
    try {
      const response = await api.post('/auth/login', credentials)
      const { token: newToken, user: userData } = response.data.data
      
      setToken(newToken)
      setUser(userData)
      
      return response
    } catch (error) {
      throw error
    }
  }

  // 注册
  const register = async (userData) => {
    try {
      const response = await api.post('/auth/register', userData)
      const { token: newToken, user: userInfo } = response.data.data
      
      setToken(newToken)
      setUser(userInfo)
      
      return response
    } catch (error) {
      throw error
    }
  }

  // 获取用户信息
  const fetchUserProfile = async () => {
    try {
      const response = await api.get('/users/profile')
      setUser(response.data.data)
      return response
    } catch (error) {
      throw error
    }
  }

  // 更新用户信息
  const updateProfile = async (profileData) => {
    try {
      const response = await api.put('/users/profile', profileData)
      setUser(response.data.data)
      return response
    } catch (error) {
      throw error
    }
  }

  // 登出
  const logout = () => {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
  }

  // 初始化用户信息
  const initUser = async () => {
    if (token.value && !user.value) {
      try {
        await fetchUserProfile()
      } catch (error) {
        // token可能已过期
        logout()
      }
    }
  }

  return {
    token,
    user,
    isLoggedIn,
    login,
    register,
    fetchUserProfile,
    updateProfile,
    logout,
    initUser
  }
})