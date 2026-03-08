<script setup lang="ts">
import type { NavigationMenuItem } from '@nuxt/ui'
import type { Conversation, Message, User } from '~/types'

const route = useRoute()
const toast = useToast()
const { user } = useUserSession()

const open = ref(false)

const links = [
  [
    {
      label: 'Home',
      icon: 'i-lucide-house',
      to: '/',
      onSelect: () => {
        open.value = false
      }
    },
    {
      label: 'Conversations',
      icon: 'i-lucide-message-circle',
      to: '/conversations',
      onSelect: () => {
        open.value = false
      }
    },
    {
      label: 'Customers',
      icon: 'i-lucide-users',
      to: '/customers',
      onSelect: () => {
        open.value = false
      }
    },
    {
      label: 'Settings',
      to: '/settings',
      icon: 'i-lucide-settings',
      defaultOpen: true,
      type: 'trigger',
      children: [
        {
          label: 'General',
          to: '/settings',
          exact: true,
          onSelect: () => {
            open.value = false
          }
        },
        {
          label: 'Members',
          to: '/settings/members',
          onSelect: () => {
            open.value = false
          }
        },
        {
          label: 'Notifications',
          to: '/settings/notifications',
          onSelect: () => {
            open.value = false
          }
        },
        {
          label: 'Security',
          to: '/settings/security',
          onSelect: () => {
            open.value = false
          }
        }
      ]
    }
  ],
  [
    {
      label: 'Feedback',
      icon: 'i-lucide-message-circle',
      to: 'https://github.com/nuxt-ui-templates/dashboard',
      target: '_blank'
    },
    {
      label: 'Help & Support',
      icon: 'i-lucide-info',
      to: 'https://github.com/nuxt-ui-templates/dashboard',
      target: '_blank'
    }
  ]
] satisfies NavigationMenuItem[][]

const groups = computed(() => [
  {
    id: 'links',
    label: 'Go to',
    items: links.flat()
  },
  {
    id: 'code',
    label: 'Code',
    items: [
      {
        id: 'source',
        label: 'View page source',
        icon: 'i-simple-icons-github',
        to: `https://github.com/nuxt-ui-templates/dashboard/blob/main/app/pages${route.path === '/' ? '/index' : route.path}.vue`,
        target: '_blank'
      }
    ]
  }
])

interface HasuraConversation {
  id: number
  name: string | null
  repo_url: string | null
  messages: {
    id: number
    message: string
    time: string
    user: {
      id: number
      email: string
      first_name: string | null
      last_name: string | null
    }
  }[]
}

const config = useRuntimeConfig()

const { data: hasuraData } = await useAsyncData<HasuraConversation[]>('conversations', async () => {
  const userId = user.value?.id
  
  const query = `
    query GetConversations($userId: Int!) {
      conversations(
        where: { 
          _exists: {
            table: messages
            where: { user_id: { eq: $userId } }
          }
        }
        order_by: { id: asc }
      ) {
        id
        name
        repo_url
        messages(order_by: { time: asc }) {
          id
          message
          time
          user {
            id
            email
            first_name
            last_name
          }
        }
      }
    }
  `
  const response = await $fetch<{ data: { conversations: HasuraConversation[] } }>(`${config.hasuraUrl}/v1/graphql`, {
    method: 'POST',
    body: { query, variables: { userId } },
    headers: {
      'Content-Type': 'application/json',
      'x-hasura-admin-secret': config.hasuraAdminSecret
    }
  })
  return response.data.conversations
})

const conversations = computed<Conversation[]>(() => {
  if (!hasuraData.value) return []
  
  return hasuraData.value.map((conv: HasuraConversation): Conversation => {
    const firstMessage = conv.messages[0]
    const msgUser = firstMessage?.user
    const participant: User = {
      id: msgUser?.id || 0,
      name: msgUser?.first_name && msgUser?.last_name 
        ? `${msgUser.first_name} ${msgUser.last_name}` 
        : msgUser?.email || '',
      email: msgUser?.email || '',
      avatar: { src: `https://i.pravatar.cc/128?u=${msgUser?.id}` },
      status: 'subscribed',
      location: ''
    }

    const messages: Message[] = conv.messages.map((msg, idx): Message => ({
      id: msg.id,
      body: msg.message,
      date: msg.time,
      isOwn: idx % 2 === 1
    }))

    const lastMsg = conv.messages[conv.messages.length - 1]
    
    return {
      id: conv.id,
      participant,
      messages,
      unreadCount: 0,
      lastMessageAt: lastMsg?.time || ''
    }
  })
})

const selectedConversation = defineModel<Conversation | null>()

onMounted(async () => {
  const cookie = useCookie('cookie-consent')
  if (cookie.value === 'accepted') {
    return
  }

  toast.add({
    title:
      'We use first-party cookies to enhance your experience on our website.',
    duration: 0,
    close: false,
    actions: [
      {
        label: 'Accept',
        color: 'neutral',
        variant: 'outline',
        onClick: () => {
          cookie.value = 'accepted'
        }
      },
      {
        label: 'Opt out',
        color: 'neutral',
        variant: 'ghost'
      }
    ]
  })
})
</script>

<template>
  <UDashboardGroup unit="rem">
    <UDashboardSidebar id="default" v-model:open="open" collapsible resizable class="bg-elevated/25"
      :ui="{ footer: 'lg:border-t lg:border-default' }">
      <template #header="{ collapsed }">
        <TeamsMenu :collapsed="collapsed" />
      </template>

      <template #default="{ collapsed }">
        <UDashboardSearchButton :collapsed="collapsed" class="bg-transparent ring-default" />

        <UNavigationMenu :collapsed="collapsed" :items="links[0]" orientation="vertical" tooltip popover />

        <UNavigationMenu :collapsed="collapsed" :items="links[1]" orientation="vertical" tooltip class="mt-auto" />
      </template>

      <template #footer="{ collapsed }">
        <UserMenu :collapsed="collapsed" />
      </template>
    </UDashboardSidebar>

    <UDashboardSearch :groups="groups" />

    <div class="flex h-full w-full overflow-hidden">
      <UDashboardPanel id="conversations" side="left" resizable class="w-full sm:w-80 lg:w-96 shrink-0">
        <UDashboardNavbar title="Conversations">
          <template #right>
            <UButton icon="i-lucide-pen-square" color="neutral" variant="ghost" />
          </template>
        </UDashboardNavbar>

        <MessagesConversationList v-if="conversations" v-model="selectedConversation" :conversations="conversations" />
      </UDashboardPanel>

      <MessagesChatInterface v-if="selectedConversation" :conversation="selectedConversation" class="flex-1"
        @close="selectedConversation = null"
        @send-message="(msg: Message) => selectedConversation?.messages.push(msg)" />

      <div v-else class="flex-1 flex items-center justify-center">
        <div class="text-center">
          <UIcon name="i-lucide-message-circle" class="size-12 text-dimmed mb-2" />
          <p class="text-dimmed">
            Select a conversation to start messaging
          </p>
        </div>
      </div>
    </div>

    <NotificationsSlideover />
  </UDashboardGroup>
</template>
