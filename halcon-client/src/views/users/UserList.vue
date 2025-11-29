<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold text-gray-900">Users</h1>
      <router-link
        to="/dashboard/users/create"
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors"
      >
        + Create User
      </router-link>
    </div>

    <div class="bg-white rounded-lg shadow overflow-hidden">
      <div v-if="usersStore.loading" class="p-8 text-center">
        <p class="text-gray-500">Loading users...</p>
      </div>

      <div v-else-if="usersStore.error" class="p-8 text-center">
        <p class="text-red-600">{{ usersStore.error }}</p>
      </div>

      <div v-else-if="usersStore.users.length === 0" class="p-8 text-center">
        <p class="text-gray-500">No users found</p>
      </div>

      <table v-else class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Username</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Full Name</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Role</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Department</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="user in usersStore.users" :key="user.id" class="hover:bg-gray-50">
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm font-medium text-gray-900">{{ user.username }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm text-gray-900">{{ user.full_name }}</div>
              <div class="text-xs text-gray-500">{{ user.email }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span :class="getRoleClass(user.role)" class="px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full">
                {{ user.role }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ user.department }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span :class="user.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'" class="px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full">
                {{ user.is_active ? 'Active' : 'Inactive' }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
              <router-link
                :to="`/dashboard/users/${user.id}/edit`"
                class="text-blue-600 hover:text-blue-900 mr-4"
              >
                Edit
              </router-link>
              <button
                @click="deleteUser(user.id)"
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
import { onMounted } from 'vue'
import { useUsersStore } from '@/stores/users'

const usersStore = useUsersStore()

onMounted(() => {
  usersStore.fetchUsers()
})

const deleteUser = async (id: number) => {
  if (confirm('Are you sure you want to delete this user?')) {
    await usersStore.deleteUser(id)
  }
}

const getRoleClass = (role: string) => {
  const classes: Record<string, string> = {
    'Admin': 'bg-purple-100 text-purple-800',
    'Sales': 'bg-blue-100 text-blue-800',
    'Purchasing': 'bg-green-100 text-green-800',
    'Warehouse': 'bg-yellow-100 text-yellow-800',
    'Route': 'bg-pink-100 text-pink-800',
  }
  return classes[role] || 'bg-gray-100 text-gray-800'
}
</script>
