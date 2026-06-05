import { defineStore } from 'pinia'
import { ref } from 'vue'

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

  function setRenderedContent(raw: string, rendered: string, timeMs: number) {
    rawHtml.value = raw
    renderedHtml.value = rendered
    renderTimeMs.value = timeMs
    htmlSize.value = formatBytes(new Blob([rendered]).size)
  }

  function setOutline(items: OutlineItem[]) {
    outline.value = items
  }

  function setViewMode(mode: 'preview' | 'source') {
    viewMode.value = mode
  }

  return { renderedHtml, rawHtml, viewMode, renderTimeMs, htmlSize, outline, setRenderedContent, setOutline, setViewMode }
})

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}
