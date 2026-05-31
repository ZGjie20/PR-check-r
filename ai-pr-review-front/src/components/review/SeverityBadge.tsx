import type { Severity } from '@/types/review';
import { SEVERITY_CONFIG } from '@/constants/severity';
import { Badge } from '@/components/ui/Badge';

interface SeverityBadgeProps {
  severity: Severity;
}

export function SeverityBadge({ severity }: SeverityBadgeProps) {
  const config = SEVERITY_CONFIG[severity];
  return <Badge className={config.badgeClass}>{config.label}</Badge>;
}
