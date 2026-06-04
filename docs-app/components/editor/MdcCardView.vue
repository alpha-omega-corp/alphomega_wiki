<script setup lang="ts">
import { NodeViewContent, nodeViewProps, NodeViewWrapper } from '@tiptap/vue-3'
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import WCard from '~/components/content/WCard.vue'

const props = defineProps(nodeViewProps)

const variants = ['default', 'info', 'warning', 'success', 'danger']

function update(key: 'icon' | 'variant', value: string) {
  props.updateAttributes({ [key]: value })
}

const isIconName = computed(() => {
  const v = String(props.node.attrs.icon ?? '')
  return /^i-[a-z0-9]+[:-]/.test(v) || v.includes(':')
})

// Show the controls whenever the caret/selection is inside this card (i.e. the
// user is editing its content), not only when the whole node is selected.
const active = ref(false)
function refreshActive() {
  const pos = typeof props.getPos === 'function' ? props.getPos() : undefined
  if (pos == null) {
    active.value = !!props.selected
    return
  }
  const { from } = props.editor.state.selection
  active.value = !!props.selected || (from >= pos && from < pos + props.node.nodeSize)
}

onMounted(() => {
  refreshActive()
  props.editor.on('selectionUpdate', refreshActive)
  props.editor.on('transaction', refreshActive)
})
onBeforeUnmount(() => {
  props.editor.off('selectionUpdate', refreshActive)
  props.editor.off('transaction', refreshActive)
})
</script>

<template>
  <NodeViewWrapper class="mdc-card-nodeview" :class="{ 'is-active': active }">
    <!-- Controls, shown while editing inside the card. -->
    <div v-if="active" class="mdc-card-editor" contenteditable="false">
      <USelect
          :model-value="node.attrs.variant"
          :items="variants"
          size="xs"
          class="w-32"
          @update:model-value="value => update('variant', value as string)"
      />
    </div>

    <WCard :icon="node.attrs.icon" :variant="node.attrs.variant">
      <!-- The icon itself is the picker trigger: click it to choose an icon. -->
      <template #icon>
        <IconPicker
            :model-value="node.attrs.icon"
            @update:model-value="value => update('icon', value)"
        >
          <template #trigger>
            <button
                type="button"
                contenteditable="false"
                title="Pick an icon"
                class="mt-0.5 inline-flex cursor-pointer rounded p-0.5 -m-0.5 hover:bg-black/5 dark:hover:bg-white/10"
            >
              <UIcon v-if="node.attrs.icon && isIconName" :name="node.attrs.icon" class="size-5"/>
              <span v-else-if="node.attrs.icon" class="text-xl leading-none">{{ node.attrs.icon }}</span>
              <UIcon v-else name="i-lucide-image-plus" class="size-5 text-muted"/>
            </button>
          </template>
        </IconPicker>
      </template>

      <NodeViewContent class="mdc-card__content"/>
    </WCard>
  </NodeViewWrapper>
</template>

<style scoped>
.mdc-card-nodeview.is-active :deep(.mdc-card) {
  outline: 2px solid var(--ui-primary, #6366f1);
  outline-offset: 2px;
}

.mdc-card-editor {
  display: flex;
  gap: 0.5rem;
  align-items: center;
  margin-bottom: 0.5rem;
}
</style>