import { type Ref } from 'vue'

// 语言别名映射（完整版，与原始 index.html 保持一致）
const LANG_ALIASES: Record<string, string> = {
  c: 'c',
  'c++': 'cpp',
  cpp: 'cpp',
  'c#': 'csharp',
  csharp: 'csharp',
  java: 'java',
  python: 'python',
  py: 'python',
  go: 'go',
  golang: 'go',
  js: 'javascript',
  javascript: 'javascript',
  ts: 'typescript',
  typescript: 'typescript',
  rb: 'ruby',
  ruby: 'ruby',
  sh: 'bash',
  bash: 'bash',
  shell: 'bash',
  zsh: 'bash',
  yml: 'yaml',
  yaml: 'yaml',
  json: 'json',
  xml: 'xml',
  html: 'xml',
  svg: 'xml',
  css: 'css',
  sql: 'sql',
  rust: 'rust',
  rs: 'rust',
  php: 'php',
  swift: 'swift',
  kotlin: 'kotlin',
  kt: 'kotlin',
  dart: 'dart',
  dockerfile: 'dockerfile',
  docker: 'dockerfile',
  makefile: 'makefile',
  mk: 'makefile',
  markdown: 'markdown',
  md: 'markdown',
}

const COPY_ICON_SVG = '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>'
const CHECK_ICON_SVG = '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>'

// Toast 回调类型
type ToastFn = (msg: string, type?: string) => void

export function useHighlight(previewRef: Ref<HTMLElement | undefined>, toast: ToastFn) {

  function highlightCode(code: HTMLElement, lang: string) {
    const mapped = LANG_ALIASES[lang] || lang
    const text = code.textContent?.trim() || ''
    if (!text) return

    // 保存原始文本，防止高亮失败后内容丢失
    const originalHTML = code.innerHTML
    try {
      const result = window.hljs?.highlight(text, { language: mapped, ignoreIllegals: true })
      if (result?.value && result.value.includes('<span')) {
        code.innerHTML = result.value
        return
      }
    } catch { /* ignore */ }

    // Fallback: 恢复原始内容后用 highlightElement
    code.innerHTML = originalHTML
    try {
      code.className = (code.className + ' language-' + mapped).trim()
      window.hljs?.highlightElement(code)
    } catch { /* ignore */ }
  }

  function addCodeCopyButtons() {
    const container = previewRef.value
    if (!container) return

    container.querySelectorAll('pre code').forEach((code) => {
      const pre = code.closest('pre')
      if (!pre || (pre as HTMLElement).dataset.ready) return
      ;(pre as HTMLElement).dataset.ready = '1'

      // 语言标签
      const langClass = Array.from(code.classList).find(c => c.startsWith('language-')) || ''
      let lang = (langClass.replace('language-', '') || 'text').toLowerCase()
      const displayLang = LANG_ALIASES[lang] || lang

      const header = document.createElement('div')
      header.className = 'code-header'

      const langLabel = document.createElement('span')
      langLabel.className = 'code-lang'
      langLabel.textContent = displayLang
      header.appendChild(langLabel)
      pre.prepend(header)

      // 复制按钮
      const copyBtn = document.createElement('button')
      copyBtn.className = 'copy-btn'
      copyBtn.title = '复制'
      copyBtn.innerHTML = COPY_ICON_SVG
      copyBtn.onclick = () => {
        navigator.clipboard.writeText(code.textContent || '').then(() => {
          copyBtn.classList.add('copied')
          copyBtn.innerHTML = CHECK_ICON_SVG
          toast('代码已复制')
          setTimeout(() => {
            copyBtn.classList.remove('copied')
            copyBtn.innerHTML = COPY_ICON_SVG
          }, 1500)
        })
      }
      pre.appendChild(copyBtn)

      // 代码高亮
      highlightCode(code as HTMLElement, lang)
    })
  }

  function attachHeadingIds() {
    const container = previewRef.value
    if (!container) return
    container.querySelectorAll('h1,h2,h3,h4,h5,h6').forEach((h, i) => {
      if (!h.id) h.id = 'h-' + i
    })
  }

  function buildOutline() {
    const container = previewRef.value
    if (!container) return []
    const items: { id: string; text: string; level: number }[] = []
    container.querySelectorAll('h1,h2,h3,h4,h5,h6').forEach((h) => {
      items.push({
        id: h.id || '',
        text: (h.textContent || '').trim(),
        level: Number(h.tagName.substring(1)),
      })
    })
    return items
  }

  return {
    addCodeCopyButtons,
    attachHeadingIds,
    buildOutline,
  }
}
