<script setup lang="ts">
import { format, isToday } from 'date-fns'
import type { Conversation } from '~/types'

const props = defineProps<{
  conversations: Conversation[]
}>()

const conversationRefs = ref<Record<number, Element | null>>({})

const selectedConversation = defineModel<Conversation | null>()

watch(selectedConversation, () => {
  if (!selectedConversation.value) {
    return
  }
  const ref = conversationRefs.value[selectedConversation.value.id]
  if (ref) {
    ref.scrollIntoView({ block: 'nearest' })
  }
})

defineShortcuts({
  arrowdown: () => {
    const index = props.conversations.findIndex((conv: Conversation) => conv.id === selectedConversation.value?.id)

    if (index === -1) {
      selectedConversation.value = props.conversations[0]
    } else if (index < props.conversations.length - 1) {
      selectedConversation.value = props.conversations[index + 1]
    }
  },
  arrowup: () => {
    const index = props.conversations.findIndex((conv: Conversation) => conv.id === selectedConversation.value?.id)

    if (index === -1) {
      selectedConversation.value = props.conversations[props.conversations.length - 1]
    } else if (index > 0) {
      selectedConversation.value = props.conversations[index - 1]
    }
  }
})

function getLastMessage(conversation: Conversation) {
  const lastMsg = conversation.messages[conversation.messages.length - 1]
  return lastMsg ? lastMsg.body : ''
}
</script>

<template>
  <div class="overflow-y-auto divide-y divide-default">
    <div
      v-for="(conversation, index) in conversations"
      :key="index"
      :ref="(el) => { conversationRefs[conversation.id] = el as Element | null }"
    >
      <div
        class="p-4 sm:px-6 text-sm cursor-pointer border-l-2 transition-colors"
        :class="[
          conversation.unreadCount > 0 ? 'text-highlighted' : 'text-toned',
          selectedConversation && selectedConversation.id === conversation.id
            ? 'border-primary bg-primary/10'
            : 'border-bg hover:border-primary hover:bg-primary/5'
        ]"
        @click="selectedConversation = conversation"
      >
        <div class="flex items-center justify-between" :class="[conversation.unreadCount > 0 && 'font-semibold']">
          <div class="flex items-center gap-3">
            <UAvatar
              v-bind="conversation.participant.avatar"
              :alt="conversation.title"
              size="sm"
            />
            {{ conversation.title }}

            <UChip v-if="conversation.unreadCount > 0" />
          </div>

          <span>{{ isToday(new Date(conversation.lastMessageAt)) ? format(new Date(conversation.lastMessageAt), 'HH:mm') : format(new Date(conversation.lastMessageAt), 'dd MMM') }}</span>
        </div>
        <p class="truncate text-dimmed" :class="[conversation.unreadCount > 0 && 'font-semibold text-highlighted']">
          {{ getLastMessage(conversation) }}
        </p>
      </div>
    </div>
  </div>
</template>
