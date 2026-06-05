import { nextTick, type Ref } from 'vue'
import { useEditorStore } from '@/stores/editor'

type ToastFn = (msg: string, type?: string) => void
type RenderFn = () => void

// 供 FindBox 组件 inject 的 API 类型
export interface FindReplaceAPI {
  findNext: () => void
  replaceOne: () => void
  replaceAll: () => void
}

export function useFindReplace(
  editorRef: Ref<HTMLTextAreaElement | undefined>,
  toast: ToastFn,
  render: RenderFn
) {
  const editorStore = useEditorStore()

  function findNext() {
    const q = editorStore.findQuery
    const editor = editorRef.value
    if (!q || !editor) return

    const from = editor.selectionEnd
    const text = editor.value
    let idx = text.indexOf(q, from)
    if (idx === -1) idx = text.indexOf(q, 0)
    if (idx === -1) {
      toast('未找到', 'err')
      return
    }
    editor.focus()
    editor.setSelectionRange(idx, idx + q.length)
  }

  function replaceOne() {
    const q = editorStore.findQuery
    const editor = editorRef.value
    if (!q || !editor) return

    const rep = editorStore.replaceQuery
    const s = editor.selectionStart
    const e = editor.selectionEnd

    if (editor.value.slice(s, e) === q) {
      editorStore.setContent(editor.value.slice(0, s) + rep + editor.value.slice(e))
      nextTick(() => {
        if (editorRef.value) {
          editorRef.value.setSelectionRange(s, s + rep.length)
        }
      })
    } else {
      findNext()
      return
    }
    render()
  }

  function replaceAll() {
    const q = editorStore.findQuery
    const editor = editorRef.value
    if (!q || !editor) return

    const rep = editorStore.replaceQuery
    const escaped = q.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
    const count = (editor.value.match(new RegExp(escaped, 'g')) || []).length
    editorStore.setContent(editor.value.split(q).join(rep))
    render()
    toast(`已替换 ${count} 处`)
  }

  return {
    findNext,
    replaceOne,
    replaceAll
  }
}
