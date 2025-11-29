import { defineStore } from 'pinia'
import { ref } from 'vue'
import apiClient from '@/api/client'

export interface User {
    id: number
    username: string
    role: 'Admin' | 'Sales' | 'Purchasing' | 'Warehouse' | 'Route'
    department: string
    full_name: string
    email: string
    is_active: boolean
    created_at: string
    updated_at: string
}

export const useUsersStore = defineStore('users', () => {
    const users = ref<User[]>([])
    const currentUser = ref<User | null>(null)
    const loading = ref(false)
    const error = ref<string | null>(null)

    const fetchUsers = async () => {
        loading.value = true
        error.value = null

        try {
            const response = await apiClient.get('/users')
            users.value = response.data
        } catch (err: any) {
            error.value = err.response?.data?.message || 'Failed to fetch users'
        } finally {
            loading.value = false
        }
    }

    const fetchUser = async (id: number) => {
        loading.value = true
        error.value = null

        try {
            const response = await apiClient.get(`/users/${id}`)
            currentUser.value = response.data
            return response.data
        } catch (err: any) {
            error.value = err.response?.data?.message || 'Failed to fetch user'
            return null
        } finally {
            loading.value = false
        }
    }

    const createUser = async (userData: Partial<User> & { password: string }) => {
        loading.value = true
        error.value = null

        try {
            const response = await apiClient.post('/users', userData)
            users.value.unshift(response.data)
            return response.data
        } catch (err: any) {
            error.value = err.response?.data?.message || 'Failed to create user'
            return null
        } finally {
            loading.value = false
        }
    }

    const updateUser = async (id: number, userData: Partial<User> & { password?: string }) => {
        loading.value = true
        error.value = null

        try {
            const response = await apiClient.put(`/users/${id}`, userData)
            const index = users.value.findIndex((u) => u.id === id)
            if (index !== -1) {
                users.value[index] = response.data
            }
            return response.data
        } catch (err: any) {
            error.value = err.response?.data?.message || 'Failed to update user'
            return null
        } finally {
            loading.value = false
        }
    }

    const deleteUser = async (id: number) => {
        loading.value = true
        error.value = null

        try {
            await apiClient.delete(`/users/${id}`)
            users.value = users.value.filter((u) => u.id !== id)
            return true
        } catch (err: any) {
            error.value = err.response?.data?.message || 'Failed to delete user'
            return false
        } finally {
            loading.value = false
        }
    }

    return {
        users,
        currentUser,
        loading,
        error,
        fetchUsers,
        fetchUser,
        createUser,
        updateUser,
        deleteUser,
    }
})
