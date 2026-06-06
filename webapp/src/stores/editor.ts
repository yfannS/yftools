import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface TabItem {
  id: string
  filename: string
  content: string
  rawHtml: string
  renderedHtml: string
  renderTimeMs: number
  htmlSize: string
}

let tabIdCounter = 0

function generateTabId(): string {
  return 'tab_' + (++tabIdCounter) + '_' + Date.now()
}

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

export const useEditorStore = defineStore('editor', () => {
  // ===== 多 Tab 管理 =====
  const tabs = ref<TabItem[]>([])
  const activeTabId = ref<string>('')

  // 当前活跃 tab
  const activeTab = computed<TabItem | undefined>(() =>
    tabs.value.find(t => t.id === activeTabId.value)
  )

  // 当前活跃 tab 的属性（兼容旧代码）
  const content = computed(() => activeTab.value?.content || '')
  const filename = computed(() => activeTab.value?.filename || 'document')
  const lineCount = computed(() => {
    const text = content.value
    return text ? text.split('\n').length : 1
  })
  const wordCount = computed(() => {
    const text = content.value
    return text.trim() ? text.trim().split(/\s+/).filter(w => w.length > 0).length : 0
  })
  const charCount = computed(() => content.value.length)

  // 查找替换状态（全局共享）
  const findBoxOpen = ref(false)
  const findQuery = ref('')
  const replaceQuery = ref('')

  // 脏标记：用户是否修改过内容（用于自动保存判断）
  const dirty = ref(false)

  // ===== Tab 操作 =====

  function createTab(name?: string, markdown?: string): TabItem {
    const tab: TabItem = {
      id: generateTabId(),
      filename: name || 'document',
      content: markdown || '',
      rawHtml: '',
      renderedHtml: '',
      renderTimeMs: 0,
      htmlSize: '0 B',
    }
    tabs.value.push(tab)
    activeTabId.value = tab.id
    saveTabsToStorage()
    return tab
  }

  function switchTab(tabId: string) {
    if (activeTabId.value === tabId) return
    activeTabId.value = tabId
    saveTabsToStorage()
  }

  function closeTab(tabId: string) {
    const idx = tabs.value.findIndex(t => t.id === tabId)
    if (idx === -1) return

    tabs.value.splice(idx, 1)

    if (activeTabId.value === tabId) {
      if (tabs.value.length === 0) {
        createTab()
      } else {
        const nextIdx = Math.min(idx, tabs.value.length - 1)
        activeTabId.value = tabs.value[nextIdx].id
      }
    }
    saveTabsToStorage()
  }

  function closeOtherTabs(tabId: string) {
    tabs.value = tabs.value.filter(t => t.id === tabId)
    activeTabId.value = tabId
    saveTabsToStorage()
  }

  // ===== 内容更新 =====

  function setContent(val: string) {
    const tab = activeTab.value
    if (tab) {
      if (tab.content !== val) {
        tab.content = val
        dirty.value = true
      }
      saveTabsToStorage()
    }
  }

  function setFilename(f: string) {
    const tab = activeTab.value
    if (tab) {
      tab.filename = f
      saveTabsToStorage()
    }
  }

  function setTabRendered(tabId: string, rawHtml: string, renderedHtml: string, timeMs: number) {
    const tab = tabs.value.find(t => t.id === tabId)
    if (tab) {
      tab.rawHtml = rawHtml
      tab.renderedHtml = renderedHtml
      tab.renderTimeMs = timeMs
      tab.htmlSize = formatBytes(new Blob([renderedHtml]).size)
    }
  }

  // ===== 持久化 =====

  function saveTabsToStorage() {
    const data = tabs.value.map(t => ({
      id: t.id,
      filename: t.filename,
      content: t.content,
    }))
    localStorage.setItem('m2h_webapp_tabs', JSON.stringify(data))
    localStorage.setItem('m2h_webapp_active_tab', activeTabId.value)
  }

  function init() {
    const savedTabs = localStorage.getItem('m2h_webapp_tabs')
    const savedActiveId = localStorage.getItem('m2h_webapp_active_tab')

    if (savedTabs) {
      try {
        const parsed = JSON.parse(savedTabs) as Array<{ id: string; filename: string; content: string }>
        if (parsed.length > 0) {
          tabs.value = parsed.map(t => ({
            id: t.id,
            filename: t.filename,
            content: t.content,
            rawHtml: '',
            renderedHtml: '',
            renderTimeMs: 0,
            htmlSize: '0 B',
          }))
          parsed.forEach(t => {
            const num = parseInt(t.id.split('_')[1])
            if (!isNaN(num) && num > tabIdCounter) tabIdCounter = num
          })
          activeTabId.value = savedActiveId || tabs.value[0].id
          if (!tabs.value.find(t => t.id === activeTabId.value)) {
            activeTabId.value = tabs.value[0].id
          }
          return
        }
      } catch {}
    }

    // 兼容旧版单文件存储
    const savedContent = localStorage.getItem('m2h_webapp_content')
    const savedFilename = localStorage.getItem('m2h_webapp_filename')
    if (savedContent) {
      createTab(savedFilename || 'document', savedContent)
      localStorage.removeItem('m2h_webapp_content')
      localStorage.removeItem('m2h_webapp_filename')
    } else {
      createTab()
    }
  }

  function markSaved() {
    dirty.value = false
  }

  return {
    tabs, activeTabId, activeTab,
    content, filename, lineCount, wordCount, charCount,
    findBoxOpen, findQuery, replaceQuery,
    dirty,
    createTab, switchTab, closeTab, closeOtherTabs,
    setContent, setFilename, setTabRendered,
    markSaved,
    init,
  }
})
