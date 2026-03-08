<script setup lang="ts">
import type { DropdownMenuItem } from '@nuxt/ui'

defineProps<{
  collapsed?: boolean
}>()

const colorMode = useColorMode()
const appConfig = useAppConfig()
const { user, clear, loggedIn } = useUserSession()

const colors = ['red', 'orange', 'amber', 'yellow', 'lime', 'green', 'emerald', 'teal', 'cyan', 'sky', 'blue', 'indigo', 'violet', 'purple', 'fuchsia', 'pink', 'rose']
const neutrals = ['slate', 'gray', 'zinc', 'neutral', 'stone']

const userDisplay = computed(() => {
  if (!user.value) {
    return {
      name: 'Guest',
      avatar: {
        src: 'https://github.com/placeholder.png',
        alt: 'Guest'
      }
    }
  }
  const fullName = user.value.firstName && user.value.lastName 
    ? `${user.value.firstName} ${user.value.lastName}`
    : user.value.firstName || user.value.email
  
  return {
    name: fullName,
    avatar: {
      src: `https://i.pravatar.cc/128?u=${user.value.id}`,
      alt: fullName
    }
  }
})

const handleLogout = async () => {
  try {
    await $fetch('/api/auth/logout', { method: 'POST' })
  } finally {
    await clear()
    navigateTo('/')
  }
}

const authItems = computed<DropdownMenuItem[][]>(() => {
  if (!loggedIn.value) {
    return [[{
      label: 'Log in',
      icon: 'i-lucide-log-in',
      to: '/login'
    }]]
  }
  
  return [[{
    type: 'label',
    label: userDisplay.value.name,
    avatar: userDisplay.value.avatar
  }], [{
    label: 'Profile',
    icon: 'i-lucide-user'
  }, {
    label: 'Billing',
    icon: 'i-lucide-credit-card'
  }, {
    label: 'Settings',
    icon: 'i-lucide-settings',
    to: '/settings'
  }], [{
    label: 'Theme',
    icon: 'i-lucide-palette',
    children: [{
      label: 'Primary',
      slot: 'chip',
      chip: appConfig.ui.colors.primary,
      content: {
        align: 'center',
        collisionPadding: 16
      },
      children: colors.map(color => ({
        label: color,
        chip: color,
        slot: 'chip',
        checked: appConfig.ui.colors.primary === color,
        type: 'checkbox',
        onSelect: (e) => {
          e.preventDefault()

          appConfig.ui.colors.primary = color
        }
      }))
    }, {
      label: 'Neutral',
      slot: 'chip',
      chip: appConfig.ui.colors.neutral === 'neutral' ? 'old-neutral' : appConfig.ui.colors.neutral,
      content: {
        align: 'end',
        collisionPadding: 16
      },
      children: neutrals.map(color => ({
        label: color,
        chip: color === 'neutral' ? 'old-neutral' : color,
        slot: 'chip',
        type: 'checkbox',
        checked: appConfig.ui.colors.neutral === color,
        onSelect: (e) => {
          e.preventDefault()

          appConfig.ui.colors.neutral = color
        }
      }))
    }]
  }, {
    label: 'Appearance',
    icon: 'i-lucide-sun-moon',
    children: [{
      label: 'Light',
      icon: 'i-lucide-sun',
      type: 'checkbox',
      checked: colorMode.value === 'light',
      onSelect(e: Event) {
        e.preventDefault()

        colorMode.preference = 'light'
      }
    }, {
      label: 'Dark',
      icon: 'i-lucide-moon',
      type: 'checkbox',
      checked: colorMode.value === 'dark',
      onUpdateChecked(checked: boolean) {
        if (checked) {
          colorMode.preference = 'dark'
        }
      },
      onSelect(e: Event) {
        e.preventDefault()
      }
    }]
  }], [{
    label: 'Templates',
    icon: 'i-lucide-layout-template',
    children: [{
      label: 'Starter',
      to: 'https://starter-template.nuxt.dev/'
    }, {
      label: 'Landing',
      to: 'https://landing-template.nuxt.dev/'
    }, {
      label: 'Docs',
      to: 'https://docs-template.nuxt.dev/'
    }, {
      label: 'SaaS',
      to: 'https://saas-template.nuxt.dev/'
    }, {
      label: 'Dashboard',
      to: 'https://dashboard-template.nuxt.dev/',
      color: 'primary',
      checked: true,
      type: 'checkbox'
    }, {
      label: 'Chat',
      to: 'https://chat-template.nuxt.dev/'
    }, {
      label: 'Portfolio',
      to: 'https://portfolio-template.nuxt.dev/'
    }, {
      label: 'Changelog',
      to: 'https://changelog-template.nuxt.dev/'
    }]
  }], [{
    label: 'Documentation',
    icon: 'i-lucide-book-open',
    to: 'https://ui.nuxt.com/docs/getting-started/installation/nuxt',
    target: '_blank'
  }, {
    label: 'GitHub repository',
    icon: 'i-simple-icons-github',
    to: 'https://github.com/nuxt-ui-templates/dashboard',
    target: '_blank'
  }, {
    label: 'Log out',
    icon: 'i-lucide-log-out',
    onSelect: () => {
      handleLogout()
    }
  }]]
})
</script>

<template>
  <UDropdownMenu
    v-if="loggedIn"
    :items="authItems"
    :content="{ align: 'center', collisionPadding: 12 }"
    :ui="{ content: collapsed ? 'w-48' : 'w-(--reka-dropdown-menu-trigger-width)' }"
  >
    <UButton
      v-bind="{
        ...userDisplay,
        label: collapsed ? undefined : userDisplay?.name,
        trailingIcon: collapsed ? undefined : 'i-lucide-chevrons-up-down'
      }"
      color="neutral"
      variant="ghost"
      block
      :square="collapsed"
      class="data-[state=open]:bg-elevated"
      :ui="{
        trailingIcon: 'text-dimmed'
      }"
    />

    <template #chip-leading="{ item }">
      <div class="inline-flex items-center justify-center shrink-0 size-5">
        <span
          class="rounded-full ring ring-bg bg-(--chip-light) dark:bg-(--chip-dark) size-2"
          :style="{
            '--chip-light': `var(--color-${(item as any).chip}-500)`,
            '--chip-dark': `var(--color-${(item as any).chip}-400)`
          }"
        />
      </div>
    </template>
  </UDropdownMenu>

  <UButton
    v-else
    to="/login"
    label="Log in"
    icon="i-lucide-log-in"
    color="neutral"
    variant="ghost"
    block
    :square="collapsed"
  />
</template>
