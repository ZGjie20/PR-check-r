import type { Severity } from '@/types/review';

export interface SeverityConfig {
  label: string;
  badgeClass: string;
  dotClass: string;
  glowClass: string;
}

export const SEVERITY_CONFIG: Record<Severity, SeverityConfig> = {
  high: {
    label: '高',
    badgeClass: 'border-red-400/40 bg-red-500/20 text-red-200',
    dotClass: 'bg-red-400 shadow-[0_0_8px_rgba(248,113,113,0.6)]',
    glowClass: 'from-red-500/20 to-red-500/5',
  },
  medium: {
    label: '中',
    badgeClass: 'border-orange-400/40 bg-orange-500/20 text-orange-200',
    dotClass: 'bg-orange-400 shadow-[0_0_8px_rgba(251,146,60,0.6)]',
    glowClass: 'from-orange-500/20 to-orange-500/5',
  },
  low: {
    label: '低',
    badgeClass: 'border-cyan-400/40 bg-cyan-500/20 text-cyan-200',
    dotClass: 'bg-cyan-400 shadow-[0_0_8px_rgba(34,211,238,0.6)]',
    glowClass: 'from-cyan-500/20 to-cyan-500/5',
  },
};

export const SEVERITY_ORDER: Severity[] = ['high', 'medium', 'low'];

export const SEVERITY_UNDERLINE_CLASS: Record<Severity, string> = {
  high: 'underline decoration-red-400 decoration-2 underline-offset-[3px] decoration-solid',
  medium:
    'underline decoration-yellow-400 decoration-2 underline-offset-[3px] decoration-solid',
  low: 'underline decoration-pink-400 decoration-2 underline-offset-[3px] decoration-solid',
};

export const SEVERITY_LEGEND: Array<{ severity: Severity; label: string; lineClassName: string }> = [
  { severity: 'high', label: '高风险', lineClassName: 'border-red-400' },
  { severity: 'medium', label: '中风险', lineClassName: 'border-yellow-400' },
  { severity: 'low', label: '低风险', lineClassName: 'border-pink-400' },
];
