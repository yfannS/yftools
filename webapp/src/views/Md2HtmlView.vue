<template>
  <div class="md2html-view">
    <ToolbarBar />
    <TabBar />
    <main class="app-main" ref="mainRef">
      <!-- 编辑器面板 -->
      <section class="app-panel editor-panel">
        <div class="panel-bar">
          <div class="panel-bar-left">
            <span class="panel-label">编辑器</span>
            <span class="panel-stat">{{ editorStore.lineCount }} 行</span>
            <span class="panel-stat">{{ editorStore.wordCount }} 词</span>
            <span class="panel-stat">{{ editorStore.charCount }} 字符</span>
          </div>
          <div class="panel-bar-right">
            <button class="panel-action-btn" @click="loadExample" title="加载示例">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="9" y1="13" x2="15" y2="13"/><line x1="9" y1="17" x2="15" y2="17"/></svg>
              <span>示例</span>
            </button>
            <button class="panel-action-btn" :class="{ active: editorStore.findBoxOpen }" @click="editorStore.findBoxOpen = !editorStore.findBoxOpen" title="查找替换 (Ctrl+F)">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
              <span>查找</span>
            </button>
            <button class="panel-action-btn" @click="fileIO.triggerImport()" title="导入 .md 文件（支持多选）">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
              <span>导入</span>
            </button>
          </div>
        </div>
        <div class="editor-body">
          <div class="line-gutter" ref="lineGutterRef"><pre>{{ lineNumbers }}</pre></div>
          <div class="editor-area">
            <div class="current-line" :style="currentLineStyle"></div>
            <textarea
              ref="editorRef"
              id="editor"
              spellcheck="false"
              :value="editorStore.content"
              @input="onInput"
              @scroll="onEditorScroll"
              @click="updateCurrentLine"
              @keyup="updateCurrentLine"
              @keydown="onEditorKeydown"
              placeholder="# 在此输入 Markdown&#10;&#10;支持 Mermaid 流程图、KaTeX 数学公式&#10;快捷键：Ctrl+S 导出 ｜ Ctrl+F 查找 ｜ Ctrl+Shift+C 复制"
            ></textarea>
            <!-- 查找替换 -->
            <div class="find-box" v-if="editorStore.findBoxOpen">
              <div class="find-row">
                <input id="findInput" v-model="editorStore.findQuery" placeholder="查找内容" @keydown.enter="findReplace.findNext()">
                <button class="toolbar-btn" @click="findReplace.findNext()">下一个</button>
              </div>
              <div class="find-row">
                <input id="replaceInput" v-model="editorStore.replaceQuery" placeholder="替换为">
                <button class="toolbar-btn" @click="findReplace.replaceOne()">替换</button>
                <button class="toolbar-btn" @click="findReplace.replaceAll()">全部</button>
              </div>
              <div class="find-row" style="justify-content:flex-end">
                <button class="toolbar-btn" @click="editorStore.findBoxOpen = false">关闭</button>
              </div>
            </div>
          </div>
        </div>
      </section>

      <div
        class="resize-handle"
        :class="{ active: panelResize.isResizing.value }"
        @mousedown="panelResize.onResizeStart($event, mainRef!)"
        title="拖拽调整面板宽度"
      ></div>

      <!-- 预览面板 -->
      <section class="app-panel preview-panel">
        <div class="panel-bar">
          <div class="panel-bar-left">
            <span class="panel-label">预览</span>
          </div>
          <div class="panel-bar-right">
            <div class="panel-tabs">
              <button class="panel-tab" :class="{ active: previewStore.viewMode === 'preview' }" @click="previewStore.setViewMode('preview')">预览</button>
              <button class="panel-tab" :class="{ active: previewStore.viewMode === 'source' }" @click="previewStore.setViewMode('source')">源码</button>
            </div>
            <span class="panel-stat" v-if="previewStore.renderTimeMs">{{ previewStore.renderTimeMs }}ms</span>
            <span class="panel-stat">{{ previewStore.htmlSize }}</span>
            <span class="panel-sep"></span>
            <button class="panel-action-btn primary" @click="exportHtml.downloadHTML()" title="导出当前 HTML (Ctrl+S)">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
              <span>导出</span>
            </button>
            <button class="panel-action-btn" v-if="editorStore.tabs.length > 1" @click="exportHtml.downloadAllHTML()" title="批量导出所有 HTML">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="10" y1="17" x2="14" y2="17"/></svg>
              <span>全部导出</span>
            </button>
            <button class="panel-action-btn" @click="exportHtml.copyHTML()" title="复制 HTML (Ctrl+Shift+C)">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
              <span>复制</span>
            </button>
            <OutlinePanel />
          </div>
        </div>
        <div class="preview-body">
          <div class="preview-pane" :class="{ active: previewStore.viewMode === 'preview' }" ref="previewPaneRef" @scroll="scrollSync.syncPreviewToEditor">
            <div id="previewRender" ref="previewRenderRef">
              <div class="empty-state" v-if="!previewStore.renderedHtml">
                <svg viewBox="0 0 24 24"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><line x1="10" y1="9" x2="8" y2="9"/></svg>
                <h3>欢迎使用 Markdown 编辑器</h3>
                <p>
                  点击顶部 <strong>「示例」</strong>按钮快速体验功能，<br>
                  或直接在左侧输入 Markdown 开始编写。<br><br>
                  支持 <strong>Mermaid 流程图、类图、时序图</strong>、<strong>KaTeX 数学公式</strong>，<br>
                  键盘快捷键：<kbd>Ctrl+S</kbd> 导出 ｜ <kbd>Ctrl+F</kbd> 查找 ｜ <kbd>Ctrl+Shift+C</kbd> 复制
                </p>
              </div>
              <div v-else class="markdown-preview" v-html="previewStore.renderedHtml"></div>
            </div>
          </div>
          <pre class="preview-pane" :class="{ active: previewStore.viewMode === 'source' }" id="sourcePane">{{ previewStore.rawHtml }}</pre>
        </div>
      </section>
    </main>

    <!-- 拖拽导入 -->
    <div class="drop-zone" :class="{ active: fileIO.isDragging.value }">释放文件以导入 Markdown</div>

    <!-- Toast 容器 -->
    <div class="toast-container" id="toast-container"></div>

    <!-- 隐藏文件输入 -->
    <input type="file" ref="fileIO.fileInputRef.value" style="display:none" accept=".md,.markdown,.txt" multiple @change="fileIO.onFileSelected" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useEditorStore } from '@/stores/editor'
import { usePreviewStore } from '@/stores/preview'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { useMarkdownRender } from '@/composables/useMarkdownRender'
import { useScrollSync } from '@/composables/useScrollSync'
import { usePanelResize } from '@/composables/usePanelResize'
import { useFileIO } from '@/composables/useFileIO'
import { useFindReplace } from '@/composables/useFindReplace'
import { useToast } from '@/composables/useToast'
import { useExportHtml } from '@/composables/useExportHtml'
import { useAutoSave } from '@/composables/useAutoSave'
import ToolbarBar from '@/components/md2html/ToolbarBar.vue'
import TabBar from '@/components/md2html/TabBar.vue'
import OutlinePanel from '@/components/md2html/OutlinePanel.vue'

const editorStore = useEditorStore()
const previewStore = usePreviewStore()
const appStore = useAppStore()
const authStore = useAuthStore()

// DOM refs
const editorRef = ref<HTMLTextAreaElement>()
const lineGutterRef = ref<HTMLElement>()
const previewRenderRef = ref<HTMLElement>()
const previewPaneRef = ref<HTMLElement>()
const mainRef = ref<HTMLElement>()

// 渲染就绪标记（用于触发自动保存）
const renderReady = ref(false)

// Composables
const { toast } = useToast()
const markdownRender = useMarkdownRender(previewRenderRef, toast)
const scrollSync = useScrollSync(editorRef, previewPaneRef, lineGutterRef)
const panelResize = usePanelResize()
const fileIO = useFileIO(toast, () => { markdownRender.renderMarkdown(); renderReady.value = true })
const findReplace = useFindReplace(editorRef, toast, () => { markdownRender.renderMarkdown(); renderReady.value = true })
const exportHtml = useExportHtml(previewRenderRef, toast)
const autoSave = useAutoSave(toast, renderReady)

// 行号
const lineNumbers = computed(() => {
  const count = Math.max(1, editorStore.lineCount)
  return Array.from({ length: count }, (_, i) => i + 1).join('\n')
})

// 当前行高亮
const currentLineTop = ref(-9999)
const currentLineStyle = computed(() => ({
  height: '22.1px', // 13px * 1.7
  transform: `translateY(${currentLineTop.value}px)`
}))

function updateCurrentLine() {
  if (!editorRef.value) return
  const textarea = editorRef.value
  const before = textarea.value.slice(0, textarea.selectionStart)
  const line = before.split('\n').length
  const lh = 13 * 1.7
  currentLineTop.value = (line - 1) * lh - textarea.scrollTop + 12
}

// 编辑器输入
function onInput(e: Event) {
  const target = e.target as HTMLTextAreaElement
  editorStore.setContent(target.value)
  updateCurrentLine()
  markdownRender.renderMarkdown()
  renderReady.value = true
}

// 编辑器滚动
function onEditorScroll() {
  if (!editorRef.value || !lineGutterRef.value) return
  lineGutterRef.value.scrollTop = editorRef.value.scrollTop
  updateCurrentLine()
  scrollSync.syncEditorToPreview()
}

// 快捷键
function onEditorKeydown(e: KeyboardEvent) {
  if (e.ctrlKey || e.metaKey) {
    if (e.key.toLowerCase() === 's') {
      e.preventDefault()
      exportHtml.downloadHTML()
    } else if (e.key.toLowerCase() === 'f') {
      e.preventDefault()
      editorStore.findBoxOpen = !editorStore.findBoxOpen
    } else if (e.shiftKey && e.key.toLowerCase() === 'c') {
      e.preventDefault()
      exportHtml.copyHTML()
    }
  }
  // Tab 插入空格
  if (e.key === 'Tab') {
    e.preventDefault()
    const textarea = e.target as HTMLTextAreaElement
    const s = textarea.selectionStart
    editorStore.setContent(
      textarea.value.slice(0, s) + '  ' + textarea.value.slice(textarea.selectionEnd)
    )
    nextTick(() => {
      if (editorRef.value) {
        editorRef.value.selectionStart = editorRef.value.selectionEnd = s + 2
      }
    })
  }
}

// 加载示例
const EXAMPLE_MD = `# Markdown → HTML 编辑器

> 欢迎使用。这是一个完全在浏览器中运行的 Markdown 编辑器。

## 核心功能

- 实时预览，编辑器与预览双向滚动同步
- 支持 **Mermaid** 流程图与 **KaTeX** 数学公式
- 查找替换、代码块一键插入、行号与高亮
- 导出带目录的完整 HTML 文件

## 代码示例

\`\`\`python
for i in range(3):
    print(f"hello {i}")
\`\`\`

## Mermaid 流程图

\`\`\`mermaid
graph LR
    A[开始] --> B{判断}
    B -->|是| C[处理]
    B -->|否| D[结束]
\`\`\`

## 数学公式

行内公式：$E = mc^2$

块级公式：

$$
\\int_{0}^{1} x^2 \\, dx = \\frac{1}{3}
$$

## 表格

| 功能 | 状态 |
|------|------|
| 实时预览 | 支持 |
| 代码高亮 | 支持 |
| 数学公式 | 支持 |
| 流程图 | 支持 |
`

function loadExample() {
  if (editorStore.content.trim()) {
    if (!confirm('加载示例将覆盖当前内容，继续吗？')) return
  }
  editorStore.setContent(EXAMPLE_MD)
  markdownRender.renderMarkdown()
  toast('已加载示例文档')
}

// 全局拖拽事件
function onGlobalDragEnter(e: DragEvent) { e.preventDefault(); fileIO.onDragEnter(e) }
function onGlobalDragOver(e: DragEvent) { e.preventDefault(); fileIO.onDragOver(e) }
function onGlobalDragLeave(e: DragEvent) { e.preventDefault(); fileIO.onDragLeave(e) }
function onGlobalDrop(e: DragEvent) { e.preventDefault(); fileIO.onDrop(e) }

// 全局快捷键
function onGlobalKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 's') {
    e.preventDefault()
    exportHtml.downloadHTML()
  }
  if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.key.toLowerCase() === 'c') {
    e.preventDefault()
    exportHtml.copyHTML()
  }
}

// 切换 tab 时重新渲染
watch(
  () => editorStore.activeTabId,
  () => {
    nextTick(() => {
      // 恢复该 tab 的渲染结果到 previewStore
      const tab = editorStore.activeTab
      if (tab?.renderedHtml) {
        previewStore.setRenderedContent(tab.rawHtml, tab.renderedHtml, tab.renderTimeMs)
      } else if (tab?.content.trim()) {
        markdownRender.renderMarkdown()
      } else {
        previewStore.setRenderedContent('', '', 0)
      }
    })
  }
)

// 初始化
onMounted(() => {
  appStore.initTheme()
  markdownRender.initWorker()

  // 恢复分栏比例
  if (mainRef.value) {
    panelResize.initRatio(mainRef.value)
  }

  // 初始化 editor store（恢复多 tab 状态）
  editorStore.init()
  if (editorStore.content.trim()) {
    markdownRender.renderMarkdown()
  }

  // 全局事件
  document.addEventListener('dragenter', onGlobalDragEnter)
  document.addEventListener('dragover', onGlobalDragOver)
  document.addEventListener('dragleave', onGlobalDragLeave)
  document.addEventListener('drop', onGlobalDrop)
  document.addEventListener('keydown', onGlobalKeydown)
})

onUnmounted(() => {
  autoSave.dispose()
  markdownRender.terminateWorker()
  document.removeEventListener('dragenter', onGlobalDragEnter)
  document.removeEventListener('dragover', onGlobalDragOver)
  document.removeEventListener('dragleave', onGlobalDragLeave)
  document.removeEventListener('drop', onGlobalDrop)
  document.removeEventListener('keydown', onGlobalKeydown)
})
</script>

<style scoped>
.md2html-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  flex: 1;
  overflow: hidden;
}

/* Main */
.app-main {
  flex: 1;
  min-height: 0;
  display: grid;
  grid-template-columns: minmax(200px, 1fr) 4px minmax(200px, 1fr);
  gap: 0;
  background: var(--border);
}

.app-panel {
  background: var(--surface);
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.resize-handle {
  background: var(--border);
  cursor: col-resize;
  transition: background 0.2s var(--ease-out), box-shadow 0.2s;
  user-select: none;
  position: relative;
  z-index: 10;
}

.resize-handle::after {
  content: '';
  position: absolute;
  inset: 0 -3px;
}

.resize-handle:hover,
.resize-handle.active {
  background: var(--accent);
  box-shadow: 0 0 8px rgba(100, 116, 139, 0.25);
}

:global(body.resize-active) {
  cursor: col-resize;
  user-select: none;
}

:global(body.resize-active *) {
  cursor: col-resize !important;
}

.panel-bar {
  height: 34px;
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 12px;
  flex-shrink: 0;
  background: var(--surface);
  gap: 8px;
  box-shadow: inset 0 -1px 0 var(--border);
}

.panel-bar-left, .panel-bar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.panel-label {
  font-size: 11px;
  font-weight: 700;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.panel-stat {
  font-size: 11px;
  color: var(--text-tertiary);
  font-family: var(--font-mono);
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}

.panel-tabs {
  display: flex;
  gap: 2px;
}

.panel-tab {
  height: 22px;
  padding: 0 8px;
  border: none;
  background: transparent;
  color: var(--text-tertiary);
  border-radius: var(--radius-sm);
  font-size: 11px;
  font-family: var(--font-sans);
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s var(--ease-out);
}

.panel-tab:hover { color: var(--text-secondary); background: var(--surface-hover); }
.panel-tab.active {
  background: var(--accent-soft);
  color: var(--accent);
  font-weight: 600;
}

.panel-sep {
  width: 1px;
  height: 16px;
  background: var(--border);
}

.panel-action-btn {
  height: 24px;
  padding: 0 8px;
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--text-secondary);
  border-radius: var(--radius-sm);
  font-size: 11px;
  font-family: var(--font-sans);
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  transition: all 0.15s var(--ease-out);
  white-space: nowrap;
}

.panel-action-btn:hover {
  border-color: var(--border-strong);
  color: var(--text);
  background: var(--surface-hover);
  box-shadow: var(--shadow-sm);
}

.panel-action-btn:active {
  transform: scale(0.97);
  transition-duration: 0.08s;
}

.panel-action-btn.primary {
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-hover) 100%);
  color: #fff;
  border-color: var(--accent);
  box-shadow: 0 1px 3px rgba(100, 116, 139, 0.15);
}

.panel-action-btn.primary:hover {
  box-shadow: 0 2px 8px rgba(100, 116, 139, 0.2);
  filter: brightness(1.05);
  color: #fff;
}

.panel-action-btn.active {
  background: var(--accent-soft);
  border-color: var(--accent);
  color: var(--accent);
  box-shadow: 0 0 0 2px var(--accent-soft);
}

/* Editor */
.editor-body {
  flex: 1;
  min-height: 0;
  display: grid;
  grid-template-columns: 48px 1fr;
  position: relative;
}

.line-gutter {
  border-right: 1px solid var(--border);
  background: var(--editor-bg);
  overflow: hidden;
  padding: 12px 0;
  font-size: 13px;
  line-height: 1.7;
  font-family: var(--font-mono);
  color: var(--line-number);
  text-align: right;
  user-select: none;
}

.line-gutter pre {
  white-space: pre;
  padding: 0 10px 0 0;
}

.editor-area {
  position: relative;
  min-height: 0;
  overflow: hidden;
  background: var(--editor-bg);
}

.current-line {
  position: absolute;
  left: 0;
  right: 0;
  background: var(--line-highlight);
  pointer-events: none;
  z-index: 1;
  border-left: 2px solid var(--accent);
  transition: top 0.05s var(--ease-out);
}

#editor {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  border: none;
  outline: none;
  background: transparent;
  color: var(--text);
  resize: none;
  font-size: 13px;
  line-height: 1.7;
  padding: 12px 14px;
  font-family: var(--font-mono);
  z-index: 2;
  overflow: auto;
  tab-size: 2;
}

#editor::placeholder { color: var(--text-tertiary); opacity: 0.7; }
#editor::selection { background: rgba(100, 116, 139, 0.2); }

/* Find box */
.find-box {
  position: absolute;
  top: 10px;
  right: 10px;
  z-index: 30;
  background: var(--surface);
  border: 1px solid var(--border-strong);
  border-radius: var(--radius);
  box-shadow: var(--shadow-lg);
  padding: 10px;
  width: 280px;
  display: flex;
  gap: 8px;
  flex-direction: column;
  animation: findBoxIn 0.15s var(--ease-out);
}

@keyframes findBoxIn {
  from { opacity: 0; transform: translateY(-4px); }
  to { opacity: 1; transform: translateY(0); }
}

.find-row { display: flex; gap: 6px; }
.find-row input {
  flex: 1;
  min-width: 0;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 5px 8px;
  background: var(--editor-bg);
  color: var(--text);
  outline: none;
  font-size: 12px;
  font-family: var(--font-sans);
}
.find-row input:focus { border-color: var(--accent); }
.find-row .toolbar-btn { height: 26px; padding: 0 8px; font-size: 11px; }

/* Preview */
.preview-body {
  flex: 1;
  min-height: 0;
  position: relative;
  overflow: hidden;
}

.preview-pane {
  position: absolute;
  inset: 0;
  overflow: auto;
  display: none;
  background: var(--preview-bg);
}

.preview-pane.active { display: block; }

#previewRender {
  padding: 20px 24px;
  max-width: 800px;
  margin: 0 auto;
  font-size: 15px;
  line-height: 1.75;
  color: var(--text);
}

#sourcePane {
  padding: 16px 20px;
  font-family: var(--font-mono);
  font-size: 12px;
  line-height: 1.7;
  color: var(--text-secondary);
  white-space: pre-wrap;
  word-break: break-word;
  background: var(--preview-bg);
}

/* Markdown preview styles */
.markdown-preview :deep(h1),
.markdown-preview :deep(h2),
.markdown-preview :deep(h3),
.markdown-preview :deep(h4),
.markdown-preview :deep(h5),
.markdown-preview :deep(h6) {
  font-weight: 600;
  line-height: 1.35;
  margin-top: 1.6em;
  margin-bottom: 0.5em;
  color: var(--text);
  scroll-margin-top: 16px;
}

.markdown-preview :deep(h1) { font-size: 1.75em; margin-top: 0; padding-bottom: 0.3em; border-bottom: 1px solid var(--border); }
.markdown-preview :deep(h2) { font-size: 1.35em; }
.markdown-preview :deep(h3) { font-size: 1.15em; }
.markdown-preview :deep(h4) { font-size: 1em; }

.markdown-preview :deep(p),
.markdown-preview :deep(ul),
.markdown-preview :deep(ol),
.markdown-preview :deep(pre),
.markdown-preview :deep(table),
.markdown-preview :deep(blockquote) {
  margin-bottom: 1em;
}

.markdown-preview :deep(a) {
  color: var(--accent);
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: border-color 0.15s;
}

.markdown-preview :deep(a:hover) { border-bottom-color: var(--accent); }

.markdown-preview :deep(strong) { font-weight: 600; color: var(--text); }
.markdown-preview :deep(em) { font-style: italic; }

.markdown-preview :deep(ul), .markdown-preview :deep(ol) { padding-left: 1.5em; }
.markdown-preview :deep(li) { margin-bottom: 0.25em; }
.markdown-preview :deep(li::marker) { color: var(--text-tertiary); }

.markdown-preview :deep(blockquote) {
  border-left: 3px solid var(--border-strong);
  padding: 8px 16px;
  color: var(--text-secondary);
  background: var(--surface-hover);
  border-radius: 0 var(--radius-sm) var(--radius-sm) 0;
}

.markdown-preview :deep(blockquote p:last-child) { margin-bottom: 0; }

.markdown-preview :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin-bottom: 1em;
  font-size: 14px;
}

.markdown-preview :deep(th),
.markdown-preview :deep(td) {
  border: 1px solid var(--border);
  padding: 8px 12px;
  text-align: left;
}

.markdown-preview :deep(th) {
  background: var(--surface-raised);
  font-weight: 600;
}

.markdown-preview :deep(tr:nth-child(even)) {
  background: var(--surface-hover);
}

.markdown-preview :deep(hr) {
  border: none;
  height: 1px;
  background: var(--border-strong);
  margin: 2em 0;
}

.markdown-preview :deep(img) {
  max-width: 100%;
  border-radius: var(--radius);
  border: 1px solid var(--border);
}

.markdown-preview :deep(del) {
  color: var(--text-tertiary);
  text-decoration: line-through;
}

.markdown-preview :deep(input[type="checkbox"]) {
  margin-right: 6px;
  accent-color: var(--accent);
  vertical-align: middle;
  position: relative;
  top: -1px;
}

.markdown-preview :deep(li) {
  margin-bottom: 0.35em;
}

.markdown-preview :deep(.mermaid) {
  position: relative;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  overflow: hidden;
  margin: 1em 0;
}

.markdown-preview :deep(.mermaid-inner) {
  padding: 16px;
  display: flex;
  justify-content: center;
  overflow: auto;
}

.markdown-preview :deep(.mermaid-toolbar) {
  position: absolute;
  top: 8px;
  right: 8px;
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s;
  z-index: 3;
}

.markdown-preview :deep(.mermaid:hover .mermaid-toolbar) {
  opacity: 1;
}

.markdown-preview :deep(.mermaid-scale-select) {
  height: 24px;
  min-width: 56px;
  padding: 0 20px 0 8px;
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--text);
  border-radius: 4px;
  font-size: 11px;
  font-family: var(--font-sans);
  cursor: pointer;
  appearance: none;
}

.markdown-preview :deep(.mermaid-download-btn) {
  height: 24px;
  padding: 0 8px;
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--text);
  border-radius: 4px;
  font-size: 11px;
  font-family: var(--font-sans);
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  transition: all 0.12s;
}

.markdown-preview :deep(.mermaid-download-btn:hover) {
  background: var(--surface-hover);
  border-color: var(--accent);
  color: var(--accent);
}

.markdown-preview :deep(.katex-display) {
  margin: 1em 0;
  overflow-x: auto;
}

.markdown-preview :deep(.code-header) {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 36px;
  display: flex;
  align-items: center;
  padding: 0 12px;
  pointer-events: none;
  z-index: 2;
}

.markdown-preview :deep(.code-lang) {
  font-size: 10px;
  color: #525252;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-family: var(--font-mono);
  user-select: none;
  pointer-events: none;
}

.markdown-preview :deep(.copy-btn) {
  position: absolute;
  top: 6px;
  right: 8px;
  width: 28px;
  height: 28px;
  padding: 0;
  border: 1px solid rgba(255,255,255,0.08);
  background: rgba(255,255,255,0.04);
  color: #525252;
  border-radius: 6px;
  font-size: 10px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  z-index: 3;
  transition: opacity 0.15s;
}

.markdown-preview :deep(pre:hover .copy-btn) { opacity: 1; }
.markdown-preview :deep(.copy-btn:hover) { background: rgba(255,255,255,0.10); color: #e5e5e5; border-color: rgba(255,255,255,0.15); }
.markdown-preview :deep(.copy-btn.copied) { color: #22c55e; border-color: rgba(34,197,94,0.25); background: rgba(34,197,94,0.08); }

.markdown-preview :deep(code) {
  font-family: var(--font-mono);
  font-size: 0.88em;
  padding: 2px 5px;
  border-radius: 4px;
  background: var(--surface-raised);
  color: var(--accent);
}

.markdown-preview :deep(pre) {
  background: var(--code-bg);
  color: #e5e5e5;
  border-radius: var(--radius);
  overflow: auto;
  position: relative;
  padding-top: 36px;
  border: 1px solid var(--border);
}

.markdown-preview :deep(pre code) {
  display: block;
  background: transparent;
  color: inherit;
  padding: 0 16px 14px;
  font-size: 12.5px;
  line-height: 1.65;
  white-space: pre;
  border-radius: 0;
}

.markdown-preview :deep(pre code.hljs) {
  background: transparent !important;
}

/* Empty state */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--text-tertiary);
  text-align: center;
  padding: 40px;
  animation: emptyFadeIn 0.4s var(--ease-out);
}

@keyframes emptyFadeIn {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}

.empty-state svg {
  width: 48px;
  height: 48px;
  stroke: var(--border-strong);
  fill: none;
  stroke-width: 1.5;
  margin-bottom: 16px;
  opacity: 0.6;
}

.empty-state h3 {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 6px;
}

.empty-state p {
  font-size: 13px;
  line-height: 1.6;
  max-width: 360px;
}

.empty-state strong {
  color: var(--text);
  font-weight: 600;
}

.empty-state kbd {
  display: inline-block;
  font-family: var(--font-mono);
  font-size: 11px;
  padding: 1px 5px;
  border: 1px solid var(--border);
  border-radius: 3px;
  background: var(--surface-raised);
  color: var(--text-secondary);
  vertical-align: middle;
  box-shadow: 0 1px 0 var(--border-strong);
}

/* Drop zone */
.drop-zone {
  position: fixed;
  inset: 6px;
  border: 2px dashed var(--accent);
  border-radius: var(--radius);
  background: var(--accent-soft);
  display: none;
  align-items: center;
  justify-content: center;
  color: var(--accent);
  font-size: 15px;
  font-weight: 500;
  z-index: 150;
  pointer-events: none;
  letter-spacing: -0.01em;
  backdrop-filter: blur(4px);
}

.drop-zone.active { display: flex; animation: dropZoneIn 0.2s var(--ease-out); }

@keyframes dropZoneIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

/* Toast */
.toast-container {
  position: fixed;
  bottom: 16px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 200;
  display: flex;
  flex-direction: column;
  gap: 6px;
  pointer-events: none;
}

/* Responsive */
@media (max-width: 900px) {
  .app-main { grid-template-columns: 1fr; }
  .resize-handle { display: none; }
  #previewRender { padding: 16px; font-size: 14px; }
}
</style>

<style>
/* Toast 全局样式（不能 scoped） */
.toast {
  background: var(--text);
  color: var(--surface);
  border-radius: var(--radius-sm);
  padding: 8px 14px;
  font-size: 12px;
  font-weight: 500;
  box-shadow: var(--shadow-lg);
  opacity: 0;
  transform: translateY(6px);
  animation: toastIn 0.25s var(--ease-spring) forwards, toastOut 0.2s 1.8s var(--ease-out) forwards;
}

@keyframes toastIn { to { opacity: 1; transform: translateY(0); } }
@keyframes toastOut { to { opacity: 0; transform: translateY(-4px); } }
</style>
