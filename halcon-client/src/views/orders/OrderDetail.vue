<template>
  <div>
    <div class="mb-6">
      <router-link to="/dashboard/orders" class="text-blue-600 hover:text-blue-700 text-sm">
        ← Back to Orders
      </router-link>
      <h1 class="text-3xl font-bold text-gray-900 mt-2">Order Details</h1>
    </div>

    <div v-if="ordersStore.loading" class="text-center py-12">
      <p class="text-gray-500">Loading order...</p>
    </div>

    <div v-else-if="order" class="space-y-6">
      <!-- Order Information -->
      <div class="bg-white rounded-lg shadow p-6">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Order Information</h2>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <span class="text-sm text-gray-500">Invoice Number:</span>
            <p class="font-medium text-gray-900">{{ order.invoice_number }}</p>
          </div>
          <div>
            <span class="text-sm text-gray-500">Status:</span>
            <span :class="getStatusClass(order.status)" class="inline-block px-3 py-1 rounded-full text-sm font-medium mt-1">
              {{ order.status }}
            </span>
          </div>
          <div>
            <span class="text-sm text-gray-500">Customer Name:</span>
            <p class="font-medium text-gray-900">{{ order.customer_name }}</p>
          </div>
          <div>
            <span class="text-sm text-gray-500">Customer Number:</span>
            <p class="font-medium text-gray-900">{{ order.customer_number }}</p>
          </div>
          <div class="col-span-2">
            <span class="text-sm text-gray-500">Delivery Address:</span>
            <p class="font-medium text-gray-900">{{ order.delivery_address || 'N/A' }}</p>
          </div>
          <div class="col-span-2">
            <span class="text-sm text-gray-500">Notes:</span>
            <p class="font-medium text-gray-900">{{ order.notes || 'N/A' }}</p>
          </div>
        </div>
      </div>

      <!-- Status Update (Warehouse and Route roles) -->
      <div v-if="canUpdateStatus" class="bg-white rounded-lg shadow p-6">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Update Status</h2>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">New Status</label>
            <select
              v-model="newStatus"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="">Select status...</option>
              <option v-for="status in availableStatuses" :key="status" :value="status">
                {{ status }}
              </option>
            </select>
          </div>

          <button
            @click="updateStatus"
            :disabled="!newStatus || ordersStore.loading"
            class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors disabled:opacity-50"
          >
            Update Status
          </button>
        </div>
      </div>

      <!-- Evidence Upload (Route role only) -->
      <div v-if="authStore.userRole === 'Route'" class="bg-white rounded-lg shadow p-6">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Upload Delivery Evidence</h2>
        <div class="space-y-4">
          <div>
            <input
              type="file"
              accept="image/*"
              @change="handleFileSelect"
              class="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
            />
          </div>

          <div v-if="selectedFile" class="flex items-center gap-4">
            <img
              :src="previewUrl"
              alt="Preview"
              class="w-32 h-32 object-cover rounded-lg border"
            />
            <div>
              <p class="text-sm text-gray-700">{{ selectedFile.name }}</p>
              <p class="text-xs text-gray-500">{{ (selectedFile.size / 1024).toFixed(2) }} KB</p>
            </div>
          </div>

          <div class="flex items-center gap-2">
            <input
              type="checkbox"
              id="markDelivered"
              v-model="markAsDelivered"
              class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
            />
            <label for="markDelivered" class="text-sm text-gray-700">
              Mark as Delivered
            </label>
          </div>

          <button
            @click="uploadEvidence"
            :disabled="!selectedFile || ordersStore.loading"
            class="bg-green-600 hover:bg-green-700 text-white px-6 py-2 rounded-lg font-medium transition-colors disabled:opacity-50"
          >
            {{ ordersStore.loading ? 'Uploading...' : 'Upload Evidence' }}
          </button>
        </div>
      </div>

      <!-- Evidence Photo Display -->
      <div v-if="order.evidence_photo_url" class="bg-white rounded-lg shadow p-6">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Delivery Evidence</h2>
        <img
          :src="`${apiUrl}${order.evidence_photo_url}`"
          alt="Delivery evidence"
          class="max-w-md rounded-lg border"
        />
      </div>

      <!-- Metadata -->
      <div class="bg-white rounded-lg shadow p-6">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Metadata</h2>
        <div class="grid grid-cols-2 gap-4 text-sm">
          <div>
            <span class="text-gray-500">Created by:</span>
            <p class="font-medium">{{ order.created_by_user?.full_name || 'N/A' }}</p>
          </div>
          <div>
            <span class="text-gray-500">Created at:</span>
            <p class="font-medium">{{ formatDate(order.created_at) }}</p>
          </div>
          <div>
            <span class="text-gray-500">Last modified by:</span>
            <p class="font-medium">{{ order.last_modified_by_user?.full_name || 'N/A' }}</p>
          </div>
          <div>
            <span class="text-gray-500">Last modified:</span>
            <p class="font-medium">{{ formatDate(order.updated_at) }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useOrdersStore, type OrderStatus } from '@/stores/orders'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const ordersStore = useOrdersStore()
const authStore = useAuthStore()

const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const order = computed(() => ordersStore.currentOrder)
const newStatus = ref<OrderStatus | ''>('')
const selectedFile = ref<File | null>(null)
const previewUrl = ref<string>('')
const markAsDelivered = ref(false)

const canUpdateStatus = computed(() => {
  const role = authStore.userRole
  return role === 'Warehouse' || role === 'Route'
})

const availableStatuses = computed(() => {
  if (!order.value) return []
  const role = authStore.userRole
  const currentStatus = order.value.status

  const transitions: Record<string, Record<string, OrderStatus[]>> = {
    'Warehouse': {
      'Ordered': ['In Process'],
      'In Process': ['In Route'],
    },
    'Route': {
      'In Route': ['Delivered'],
    },
  }

  return transitions[role || '']?.[currentStatus] || []
})

onMounted(async () => {
  const id = Number(route.params.id)
  await ordersStore.fetchOrder(id)
})

const updateStatus = async () => {
  if (!order.value || !newStatus.value) return
  
  await ordersStore.updateOrder(order.value.id, { status: newStatus.value })
  newStatus.value = ''
}

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    selectedFile.value = file
    previewUrl.value = URL.createObjectURL(file)
  }
}

const uploadEvidence = async () => {
  if (!order.value || !selectedFile.value) return

  const status = markAsDelivered.value ? 'Delivered' : undefined
  const result = await ordersStore.uploadEvidence(order.value.id, selectedFile.value, status)
  
  if (result) {
    selectedFile.value = null
    previewUrl.value = ''
    markAsDelivered.value = false
    await ordersStore.fetchOrder(order.value.id)
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
