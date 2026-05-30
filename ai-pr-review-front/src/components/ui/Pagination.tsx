import { Button } from '@/components/ui/Button';

interface PaginationProps {
  page: number;
  total: number;
  limit: number;
  onPageChange: (page: number) => void;
}

export function Pagination({ page, total, limit, onPageChange }: PaginationProps) {
  const totalPages = Math.max(1, Math.ceil(total / limit));
  const hasPrev = page > 1;
  const hasNext = page * limit < total;

  if (total === 0) return null;

  return (
    <div className="glass-card flex items-center justify-between gap-4 px-5 py-4">
      <p className="text-sm text-white/50">
        共 <span className="font-medium text-sakura-300">{total}</span> 条 · 第{' '}
        <span className="font-medium text-white/80">{page}</span> / {totalPages} 页
      </p>
      <div className="flex gap-2">
        <Button variant="secondary" disabled={!hasPrev} onClick={() => onPageChange(page - 1)}>
          ← 上一页
        </Button>
        <Button variant="secondary" disabled={!hasNext} onClick={() => onPageChange(page + 1)}>
          下一页 →
        </Button>
      </div>
    </div>
  );
}
