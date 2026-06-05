// 语言别名映射（与 index.html 中 LANG_ALIASES 保持一致）
export const LANG_ALIASES: Record<string, string> = {
  js: 'javascript',
  ts: 'typescript',
  py: 'python',
  rb: 'ruby',
  sh: 'bash',
  yml: 'yaml',
  md: 'markdown',
  'c++': 'cpp',
  'c#': 'csharp',
  'f#': 'fsharp',
  objc: 'objectivec',
  tsx: 'typescript',
  jsx: 'javascript',
}

// localStorage 键名前缀
export const STORAGE_PREFIX = 'm2h_webapp_'

// 渲染配置
export const RENDER_DEBOUNCE_MS = 120
export const CHUNK_RENDER_THRESHOLD = 14000
export const CHUNK_SIZE = 36

// 快捷键
export const SHORTCUTS = {
  SAVE: { key: 's', ctrl: true },
  FIND: { key: 'f', ctrl: true },
  COPY_HTML: { key: 'C', ctrl: true, shift: true },
} as const

// 支持的文件类型
export const SUPPORTED_FILE_TYPES = ['.md', '.markdown', '.txt']
