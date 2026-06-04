import { icons as lucide } from '@iconify-json/lucide'

// Names are read once from the locally-installed collection (no network), and
// the array is reused across requests.
const names = Object.keys(lucide.icons).sort()

export default defineEventHandler(() => names)