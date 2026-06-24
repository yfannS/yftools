import { defineStore } from 'pinia'
import { ref } from 'vue'
import DOMPurify from 'dompurify'

// DOMPurify 配置：允许渲染后的 Markdown HTML，禁止所有危险属性/标签
const SANITIZE_OPTS: Parameters<typeof DOMPurify.sanitize>[1] = {
  ADD_TAGS: ['mermaid'],
  ADD_ATTR: ['class', 'id', 'data-mermaid-source'],
  FORBID_TAGS: ['script', 'iframe', 'object', 'embed', 'form', 'input', 'textarea', 'select', 'button'],
  FORBID_ATTR: ['on*', 'formaction', 'formmethod', 'action', 'xmlns'],
  ALLOW_DATA_ATTR: false,
  SANITIZE_DOM: true,
}

export interface OutlineItem {
  id: string
  text: string
  level: number // 1-6
}

export const usePreviewStore = defineStore('preview', () => {
  const renderedHtml = ref('')
  const rawHtml = ref('')
  const viewMode = ref<'preview' | 'source'>('preview')
  const renderTimeMs = ref(0)
  const htmlSize = ref('0 B')
  const outline = ref<OutlineItem[]>([])
  const previewScale = ref(1)

  function setRenderedContent(raw: string, rendered: string, timeMs: number) {
    rawHtml.value = raw
    renderedHtml.value = DOMPurify.sanitize(rendered, SANITIZE_OPTS)
    renderTimeMs.value = timeMs
    htmlSize.value = formatBytes(new Blob([renderedHtml.value]).size)
  }

  function setOutline(items: OutlineItem[]) {
    outline.value = items
  }

  function setViewMode(mode: 'preview' | 'source') {
    viewMode.value = mode
  }

  function setPreviewScale(scale: number) {
    previewScale.value = Math.min(1.25, Math.max(0.85, scale))
  }

  return {
    renderedHtml, rawHtml, viewMode, renderTimeMs, htmlSize, outline, previewScale,
    setRenderedContent, setOutline, setViewMode, setPreviewScale,
  }
})

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}
