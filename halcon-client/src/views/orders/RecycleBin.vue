<template>
  <div>
    <h1 class="text-3xl font-bold text-gray-900 mb-6">Recycle Bin</h1>

    <div class="bg-white rounded-lg shadow overflow-hidden">
      <div v-if="ordersStore.loading" class="p-8 text-center">
        <p class="text-gray-500">Loading deleted orders...</p>
      </div>

      <div v-else-if="ordersStore.orders.length === 0" class="p-8 text-center">
        <p class="text-gray-500">No deleted orders</p>
      </div>

      <table v-else class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Invoice</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Customer</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Deleted</th>
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
              {{ formatDate(order.updated_at) }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
              <button
                @click="restoreOrder(order.id)"
                class="text-green-600 hover:text-green-900"
              >
                Restore
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useOrdersStore } from '@/stores/orders'

const ordersStore = useOrdersStore()

onMounted(() => {
  ordersStore.fetchOrders({ include_deleted: true })
})

const restoreOrder = async (id: number) => {
  if (confirm('Are you sure you want to restore this order?')) {
    const result = await ordersStore.restoreOrder(id)
    if (result) {
      ordersStore.fetchOrders({ include_deleted: true })
    }
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
