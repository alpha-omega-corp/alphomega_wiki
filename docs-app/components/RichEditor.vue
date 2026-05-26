<template>
  <div class="border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden">
    <!-- Toolbar -->
    <div class="flex flex-wrap items-center gap-1 p-2 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800">
      <UButtonGroup>
        <UButton
          icon="i-heroicons-bold"
          variant="ghost"
          color="neutral"
          size="xs"
          :class="{ 'bg-gray-200 dark:bg-gray-700': editor?.isActive('bold') }"
          @click="editor?.chain().focus().toggleBold().run()"
        />
        <UButton
          icon="i-heroicons-italic"
          variant="ghost"
          color="neutral"
          size="xs"
          :class="{ 'bg-gray-200 dark:bg-gray-700': editor?.isActive('italic') }"
          @click="editor?.chain().focus().toggleItalic().run()"
        />
        <UButton
          icon="i-heroicons-underline"
          variant="ghost"
          color="neutral"
          size="xs"
          :class="{ 'bg-gray-200 dark:bg-gray-700': editor?.isActive('underline') }"
          @click="editor?.chain().focus().toggleUnderline().run()"
        />
        <UButton
          icon="i-heroicons-strikethrough"
          variant="ghost"
          color="neutral"
          size="xs"
          :class="{ 'bg-gray-200 dark:bg-gray-700': editor?.isActive('strike') }"
          @click="editor?.chain().focus().toggleStrike().run()"
        />
      </UButtonGroup>

      <UDivider orientation="vertical" class="h-6" />

      <UButtonGroup>
        <UButton
          label="H1"
          variant="ghost"
          color="neutral"
          size="xs"
          class="font-bold"
          :class="{ 'bg-gray-200 dark:bg-gray-700': editor?.isActive('heading', { level: 1 }) }"
          @click="editor?.chain().focus().toggleHeading({ level: 1 }).run()"
        />
        <UButton
          label="H2"
          variant="ghost"
          color="neutral"
          size="xs"
          class="font-bold"
          :class="{ 'bg-gray-200 dark:bg-gray-700': editor?.isActive('heading', { level: 2 }) }"
          @click="editor?.chain().focus().toggleHeading({ level: 2 }).run()"
        />
        <UButton
          label="H3"
          variant="ghost"
          color="neutral"
          size="xs"
          class="font-bold"
          :class="{ 'bg-gray-200 dark:bg-gray-700': editor?.isActive('heading', { level: 3 }) }"
          @click="editor?.chain().focus().toggleHeading({ level: 3 }).run()"
        />
      </UButtonGroup>

      <UDivider orientation="vertical" class="h-6" />

      <UButtonGroup>
        <UButton
          icon="i-heroicons-list-bullet"
          variant="ghost"
          color="neutral"
          size="xs"
          :class="{ 'bg-gray-200 dark:bg-gray-700': editor?.isActive('bulletList') }"
          @click="editor?.chain().focus().toggleBulletList().run()"
        />
        <UButton
          icon="i-heroicons-numbered-list"
          variant="ghost"
          color="neutral"
          size="xs"
          :class="{ 'bg-gray-200 dark:bg-gray-700': editor?.isActive('orderedList') }"
          @click="editor?.chain().focus().toggleOrderedList().run()"
        />
        <UButton
          icon="i-heroicons-code-bracket"
          variant="ghost"
          color="neutral"
          size="xs"
          :class="{ 'bg-gray-200 dark:bg-gray-700': editor?.isActive('codeBlock') }"
          @click="editor?.chain().focus().toggleCodeBlock().run()"
        />
        <UButton
          icon="i-heroicons-chat-bubble-left"
          variant="ghost"
          color="neutral"
          size="xs"
          :class="{ 'bg-gray-200 dark:bg-gray-700': editor?.isActive('blockquote') }"
          @click="editor?.chain().focus().toggleBlockquote().run()"
        />
      </UButtonGroup>

      <UDivider orientation="vertical" class="h-6" />

      <UButtonGroup>
        <UButton
          icon="i-heroicons-arrow-uturn-left"
          variant="ghost"
          color="neutral"
          size="xs"
          :disabled="!editor?.can().undo()"
          @click="editor?.chain().focus().undo().run()"
        />
        <UButton
          icon="i-heroicons-arrow-uturn-right"
          variant="ghost"
          color="neutral"
          size="xs"
          :disabled="!editor?.can().redo()"
          @click="editor?.chain().focus().redo().run()"
        />
      </UButtonGroup>
    </div>

    <!-- Editor Content -->
    <div class="bg-white dark:bg-gray-900">
      <EditorContent
        :editor="editor"
        class="prose dark:prose-invert max-w-none p-6 min-h-[400px] focus:outline-none"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Underline from '@tiptap/extension-underline'

const props = defineProps<{
  modelValue?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const editor = useEditor({
  content: props.modelValue || defaultContent,
  extensions: [
    StarterKit,
    Underline,
  ],
  editorProps: {
    attributes: {
      class: 'focus:outline-none',
    },
  },
  onUpdate: ({ editor }) => {
    emit('update:modelValue', editor.getHTML())
  },
})

onBeforeUnmount(() => {
  editor.value?.destroy()
})

const defaultContent = `
<h1>Welcome to the Editor</h1>
<p>This is a rich text editor powered by <strong>TipTap</strong>. You can use the toolbar above to format your content.</p>
<h2>Features</h2>
<ul>
  <li>Bold, italic, underline, and strikethrough</li>
  <li>Headings (H1, H2, H3)</li>
  <li>Bullet and numbered lists</li>
  <li>Code blocks and blockquotes</li>
  <li>Undo / Redo</li>
</ul>
<blockquote>
  <p>Start writing your documentation here...</p>
</blockquote>
`
</script>

<style>
.tiptap {
  min-height: 400px;
  outline: none;
}

.tiptap p.is-editor-empty:first-child::before {
  color: #adb5bd;
  content: attr(data-placeholder);
  float: left;
  height: 0;
  pointer-events: none;
}
</style>
