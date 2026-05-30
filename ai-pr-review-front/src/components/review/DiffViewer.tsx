import { useMemo, useState } from 'react';
import {
  SEVERITY_LEGEND,
  SEVERITY_UNDERLINE_CLASS,
} from '@/constants/severity';
import type { ReviewResultBySeverity } from '@/types/review';
import {
  buildIssueLineMap,
  countHighlightedLines,
  parseDiffWithIssues,
  type DiffDisplayLine,
} from '@/utils/diffHighlight';

interface DiffViewerProps {
  diff: string;
  reviewResult?: ReviewResultBySeverity;
}

function DiffLine({ line }: { line: DiffDisplayLine }) {
  if (line.changeType === 'meta') {
    return (
      <div className="whitespace-pre-wrap text-white/45">{line.text || ' '}</div>
    );
  }

  const prefix = line.text[0] ?? ' ';
  const content = line.text.slice(1);
  const prefixClass =
    prefix === '+'
      ? 'text-emerald-400'
      : prefix === '-'
        ? 'text-red-400/80'
        : 'text-white/50';

  const underlineClass = line.severity
    ? SEVERITY_UNDERLINE_CLASS[line.severity]
    : '';
  const title =
    line.severity && line.file && line.newLine
      ? `${line.file}:${line.newLine} · ${SEVERITY_LEGEND.find((item) => item.severity === line.severity)?.label ?? line.severity}`
      : undefined;

  return (
    <div className="whitespace-pre-wrap">
      <span className={prefixClass}>{prefix}</span>
      <span className={underlineClass} title={title}>
        {content}
      </span>
    </div>
  );
}

export function DiffViewer({ diff, reviewResult }: DiffViewerProps) {
  const [expanded, setExpanded] = useState(false);

  const { lines, highlightedCount, hasIssues } = useMemo(() => {
    const issueMap = reviewResult
      ? buildIssueLineMap(reviewResult, diff)
      : new Map();
    const parsedLines = parseDiffWithIssues(diff, issueMap);
    const highlighted = countHighlightedLines(parsedLines);
    const issueCount =
      (reviewResult?.high.length ?? 0) +
      (reviewResult?.medium.length ?? 0) +
      (reviewResult?.low.length ?? 0);

    return {
      lines: parsedLines,
      highlightedCount: highlighted,
      hasIssues: issueCount > 0,
    };
  }, [diff, reviewResult]);

  if (!diff) return null;

  const lineCount = diff.split('\n').length;

  return (
    <div className="glass-card overflow-hidden">
      <button
        type="button"
        onClick={() => setExpanded((value) => !value)}
        className="flex w-full items-center justify-between px-5 py-4 text-left transition-colors hover:bg-white/5"
      >
        <div className="flex items-center gap-2">
          <span>📝</span>
          <span className="text-sm font-medium text-white/70">原始 Diff</span>
          <span className="rounded-full bg-white/10 px-2 py-0.5 text-[11px] text-white/40">
            {lineCount} 行
          </span>
          {hasIssues && highlightedCount > 0 && (
            <span className="rounded-full bg-red-500/15 px-2 py-0.5 text-[11px] text-red-200">
              {highlightedCount} 处已标注
            </span>
          )}
        </div>
        <span className="text-xs text-white/40">{expanded ? '收起 ▲' : '展开 ▼'}</span>
      </button>

      {expanded && (
        <>
          {hasIssues && (
            <div className="flex flex-wrap items-center gap-4 border-t border-white/10 bg-white/[0.02] px-5 py-3">
              <span className="text-[11px] text-white/40">风险标注：</span>
              {SEVERITY_LEGEND.map((item) => (
                <span
                  key={item.severity}
                  className="flex items-center gap-1.5 text-[11px] text-white/60"
                >
                  <span
                    className={`inline-block h-0 w-6 shrink-0 border-b-2 ${item.lineClassName}`}
                    aria-hidden
                  />
                  {item.label}
                </span>
              ))}
            </div>
          )}

          <pre className="max-h-96 overflow-auto border-t border-white/10 bg-black/40 p-4 font-mono text-xs leading-relaxed text-emerald-300/80">
            {lines.map((line, index) => (
              <DiffLine key={`${index}-${line.text}`} line={line} />
            ))}
          </pre>
        </>
      )}
    </div>
  );
}
