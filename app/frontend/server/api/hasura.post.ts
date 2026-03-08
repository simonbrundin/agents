export default defineEventHandler(async (event) => {
  const hasuraUrl = process.env.HASURA_URL || 'http://localhost:8080'
  const hasuraSecret = process.env.HASURA_ADMIN_SECRET || 'hasura-dev-secret'
  
  const body = await readBody(event)
  
  const response = await $fetch(`${hasuraUrl}/v1/graphql`, {
    method: 'POST',
    body,
    headers: {
      'Content-Type': 'application/json',
      'x-hasura-admin-secret': hasuraSecret
    }
  })
  
  return response
})
