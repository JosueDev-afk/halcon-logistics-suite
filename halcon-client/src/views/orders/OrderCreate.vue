<template>
  <div>
    <div class="mb-6">
      <router-link to="/dashboard/orders" class="text-blue-600 hover:text-blue-700 text-sm">
        ← Back to Orders
      </router-link>
      <h1 class="text-3xl font-bold text-gray-900 mt-2">Create New Order</h1>
    </div>

    <div class="bg-white rounded-lg shadow p-6 max-w-2xl">
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div>
          <label for="invoiceNumber" class="block text-sm font-medium text-gray-700 mb-2">
            Invoice Number *
          </label>
          <input
            id="invoiceNumber"
            v-model="form.invoice_number"
            type="text"
            required
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        <div>
          <label for="customerName" class="block text-sm font-medium text-gray-700 mb-2">
            Customer Name *
          </label>
          <input
            id="customerName"
            v-model="form.customer_name"
            type="text"
            required
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        <div>
          <label for="customerNumber" class="block text-sm font-medium text-gray-700 mb-2">
            Customer Number *
          </label>
          <input
            id="customerNumber"
            v-model="form.customer_number"
            type="text"
            required
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        <div>
          <label for="deliveryAddress" class="block text-sm font-medium text-gray-700 mb-2">
            Delivery Address
          </label>
          <textarea
            id="deliveryAddress"
            v-model="form.delivery_address"
            rows="3"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          ></textarea>
        </div>

        <div>
          <label for="notes" class="block text-sm font-medium text-gray-700 mb-2">
            Notes
          </label>
          <textarea
            id="notes"
            v-model="form.notes"
            rows="3"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          ></textarea>
        </div>

        <div v-if="ordersStore.error" class="p-4 bg-red-50 border border-red-200 rounded-lg">
          <p class="text-sm text-red-600">{{ ordersStore.error }}</p>
        </div>

        <div class="flex gap-4">
          <button
            type="submit"
            :disabled="ordersStore.loading"
            class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors disabled:opacity-50"
          >
            {{ ordersStore.loading ? 'Creating...' : 'Create Order' }}
          </button>
          <router-link
            to="/dashboard/orders"
            class="bg-gray-200 hover:bg-gray-300 text-gray-700 px-6 py-2 rounded-lg font-medium transition-colors inline-block"
          >
            Cancel
          </router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useOrdersStore } from '@/stores/orders'

const router = useRouter()
const ordersStore = useOrdersStore()

const form = ref({
  invoice_number: '',
  customer_name: '',
  customer_number: '',
  delivery_address: '',
  notes: '',
})

const handleSubmit = async () => {
  const order = await ordersStore.createOrder(form.value)
  if (order) {
    router.push('/dashboard/orders')
  }
}
</script>
