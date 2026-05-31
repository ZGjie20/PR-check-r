import { Card, CardBody } from '@/components/ui/Card';

interface PrChangeSummaryProps {
  summary?: string;
}

export function PrChangeSummary({ summary }: PrChangeSummaryProps) {
  const text = summary?.trim();
  if (!text) {
    return null;
  }

  return (
    <section>
      <h2 className="mb-5 flex items-center gap-2 text-lg font-semibold">
        <span>📝</span>
        <span className="gradient-text">PR 变更总结</span>
      </h2>
      <Card>
        <CardBody>
          <p className="whitespace-pre-wrap text-sm leading-relaxed text-white/80">
            {text}
          </p>
        </CardBody>
      </Card>
    </section>
  );
}
