import type { ReviewIssueDetail, Severity } from '@/types/review';
import { SeverityBadge } from '@/components/review/SeverityBadge';
import { Card, CardBody } from '@/components/ui/Card';

interface IssueCardProps {
  issue: ReviewIssueDetail;
  severity: Severity;
}

export function IssueCard({ issue, severity }: IssueCardProps) {
  return (
    <Card className="group">
      <CardBody>
        <div className="flex items-start gap-3">
          <div className="mt-1 flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-white/5 text-xs transition-colors group-hover:bg-white/10">
            📄
          </div>
          <div className="min-w-0 flex-1">
            <div className="flex flex-wrap items-center gap-2">
              <SeverityBadge severity={severity} />
              <span className="rounded-md bg-white/5 px-2 py-0.5 font-mono text-xs text-white/50">
                {issue.file}:{issue.line}
              </span>
            </div>
            <p className="mt-3 text-sm leading-relaxed text-white/80">{issue.message}</p>
            {issue.suggestion && (
              <div className="mt-3 rounded-xl border border-dream-400/20 bg-dream-500/10 px-4 py-3">
                <p className="mb-1 text-[11px] font-medium tracking-wider text-dream-300/70">
                  💡 修复建议
                </p>
                <p className="text-sm text-dream-100/80">{issue.suggestion}</p>
              </div>
            )}
          </div>
        </div>
      </CardBody>
    </Card>
  );
}
