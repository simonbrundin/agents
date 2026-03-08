<script setup lang="ts">
import { format, isToday } from 'date-fns'
import type { Conversation, Message } from '~/types'

const props = defineProps<{
  conversation: Conversation
  userId?: number
}>()

const emits = defineEmits(['close', 'send-message'])

const config = useRuntimeConfig()
const toast = useToast()

const userId = props.userId || 1

const messagesContainer = ref<HTMLElement | null>(null)
const newMessage = ref('')
const loading = ref(false)

onMounted(() => {
  scrollToBottom()
})

watch(() => props.conversation.messages, () => {
  nextTick(() => {
    scrollToBottom()
  })
}, { deep: true })

function scrollToBottom() {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

function formatMessageTime(date: string) {
  const d = new Date(date)
  if (isToday(d)) {
    return format(d, 'HH:mm')
  }
  return format(d, 'dd MMM HH:mm')
}

async function onSubmit() {
  if (!newMessage.value.trim()) return

  loading.value = true

  const mutation = `
    mutation InsertMessage($conversationId: Int!, $userId: Int!, $message: String!) {
      insert_messages_one(object: {
        conversation_id: $conversationId,
        user_id: $userId,
        message: $message
      }) {
        id
        message
        time
      }
    }
  `

  const hasuraUrl = 'http://localhost:8080'
  const hasuraSecret = 'hasura-dev-secret'

  console.log('Config hasuraUrl:', config.public.hasuraUrl)
  console.log('Config hasuraAdminSecret:', config.public.hasuraAdminSecret)
  console.log('Using hasuraUrl:', hasuraUrl)
  console.log('Conversation ID:', props.conversation.id)
  console.log('User ID:', userId)

  try {
    console.log('Making fetch to:', `${hasuraUrl}/v1/graphql`)
    const res = await fetch(`${hasuraUrl}/v1/graphql`, {
      method: 'POST',
      body: JSON.stringify({
        query: mutation,
        variables: {
          conversationId: props.conversation.id,
          userId,
          message: newMessage.value
        }
      }),
      headers: {
        'Content-Type': 'application/json',
        'x-hasura-admin-secret': hasuraSecret,
        'Origin': 'http://localhost:3000'
      }
    })
    const response = await res.json()
    console.log('Response status:', res.status)
    console.log('Response:', response)

    const newMsg: Message = {
      id: props.conversation.messages.length + 1,
      body: newMessage.value,
      date: new Date().toISOString(),
      isOwn: true
    }

    emits('send-message', newMsg)
    newMessage.value = ''

    toast.add({
      title: 'Message sent',
      description: 'Your message has been sent successfully',
      icon: 'i-lucide-check-circle',
      color: 'success'
    })
  } catch (err) {
    console.error('Failed to send message:', err)
    toast.add({
      title: 'Failed to send message',
      description: 'Please try again',
      icon: 'i-lucide-x-circle',
      color: 'error'
    })
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <UDashboardPanel id="chat">
    <UDashboardNavbar :title="conversation.participant.name" :toggle="false">
      <template #leading>
        <UButton
          icon="i-lucide-x"
          color="neutral"
          variant="ghost"
          class="-ms-1.5"
          @click="emits('close')"
        />
      </template>

      <template #right>
        <UTooltip text="Voice call">
          <UButton
            icon="i-lucide-phone"
            color="neutral"
            variant="ghost"
          />
        </UTooltip>

        <UTooltip text="Video call">
          <UButton icon="i-lucide-video" color="neutral" variant="ghost" />
        </UTooltip>

        <UTooltip text="More options">
          <UButton
            icon="i-lucide-ellipsis-vertical"
            color="neutral"
            variant="ghost"
          />
        </UTooltip>
      </template>
    </UDashboardNavbar>

    <div class="flex items-center gap-4 p-4 sm:px-6 border-b border-default">
      <UAvatar
        v-bind="conversation.participant.avatar"
        :alt="conversation.participant.name"
        size="lg"
      />

      <div class="min-w-0">
        <p class="font-semibold text-highlighted">
          {{ conversation.participant.name }}
        </p>
        <p class="text-muted text-sm">
          {{ conversation.participant.email }}
        </p>
      </div>
    </div>

    <div
      ref="messagesContainer"
      class="flex-1 p-4 sm:p-6 overflow-y-auto space-y-4"
    >
      <div
        v-for="message in conversation.messages"
        :key="message.id"
        class="flex"
        :class="[message.isOwn ? 'justify-end' : 'justify-start']"
      >
        <div
          class="max-w-[75%] rounded-2xl px-4 py-2.5"
          :class="[
            message.isOwn
              ? 'bg-primary text-primary-fg rounded-br-sm'
              : 'bg-muted text-muted-fg rounded-bl-sm'
          ]"
        >
          <p class="whitespace-pre-wrap">
            {{ message.body }}
          </p>
          <p
            class="text-xs mt-1.5 opacity-70"
            :class="[message.isOwn ? 'text-right' : 'text-left']"
          >
            {{ formatMessageTime(message.date) }}
          </p>
        </div>
      </div>
    </div>

    <div class="pb-4 px-4 sm:px-6 shrink-0">
      <UCard variant="subtle" class="mt-auto">
        <form @submit.prevent="onSubmit">
          <div class="flex items-end gap-2">
            <UTextarea
              v-model="newMessage"
              color="neutral"
              variant="none"
              required
              autoresize
              placeholder="Type a message..."
              :rows="1"
              :disabled="loading"
              class="flex-1"
              :ui="{ base: 'resize-none' }"
              @keydown.enter.exact.prevent="onSubmit"
            />

            <UButton
              type="submit"
              color="primary"
              :loading="loading"
              icon="i-lucide-send"
            />
          </div>
        </form>
      </UCard>
    </div>
  </UDashboardPanel>
</template>
