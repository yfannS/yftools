import type { Ref } from 'vue'

export function useScrollSync(
  editorRef: Ref<HTMLTextAreaElement | undefined>,
  previewRef: Ref<HTMLElement | undefined>,
  lineGutterRef?: Ref<HTMLElement | undefined>
) {
  let isSyncingFromEditor = false
  let isSyncingFromPreview = false

  function syncEditorToPreview() {
    if (isSyncingFromPreview || !editorRef.value || !previewRef.value) return
    isSyncingFromEditor = true

    const editor = editorRef.value
    const preview = previewRef.value
    const eMax = editor.scrollHeight - editor.clientHeight
    const pMax = preview.scrollHeight - preview.clientHeight
    const ratio = eMax > 0 ? editor.scrollTop / eMax : 0
    preview.scrollTop = ratio * pMax

    if (lineGutterRef?.value) {
      lineGutterRef.value.scrollTop = editor.scrollTop
    }

    isSyncingFromEditor = false
  }

  function syncPreviewToEditor() {
    if (isSyncingFromEditor || !editorRef.value || !previewRef.value) return
    isSyncingFromPreview = true

    const editor = editorRef.value
    const preview = previewRef.value
    const eMax = editor.scrollHeight - editor.clientHeight
    const pMax = preview.scrollHeight - preview.clientHeight
    const ratio = pMax > 0 ? preview.scrollTop / pMax : 0
    editor.scrollTop = ratio * eMax

    if (lineGutterRef?.value) {
      lineGutterRef.value.scrollTop = editor.scrollTop
    }

    isSyncingFromPreview = false
  }

  return {
    syncEditorToPreview,
    syncPreviewToEditor
  }
}
