import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useEditorStore = defineStore('editor', () => {
  const content = ref('')
  const filename = ref('document')
  const cursorLine = ref(1)
  const lineCount = ref(1)
  const wordCount = ref(0)
  const charCount = ref(0)
  const findBoxOpen = ref(false)
  const findQuery = ref('')
  const replaceQuery = ref('')

  function setContent(val: string) {
    content.value = val
    updateStats()
    // 自动保存到本地
    localStorage.setItem('m2h_webapp_content', val)
  }

  function updateStats() {
    const text = content.value
    charCount.value = text.length
    lineCount.value = text ? text.split('\n').length : 1
    wordCount.value = text
      ? text.trim().split(/\s+/).filter(w => w.length > 0).length
      : 0
  }

  function init() {
    const saved = localStorage.getItem('m2h_webapp_content')
    if (saved) content.value = saved

    const savedFilename = localStorage.getItem('m2h_webapp_filename')
    if (savedFilename) filename.value = savedFilename

    updateStats()
  }

  return {
    content, filename, cursorLine, lineCount, wordCount, charCount,
    findBoxOpen, findQuery, replaceQuery,
    setContent, setFilename: (f: string) => { filename.value = f; localStorage.setItem('m2h_webapp_filename', f) },
    updateStats, init,
  }
})
