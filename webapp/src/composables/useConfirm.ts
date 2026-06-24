const CONFIRM_STYLE_ID = 'confirm-global-style'

interface ConfirmOptions {
  title?: string
  message: string
  /** 确认按钮文字，默认"确定" */
  confirmText?: string
  /** 取消按钮文字，默认"取消" */
  cancelText?: string
  /** 危险操作时确认按钮显示红色 */
  danger?: boolean
}

function ensureConfirmStyle() {
  if (document.getElementById(CONFIRM_STYLE_ID)) return

  const style = document.createElement('style')
  style.id = CONFIRM_STYLE_ID
  style.textContent = `
    .confirm-overlay {
      position: fixed;
      inset: 0;
      background: rgba(0, 0, 0, 0.25);
      display: flex;
      align-items: center;
      justify-content: center;
      z-index: 4000;
      padding: 20px;
      backdrop-filter: blur(6px);
      -webkit-backdrop-filter: blur(6px);
      animation: confirmOverlayIn 0.18s ease both;
    }

    @keyframes confirmOverlayIn {
      from { opacity: 0; }
      to   { opacity: 1; }
    }

    .confirm-dialog {
      width: min(400px, 100%);
      background: var(--surface, #fff);
      border: 1px solid var(--border, #e2e8f0);
      border-radius: 12px;
      box-shadow: 0 20px 60px rgba(0, 0, 0, 0.16), 0 4px 16px rgba(0, 0, 0, 0.08);
      padding: 24px;
      animation: confirmDialogIn 0.22s cubic-bezier(0.16, 1, 0.3, 1) both;
    }

    @keyframes confirmDialogIn {
      from { opacity: 0; transform: scale(0.93) translateY(8px); }
      to   { opacity: 1; transform: scale(1)    translateY(0);    }
    }

    .confirm-title {
      font-size: 16px;
      font-weight: 700;
      color: var(--text, #0f172a);
      margin-bottom: 8px;
    }

    .confirm-message {
      font-size: 13.5px;
      line-height: 1.6;
      color: var(--text-secondary, #475569);
      margin-bottom: 20px;
    }

    .confirm-actions {
      display: flex;
      justify-content: flex-end;
      gap: 8px;
    }

    .confirm-btn {
      height: 32px;
      padding: 0 16px;
      border-radius: 7px;
      font-size: 13px;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.15s;
      font-family: var(--font-sans, system-ui, sans-serif);
      outline: none;
      white-space: nowrap;
    }

    .confirm-btn:focus-visible {
      box-shadow: 0 0 0 3px var(--accent-soft, rgba(99,102,241,0.15));
    }

    .confirm-btn-cancel {
      background: transparent;
      border: 1px solid var(--border, #e2e8f0);
      color: var(--text-secondary, #475569);
    }

    .confirm-btn-cancel:hover {
      background: var(--surface-hover, #f8fafc);
      border-color: var(--border-strong, #94a3b8);
      color: var(--text, #0f172a);
    }

    .confirm-btn-ok {
      background: var(--accent, #6366f1);
      border: 1px solid var(--accent, #6366f1);
      color: #fff;
    }

    .confirm-btn-ok:hover {
      background: var(--accent-hover, #4f46e5);
      border-color: var(--accent-hover, #4f46e5);
    }

    .confirm-btn-ok.danger {
      background: #dc2626;
      border-color: #dc2626;
    }

    .confirm-btn-ok.danger:hover {
      background: #b91c1c;
      border-color: #b91c1c;
    }
  `
  document.head.appendChild(style)
}

/**
 * 全局命令式确认对话框，替代原生 confirm()。
 * 返回 Promise<boolean>，用户确认为 true，取消为 false。
 *
 * @example
 * const ok = await showConfirm({ message: '确定删除这条记录？', danger: true })
 * if (!ok) return
 */
export function showConfirm(opts: ConfirmOptions): Promise<boolean> {
  if (typeof document === 'undefined') return Promise.resolve(false)
  ensureConfirmStyle()

  return new Promise<boolean>((resolve) => {
    const overlay = document.createElement('div')
    overlay.className = 'confirm-overlay'

    const title = opts.title ? `<div class="confirm-title">${escHtml(opts.title)}</div>` : ''
    overlay.innerHTML = `
      <div class="confirm-dialog" role="dialog" aria-modal="true">
        ${title}
        <div class="confirm-message">${escHtml(opts.message)}</div>
        <div class="confirm-actions">
          <button class="confirm-btn confirm-btn-cancel">${escHtml(opts.cancelText ?? '取消')}</button>
          <button class="confirm-btn confirm-btn-ok ${opts.danger ? 'danger' : ''}">${escHtml(opts.confirmText ?? '确定')}</button>
        </div>
      </div>
    `

    function cleanup(result: boolean) {
      document.removeEventListener('keydown', onKeydown)
      overlay.remove()
      resolve(result)
    }

    function onKeydown(e: KeyboardEvent) {
      if (e.key === 'Escape') { e.preventDefault(); cleanup(false) }
      if (e.key === 'Enter')  { e.preventDefault(); cleanup(true)  }
    }

    overlay.querySelector('.confirm-btn-cancel')?.addEventListener('click', () => cleanup(false))
    overlay.querySelector('.confirm-btn-ok')?.addEventListener('click',     () => cleanup(true))
    overlay.addEventListener('click', (e) => { if (e.target === overlay) cleanup(false) })

    document.addEventListener('keydown', onKeydown)
    document.body.appendChild(overlay)

    // 自动聚焦确认按钮，方便键盘操作
    setTimeout(() => (overlay.querySelector('.confirm-btn-ok') as HTMLElement)?.focus(), 50)
  })
}

function escHtml(s: string): string {
  return s
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}
