import type { Ref } from 'vue'

export function useKaTeX(previewRef: Ref<HTMLElement | undefined>) {

  function processMath() {
    const container = previewRef.value
    if (!container) return

    try {
      ;(window as any).renderMathInElement(container, {
        delimiters: [
          { left: '$$', right: '$$', display: true },
          { left: '$', right: '$', display: false },
          { left: '\\(', right: '\\)', display: false },
          { left: '\\[', right: '\\]', display: true }
        ],
        throwOnError: false
      })
    } catch (_) { /* ignore */ }
  }

  return {
    processMath
  }
}
