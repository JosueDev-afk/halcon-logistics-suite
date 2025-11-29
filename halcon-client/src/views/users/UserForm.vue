<template>
  <div>
    <div class="mb-6">
      <router-link to="/dashboard/users" class="text-blue-600 hover:text-blue-700 text-sm">
        ← Back to Users
      </router-link>
      <h1 class="text-3xl font-bold text-gray-900 mt-2">{{ isEdit ? 'Edit User' : 'Create User' }}</h1>
    </div>

    <div class="bg-white rounded-lg shadow p-6 max-w-2xl">
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div>
          <label for="username" class="block text-sm font-medium text-gray-700 mb-2">
            Username *
          </label>
          <input
            id="username"
            v-model="form.username"
            type="text"
            required
            :disabled="isEdit"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
          />
        </div>

        <div>
          <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
            Password {{ isEdit ? '(leave blank to keep current)' : '*' }}
          </label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            :required="!isEdit"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        <div>
          <label for="fullName" class="block text-sm font-medium text-gray-700 mb-2">
            Full Name
          </label>
          <input
            id="fullName"
            v-model="form.full_name"
            type="text"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        <div>
          <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
            Email
          </label>
          <input
            id="email"
            v-model="form.email"
            type="email"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        <div>
          <label for="role" class="block text-sm font-medium text-gray-700 mb-2">
            Role *
          </label>
          <select
            id="role"
            v-model="form.role"
            required
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          >
            <option value="">Select role...</option>
            <option value="Admin">Admin</option>
            <option value="Sales">Sales</option>
            <option value="Purchasing">Purchasing</option>
            <option value="Warehouse">Warehouse</option>
            <option value="Route">Route</option>
          </select>
        </div>

        <div>
          <label for="department" class="block text-sm font-medium text-gray-700 mb-2">
            Department
          </label>
          <input
            id="department"
            v-model="form.department"
            type="text"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        <div v-if="isEdit" class="flex items-center">
          <input
            id="isActive"
            v-model="form.is_active"
            type="checkbox"
            class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
          />
          <label for="isActive" class="ml-2 text-sm text-gray-700">
            Active
          </label>
        </div>

        <div v-if="usersStore.error" class="p-4 bg-red-50 border border-red-200 rounded-lg">
          <p class="text-sm text-red-600">{{ usersStore.error }}</p>
        </div>

        <div class="flex gap-4">
          <button
            type="submit"
            :disabled="usersStore.loading"
            class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors disabled:opacity-50"
          >
            {{ usersStore.loading ? 'Saving...' : (isEdit ? 'Update User' : 'Create User') }}
          </button>
          <router-link
            to="/dashboard/users"
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
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUsersStore } from '@/stores/users'

const route = useRoute()
const router = useRouter()
const usersStore = useUsersStore()

const isEdit = computed(() => !!route.params.id)

const form = ref({
  username: '',
  password: '',
  full_name: '',
  email: '',
  role: '',
  department: '',
  is_active: true,
})

onMounted(async () => {
  if (isEdit.value) {
    const id = Number(route.params.id)
    const user = await usersStore.fetchUser(id)
    if (user) {
      form.value = {
        username: user.username,
        password: '',
        full_name: user.full_name,
        email: user.email,
        role: user.role,
        department: user.department,
        is_active: user.is_active,
      }
    }
  }
})

const handleSubmit = async () => {
  const userData: any = { ...form.value }
  
  // Remove password if empty in edit mode
  if (isEdit.value && !userData.password) {
    delete userData.password
  }

  let result
  if (isEdit.value) {
    const id = Number(route.params.id)
    result = await usersStore.updateUser(id, userData)
  } else {
    result = await usersStore.createUser(userData)
  }

  if (result) {
    router.push('/dashboard/users')
  }
}
</script>
