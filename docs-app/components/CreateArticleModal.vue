<template>
  <UModal v-model:open="open" title="New Page">
    <template #body>
      <div class="space-y-4">
        <UFormField label="Title" name="title" required>
          <UInput
            v-model="form.title"
            placeholder="Page title"
            class="w-full"
            @keydown.enter="submit"
          />
        </UFormField>

        <UFormField label="Parent page" name="parentId">
          <USelect
            v-model="form.parentId"
            :items="parentItems"
            value-key="value"
            placeholder="None (top level)"
            class="w-full"
          />
        </UFormField>
      </div>
    </template>

    <template #footer>
      <div class="flex justify-end gap-2 w-full">
        <UButton variant="ghost" color="neutral" @click="open = false">
          Cancel
        </UButton>
        <UButton :loading="loading" :disabled="!form.title.trim()" @click="submit">
          Create Page
        </UButton>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import type { NavPage } from '~/types/navigation'

const props = defineProps<{
  pages: NavPage[]
}>()

const emit = defineEmits<{
  created: []
}>()

const open = defineModel<boolean>('open', { default: false })
const parentId = defineModel<number | null>('parentId', { default: null })

const loading = ref(false)
const form = reactive({ title: '', parentId: null as number | null })

const parentItems = computed(() =>
  props.pages.map(p => ({ label: p.title, value: p.id }))
)

watch(parentId, (val) => { form.parentId = val }, { immediate: true })
watch(open, (val) => { if (!val) { form.title = ''; form.parentId = parentId.value } })

async function submit() {
  if (!form.title.trim()) return
  loading.value = true
  try {
    await $fetch('/api/articles', {
      method: 'POST',
      body: { title: form.title.trim(), parent_id: form.parentId },
    })
    open.value = false
    emit('created')
  }
  finally {
    loading.value = false
  }
}
</script>