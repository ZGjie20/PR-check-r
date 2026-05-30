import { useState } from 'react';

interface DiffViewerProps {
  diff: string;
}

export function DiffViewer({ diff }: DiffViewerProps) {
  const [expanded, setExpanded] = useState(false);

  if (!diff) return null;

  const lineCount = diff.split('\n').length;

  return (
    <div className="glass-card overflow-hidden">
      <button
        type="button"
        onClick={() => setExpanded((v) => !v)}
        className="flex w-full items-center justify-between px-5 py-4 text-left transition-colors hover:bg-white/5"
      >
        <div className="flex items-center gap-2">
          <span>📝</span>
          <span className="text-sm font-medium text-white/70">原始 Diff</span>
          <span className="rounded-full bg-white/10 px-2 py-0.5 text-[11px] text-white/40">
            {lineCount} 行
          </span>
        </div>
        <span className="text-xs text-white/40">{expanded ? '收起 ▲' : '展开 ▼'}</span>
      </button>
      {expanded && (
        <pre className="max-h-96 overflow-auto border-t border-white/10 bg-black/40 p-4 font-mono text-xs leading-relaxed text-emerald-300/80">
          {diff}
        </pre>
      )}
    </div>
  );
}
