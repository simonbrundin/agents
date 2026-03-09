<script setup lang="ts">
const config = useRuntimeConfig()
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)
const devLoading = ref(false)

const isDev = computed(() => process.env.NODE_ENV === 'development')

const handleLogin = async () => {
  error.value = ''
  loading.value = true

  try {
    await $fetch('/api/auth/login', {
      method: 'POST',
      body: { email: email.value, password: password.value }
    })
    navigateTo('/conversations')
  } catch (e: any) {
    error.value = e.data?.message || 'Login failed'
  } finally {
    loading.value = false
  }
}

const handleDevLogin = async () => {
  error.value = ''
  devLoading.value = true

  try {
    await $fetch('/api/auth/dev-login', {
      method: 'POST'
    })
    navigateTo('/conversations')
  } catch (e: any) {
    error.value = e.data?.message || 'Dev login failed'
  } finally {
    devLoading.value = false
  }
}
</script>

<template>
  <div
    class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900"
  >
    <UCard class="w-full max-w-md">
      <template #header>
        <h1 class="text-xl font-semibold">
          Sign in
        </h1>
      </template>

      <form class="space-y-4" @submit.prevent="handleLogin">
        <UFormField label="Email">
          <UInput
            v-model="email"
            type="email"
            placeholder="you@example.com"
            required
            block
          />
        </UFormField>

        <UFormField label="Password">
          <UInput
            v-model="password"
            type="password"
            placeholder="••••••••"
            required
            block
          />
        </UFormField>

        <UAlert v-if="error" color="error" variant="subtle">
          {{ error }}
        </UAlert>

        <UButton type="submit" :loading="loading" block>
          Sign in
        </UButton>
      </form>

      <div
        v-if="isDev"
        class="mt-6 pt-6 border-t border-gray-200 dark:border-gray-700"
      >
        <UButton
          color="secondary"
          variant="outline"
          :loading="devLoading"
          block
          @click="handleDevLogin"
        >
          Dev Login
        </UButton>
        <p class="text-xs text-dimmed mt-2 text-center">
          Auto-login for development
        </p>
      </div>
    </UCard>
  </div>
</template>
