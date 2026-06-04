// POST /api/auth/logout — clears the local session. (The AOAuth access token
// is short-lived; a full RP-initiated logout at the provider could be added
// here later.)
export default defineEventHandler(async (event) => {
  const session = await authSession(event)
  await session.clear()
  return { ok: true }
})
