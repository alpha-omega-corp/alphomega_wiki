// backendUrl resolves the base URL of the Go wiki backend. Defaults to the
// local dev address so `npm run dev` works with no setup; override with
// BACKEND_URL (e.g. http://backend:1000) when running in containers.
export function backendUrl(): string {
  return process.env.BACKEND_URL || 'http://localhost:1000'
}
