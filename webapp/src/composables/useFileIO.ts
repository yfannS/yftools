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
    const file = (e.target as HTMLInputElement).files?.[0]
    if (file) readFile(file)
    // 重置，否则同名文件再次选择不触发 change 事件
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
      editorStore.setContent(content)
      editorStore.setFilename(file.name.replace(/\.[^/.]+$/, ''))
      render()
      toast('已导入：' + file.name)
    }
    reader.readAsText(file)
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
    const file = e.dataTransfer?.files?.[0]
    if (file) readFile(file)
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
