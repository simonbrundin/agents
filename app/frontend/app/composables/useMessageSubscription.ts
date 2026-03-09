import { ref, reactive, onUnmounted, watch, type Ref } from 'vue'
import { createClient } from 'graphql-ws'

interface Message {
  id: number
  message: string
  time: string
  user: {
    id: number
    email: string
    first_name: string | null
    last_name: string | null
  }
}

export const useMessageSubscription = (conversationId: number | Ref<number>) => {
  const config = useRuntimeConfig()
  const messages = reactive<Message[]>([])
  const isLoading = ref(true)
  const error = ref<Error | null>(null)

  const hasuraUrl = config.public.hasuraUrl || 'http://localhost:8080'
  const wsUrl = hasuraUrl.replace(/^http/, 'ws') + '/v1/graphql'

  let wsClient: ReturnType<typeof createClient> | null = null
  let dispose: (() => void) | null = null

  const subscriptionQuery = `
    subscription OnNewMessage($conversationId: Int!) {
      messages(
        where: { conversation_id: { _eq: $conversationId } },
        order_by: { time: desc },
        limit: 10
      ) {
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
  `

  const getConversationId = () => {
    if (typeof conversationId === 'number') {
      return conversationId
    }
    return conversationId.value
  }

  const subscribe = () => {
    if (dispose) {
      dispose()
    }
    if (wsClient) {
      wsClient.dispose()
    }

    const currentId = getConversationId()
    if (!currentId) return

    wsClient = createClient({
      url: wsUrl,
      connectionParams: {
        'x-hasura-admin-secret': config.public.hasuraAdminSecret || 'hasura-dev-secret'
      },
      on: {
        error: (err) => {
          console.error('WebSocket error:', err)
          error.value = err
        }
      }
    })

    dispose = wsClient.subscribe(
      {
        query: subscriptionQuery,
        variables: { conversationId: currentId }
      },
      {
        next: (data: unknown) => {
          if (data && typeof data === 'object' && 'data' in data) {
            const responseData = data as { data?: { messages?: Message[] } }
            if (responseData.data?.messages) {
              const newMessages = responseData.data.messages
              const existingIds = new Set(messages.map(m => m.id))

              newMessages.forEach((msg) => {
                if (!existingIds.has(msg.id)) {
                  messages.unshift(msg)
                }
              })

              messages.sort((a, b) => new Date(b.time).getTime() - new Date(a.time).getTime())
            }
          }
          isLoading.value = false
        },
        error: (err: unknown) => {
          console.error('Subscription error:', err)
          error.value = err instanceof Error ? err : new Error(String(err))
          isLoading.value = false
        },
        complete: () => {
          console.log('Subscription complete')
        }
      }
    )
  }

  subscribe()

  if (typeof conversationId !== 'number') {
    watch(conversationId, () => {
      messages.length = 0
      subscribe()
    })
  }

  onUnmounted(() => {
    if (dispose) {
      dispose()
    }
    if (wsClient) {
      wsClient.dispose()
    }
  })

  return {
    messages,
    isLoading,
    error,
    reconnect: subscribe
  }
}
