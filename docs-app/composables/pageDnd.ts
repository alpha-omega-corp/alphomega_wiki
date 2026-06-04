import type { InjectionKey, Ref } from 'vue'

// Where a dragged page lands relative to the page it was dropped on.
export type DropPosition = 'before' | 'after' | 'inside'

// Shared drag-and-drop state/handlers, provided by PageList and injected by the
// recursive PageItem so every node in the tree talks to one coordinator.
export interface PageDnd {
  draggingId: Ref<number | null>
  start: (id: number) => void
  end: () => void
  move: (draggedId: number, targetId: number, position: DropPosition) => Promise<void>
}

export const PageDndKey = Symbol('pageDnd') as InjectionKey<PageDnd>