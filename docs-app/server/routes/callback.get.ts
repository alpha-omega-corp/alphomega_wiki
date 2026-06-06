// GET /callback — OAuth redirect target (registered as the client's
// redirect_uri). Validates state, exchanges the authorization code for tokens
// (server-side, using the client secret + PKCE verifier), fetches the user
// profile, and stores everything in the sealed session.
//
// This is a Nitro route (not a Nuxt page) so it runs purely on the server.
import {internalIssuerUrl} from "#server/utils/session";

interface TokenResponse {
    access_token: string
    refresh_token?: string
    expires_in?: number
    id_token?: string
}

interface UserInfo {
    sub: string
    name?: string
    email?: string
    email_verified?: boolean
}

export default defineEventHandler(async (event) => {
    const oauth = oauthConfig()
    const query = getQuery(event)
    const session = await authSession(event)

    // The authorization server reported an error (e.g. access_denied).
    if (query.error) {
        return sendRedirect(event, `/?auth_error=${encodeURIComponent(String(query.error))}`)
    }

    const code = String(query.code || '')
    const state = String(query.state || '')
    const {state: expectedState, codeVerifier, returnTo = '/'} = session.data

    if (!code || !state || !expectedState || state !== expectedState) {
        return sendRedirect(event, '/?auth_error=invalid_state')
    }

    try {
        const token = await $fetch<TokenResponse>(`${internalIssuerUrl()}/token`, {
            method: 'POST',
            headers: {
                Authorization:
                    'Basic ' + Buffer.from(`${oauth.clientId}:${oauth.clientSecret}`).toString('base64'),
            },
            body: new URLSearchParams({
                grant_type: 'authorization_code',
                code,
                redirect_uri: oauth.redirectUri,
                code_verifier: codeVerifier || '',
            }),
        })

        const userinfo = await $fetch<UserInfo>(`${internalIssuerUrl()}/userinfo`, {
            headers: {Authorization: `Bearer ${token.access_token}`},
        }).catch(() => null)

        await session.update({
            user: {
                sub: userinfo?.sub || '',
                name: userinfo?.name,
                email: userinfo?.email,
                email_verified: userinfo?.email_verified,
            },
            accessToken: token.access_token,
            refreshToken: token.refresh_token,
            expiresAt: token.expires_in ? Date.now() + token.expires_in * 1000 : undefined,
            // Clear the one-time flow values.
            state: undefined,
            codeVerifier: undefined,
            returnTo: undefined,
        })

        return sendRedirect(event, returnTo)
    } catch {
        return sendRedirect(event, '/?auth_error=token_exchange_failed')
    }
})
