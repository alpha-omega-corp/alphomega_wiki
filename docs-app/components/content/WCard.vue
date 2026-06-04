<template>
  <div class="mdc-card rounded-lg border p-2" :class="styles.root">
    <div class="flex align-middle gap-3">
      <!-- Icon sits inline with the content. The #icon slot lets the editor swap in a clickable picker. -->
      <div v-if="$slots.icon || icon" class="mt-0.5">
        <slot name="icon">
          <UIcon v-if="icon && isIconName" :name="icon" class="mt-0.5 size-5" :class="styles.icon"/>
          <span v-else-if="icon" class="text-xl leading-none" aria-hidden="true">{{ icon }}</span>
        </slot>
      </div>

      <div class="mdc-card__body min-w-0 flex-1 text-sm leading-relaxed text-gray-700 dark:text-gray-300 font-bold">
        <!-- Default slot: rendered markdown placed between ::card and :: -->
        <slot/>
      </div>
    </div>

    <footer v-if="$slots.footer" class="mt-3 border-t border-black/10 pt-3 text-sm opacity-80 dark:border-white/10">
      <slot name="footer"/>
    </footer>
  </div>
</template>

<script setup lang="ts">
/**
 * WCard component for Nuxt Content (MDC).
 *
 * Components in components/content/ are auto-registered as global components and
 * become usable inside markdown rendered by <ContentRenderer>.
 *
 *   ::card{icon="i-lucide-triangle-alert" variant="warning"}
 *   Body content. **Markdown** works here.
 *   ::
 */
type Variant = 'default' | 'info' | 'warning' | 'success' | 'danger'

const props = withDefaults(
    defineProps<{
      /** Leading icon: an icon token like `i-lucide-flame` (rendered as an icon) or plain text/emoji. */
      icon?: string
      /** Visual style. */
      variant?: Variant
    }>(),
    {
      icon: '',
      variant: 'default',
    },
)

// Treat values like `i-lucide-flame` or `lucide:flame` as icon names; everything
// else (emoji, plain text) is rendered literally.
const isIconName = computed(() => /^i-[a-z0-9]+[:-]/.test(props.icon) || props.icon.includes(':'))

// Per-variant tint. Full literal class strings so Tailwind's scanner keeps them.
const VARIANTS: Record<Variant, { root: string; icon: string }> = {
  default: {
    root: 'bg-gray-50 border-gray-200 dark:bg-gray-900/50 dark:border-gray-800',
    icon: 'text-gray-500 dark:text-gray-400',
  },
  info: {
    root: 'bg-blue-50 border-blue-200 dark:bg-blue-950/40 dark:border-blue-900',
    icon: 'text-blue-500 dark:text-blue-400',
  },
  success: {
    root: 'bg-green-50 border-green-200 dark:bg-green-950/40 dark:border-green-900',
    icon: 'text-green-600 dark:text-green-400',
  },
  warning: {
    root: 'bg-amber-50 border-amber-200 dark:bg-amber-950/40 dark:border-amber-900',
    icon: 'text-amber-500 dark:text-amber-400',
  },
  danger: {
    root: 'bg-red-50 border-red-200 dark:bg-red-950/40 dark:border-red-900',
    icon: 'text-red-500 dark:text-red-400',
  },
}

const styles = computed(() => VARIANTS[props.variant] ?? VARIANTS.default)
</script>

<style scoped>
/* Tidy up spacing of markdown rendered inside the body. */
.mdc-card__body :deep(> :first-child) {
  margin-top: 0;
}

.mdc-card__body :deep(> :last-child) {
  margin-bottom: 0;
}
</style>