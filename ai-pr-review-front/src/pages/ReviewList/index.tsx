import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useReviewList } from '@/hooks/useReview';
import { PageHeader } from '@/components/layout/PageHeader';
import { PageDecorBanner, DECOR_IMAGES } from '@/components/layout/PageDecor';
import { Card, CardBody } from '@/components/ui/Card';
import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';
import { Pagination } from '@/components/ui/Pagination';
import { Spinner } from '@/components/ui/Spinner';
import { EmptyState } from '@/components/common/EmptyState';
import { ErrorAlert } from '@/components/common/ErrorAlert';
import { formatDate } from '@/utils/format';

const PAGE_SIZE = 20;

export function ReviewListPage() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [prNumberInput, setPrNumberInput] = useState('');
  const [prNumberFilter, setPrNumberFilter] = useState<number | undefined>();

  const { data, loading, error } = useReviewList(page, PAGE_SIZE, prNumberFilter);

  function handleFilter() {
    const num = Number.parseInt(prNumberInput, 10);
    setPrNumberFilter(prNumberInput && num > 0 ? num : undefined);
    setPage(1);
  }

  function handleClearFilter() {
    setPrNumberInput('');
    setPrNumberFilter(undefined);
    setPage(1);
  }

  return (
    <div className="space-y-6">
      <PageDecorBanner
        image={DECOR_IMAGES.campfire}
        objectPosition="22% 28%"
      />

      <PageHeader
        title="审查历史"
        subtitle="浏览所有 PR 审查记录"
        icon="📜"
      />

      <div className="glass-card flex flex-col gap-3 p-4 sm:flex-row sm:items-end">
        <div className="flex-1">
          <label htmlFor="pr-filter" className="mb-2 block text-xs text-white/50">
            🔎 PR 编号筛选
          </label>
          <Input
            id="pr-filter"
            type="number"
            placeholder="如 123"
            value={prNumberInput}
            onChange={(e) => setPrNumberInput(e.target.value)}
          />
        </div>
        <div className="flex gap-2">
          <Button variant="secondary" onClick={handleFilter}>
            筛选
          </Button>
          {prNumberFilter && (
            <Button variant="ghost" onClick={handleClearFilter}>
              清除
            </Button>
          )}
        </div>
      </div>

      {error && <ErrorAlert message={error} />}

      {loading ? (
        <div className="flex justify-center py-20">
          <Spinner size="lg" />
        </div>
      ) : !data || data.items.length === 0 ? (
        <EmptyState
          title="暂无审查记录"
          description="提交第一个 PR 链接，开启你的代码审查之旅吧"
          icon="🌙"
        />
      ) : (
        <>
          <Card>
            <CardBody className="p-0">
              <div className="overflow-x-auto">
                <table className="w-full text-left text-sm">
                  <thead>
                    <tr className="border-b border-white/10 text-white/40">
                      <th className="px-6 py-4 font-medium">PR 标题</th>
                      <th className="px-6 py-4 font-medium">编号</th>
                      <th className="px-6 py-4 font-medium">问题数</th>
                      <th className="hidden px-6 py-4 font-medium sm:table-cell">
                        创建时间
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    {data.items.map((item) => (
                      <tr
                        key={item.id}
                        onClick={() => navigate(`/reviews/${item.id}`)}
                        className="table-row-glass"
                      >
                        <td className="px-6 py-4">
                          <div className="flex items-center gap-2">
                            <span className="text-xs opacity-40">📌</span>
                            <span className="font-medium text-white/90">
                              {item.pr_title}
                            </span>
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <span className="rounded-md bg-dream-500/20 px-2 py-0.5 text-xs text-dream-200">
                            #{item.pr_number}
                          </span>
                        </td>
                        <td className="px-6 py-4">
                          <span
                            className={`rounded-full px-2.5 py-0.5 text-xs font-medium ${
                              item.total_issues > 0
                                ? 'bg-sakura-500/20 text-sakura-200'
                                : 'bg-emerald-500/20 text-emerald-200'
                            }`}
                          >
                            {item.total_issues} 个问题
                          </span>
                        </td>
                        <td className="hidden px-6 py-4 text-white/40 sm:table-cell">
                          {formatDate(item.created_at)}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </CardBody>
          </Card>
          <Pagination
            page={data.page}
            total={data.total}
            limit={data.limit}
            onPageChange={setPage}
          />
        </>
      )}
    </div>
  );
}
