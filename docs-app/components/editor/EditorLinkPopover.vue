<script setup lang="ts">
import type { Editor } from '@tiptap/vue-3'

const props = defineProps<{
  editor: Editor
  autoOpen?: boolean
}>()

const href = ref('')
const isOpen = ref(false)

const isLinkActive = computed(() => props.editor.isActive('link'))

watch(isLinkActive, (active) => {
  if (props.autoOpen && active) {
    href.value = props.editor.getAttributes('link').href ?? ''
    isOpen.value = true
  }
})

watch(isOpen, (open) => {
  if (open) {
    href.value = props.editor.getAttributes('link').href ?? ''
  }
})

function apply() {
  if (!href.value) {
    props.editor.chain().focus().extendMarkRange('link').unsetLink().run()
  } else {
    props.editor.chain().focus().extendMarkRange('link').setLink({ href: href.value }).run()
  }
  isOpen.value = false
}

function remove() {
  props.editor.chain().focus().extendMarkRange('link').unsetLink().run()
  isOpen.value = false
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter') apply()
  if (e.key === 'Escape') isOpen.value = false
}
</script>

<template>
  <UPopover v-model:open="isOpen">
    <UButton
      icon="i-lucide-link"
      color="neutral"
      :variant="isLinkActive ? 'soft' : 'ghost'"
      size="sm"
      @click="isOpen = !isOpen"
    />

    <template #content>
      <div class="flex items-center gap-1.5 p-2 w-72">
        <UInput
          v-model="href"
          placeholder="https://"
          size="sm"
          class="flex-1"
          autofocus
          @keydown="onKeydown"
        />
        <UButton icon="i-lucide-check" size="sm" color="primary" variant="ghost" @click="apply" />
        <UButton v-if="isLinkActive" icon="i-lucide-unlink" size="sm" color="neutral" variant="ghost" @click="remove" />
      </div>
    </template>
  </UPopover>
</template>