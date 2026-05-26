import { Node, mergeAttributes } from '@tiptap/core'
import { VueNodeViewRenderer, NodeViewWrapper } from '@tiptap/vue-3'
import { defineComponent, h, ref } from 'vue'

const ImageUploadView = defineComponent({
  name: 'ImageUploadView',
  props: {
    editor: { type: Object, required: true },
    node: { type: Object, required: true },
    getPos: { type: Function, required: true },
    selected: { type: Boolean, default: false },
    deleteNode: { type: Function, required: true },
  },
  setup(props) {
    const fileInput = ref<HTMLInputElement | null>(null)
    const isDragging = ref(false)

    const replaceWithImage = (src: string) => {
      const pos = typeof props.getPos === 'function' ? props.getPos() : 0
      props.editor
        .chain()
        .focus()
        .deleteRange({ from: pos, to: pos + props.node.nodeSize })
        .insertContentAt(pos, { type: 'image', attrs: { src } })
        .run()
    }

    const handleFile = (file: File) => {
      if (!file.type.startsWith('image/')) return
      const reader = new FileReader()
      reader.onload = e => replaceWithImage(e.target?.result as string)
      reader.readAsDataURL(file)
    }

    return () =>
      h(NodeViewWrapper, { class: 'image-upload-wrapper' }, () =>
        h(
          'div',
          {
            class: [
              'image-upload-placeholder flex flex-col items-center justify-center gap-2',
              'border-2 border-dashed rounded-lg p-8 cursor-pointer transition-colors',
              'text-muted hover:text-default hover:border-default',
              props.selected ? 'border-primary' : 'border-muted',
              isDragging.value ? 'bg-elevated border-primary' : 'bg-muted/30',
            ].join(' '),
            onClick: () => fileInput.value?.click(),
            onDragover: (e: DragEvent) => { e.preventDefault(); isDragging.value = true },
            onDragleave: () => { isDragging.value = false },
            onDrop: (e: DragEvent) => {
              e.preventDefault()
              isDragging.value = false
              const file = e.dataTransfer?.files?.[0]
              if (file) handleFile(file)
            },
          },
          [
            h('input', {
              ref: fileInput,
              type: 'file',
              accept: 'image/*',
              style: 'display:none',
              onChange: (e: Event) => {
                const file = (e.target as HTMLInputElement).files?.[0]
                if (file) handleFile(file)
              },
            }),
            h('div', { class: 'i-lucide-image-plus text-3xl opacity-60' }),
            h('p', { class: 'text-sm' }, 'Click or drag to upload an image'),
          ]
        )
      )
  },
})

export const ImageUpload = Node.create({
  name: 'imageUpload',
  group: 'block',
  selectable: true,
  draggable: true,
  atom: true,

  parseHTML() {
    return [{ tag: 'image-upload' }]
  },

  renderHTML({ HTMLAttributes }) {
    return ['image-upload', mergeAttributes(HTMLAttributes)]
  },

  addNodeView() {
    return VueNodeViewRenderer(ImageUploadView)
  },
})