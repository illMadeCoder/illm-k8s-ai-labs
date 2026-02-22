/**
 * Tooltip utility for graph nodes.
 * Reads data-tooltip-* attributes, renders a position:fixed tooltip near the node
 * in viewport coords, flips at edges, hides during drag.
 */

export interface TooltipOptions {
  /** Offset from cursor/node in px */
  offset?: number;
  /** Additional CSS class for the tooltip element */
  className?: string;
  /** Return true to suppress tooltip for a given target element */
  shouldSuppress?: (target: HTMLElement) => boolean;
}

export function initTooltip(
  container: HTMLElement,
  selector: string,
  opts: TooltipOptions = {},
): { cleanup: () => void } {
  const offset = opts.offset ?? 12;

  // Create tooltip element
  const tip = document.createElement('div');
  tip.style.cssText =
    'position:fixed;z-index:9999;pointer-events:none;opacity:0;transition:opacity 0.15s;' +
    'background:var(--bg-card);border:1px solid var(--border);border-radius:6px;' +
    'padding:8px 12px;font-size:0.75rem;font-family:var(--font-mono);' +
    'color:var(--text);box-shadow:var(--shadow-lg);max-width:280px;line-height:1.5;';
  if (opts.className) tip.classList.add(opts.className);
  document.body.appendChild(tip);

  let visible = false;

  function show(target: HTMLElement, e: MouseEvent) {
    if (opts.shouldSuppress?.(target)) { hide(); return; }

    // Build content from data-tooltip-* attributes
    const attrs = target.dataset;
    const title = attrs.tooltipTitle;
    const body = attrs.tooltipBody;
    const meta = attrs.tooltipMeta;

    if (!title && !body) return;

    let html = '';
    if (title) {
      const color = attrs.tooltipColor || 'var(--accent)';
      html += `<div style="font-weight:600;color:${color};margin-bottom:4px">${escapeHtml(title)}</div>`;
    }
    if (body) {
      html += `<div style="color:var(--text)">${escapeHtml(body)}</div>`;
    }
    if (meta) {
      html += `<div style="color:var(--text-muted);margin-top:4px;font-size:0.7rem">${escapeHtml(meta)}</div>`;
    }

    tip.innerHTML = html;
    tip.style.opacity = '1';
    visible = true;

    positionTip(e);
  }

  function positionTip(e: MouseEvent) {
    if (!visible) return;

    const vw = window.innerWidth;
    const vh = window.innerHeight;
    const tipRect = tip.getBoundingClientRect();

    let left = e.clientX + offset;
    let top = e.clientY + offset;

    // Flip horizontally if overflowing right
    if (left + tipRect.width > vw - 8) {
      left = e.clientX - tipRect.width - offset;
    }

    // Flip vertically if overflowing bottom
    if (top + tipRect.height > vh - 8) {
      top = e.clientY - tipRect.height - offset;
    }

    // Clamp to viewport
    left = Math.max(8, Math.min(left, vw - tipRect.width - 8));
    top = Math.max(8, Math.min(top, vh - tipRect.height - 8));

    tip.style.left = `${left}px`;
    tip.style.top = `${top}px`;
  }

  function hide() {
    tip.style.opacity = '0';
    visible = false;
  }

  // Event delegation on container
  function onMouseOver(e: MouseEvent) {
    const target = (e.target as HTMLElement).closest(selector) as HTMLElement | null;
    if (target) {
      show(target, e);
    }
  }

  function onMouseMove(e: MouseEvent) {
    if (visible) {
      const target = (e.target as HTMLElement).closest(selector) as HTMLElement | null;
      if (target) {
        if (opts.shouldSuppress?.(target)) { hide(); return; }
        positionTip(e);
      } else {
        hide();
      }
    }
  }

  function onMouseOut(e: MouseEvent) {
    const target = (e.target as HTMLElement).closest(selector) as HTMLElement | null;
    if (target) {
      const related = e.relatedTarget as HTMLElement | null;
      if (!related || !target.contains(related)) {
        hide();
      }
    }
  }

  container.addEventListener('mouseover', onMouseOver);
  container.addEventListener('mousemove', onMouseMove);
  container.addEventListener('mouseout', onMouseOut);

  function cleanup() {
    container.removeEventListener('mouseover', onMouseOver);
    container.removeEventListener('mousemove', onMouseMove);
    container.removeEventListener('mouseout', onMouseOut);
    tip.remove();
  }

  return { cleanup };
}

function escapeHtml(str: string): string {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}
