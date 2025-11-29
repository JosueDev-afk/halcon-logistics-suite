import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import apiClient from '@/api/client'

export interface User {
    id: number
    username: string
    role: 'Admin' | 'Sales' | 'Purchasing' | 'Warehouse' | 'Route'
    department: string
    full_name: string
    email: string
}

export const useAuthStore = defineStore('auth', () => {
    const user = ref<User | null>(null)
    const token = ref<string | null>(null)
    const loading = ref(false)
    const error = ref<string | null>(null)

    const isAuthenticated = computed(() => !!token.value)
    const userRole = computed(() => user.value?.role)

    // Initialize from localStorage
    const initAuth = () => {
        const savedToken = localStorage.getItem('token')
        const savedUser = localStorage.getItem('user')

        if (savedToken && savedUser) {
            token.value = savedToken
            user.value = JSON.parse(savedUser)
        }
    }

    const login = async (username: string, password: string) => {
        loading.value = true
        error.value = null

        try {
            const response = await apiClient.post('/auth/login', { username, password })
            const { token: authToken, user: userData } = response.data

            token.value = authToken
            user.value = userData

            localStorage.setItem('token', authToken)
            localStorage.setItem('user', JSON.stringify(userData))

            return true
        } catch (err: any) {
            error.value = err.response?.data?.message || 'Login failed'
            return false
        } finally {
            loading.value = false
        }
    }

    const logout = () => {
        user.value = null
        token.value = null
        localStorage.removeItem('token')
        localStorage.removeItem('user')
    }

    const fetchCurrentUser = async () => {
        try {
            const response = await apiClient.get('/auth/me')
            user.value = response.data
            localStorage.setItem('user', JSON.stringify(response.data))
        } catch (err) {
            logout()
        }
    }

    return {
        user,
        token,
        loading,
        error,
        isAuthenticated,
        userRole,
        initAuth,
        login,
        logout,
        fetchCurrentUser,
    }
})
