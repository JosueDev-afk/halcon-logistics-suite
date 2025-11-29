import { defineStore } from 'pinia'
import { ref } from 'vue'
import apiClient from '@/api/client'

export type OrderStatus = 'Ordered' | 'In Process' | 'In Route' | 'Delivered'

export interface Order {
    id: number
    invoice_number: string
    customer_name: string
    customer_number: string
    status: OrderStatus
    delivery_address: string
    notes: string
    evidence_photo_url: string
    is_deleted: boolean
    created_by: number
    last_modified_by: number
    created_at: string
    updated_at: string
    created_by_user?: {
        id: number
        username: string
        full_name: string
    }
    last_modified_by_user?: {
        id: number
        username: string
        full_name: string
    }
}

export interface OrderFilters {
    invoice_number?: string
    customer_name?: string
    customer_number?: string
    status?: string
    include_deleted?: boolean
}

export const useOrdersStore = defineStore('orders', () => {
    const orders = ref<Order[]>([])
    const currentOrder = ref<Order | null>(null)
    const loading = ref(false)
    const error = ref<string | null>(null)

    const fetchOrders = async (filters?: OrderFilters) => {
        loading.value = true
        error.value = null

        try {
            const params = new URLSearchParams()
            if (filters?.invoice_number) params.append('invoice_number', filters.invoice_number)
            if (filters?.customer_name) params.append('customer_name', filters.customer_name)
            if (filters?.customer_number) params.append('customer_number', filters.customer_number)
            if (filters?.status) params.append('status', filters.status)
            if (filters?.include_deleted) params.append('include_deleted', 'true')

            const response = await apiClient.get(`/orders?${params.toString()}`)
            orders.value = response.data
        } catch (err: any) {
            error.value = err.response?.data?.message || 'Failed to fetch orders'
        } finally {
            loading.value = false
        }
    }

    const fetchOrder = async (id: number) => {
        loading.value = true
        error.value = null

        try {
            const response = await apiClient.get(`/orders/${id}`)
            currentOrder.value = response.data
            return response.data
        } catch (err: any) {
            error.value = err.response?.data?.message || 'Failed to fetch order'
            return null
        } finally {
            loading.value = false
        }
    }

    const createOrder = async (orderData: Partial<Order>) => {
        loading.value = true
        error.value = null

        try {
            const response = await apiClient.post('/orders', orderData)
            orders.value.unshift(response.data)
            return response.data
        } catch (err: any) {
            error.value = err.response?.data?.message || 'Failed to create order'
            return null
        } finally {
            loading.value = false
        }
    }

    const updateOrder = async (id: number, orderData: Partial<Order>) => {
        loading.value = true
        error.value = null

        try {
            const response = await apiClient.put(`/orders/${id}`, orderData)
            const index = orders.value.findIndex((o) => o.id === id)
            if (index !== -1) {
                orders.value[index] = response.data
            }
            if (currentOrder.value?.id === id) {
                currentOrder.value = response.data
            }
            return response.data
        } catch (err: any) {
            error.value = err.response?.data?.message || 'Failed to update order'
            return null
        } finally {
            loading.value = false
        }
    }

    const deleteOrder = async (id: number) => {
        loading.value = true
        error.value = null

        try {
            await apiClient.delete(`/orders/${id}`)
            orders.value = orders.value.filter((o) => o.id !== id)
            return true
        } catch (err: any) {
            error.value = err.response?.data?.message || 'Failed to delete order'
            return false
        } finally {
            loading.value = false
        }
    }

    const restoreOrder = async (id: number) => {
        loading.value = true
        error.value = null

        try {
            const response = await apiClient.post(`/orders/${id}/restore`)
            return response.data
        } catch (err: any) {
            error.value = err.response?.data?.message || 'Failed to restore order'
            return null
        } finally {
            loading.value = false
        }
    }

    const uploadEvidence = async (id: number, file: File, status?: OrderStatus) => {
        loading.value = true
        error.value = null

        try {
            const formData = new FormData()
            formData.append('photo', file)
            if (status) {
                formData.append('status', status)
            }

            const response = await apiClient.post(`/orders/${id}/evidence`, formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            })
            return response.data
        } catch (err: any) {
            error.value = err.response?.data?.message || 'Failed to upload evidence'
            return null
        } finally {
            loading.value = false
        }
    }

    return {
        orders,
        currentOrder,
        loading,
        error,
        fetchOrders,
        fetchOrder,
        createOrder,
        updateOrder,
        deleteOrder,
        restoreOrder,
        uploadEvidence,
    }
})
