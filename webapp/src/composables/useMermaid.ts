import { type Ref } from 'vue'

type ToastFn = (msg: string, type?: string) => void

// 所有 Mermaid 图表类型关键字（不区分大小写）
const MERMAID_KEYWORDS = [
  'graph', 'flowchart',
  'sequenceDiagram', 'sequencediagram',
  'classDiagram', 'classdiagram',
  'stateDiagram', 'statediagram',
  'erDiagram', 'erdiagram',
  'gantt',
  'pie',
  'journey',
  'gitgraph', 'gitGraph',
  'mindmap',
  'timeline',
  'quadrantChart', 'quadrantchart',
  'sankey',
  'block-beta', 'block',
  'xychart', 'xyChart',
  'requirementDiagram', 'requirementdiagram',
  'c4Context', 'c4context',
]

const KNOWN_MERMAID_LANGS = new Set([
  'mermaid', 'seq', 'sequence', 'sequencediagram',
  'flowchart', 'class', 'classdiagram', 'classdiagram',
  'state', 'statediagram', 'er', 'erdiagram',
  'graph', 'gantt', 'pie', 'journey', 'gitgraph', 'git',
  'mindmap', 'timeline', 'sankey', 'block',
])

export function useMermaid(previewRef: Ref<HTMLElement | undefined>, toast: ToastFn, renderToken: Ref<number>) {

  function looksLikeMermaid(code: HTMLElement): boolean {
    const lang = Array.from(code.classList)
      .find(c => c.startsWith('language-'))
      ?.replace('language-', '')
      ?.toLowerCase() || ''
    if (KNOWN_MERMAID_LANGS.has(lang)) return true
    const raw = (code.textContent || '').trim()
    return MERMAID_KEYWORDS.some(kw => {
      if (raw.toLowerCase().startsWith(kw.toLowerCase())) return true
      const kwFirstWord = kw.split(/\s/)[0]
      return raw.toLowerCase().startsWith(kwFirstWord.toLowerCase())
    })
  }

  function normalizeMermaidSource(code: HTMLElement): string {
    const raw = (code.textContent || '').trim()
    const isFlowchartFence = code.classList.contains('language-flowchart')
    const isSequenceFence =
      code.classList.contains('language-seq') ||
      code.classList.contains('language-sequence')

    // sequenceDiagram shorthand
    if (isSequenceFence) {
      if (/^sequenceDiagram/i.test(raw)) return raw
      return 'sequenceDiagram\n' + raw
    }

    // flowchart shorthand
    if (isFlowchartFence) {
      if (/^(flowchart|graph)\s+/i.test(raw)) return raw
      const lines = raw.split('\n')
      const head = (lines.shift() || 'TB').trim()
      return 'flowchart ' + (head || 'TB') + '\n' + lines.join('\n')
    }

    return raw
  }

  function processMermaid() {
    const container = previewRef.value
    if (!container) return

    // 查找所有代码块，判断是否为 mermaid
    container.querySelectorAll('pre code').forEach(code => {
      if (!looksLikeMermaid(code as HTMLElement)) return
      // 跳过已经处理过的
      if (code.closest('.mermaid')) return
      const pre = code.closest('pre')
      if (!pre || (pre as HTMLElement).dataset.mermaidProcessed) return
      ;(pre as HTMLElement).dataset.mermaidProcessed = '1'
      const wrap = document.createElement('div')
      wrap.className = 'mermaid'
      const source = normalizeMermaidSource(code as HTMLElement)
      wrap.textContent = source
      wrap.setAttribute('data-mermaid-source', source)
      pre.replaceWith(wrap)
    })

    const nodes = container.querySelectorAll('.mermaid')
    if (nodes.length) {
      const currentToken = renderToken.value
      const api = window.mermaid
      if (api) {
        api.run({ nodes: Array.from(nodes) })
          .then(() => {
            if (renderToken.value === currentToken) {
              addMermaidDownloadButtons()
            }
          })
          .catch(() => {})
      }
    }
  }

  function addMermaidDownloadButtons() {
    const container = previewRef.value
    if (!container) return

    container.querySelectorAll('.mermaid').forEach((containerEl, index) => {
      const svg = containerEl.querySelector('svg')
      if (!svg || (containerEl as HTMLElement).dataset.downloadReady) return
      ;(containerEl as HTMLElement).dataset.downloadReady = '1'

      // 创建 mermaid-inner 容器
      const inner = document.createElement('div')
      inner.className = 'mermaid-inner'
      while (containerEl.firstChild) inner.appendChild(containerEl.firstChild)
      containerEl.appendChild(inner)

      // 工具栏
      const toolbar = document.createElement('div')
      toolbar.className = 'mermaid-toolbar'

      // 缩放选择器
      const scaleSelect = document.createElement('select')
      scaleSelect.className = 'mermaid-scale-select'
      scaleSelect.title = 'PNG 清晰度倍率'
      ;[2, 4, 6].forEach(scale => {
        const opt = document.createElement('option')
        opt.value = String(scale)
        opt.textContent = scale + 'x'
        if (scale === 2) opt.selected = true
        scaleSelect.appendChild(opt)
      })

      // SVG 下载按钮
      const svgBtn = document.createElement('button')
      svgBtn.className = 'mermaid-download-btn'
      svgBtn.innerHTML = '<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg> SVG'
      svgBtn.onclick = (e: Event) => {
        e.stopPropagation()
        downloadMermaidSVG(svg, index)
      }

      // PNG 下载按钮
      const pngBtn = document.createElement('button')
      pngBtn.className = 'mermaid-download-btn'
      pngBtn.innerHTML = '<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg> PNG'
      pngBtn.onclick = (e: Event) => {
        e.stopPropagation()
        const scale = Math.max(1, Number(scaleSelect.value) || 1)
        downloadMermaidPNG(svg, index, scale)
      }

      toolbar.appendChild(scaleSelect)
      toolbar.appendChild(svgBtn)
      toolbar.appendChild(pngBtn)
      containerEl.appendChild(toolbar)
    })
  }

  function downloadMermaidSVG(svg: SVGSVGElement, index: number) {
    const bbox = svg.getBBox()
    const pad = 25
    const vx = bbox.x - pad
    const vy = bbox.y - pad
    const vw = bbox.width + pad * 2
    const vh = bbox.height + pad * 2

    const clone = svg.cloneNode(true) as SVGSVGElement
    clone.setAttribute('xmlns', 'http://www.w3.org/2000/svg')
    clone.setAttribute('viewBox', `${vx} ${vy} ${vw} ${vh}`)
    clone.setAttribute('width', String(vw))
    clone.setAttribute('height', String(vh))
    clone.removeAttribute('style')

    const svgData = new XMLSerializer().serializeToString(clone)
    const blob = new Blob([svgData], { type: 'image/svg+xml;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `diagram-${index + 1}.svg`
    document.body.appendChild(a)
    a.click()
    a.remove()
    URL.revokeObjectURL(url)
    toast('SVG 已下载')
  }

  function downloadMermaidPNG(svg: SVGSVGElement, index: number, scale: number) {
    const bbox = svg.getBBox()
    const pad = 25
    const vw = Math.ceil(bbox.width + pad * 2)
    const vh = Math.ceil(bbox.height + pad * 2)
    const ratio = Math.max(1, scale || 1)

    const clone = svg.cloneNode(true) as SVGSVGElement
    clone.setAttribute('xmlns', 'http://www.w3.org/2000/svg')
    clone.setAttribute('xmlns:xlink', 'http://www.w3.org/1999/xlink')
    clone.setAttribute('viewBox', `${bbox.x - pad} ${bbox.y - pad} ${vw} ${vh}`)
    clone.setAttribute('width', String(vw))
    clone.setAttribute('height', String(vh))
    clone.removeAttribute('style')

    const svgData = new XMLSerializer().serializeToString(clone)
    const svgBase64 = 'data:image/svg+xml;base64,' + btoa(unescape(encodeURIComponent(svgData)))

    const canvas = document.createElement('canvas')
    canvas.width = vw * ratio
    canvas.height = vh * ratio
    const ctx = canvas.getContext('2d')!

    const img = new Image()
    img.onload = function () {
      ctx.fillStyle = '#ffffff'
      ctx.fillRect(0, 0, canvas.width, canvas.height)
      ctx.drawImage(img, 0, 0, canvas.width, canvas.height)

      canvas.toBlob(function (blob) {
        if (!blob) { toast('PNG 导出失败', 'err'); return }
        const url = URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = `diagram-${index + 1}.png`
        document.body.appendChild(a)
        a.click()
        a.remove()
        URL.revokeObjectURL(url)
        toast(`PNG 已下载（${ratio}x）`)
      }, 'image/png')
    }
    img.onerror = function () {
      toast('PNG 导出失败，请使用 SVG 格式', 'err')
    }
    img.src = svgBase64
  }

  return {
    processMermaid,
    addMermaidDownloadButtons,
    downloadMermaidSVG,
    downloadMermaidPNG
  }
}
