<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
    <div class="max-w-md w-full">
      <!-- Header -->
      <div class="text-center mb-8">
        <h1 class="text-4xl font-bold text-gray-900 mb-2">🦅 Halcon Logistics</h1>
        <p class="text-gray-600">Track your order status</p>
      </div>

      <!-- Tracking Form -->
      <div class="bg-white rounded-2xl shadow-xl p-8">
        <form @submit.prevent="trackOrder" class="space-y-6">
          <div>
            <label for="customerNumber" class="block text-sm font-medium text-gray-700 mb-2">
              Customer Number
            </label>
            <input
              id="customerNumber"
              v-model="customerNumber"
              type="text"
              required
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Enter your customer number"
            />
          </div>

          <div>
            <label for="invoiceNumber" class="block text-sm font-medium text-gray-700 mb-2">
              Invoice Number
            </label>
            <input
              id="invoiceNumber"
              v-model="invoiceNumber"
              type="text"
              required
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Enter your invoice number"
            />
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full bg-blue-600 hover:bg-blue-700 text-white font-semibold py-3 px-6 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ loading ? 'Tracking...' : 'Track Order' }}
          </button>
        </form>

        <!-- Error Message -->
        <div v-if="error" class="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
          <p class="text-sm text-red-600">{{ error }}</p>
        </div>

        <!-- Tracking Result -->
        <div v-if="trackingResult && trackingResult.found" class="mt-6 space-y-4">
          <div class="border-t pt-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Order Details</h3>
            
            <div class="space-y-3">
              <div>
                <span class="text-sm text-gray-500">Invoice Number:</span>
                <p class="font-medium text-gray-900">{{ trackingResult.invoice_number }}</p>
              </div>

              <div>
                <span class="text-sm text-gray-500">Customer Name:</span>
                <p class="font-medium text-gray-900">{{ trackingResult.customer_name }}</p>
              </div>

              <div>
                <span class="text-sm text-gray-500">Status:</span>
                <span :class="getStatusClass(trackingResult.status)" class="inline-block px-3 py-1 rounded-full text-sm font-medium mt-1">
                  {{ trackingResult.status }}
                </span>
              </div>

              <div>
                <span class="text-sm text-gray-500">Delivery Address:</span>
                <p class="font-medium text-gray-900">{{ trackingResult.delivery_address || 'N/A' }}</p>
              </div>

              <div v-if="trackingResult.evidence_photo_url" class="mt-4">
                <span class="text-sm text-gray-500 block mb-2">Delivery Evidence:</span>
                <img
                  :src="`${apiUrl}${trackingResult.evidence_photo_url}`"
                  alt="Delivery evidence"
                  class="w-full rounded-lg border border-gray-200"
                />
              </div>

              <div class="text-xs text-gray-400 mt-4">
                Last updated: {{ formatDate(trackingResult.updated_at) }}
              </div>
            </div>
          </div>
        </div>

        <div v-else-if="trackingResult && !trackingResult.found" class="mt-6">
          <div class="p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
            <p class="text-sm text-yellow-800">No order found with the provided information.</p>
          </div>
        </div>

        <!-- Login Link -->
        <div class="mt-6 text-center">
          <router-link to="/login" class="text-sm text-blue-600 hover:text-blue-700">
            Employee Login →
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'

const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const customerNumber = ref('')
const invoiceNumber = ref('')
const loading = ref(false)
const error = ref<string | null>(null)
const trackingResult = ref<any>(null)

const trackOrder = async () => {
  loading.value = true
  error.value = null
  trackingResult.value = null

  try {
    const response = await axios.get(`${apiUrl}/api/track`, {
      params: {
        customer_number: customerNumber.value,
        invoice_number: invoiceNumber.value,
      },
    })
    trackingResult.value = response.data
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to track order'
  } finally {
    loading.value = false
  }
}

const getStatusClass = (status: string) => {
  const classes: Record<string, string> = {
    'Ordered': 'bg-blue-100 text-blue-800',
    'In Process': 'bg-yellow-100 text-yellow-800',
    'In Route': 'bg-purple-100 text-purple-800',
    'Delivered': 'bg-green-100 text-green-800',
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}
</script>
