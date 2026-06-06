import type { H3Event } from 'h3'

// SessionData is what we seal into the encrypted session cookie. Tokens live
// here (server-side only) and are never sent to the browser; /api/auth/session
// exposes just `user`.
export interface SessionData {
  user?: {
    sub: string
    name?: string
    email?: string
    email_verified?: boolean
  }
  accessToken?: string
  refreshToken?: string
  expiresAt?: number

  // Transient values used only between /api/auth/login and /callback.
  state?: string
  codeVerifier?: string
  returnTo?: string
}

// oauthConfig resolves the OAuth settings, preferring runtime process.env over
// the build-time runtimeConfig defaults. Reading process.env at runtime keeps
// the same AOAUTH_* names working regardless of build caching or Nuxt's NUXT_
// override convention.
export function oauthConfig() {
  const c = (useRuntimeConfig().oauth ?? {}) as Record<string, string>
  return {
    issuer: process.env.AOAUTH_ISSUER || c.issuer || 'http://localhost:8080',
    clientId: process.env.AOAUTH_CLIENT_ID || c.clientId || 'demo-client',
    clientSecret: process.env.AOAUTH_CLIENT_SECRET || c.clientSecret || 'demo-secret',
    redirectUri: process.env.AOAUTH_REDIRECT_URI || c.redirectUri || 'http://localhost:3000/callback',
  }
}

// internalIssuerUrl returns the AOAuth base URL for server-side calls
// (token exchange, userinfo, token refresh). When running in a container,
// AOAUTH_INTERNAL_URL should point to the auth server's container-reachable
// address while AOAUTH_ISSUER stays as the browser-facing URL.
export function internalIssuerUrl(): string {
  return process.env.AOAUTH_INTERNAL_URL || oauthConfig().issuer
}

// bearer builds an Authorization header for the upstream Go backend, or {} when
// there is no token.
export function bearer(token?: string): Record<string, string> {
  return token ? { Authorization: `Bearer ${token}` } : {}
}

interface RefreshResponse {
  access_token: string
  refresh_token?: string
  expires_in?: number
}

// refreshTokens exchanges a refresh token for a fresh access token at AOAuth.
async function refreshTokens(refreshToken: string): Promise<RefreshResponse> {
  const oauth = oauthConfig()
  return await $fetch<RefreshResponse>(`${internalIssuerUrl()}/token`, {
    method: 'POST',
    headers: {
      Authorization:
        'Basic ' + Buffer.from(`${oauth.clientId}:${oauth.clientSecret}`).toString('base64'),
    },
    body: new URLSearchParams({ grant_type: 'refresh_token', refresh_token: refreshToken }),
  })
}

// requireAuth enforces an authenticated session for a request and returns a
// valid access token, silently refreshing (and rotating) it when it has expired.
// Throws a 401 when there is no session or the refresh fails.
export async function requireAuth(event: H3Event) {
  const session = await authSession(event)
  const data = session.data

  if (!data.user) {
    throw createError({ statusCode: 401, statusMessage: 'Authentication required' })
  }

  let accessToken = data.accessToken
  // Refresh ~30s before expiry to avoid races.
  const expiringSoon = data.expiresAt !== undefined && Date.now() > data.expiresAt - 30_000
  if (expiringSoon && data.refreshToken) {
    const refreshed = await refreshTokens(data.refreshToken).catch(() => null)
    if (!refreshed) {
      await session.clear()
      throw createError({ statusCode: 401, statusMessage: 'Session expired, please sign in again' })
    }
    accessToken = refreshed.access_token
    await session.update({
      accessToken: refreshed.access_token,
      // AOAuth rotates refresh tokens — keep the new one.
      refreshToken: refreshed.refresh_token ?? data.refreshToken,
      expiresAt: refreshed.expires_in ? Date.now() + refreshed.expires_in * 1000 : undefined,
    })
  }

  return { user: data.user, accessToken }
}

// authSession opens the sealed (AEAD-encrypted, httpOnly) session for this
// request. sameSite=lax so the cookie survives the top-level redirect back from
// the authorization server.
export function authSession(event: H3Event) {
  const password =
    process.env.AUTH_SESSION_PASSWORD ||
    (useRuntimeConfig().session?.password as string)
  return useSession<SessionData>(event, {
    password,
    name: 'docs_session',
    // sameSite=lax so the cookie survives the OAuth redirect; Secure in
    // production, relaxed in dev so it works over http://localhost.
    cookie: { sameSite: 'lax', secure: !import.meta.dev },
  })
}
