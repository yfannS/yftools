/**
 * Mermaid 初始化配置
 * 由 index.html 底部 <script src="..."> 加载
 */
if (window.mermaid) {
  window.mermaid.initialize({
    startOnLoad: false,
    securityLevel: 'strict',
    theme: 'default',
  });
}
