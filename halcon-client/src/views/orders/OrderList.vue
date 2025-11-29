<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold text-gray-900">Orders</h1>
      <router-link
        v-if="authStore.userRole === 'Sales'"
        to="/dashboard/orders/create"
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors"
      >
        + Create Order
      </router-link>
    </div>

    <!-- Filters -->
    <div class="bg-white rounded-lg shadow p-6 mb-6">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <input
          v-model="filters.invoice_number"
          type="text"
          placeholder="Invoice Number"
          class="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
        <input
          v-model="filters.customer_name"
          type="text"
          placeholder="Customer Name"
          class="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
        <input
          v-model="filters.customer_number"
          type="text"
          placeholder="Customer Number"
          class="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
        <select
          v-model="filters.status"
          class="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        >
          <option value="">All Statuses</option>
          <option value="Ordered">Ordered</option>
          <option value="In Process">In Process</option>
          <option value="In Route">In Route</option>
          <option value="Delivered">Delivered</option>
        </select>
      </div>
      <div class="mt-4 flex gap-2">
        <button
          @click="searchOrders"
          class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors"
        >
          Search
        </button>
        <button
          @click="clearFilters"
          class="bg-gray-200 hover:bg-gray-300 text-gray-700 px-6 py-2 rounded-lg font-medium transition-colors"
        >
          Clear
        </button>
      </div>
    </div>

    <!-- Orders Table -->
    <div class="bg-white rounded-lg shadow overflow-hidden">
      <div v-if="ordersStore.loading" class="p-8 text-center">
        <p class="text-gray-500">Loading orders...</p>
      </div>

      <div v-else-if="ordersStore.error" class="p-8 text-center">
        <p class="text-red-600">{{ ordersStore.error }}</p>
      </div>

      <div v-else-if="ordersStore.orders.length === 0" class="p-8 text-center">
        <p class="text-gray-500">No orders found</p>
      </div>

      <table v-else class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Invoice</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Customer</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Created</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="order in ordersStore.orders" :key="order.id" class="hover:bg-gray-50">
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm font-medium text-gray-900">{{ order.invoice_number }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm text-gray-900">{{ order.customer_name }}</div>
              <div class="text-xs text-gray-500">{{ order.customer_number }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span :class="getStatusClass(order.status)" class="px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full">
                {{ order.status }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ formatDate(order.created_at) }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
              <router-link
                :to="`/dashboard/orders/${order.id}`"
                class="text-blue-600 hover:text-blue-900 mr-4"
              >
                View
              </router-link>
              <button
                v-if="authStore.userRole === 'Admin' || authStore.userRole === 'Sales'"
                @click="deleteOrder(order.id)"
                class="text-red-600 hover:text-red-900"
              >
                Delete
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useOrdersStore } from '@/stores/orders'
import { useAuthStore } from '@/stores/auth'

const ordersStore = useOrdersStore()
const authStore = useAuthStore()

const filters = ref({
  invoice_number: '',
  customer_name: '',
  customer_number: '',
  status: '',
})

onMounted(() => {
  ordersStore.fetchOrders()
})

const searchOrders = () => {
  const cleanFilters = Object.fromEntries(
    Object.entries(filters.value).filter(([_, v]) => v !== '')
  )
  ordersStore.fetchOrders(cleanFilters)
}

const clearFilters = () => {
  filters.value = {
    invoice_number: '',
    customer_name: '',
    customer_number: '',
    status: '',
  }
  ordersStore.fetchOrders()
}

const deleteOrder = async (id: number) => {
  if (confirm('Are you sure you want to delete this order?')) {
    await ordersStore.deleteOrder(id)
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
  return new Date(dateString).toLocaleDateString()
}
</script>
