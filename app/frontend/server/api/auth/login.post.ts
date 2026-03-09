export default defineEventHandler(async (event) => {
  const body = await readBody(event)
  const { email, password } = body

  if (!email || !password) {
    throw createError({
      statusCode: 400,
      message: 'Email and password are required'
    })
  }

  // Use localhost when running locally, backend service name when in Docker
  const isDocker = process.env.DOCKER === 'true' || process.env.NUXT_PUBLIC_API_BASE?.includes('backend')
  const backendUrl = isDocker ? 'http://backend:8000' : 'http://localhost:8000'

  try {
    const formData = new URLSearchParams()
    formData.append('username', email)
    formData.append('password', password)

    const tokenResponse = await $fetch<{ access_token: string, token_type: string }>(`${backendUrl}/api/auth/login`, {
      method: 'POST',
      body: formData,
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
      }
    })

    const userResponse = await $fetch<{ id: number, email: string, first_name: string | null, last_name: string | null }>(`${backendUrl}/api/auth/me`, {
      headers: {
        Authorization: `Bearer ${tokenResponse.access_token}`
      }
    })

    await setUserSession(event, {
      user: {
        id: userResponse.id,
        email: userResponse.email,
        firstName: userResponse.first_name,
        lastName: userResponse.last_name
      },
      loggedInAt: Date.now()
    })

    return { success: true }
  } catch (error: any) {
    throw createError({
      statusCode: 401,
      message: error.data?.detail || 'Invalid credentials'
    })
  }
})
