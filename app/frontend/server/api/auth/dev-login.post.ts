export default defineEventHandler(async (event) => {
  // Dev login - create a mock user session for development
  // This bypasses the backend database which may not be set up
  await setUserSession(event, {
    user: {
      id: 1,
      email: 'dev@example.com',
      firstName: 'Dev',
      lastName: 'User'
    },
    loggedInAt: Date.now()
  })

  return { success: true }
})
