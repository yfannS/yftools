import { ref } from 'vue'
import { useEditorStore } from '@/stores/editor'
import { SUPPORTED_FILE_TYPES } from '@/utils/constants'

type ToastFn = (msg: string, type?: string) => void
type RenderFn = () => void

export function useFileIO(toast: ToastFn, render: RenderFn) {
  const fileInputRef = ref<HTMLInputElement>()
  const dragCounter = ref(0)
  const isDragging = ref(false)

  const editorStore = useEditorStore()

  function triggerImport() {
    fileInputRef.value?.click()
  }

  function onFileSelected(e: Event) {
    const files = (e.target as HTMLInputElement).files
    if (!files || files.length === 0) return
    if (files.length === 1) {
      readFile(files[0])
    } else {
      readMultipleFiles(Array.from(files))
    }
    if (fileInputRef.value) fileInputRef.value.value = ''
  }

  function onFileDropped(file: File) {
    readFile(file)
  }

  function readFile(file: File) {
    const ext = '.' + file.name.split('.').pop()?.toLowerCase()
    if (!SUPPORTED_FILE_TYPES.includes(ext)) {
      toast('只支持 .md / .txt 文件', 'err')
      return
    }
    const reader = new FileReader()
    reader.onload = (e) => {
      const content = e.target?.result as string
      const name = file.name.replace(/\.[^/.]+$/, '')
      // 如果当前 tab 为空，直接写入；否则创建新 tab
      if (!editorStore.content.trim()) {
        editorStore.setContent(content)
        editorStore.setFilename(name)
      } else {
        editorStore.createTab(name, content)
      }
      render()
      toast('已导入：' + file.name)
    }
    reader.readAsText(file)
  }

  function readMultipleFiles(files: File[]) {
    const validFiles = files.filter(f => {
      const ext = '.' + f.name.split('.').pop()?.toLowerCase()
      return SUPPORTED_FILE_TYPES.includes(ext)
    })
    if (validFiles.length === 0) {
      toast('没有支持的文件类型（.md / .txt）', 'err')
      return
    }

    let loaded = 0
    validFiles.forEach((file, idx) => {
      const reader = new FileReader()
      reader.onload = (e) => {
        const content = e.target?.result as string
        const name = file.name.replace(/\.[^/.]+$/, '')
        if (idx === 0 && !editorStore.content.trim()) {
          editorStore.setContent(content)
          editorStore.setFilename(name)
        } else {
          editorStore.createTab(name, content)
        }
        loaded++
        if (loaded === validFiles.length) {
          render()
          toast(`已导入 ${validFiles.length} 个文件`)
        }
      }
      reader.readAsText(file)
    })

    if (validFiles.length < files.length) {
      toast(`${files.length - validFiles.length} 个文件格式不支持，已跳过`, 'err')
    }
  }

  // 拖拽事件处理
  function onDragEnter(e: DragEvent) {
    e.preventDefault()
    dragCounter.value++
    isDragging.value = true
  }

  function onDragOver(e: DragEvent) {
    e.preventDefault()
  }

  function onDragLeave(e: DragEvent) {
    e.preventDefault()
    dragCounter.value--
    if (dragCounter.value <= 0) {
      dragCounter.value = 0
      isDragging.value = false
    }
  }

  function onDrop(e: DragEvent) {
    e.preventDefault()
    dragCounter.value = 0
    isDragging.value = false
    const files = e.dataTransfer?.files
    if (!files || files.length === 0) return
    if (files.length === 1) {
      readFile(files[0])
    } else {
      readMultipleFiles(Array.from(files))
    }
  }

  return {
    fileInputRef,
    isDragging,
    triggerImport,
    onFileSelected,
    onFileDropped,
    onDragEnter,
    onDragOver,
    onDragLeave,
    onDrop
  }
}
