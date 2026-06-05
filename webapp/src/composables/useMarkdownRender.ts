import { watch, nextTick, type Ref } from 'vue'
import { useMarkdownWorker } from './useMarkdownWorker'
import { useMermaid } from './useMermaid'
import { useKaTeX } from './useKaTeX'
import { useHighlight } from './useHighlight'
import { useEditorStore } from '@/stores/editor'
import { usePreviewStore } from '@/stores/preview'

type ToastFn = (msg: string, type?: string) => void

export function useMarkdownRender(previewRef: Ref<HTMLElement | undefined>, toast: ToastFn) {
  const editorStore = useEditorStore()
  const previewStore = usePreviewStore()

  const { worker, renderToken, isRendering, initWorker, renderMarkdown, terminateWorker } = useMarkdownWorker()
  const { processMermaid } = useMermaid(previewRef, toast, renderToken)
  const { processMath } = useKaTeX(previewRef)
  const { addCodeCopyButtons, attachHeadingIds, buildOutline } = useHighlight(previewRef, toast)

  // 当 previewStore 的 renderedHtml 变化时，执行后处理
  watch(
    () => previewStore.renderedHtml,
    async () => {
      await nextTick()
      postProcess()
    }
  )

  function postProcess() {
    const container = previewRef.value
    if (!container || !previewStore.renderedHtml) return

    attachHeadingIds()
    processMermaid()
    processMath()
    addCodeCopyButtons()
    previewStore.setOutline(buildOutline())
  }

  return {
    worker,
    renderToken,
    isRendering,
    initWorker,
    renderMarkdown,
    terminateWorker,
    postProcess,
  }
}
