import { useEditorStore } from '@/stores/editor'
import { usePreviewStore } from '@/stores/preview'
import { useAppStore } from '@/stores/app'
import { nextTick, type Ref } from 'vue'

type ToastFn = (msg: string, type?: string) => void

export function useExportHtml(previewRef: Ref<HTMLElement | undefined>, toast: ToastFn) {
  const editorStore = useEditorStore()
  const previewStore = usePreviewStore()
  const appStore = useAppStore()

  function escapeHtml(s: string): string {
    return s
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
  }

  function buildSidebar(): string {
    const container = previewRef.value
    if (!container) return ''

    const hs = container.querySelectorAll('h1,h2,h3,h4,h5,h6')
    if (!hs.length) return ''

    let list = ''
    hs.forEach((h) => {
      const id = h.id || ''
      const lv = Number(h.tagName.substring(1))
      const text = (h.textContent || '').trim()
      list += `<li class="sb-item lv-${lv}"><a href="#${id}" data-target="${id}">${escapeHtml(text)}</a></li>`
    })

    return `<aside class="sidebar" id="sidebar"><div class="sb-header"><span class="sb-title">目录</span><button class="sb-toggle" id="sbToggle" type="button" title="折叠/展开" aria-label="折叠/展开">«</button></div><nav class="sb-nav"><ul class="sb-list">${list}</ul></nav></aside><button class="sb-fab" id="sbFab" type="button" title="显示目录" aria-label="显示目录">☰</button>`
  }

  function buildFullHTML(): string {
    const container = previewRef.value
    if (!container) return ''

    const bodyRaw = container.innerHTML
    const cleaner = document.createElement('div')
    cleaner.innerHTML = bodyRaw
    // 移除工具栏
    cleaner.querySelectorAll('.mermaid-toolbar').forEach(el => el.remove())
    // 将渲染后的 SVG 替换为 Mermaid 源码，以便导出后客户端重新渲染
    cleaner.querySelectorAll('.mermaid').forEach(el => {
      const source = el.getAttribute('data-mermaid-source')
      if (source) {
        el.textContent = source
        el.removeAttribute('data-processed')
        el.removeAttribute('data-mermaid-source')
        el.removeAttribute('data-download-ready')
      }
    })

    const body = cleaner.innerHTML
    const sidebar = buildSidebar()
    const dark = appStore.theme === 'dark'
    const scTag = '</' + 'script>' // 安全的 </script> 闭合标签
    const bg = dark ? '#0a0a0a' : '#fafafa'
    const text = dark ? '#e5e5e5' : '#171717'
    const border = dark ? '#262626' : '#e5e5e5'
    const surface = dark ? '#111111' : '#ffffff'
    const mermaidTheme = dark ? 'dark' : 'default'

    const lbStyle = [
      '/* lightbox */',
      '.lb-overlay{position:fixed;inset:0;background:rgba(0,0,0,0.88);z-index:9999;display:none}',
      '.lb-overlay.active{display:block}',
      '.lb-stage{position:fixed;inset:0;display:flex;align-items:center;justify-content:center;padding:24px;box-sizing:border-box;overflow:hidden;cursor:zoom-out}',
      '.lb-media{max-width:90vw;max-height:90vh;transform-origin:center center;transition:transform .08s ease-out;cursor:grab;user-select:none;box-shadow:0 12px 48px rgba(0,0,0,0.35)}',
      '.lb-img{border-radius:8px;object-fit:contain}',
      '.lb-svg{background:#fff;border-radius:8px;padding:12px;box-sizing:border-box}',
      '.lb-close{position:fixed;top:16px;right:16px;width:36px;height:36px;border:none;background:rgba(255,255,255,0.12);color:#fff;border-radius:50%;cursor:pointer;display:flex;align-items:center;justify-content:center;font-size:18px;line-height:1;z-index:10000}',
      '.lb-close:hover{background:rgba(255,255,255,0.25)}',
      '.lb-hint{position:fixed;bottom:16px;left:50%;transform:translateX(-50%);color:rgba(255,255,255,0.7);font-size:12px;pointer-events:none;user-select:none;z-index:10000}',
      'body.lb-open{overflow:hidden}',
      '.mm-toolbar{position:absolute;top:8px;right:8px;display:flex;gap:4px;opacity:0;transition:opacity .2s;z-index:3}',
      '.mermaid:hover .mm-toolbar{opacity:1}',
      `.mm-btn{height:24px;padding:0 8px;border:1px solid ${border};background:${surface};color:${text};border-radius:4px;font-size:11px;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;cursor:pointer;display:inline-flex;align-items:center;gap:4px}`,
      `.mm-btn:hover{background:${dark ? '#1a1a1a' : '#f5f5f5'};border-color:#2563eb;color:#2563eb}`,
      `.mm-scale{height:24px;min-width:56px;padding:0 20px 0 8px;border:1px solid ${border};background:${surface};color:${text};border-radius:4px;font-size:11px;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;cursor:pointer;appearance:none}`,
      '/* sidebar */',
      `.sidebar{position:fixed;top:0;left:0;bottom:0;width:240px;background:${surface};border-right:1px solid ${border};display:flex;flex-direction:column;z-index:50;transition:transform .25s ease;box-shadow:2px 0 8px rgba(0,0,0,.04)}`,
      `body:not(.sb-open) .sidebar{transform:translateX(-100%)}`,
      '.sb-header{height:44px;display:flex;align-items:center;justify-content:space-between;padding:0 12px;border-bottom:1px solid ${border};flex-shrink:0}',
      `.sb-title{font-size:12px;font-weight:600;color:${text};letter-spacing:.04em;text-transform:uppercase}`,
      `.sb-toggle{width:24px;height:24px;border:none;background:transparent;color:${text === '#e5e5e5' ? '#a3a3a3' : '#525252'};cursor:pointer;font-size:14px;border-radius:4px}`,
      `.sb-toggle:hover{background:${dark ? '#1a1a1a' : '#f5f5f5'};color:${text}}`,
      '.sb-nav{flex:1;overflow-y:auto;padding:8px 0}',
      '.sb-list{list-style:none;margin:0;padding:0}',
      '.sb-item{margin:0;line-height:1}',
      `.sb-item a{display:block;padding:7px 14px 7px 14px;color:${text === '#e5e5e5' ? '#a3a3a3' : '#525252'};text-decoration:none;font-size:13px;border-left:2px solid transparent;transition:background .12s,color .12s,border-color .12s;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}`,
      `.sb-item a:hover{color:${text};background:${dark ? '#171717' : '#fafafa'}}`,
      `.sb-item.lv-1 a{font-weight:600}`,
      '.sb-item.lv-2 a{padding-left:28px}',
      '.sb-item.lv-3 a{padding-left:42px;font-size:12.5px}',
      '.sb-item.lv-4 a{padding-left:56px;font-size:12.5px}',
      '.sb-item.lv-5 a{padding-left:70px;font-size:12px}',
      '.sb-item.lv-6 a{padding-left:84px;font-size:12px}',
      `.sb-item.active a{color:#2563eb;background:${dark ? 'rgba(37,99,235,.12)' : 'rgba(37,99,235,.08)'};border-left-color:#2563eb}`,
      `.sb-fab{position:fixed;top:12px;left:12px;width:34px;height:34px;border:1px solid ${border};background:${surface};color:${text};border-radius:6px;cursor:pointer;display:none;align-items:center;justify-content:center;font-size:16px;z-index:49;box-shadow:0 2px 6px rgba(0,0,0,.06)}`,
      `body:not(.sb-open) .sb-fab{display:flex}`,
      `.sb-fab:hover{border-color:#2563eb;color:#2563eb}`,
      `body.sb-open .main-wrap{margin-left:240px}`,
      '@media (max-width:900px){.sidebar{width:80vw;max-width:280px}body.sb-open .main-wrap{margin-left:0}}',
    ].join('')

    const lbScript = [
      '<script src="https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.min.js"><' + '/script>',
      '<script>',
      '(function(){',
      'var overlay=document.createElement("div");',
      'overlay.className="lb-overlay";',
      'overlay.innerHTML="<div class=\\"lb-stage\\"></div><button class=\\"lb-close\\" aria-label=\\"关闭\\">&times;</button><div class=\\"lb-hint\\">滚轮缩放 · 拖拽移动 · 点击空白/ESC 关闭</div>";',
      'document.body.appendChild(overlay);',
      'var stage=overlay.querySelector(".lb-stage");',
      'var closeBtn=overlay.querySelector(".lb-close");',
      'var activeEl=null,scale=1,tx=0,ty=0,isDragging=false,sx=0,sy=0;',
      'function apply(){if(!activeEl)return;activeEl.style.transform="translate("+tx+"px,"+ty+"px) scale("+scale+")";}',
      'function mount(el){stage.innerHTML="";activeEl=el;scale=1;tx=0;ty=0;isDragging=false;stage.appendChild(el);apply();overlay.classList.add("active");document.body.classList.add("lb-open");}',
      'function openImage(img){var el=new Image();el.src=img.currentSrc||img.src;el.alt=img.alt||"";el.className="lb-media lb-img";mount(el);}',
      'function openSvg(svg){var el=svg.cloneNode(true);el.classList.add("lb-media","lb-svg");el.removeAttribute("width");el.removeAttribute("height");mount(el);}',
      'function close(){overlay.classList.remove("active");document.body.classList.remove("lb-open");stage.innerHTML="";activeEl=null;scale=1;tx=0;ty=0;isDragging=false;}',
      'document.querySelectorAll("img").forEach(function(el){if(el.closest(".lb-overlay"))return;el.style.cursor="zoom-in";el.addEventListener("click",function(e){e.preventDefault();e.stopPropagation();openImage(el);});});',
      'overlay.addEventListener("click",function(e){if(e.target===overlay||e.target===stage)close();});',
      'closeBtn.addEventListener("click",function(e){e.stopPropagation();close();});',
      'document.addEventListener("keydown",function(e){if(e.key==="Escape"&&overlay.classList.contains("active"))close();});',
      'overlay.addEventListener("wheel",function(e){if(!activeEl)return;e.preventDefault();var delta=e.deltaY<0?1.1:0.9;scale=Math.min(Math.max(scale*delta,0.25),8);apply();},{passive:false});',
      'stage.addEventListener("mousedown",function(e){if(!activeEl||e.target!==activeEl)return;isDragging=true;sx=e.clientX-tx;sy=e.clientY-ty;activeEl.style.cursor="grabbing";e.preventDefault();});',
      'document.addEventListener("mousemove",function(e){if(!isDragging||!activeEl)return;tx=e.clientX-sx;ty=e.clientY-sy;apply();});',
      'document.addEventListener("mouseup",function(){if(activeEl)activeEl.style.cursor="grab";isDragging=false;});',
      'function toast(msg){var t=document.createElement("div");t.style.cssText="position:fixed;bottom:24px;left:50%;transform:translateX(-50%);background:#171717;color:#fff;padding:8px 16px;border-radius:6px;font-size:13px;z-index:99999;opacity:0;transition:opacity .2s";t.textContent=msg;document.body.appendChild(t);requestAnimationFrame(function(){t.style.opacity="1"});setTimeout(function(){t.style.opacity="0";setTimeout(function(){t.remove()},200)},2000);}',
      'function downloadSvg(svg,index){var bbox=svg.getBBox();var pad=25;var vx=bbox.x-pad;var vy=bbox.y-pad;var vw=Math.ceil(bbox.width+pad*2);var vh=Math.ceil(bbox.height+pad*2);var clone=svg.cloneNode(true);clone.setAttribute("xmlns","http://www.w3.org/2000/svg");clone.setAttribute("viewBox",vx+" "+vy+" "+vw+" "+vh);clone.setAttribute("width",vw);clone.setAttribute("height",vh);clone.removeAttribute("style");var data=new XMLSerializer().serializeToString(clone);var blob=new Blob([data],{type:"image/svg+xml;charset=utf-8"});var url=URL.createObjectURL(blob);var a=document.createElement("a");a.href=url;a.download="diagram-"+(index+1)+".svg";document.body.appendChild(a);a.click();a.remove();URL.revokeObjectURL(url);toast("SVG 已下载");}',
      'function downloadPng(svg,index,ratio){var bbox=svg.getBBox();var pad=25;var vw=Math.ceil(bbox.width+pad*2);var vh=Math.ceil(bbox.height+pad*2);var sc=Math.max(2,ratio||2);var clone=svg.cloneNode(true);clone.setAttribute("xmlns","http://www.w3.org/2000/svg");clone.setAttribute("xmlns:xlink","http://www.w3.org/1999/xlink");clone.setAttribute("viewBox",(bbox.x-pad)+" "+(bbox.y-pad)+" "+vw+" "+vh);clone.setAttribute("width",vw);clone.setAttribute("height",vh);clone.removeAttribute("style");var data=new XMLSerializer().serializeToString(clone);var base64="data:image/svg+xml;base64,"+btoa(unescape(encodeURIComponent(data)));var canvas=document.createElement("canvas");canvas.width=vw*sc;canvas.height=vh*sc;var ctx=canvas.getContext("2d");var img=new Image();img.onload=function(){ctx.fillStyle="#ffffff";ctx.fillRect(0,0,canvas.width,canvas.height);ctx.drawImage(img,0,0,canvas.width,canvas.height);canvas.toBlob(function(blob){if(!blob){toast("PNG 导出失败");return;}var url=URL.createObjectURL(blob);var a=document.createElement("a");a.href=url;a.download="diagram-"+(index+1)+".png";document.body.appendChild(a);a.click();a.remove();URL.revokeObjectURL(url);toast("PNG 已下载（"+sc+"x）");},"image/png");};img.onerror=function(){toast("PNG 导出失败，请使用 SVG 格式");};img.src=base64;}',
      'function bindMermaidInteractions(){',
      '  document.querySelectorAll(".mermaid").forEach(function(container,idx){',
      '    var svg=container.querySelector("svg");',
      '    if(!svg)return;',
      '    container.style.position="relative";',
      '    var overlay_div=document.createElement("div");',
      '    overlay_div.style.cssText="position:absolute;inset:0;cursor:zoom-in;z-index:1";',
      '    overlay_div.setAttribute("tabindex","0");',
      '    overlay_div.setAttribute("role","button");',
      '    overlay_div.setAttribute("aria-label","点击放大查看流程图");',
      '    overlay_div.addEventListener("click",function(e){e.preventDefault();e.stopPropagation();openSvg(svg);});',
      '    container.appendChild(overlay_div);',
      '    var tb=document.createElement("div");',
      '    tb.className="mm-toolbar";',
      '    tb.innerHTML="<select class=\\"mm-scale\\"><option value=\\"2\\" selected>2x</option><option value=\\"4\\">4x</option><option value=\\"6\\">6x</option></select><button class=\\"mm-btn mm-svg\\">SVG</button><button class=\\"mm-btn mm-png\\">PNG</button>";',
      '    container.appendChild(tb);',
      '    var sel=tb.querySelector(".mm-scale");',
      '    tb.querySelector(".mm-svg").addEventListener("click",function(e){e.stopPropagation();downloadSvg(svg,idx);});',
      '    tb.querySelector(".mm-png").addEventListener("click",function(e){e.stopPropagation();var scl=Math.max(1,Number(sel.value)||1);downloadPng(svg,idx,scl);});',
      '  });',
      '}',
      'if(typeof mermaid!=="undefined"){',
      '  mermaid.initialize({startOnLoad:false,securityLevel:"loose",theme:"' + mermaidTheme + '"});',
      '  mermaid.run().then(function(){bindMermaidInteractions();}).catch(function(){bindMermaidInteractions();});',
      '}else{',
      '  bindMermaidInteractions();',
      '}',
      '})();',
      scTag,
      '<script>',
      '(function(){',
      'var sb=document.getElementById("sidebar");',
      'var fab=document.getElementById("sbFab");',
      'var tg=document.getElementById("sbToggle");',
      'if(!sb)return;',
      'function setOpen(open){',
      '  if(open){document.body.classList.add("sb-open");document.body.classList.remove("sb-collapsed");}',
      '  else{document.body.classList.remove("sb-open");document.body.classList.add("sb-collapsed");}',
      '}',
      'setOpen(true);',
      'if(tg){tg.addEventListener("click",function(){var isOpen=document.body.classList.contains("sb-open");setOpen(!isOpen);});}',
      'if(fab){fab.addEventListener("click",function(){setOpen(true);});}',
      'var links=sb.querySelectorAll(".sb-item a");',
      'links.forEach(function(a){',
      '  a.addEventListener("click",function(e){',
      '    var id=a.getAttribute("data-target");',
      '    var el=id?document.getElementById(id):null;',
      '    if(el){e.preventDefault();el.scrollIntoView({behavior:"smooth",block:"start"});history.replaceState(null,"","#"+id);}',
      '  });',
      '});',
      'var headings=Array.prototype.slice.call(document.querySelectorAll("h1,h2,h3,h4,h5,h6")).filter(function(h){return h.id});',
      'function updateActive(){',
      '  if(!headings.length)return;',
      '  var top=window.scrollY+80;',
      '  var current=headings[0].id;',
      '  for(var i=0;i<headings.length;i++){if(headings[i].offsetTop<=top)current=headings[i].id;}',
      '  links.forEach(function(a){var li=a.parentElement;if(a.getAttribute("data-target")===current)li.classList.add("active");else li.classList.remove("active");});',
      '}',
      'window.addEventListener("scroll",updateActive,{passive:true});',
      'updateActive();',
      'window.addEventListener("hashchange",updateActive);',
      '})();',
      scTag,
    ].join('')

    return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>${escapeHtml(editorStore.filename || 'markdown-output')}</title>
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/katex.min.css" />
<link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/styles/atom-one-dark.min.css" />
<style>
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;background:${bg};color:${text};line-height:1.75;margin:0;padding:0;scroll-behavior:smooth}
.main-wrap{max-width:860px;margin:0 auto;padding:32px 20px 60px;transition:margin-left .25s ease}
h1{font-size:1.75em;border-bottom:1px solid ${border};padding-bottom:.3em;margin:0 0 .6em}
h2{font-size:1.35em;margin-top:1.6em}
h3{font-size:1.15em}
table{display:block;overflow:auto;border-collapse:collapse;width:100%}
th,td{border:1px solid ${border};padding:8px 12px}
th{background:${surface};font-weight:600}
tr:nth-child(even){background:${surface}}
del{color:#a3a3a3;text-decoration:line-through}
input[type="checkbox"]{margin-right:6px;accent-color:#2563eb;vertical-align:middle}
pre{background:#0a0a0a;color:#e5e5e5;padding:36px 16px 12px;border-radius:8px;overflow:auto;position:relative}
code{font-family:"SF Mono","Fira Code",Consolas,monospace;font-size:.88em}
pre code{display:block;background:transparent;padding:0 16px 14px;font-size:12.5px;line-height:1.65;white-space:pre}
pre code.hljs{background:transparent!important}
.code-header{position:absolute;top:0;left:0;right:0;height:36px;display:flex;align-items:center;justify-content:space-between;padding:0 12px;pointer-events:none;z-index:2}
.code-lang{font-size:10px;color:#525252;text-transform:uppercase;letter-spacing:0.05em;font-family:"SF Mono","Fira Code",Consolas,monospace;user-select:none;pointer-events:none}
.copy-btn{position:absolute;top:6px;right:8px;width:28px;height:28px;padding:0;border:1px solid rgba(255,255,255,0.08);background:rgba(255,255,255,0.04);color:#525252;border-radius:6px;font-size:10px;cursor:pointer;display:inline-flex;align-items:center;justify-content:center;opacity:0;z-index:3;transition:opacity .15s}
pre:hover .copy-btn{opacity:1}
.copy-btn:hover{background:rgba(255,255,255,0.10);color:#e5e5e5;border-color:rgba(255,255,255,0.15)}
.copy-btn.copied{color:#22c55e;border-color:rgba(34,197,94,0.25);background:rgba(34,197,94,0.08)}
blockquote{border-left:3px solid ${border};padding:8px 16px;margin:0 0 1em;color:#525252;background:${surface};border-radius:0 6px 6px 0}
img{max-width:100%;border-radius:8px;border:1px solid ${border}}
a{color:#2563eb;text-decoration:none}
a:hover{text-decoration:underline}
.mermaid{position:relative;background:${surface};border:1px solid ${border};border-radius:8px;overflow:hidden;margin:1em 0;padding:16px;display:flex;justify-content:center;align-items:center;overflow:auto}
.mermaid svg{max-width:100%;height:auto}
.mermaid-inner{padding:16px;display:flex;justify-content:center;align-items:center;overflow:auto;width:100%}
${lbStyle}
</style>
</head>
<body class="sb-open">
${sidebar}
<main class="main-wrap">
${body}
</main>
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/highlight.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/c.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/cpp.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/java.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/python.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/go.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/javascript.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/typescript.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/ruby.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/bash.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/yaml.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/json.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/xml.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/css.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/sql.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/rust.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/csharp.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/php.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/swift.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/kotlin.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/dart.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/dockerfile.min.js">${scTag}
<script>
document.querySelectorAll("pre code").forEach(function(el){hljs.highlightElement(el);});
${scTag}
<script src="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/katex.min.js">${scTag}
<script src="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/contrib/auto-render.min.js">${scTag}
<script>
renderMathInElement(document.body,{delimiters:[{left:"$$",right:"$$",display:true},{left:"$",right:"$",display:false},{left:"\\\\(",right:"\\\\)",display:false},{left:"\\\\[",right:"\\\\]",display:true}],throwOnError:false});
${scTag}
${lbScript}
</body>
</html>`
  }

  function downloadHTML() {
    try {
      const html = buildFullHTML()
      const blob = new Blob([html], { type: 'text/html;charset=utf-8' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      const name = (editorStore.filename || 'markdown-output').replace(/\.html$/i, '')
      a.href = url
      a.download = name + '.html'
      document.body.appendChild(a)
      a.click()
      a.remove()
      URL.revokeObjectURL(url)
      toast('HTML 已导出')
    } catch (e) {
      console.error('[exportHtml] downloadHTML error:', e)
      toast('导出失败: ' + (e instanceof Error ? e.message : String(e)), 'err')
    }
  }

  function downloadAllHTML() {
    const tabs = editorStore.tabs
    if (tabs.length <= 1) {
      downloadHTML()
      return
    }

    const tabsWithContent = tabs.filter(t => t.content.trim())
    if (tabsWithContent.length === 0) {
      toast('没有可导出的内容', 'err')
      return
    }

    // 逐个下载，间隔 200ms 避免浏览器拦截
    tabsWithContent.forEach((tab, idx) => {
      setTimeout(() => {
        // 临时切换到该 tab 的内容构建 HTML
        const savedActiveId = editorStore.activeTabId
        editorStore.switchTab(tab.id)

        // 需要等 tab 切换后构建
        nextTick(() => {
          try {
            const html = buildFullHTML()
            const blob = new Blob([html], { type: 'text/html;charset=utf-8' })
            const url = URL.createObjectURL(blob)
            const a = document.createElement('a')
            const name = (tab.filename || 'document').replace(/\.html$/i, '')
            a.href = url
            a.download = name + '.html'
            document.body.appendChild(a)
            a.click()
            a.remove()
            URL.revokeObjectURL(url)
          } catch (e) {
            console.error('[exportHtml] batch download error:', e)
          }

          // 恢复原来的活跃 tab
          if (savedActiveId !== tab.id) {
            editorStore.switchTab(savedActiveId)
          }

          if (idx === tabsWithContent.length - 1) {
            toast(`已导出 ${tabsWithContent.length} 个 HTML 文件`)
          }
        })
      }, idx * 200)
    })
  }

  function copyHTML() {
    try {
      const html = buildFullHTML()
      navigator.clipboard.writeText(html).then(() => toast('HTML 源码已复制')).catch((e) => {
        console.error('[exportHtml] copyHTML error:', e)
        toast('复制失败: ' + String(e), 'err')
      })
    } catch (e) {
      console.error('[exportHtml] copyHTML error:', e)
      toast('复制失败: ' + (e instanceof Error ? e.message : String(e)), 'err')
    }
  }

  return {
    buildFullHTML,
    downloadHTML,
    downloadAllHTML,
    copyHTML
  }
}
