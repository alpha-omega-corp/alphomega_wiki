<script setup lang="ts">
import type {Page} from "~/types/page";

// Controlled by PageList: `open` toggles the modal, `parentId` is the page the
// new page will be nested under (null = top level).
const open = defineModel<boolean>('open', {default: false})
const parentId = defineModel<number | null>('parentId', {default: null})

const title = ref('')
const slug = ref('')
const submitting = ref(false)
// Once the user edits the slug by hand, stop overwriting it from the title.
const slugEdited = ref(false)

function slugify(value: string) {
  return value
      .toLowerCase()
      .trim()
      .normalize('NFKD')
      .replace(/[̀-ͯ]/g, '') // strip accents
      .replace(/[^a-z0-9]+/g, '-')     // non-alphanumerics -> hyphen
      .replace(/^-+|-+$/g, '')         // trim leading/trailing hyphens
}

// Mirror the title into the slug on every keystroke, until the user overrides
// the slug by hand.
function onTitleInput(value: string) {
  title.value = value
  if (!slugEdited.value) slug.value = slugify(value)
}

function onSlugInput(value: string) {
  slugEdited.value = true
  slug.value = slugify(value)
}

watch(open, (value) => {
  if (value) {
    title.value = ''
    slug.value = ''
    slugEdited.value = false
  }
})

async function create() {
  if (!title.value.trim()) return
  submitting.value = true
  try {
    await $fetch<Page>('/api/pages', {
      method: 'POST',
      body: {
        title: title.value.trim(),
        slug: slug.value,
        parent_id: parentId.value,
      },
    })
    // Refresh the shared 'pages' tree so the sidebar shows the new page.
    await refreshNuxtData('pages')
    open.value = false
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <UModal v-model:open="open" title="New page">
    <template #body>
      <UForm :state="{title, slug}" class="space-y-4" @submit="create">
        <UFormField label="Title" name="title" required>
          <UInput
              :model-value="title"
              placeholder="Page title"
              autofocus
              class="w-full"
              @update:model-value="onTitleInput"
              @keydown.enter.prevent="create"
          />
        </UFormField>

        <UFormField label="Slug" name="slug" help="Auto-generated from the title">
          <UInput
              :model-value="slug"
              placeholder="page-slug"
              class="w-full"
              :ui="{
                base: 'pl-13',
                leading: 'pointer-events-none'
              }"
              @update:model-value="onSlugInput"
              @keydown.enter.prevent="create"
          >
            <template #leading>
              <p class="text-sm text-muted">
                docs /
              </p>
            </template>
          </UInput>
        </UFormField>
      </UForm>
    </template>

    <template #footer>
      <div class="flex justify-end gap-2 w-full">
        <UButton label="Cancel" color="neutral" variant="ghost" @click="open = false"/>
        <UButton
            label="Create"
            color="primary"
            :loading="submitting"
            :disabled="!title.trim()"
            @click="create"
        />
      </div>
    </template>
  </UModal>
</template>

<style scoped>

</style>