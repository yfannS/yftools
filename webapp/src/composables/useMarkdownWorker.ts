import { ref, onUnmounted, shallowRef } from 'vue'
import { useEditorStore } from '@/stores/editor'
import { usePreviewStore } from '@/stores/preview'
import { RENDER_DEBOUNCE_MS } from '@/utils/constants'

interface WorkerMessage {
  id: number
  markdown: string
}

interface WorkerResponse {
  id: number
  ok: boolean
  html: string
  error?: string
  workerMs: number
}

export function useMarkdownWorker() {
  const editorStore = useEditorStore()
  const previewStore = usePreviewStore()

  const worker = shallowRef<Worker | null>(null)
  const renderToken = ref(0)
  const isRendering = ref(false)
  let debounceTimer: ReturnType<typeof setTimeout> | null = null

  function initWorker() {
    const code = String.raw`
      self.importScripts('https://cdn.jsdelivr.net/npm/marked@12/marked.min.js');
      marked.use({ gfm: true, breaks: true });

      function escapeHtml(s) {
        return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
      }

      function protectMath(text) {
        var blocks = [];
        text = text.replace(/\$\$([\s\S]*?)\$\$/g, function(m) {
          blocks.push(m);
          return 'MATHPLACEHOLDER' + (blocks.length - 1) + 'XEND';
        });
        text = text.replace(/\$([^\$\n]+?)\$/g, function(m) {
          blocks.push(m);
          return 'MATHPLACEHOLDER' + (blocks.length - 1) + 'XEND';
        });
        text = text.replace(/\\\[([\s\S]*?)\\\]/g, function(m) {
          blocks.push(m);
          return 'MATHPLACEHOLDER' + (blocks.length - 1) + 'XEND';
        });
        text = text.replace(/\\\(([\s\S]*?)\\\)/g, function(m) {
          blocks.push(m);
          return 'MATHPLACEHOLDER' + (blocks.length - 1) + 'XEND';
        });
        return { text: text, blocks: blocks };
      }

      function restoreMath(html, blocks) {
        for (var i = 0; i < blocks.length; i++) {
          html = html.replace('MATHPLACEHOLDER' + i + 'XEND', escapeHtml(blocks[i]));
        }
        return html;
      }

      self.onmessage = function(ev) {
        var data = ev.data;
        var t0 = performance.now();
        try {
          var result = protectMath(data.markdown || '');
          var html = marked.parse(result.text);
          html = restoreMath(html, result.blocks);
          self.postMessage({ id: data.id, ok: true, html: html, workerMs: performance.now() - t0 });
        } catch(err) {
          self.postMessage({ id: data.id, ok: false, error: String(err), workerMs: performance.now() - t0 });
        }
      };
    `
    const blob = new Blob([code], { type: 'application/javascript' })
    worker.value = new Worker(URL.createObjectURL(blob))

    worker.value.onmessage = (ev: MessageEvent<WorkerResponse>) => {
      const d = ev.data
      if (d.id !== renderToken.value) return
      isRendering.value = false

      if (!d.ok) {
        previewStore.setRenderedContent(
          `<pre style="color:#dc2626;padding:20px">${escapeHtml(d.error || 'Unknown error')}</pre>`,
          '',
          0
        )
        return
      }

      const startTime = performance.now()

      const afterPaint = () => {
        const renderMs = Math.round(performance.now() - startTime + d.workerMs)
        previewStore.setRenderedContent(d.html, d.html, renderMs)
        // 保存渲染结果到当前 tab
        const tabId = editorStore.activeTabId
        if (tabId) {
          editorStore.setTabRendered(tabId, d.html, d.html, renderMs)
        }
      }

      afterPaint()
    }
  }

  function renderMarkdown() {
    if (!worker.value) initWorker()

    clearTimeout(debounceTimer!)
    debounceTimer = setTimeout(() => {
      const md = editorStore.content
      if (!md.trim()) {
        previewStore.setRenderedContent('', '', 0)
        return
      }
      renderToken.value++
      isRendering.value = true
      worker.value!.postMessage({
        id: renderToken.value,
        markdown: md
      } as WorkerMessage)
    }, RENDER_DEBOUNCE_MS)
  }

  function terminateWorker() {
    if (worker.value) {
      worker.value.terminate()
      worker.value = null
    }
    if (debounceTimer) clearTimeout(debounceTimer)
  }

  onUnmounted(() => {
    terminateWorker()
  })

  return {
    worker,
    renderToken,
    isRendering,
    initWorker,
    renderMarkdown,
    terminateWorker
  }
}

function escapeHtml(s: string): string {
  return s
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}
