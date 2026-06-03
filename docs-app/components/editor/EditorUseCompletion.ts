import type { EditorCustomHandlers } from '@nuxt/ui'
import type { Editor } from '@tiptap/vue-3'
import { Extension } from '@tiptap/core'

type AiAction = 'fix' | 'extend' | 'reduce' | 'simplify' | 'continue' | 'summarize' | 'translate'

export function useEditorCompletion(editorRef: Readonly<{ value: { editor: Editor } | null }>) {
  const isLoading = ref(false)

  const getSelectedText = (editor: Editor) => {
    const { from, to, empty } = editor.state.selection
    if (empty) return editor.state.doc.textContent
    return editor.state.doc.textBetween(from, to, ' ')
  }

  const runAi = async (editor: Editor, action: AiAction, extra?: Record<string, string>) => {
    const selectedText = getSelectedText(editor)
    if (!selectedText || isLoading.value) return

    isLoading.value = true
    try {
      const res = await $fetch<{ text: string }>('/api/ai/complete', {
        method: 'POST',
        body: { action, text: selectedText, ...extra },
      })

      const { from, to, empty } = editor.state.selection
      if (empty) {
        editor.chain().focus().insertContentAt(editor.state.doc.content.size - 1, ` ${res.text}`).run()
      } else {
        editor.chain().focus().insertContentAt({ from, to }, res.text).run()
      }
    } catch {
      // silently fail — AI endpoint may not be configured
    } finally {
      isLoading.value = false
    }
  }

  const makeHandler = (action: AiAction, extra?: (cmd: any) => Record<string, string>) => ({
    canExecute: (editor: Editor) => !isLoading.value && editor.isEditable,
    execute: (editor: Editor, cmd?: any) => {
      runAi(editor, action, extra?.(cmd))
      return editor.chain()
    },
    isActive: () => false,
    isDisabled: () => isLoading.value,
  })

  const handlers = {
    aiFix: makeHandler('fix'),
    aiExtend: makeHandler('extend'),
    aiReduce: makeHandler('reduce'),
    aiSimplify: makeHandler('simplify'),
    aiContinue: makeHandler('continue'),
    aiSummarize: makeHandler('summarize'),
    aiTranslate: makeHandler('translate', cmd => ({ language: cmd?.language ?? 'English' })),
  } satisfies EditorCustomHandlers

  const extension = Extension.create({ name: 'aiCompletion' })

  return { extension, handlers, isLoading }
}