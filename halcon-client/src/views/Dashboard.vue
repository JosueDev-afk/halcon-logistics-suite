<template>
  <div class="flex h-screen bg-gray-100">
    <!-- Sidebar -->
    <aside class="w-64 bg-white shadow-lg">
      <div class="p-6">
        <h1 class="text-2xl font-bold text-gray-900">🦅 Halcon</h1>
        <p class="text-sm text-gray-500 mt-1">{{ authStore.user?.role }}</p>
      </div>

      <nav class="mt-6">
        <router-link
          v-for="item in menuItems"
          :key="item.path"
          :to="item.path"
          class="flex items-center px-6 py-3 text-gray-700 hover:bg-blue-50 hover:text-blue-600 transition-colors"
          active-class="bg-blue-50 text-blue-600 border-r-4 border-blue-600"
        >
          <span class="text-xl mr-3">{{ item.icon }}</span>
          <span class="font-medium">{{ item.label }}</span>
        </router-link>
      </nav>

      <div class="absolute bottom-0 w-64 p-6 border-t">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-900">{{ authStore.user?.full_name }}</p>
            <p class="text-xs text-gray-500">{{ authStore.user?.username }}</p>
          </div>
          <button
            @click="handleLogout"
            class="text-red-600 hover:text-red-700 text-sm font-medium"
          >
            Logout
          </button>
        </div>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 overflow-y-auto">
      <div class="p-8">
        <RouterView />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, RouterView } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const menuItems = computed(() => {
  const role = authStore.user?.role
  const items = [
    { path: '/dashboard/orders', label: 'Orders', icon: '📦', roles: ['Admin', 'Sales', 'Purchasing', 'Warehouse', 'Route'] },
  ]

  if (role === 'Sales') {
    items.push({ path: '/dashboard/orders/create', label: 'Create Order', icon: '➕', roles: ['Sales'] })
  }

  if (role === 'Admin' || role === 'Sales') {
    items.push({ path: '/dashboard/recycle-bin', label: 'Recycle Bin', icon: '🗑️', roles: ['Admin', 'Sales'] })
  }

  if (role === 'Admin') {
    items.push({ path: '/dashboard/users', label: 'Users', icon: '👥', roles: ['Admin'] })
  }

  return items.filter(item => item.roles.includes(role || ''))
})

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>
