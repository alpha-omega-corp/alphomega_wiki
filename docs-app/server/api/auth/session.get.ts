// GET /api/auth/session — returns the current user (or null). Tokens are never
// included; they stay sealed in the session cookie.
export default defineEventHandler(async (event) => {
  const session = await authSession(event)
  return { user: session.data.user ?? null }
})
