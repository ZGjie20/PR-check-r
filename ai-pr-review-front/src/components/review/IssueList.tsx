import type { ReviewResultBySeverity } from '@/types/review';
import { SEVERITY_CONFIG, SEVERITY_ORDER } from '@/constants/severity';
import { IssueCard } from '@/components/review/IssueCard';

interface IssueListProps {
  reviewResult: ReviewResultBySeverity;
}

export function IssueList({ reviewResult }: IssueListProps) {
  const hasIssues = SEVERITY_ORDER.some(
    (severity) => reviewResult[severity].length > 0,
  );

  if (!hasIssues) {
    return (
      <div className="glass-card flex flex-col items-center py-12 text-center">
        <span className="mb-3 text-4xl">🎉</span>
        <p className="gradient-text text-base font-medium">未发现问题</p>
        <p className="mt-1 text-sm text-white/40">代码审查通过，干得漂亮！</p>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {SEVERITY_ORDER.map((severity) => {
        const issues = reviewResult[severity];
        if (issues.length === 0) return null;
        const config = SEVERITY_CONFIG[severity];

        return (
          <section key={severity}>
            <div className="mb-4 flex items-center gap-3">
              <span className={`h-2.5 w-2.5 rounded-full ${config.dotClass}`} />
              <h3 className="text-sm font-semibold text-white/70">
                {config.label}风险
              </h3>
              <span className="rounded-full bg-white/10 px-2 py-0.5 text-xs text-white/50">
                {issues.length}
              </span>
            </div>
            <div className="space-y-3">
              {issues.map((issue, index) => (
                <IssueCard
                  key={`${issue.file}-${issue.line}-${index}`}
                  issue={issue}
                  severity={severity}
                />
              ))}
            </div>
          </section>
        );
      })}
    </div>
  );
}
