interface ReviewSummaryProps {
  totalIssues: number;
  highIssues: number;
  mediumIssues: number;
  lowIssues: number;
}

const STAT_CONFIG = [
  { key: 'total', label: '总计', icon: '📊', color: 'text-white', glow: 'from-white/10' },
  { key: 'high', label: '高风险', icon: '🔴', color: 'text-red-300', glow: 'from-red-500/20' },
  { key: 'medium', label: '中风险', icon: '🟠', color: 'text-orange-300', glow: 'from-orange-500/20' },
  { key: 'low', label: '低风险', icon: '🔵', color: 'text-cyan-300', glow: 'from-cyan-500/20' },
] as const;

export function ReviewSummary({
  totalIssues,
  highIssues,
  mediumIssues,
  lowIssues,
}: ReviewSummaryProps) {
  const values: Record<string, number> = {
    total: totalIssues,
    high: highIssues,
    medium: mediumIssues,
    low: lowIssues,
  };

  return (
    <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
      {STAT_CONFIG.map((stat) => (
        <div key={stat.key} className="stat-card group">
          <div className="relative z-10">
            <div className="mb-2 flex items-center justify-center gap-1.5">
              <span className="text-sm">{stat.icon}</span>
              <p className="text-xs text-white/50">{stat.label}</p>
            </div>
            <p className={`text-3xl font-bold ${stat.color}`}>{values[stat.key]}</p>
          </div>
        </div>
      ))}
    </div>
  );
}
