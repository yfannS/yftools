import { nextTick, watch, type Ref } from 'vue'
import { useEditorStore } from '@/stores/editor'

type ToastFn = (msg: string, type?: string) => void
type RenderFn = () => void

// 供 FindBox 组件 inject 的 API 类型
export interface FindReplaceAPI {
  findNext: () => void
  findPrev: () => void
  replaceOne: () => void
  replaceAll: () => void
}

export function useFindReplace(
  editorRef: Ref<HTMLTextAreaElement | undefined>,
  toast: ToastFn,
  render: RenderFn
) {
  const editorStore = useEditorStore()

  function escapeRegExp(s: string): string {
    return s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  }

  function buildFindRegExp(global = true): RegExp | null {
    const q = editorStore.findQuery
    if (!q) return null
    const boundary = editorStore.findWholeWord ? '\\b' : ''
    const flags = (global ? 'g' : '') + (editorStore.findCaseSensitive ? '' : 'i')
    return new RegExp(boundary + escapeRegExp(q) + boundary, flags)
  }

  function getMatches(text: string) {
    const re = buildFindRegExp(true)
    if (!re) return []
    return Array.from(text.matchAll(re)).map(match => ({
      index: match.index || 0,
      text: match[0],
    }))
  }

  function refreshMatches() {
    const editor = editorRef.value
    const matches = getMatches(editor?.value || editorStore.content)
    const s = editor?.selectionStart ?? 0
    const e = editor?.selectionEnd ?? 0
    const current = matches.findIndex(match => match.index === s && match.index + match.text.length === e)
    editorStore.setFindStats(current >= 0 ? current + 1 : 0, matches.length)
  }

  function selectMatch(index: number, matches: Array<{ index: number; text: string }>) {
    const editor = editorRef.value
    if (!editor || !matches.length) return
    const match = matches[index]
    editor.focus()
    editor.setSelectionRange(match.index, match.index + match.text.length)
    editorStore.setFindStats(index + 1, matches.length)
  }

  function findNext() {
    const editor = editorRef.value
    const q = editorStore.findQuery
    if (!q || !editor) return

    const matches = getMatches(editor.value)
    if (!matches.length) {
      editorStore.setFindStats(0, 0)
      toast('未找到', 'err')
      return
    }
    const from = editor.selectionEnd
    const idx = matches.findIndex(match => match.index >= from)
    selectMatch(idx === -1 ? 0 : idx, matches)
  }

  function findPrev() {
    const editor = editorRef.value
    const q = editorStore.findQuery
    if (!q || !editor) return

    const matches = getMatches(editor.value)
    if (!matches.length) {
      editorStore.setFindStats(0, 0)
      toast('未找到', 'err')
      return
    }
    const from = editor.selectionStart
    let idx = matches.findIndex(match => match.index >= from) - 1
    if (idx < 0) idx = matches.length - 1
    selectMatch(idx, matches)
  }

  function replaceOne() {
    const q = editorStore.findQuery
    const editor = editorRef.value
    if (!q || !editor) return

    const rep = editorStore.replaceQuery
    const s = editor.selectionStart
    const e = editor.selectionEnd

    const selected = editor.value.slice(s, e)
    const re = buildFindRegExp(false)
    if (re && re.test(selected) && selected.length > 0) {
      editorStore.setContent(editor.value.slice(0, s) + rep + editor.value.slice(e))
      nextTick(() => {
        if (editorRef.value) {
          editorRef.value.setSelectionRange(s, s + rep.length)
          refreshMatches()
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
    const re = buildFindRegExp(true)
    if (!re) return
    const count = getMatches(editor.value).length
    editorStore.setContent(editor.value.replace(re, rep))
    render()
    nextTick(refreshMatches)
    toast(`已替换 ${count} 处`)
  }

  watch(
    () => [
      editorStore.content,
      editorStore.findQuery,
      editorStore.findCaseSensitive,
      editorStore.findWholeWord,
    ],
    refreshMatches,
  )

  return {
    refreshMatches,
    findNext,
    findPrev,
    replaceOne,
    replaceAll
  }
}
