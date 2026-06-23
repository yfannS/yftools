/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

/** highlight.js 实际使用的最小类型声明（CDN 引入） */
interface HljsHighlightResult {
  value: string
}

interface Hljs {
  highlight(code: string, options: { language: string; ignoreIllegals: boolean }): HljsHighlightResult
  highlightElement(target: HTMLElement): void
}

/** KaTeX auto-render 实际使用的最小类型声明（CDN 引入） */
interface KatexRenderOptions {
  delimiters: Array<{ left: string; right: string; display: boolean }>
  throwOnError: boolean
}

/** Mermaid 实际使用的最小类型声明（CDN 引入） */
interface MermaidApi {
  run(options: { nodes: Element[] }): Promise<unknown>
}

interface Window {
  hljs?: Hljs
  renderMathInElement?: (element: HTMLElement, options: KatexRenderOptions) => void
  mermaid?: MermaidApi
}
