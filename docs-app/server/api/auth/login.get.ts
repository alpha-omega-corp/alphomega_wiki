import {createHash, randomBytes} from 'node:crypto'

// GET /api/auth/login — start the OAuth 2.0 Authorization Code + PKCE flow.
// Generates state + PKCE, stashes them in the session, and redirects the
// browser to AOAuth's /authorize endpoint.
export default defineEventHandler(async (event) => {
    const oauth = oauthConfig()
    console.log(oauth)

    const codeVerifier = randomBytes(32).toString('base64url')
    const codeChallenge = createHash('sha256').update(codeVerifier).digest('base64url')
    const state = randomBytes(16).toString('base64url')

    // Only allow internal return paths to avoid open-redirects.
    const requested = String(getQuery(event).returnTo || '/')
    const returnTo = requested.startsWith('/') ? requested : '/'

    const session = await authSession(event)
    await session.update({state, codeVerifier, returnTo})

    const url = new URL(`${oauth.issuer}/authorize`)
    url.searchParams.set('response_type', 'code')
    url.searchParams.set('client_id', oauth.clientId)
    url.searchParams.set('redirect_uri', oauth.redirectUri)
    url.searchParams.set('scope', 'openid profile email')
    url.searchParams.set('state', state)
    url.searchParams.set('code_challenge', codeChallenge)
    url.searchParams.set('code_challenge_method', 'S256')

    return sendRedirect(event, url.toString())
})
