import { Link, useParams } from 'react-router-dom';
import { useReviewDetail } from '@/hooks/useReview';
import { ReviewSummary } from '@/components/review/ReviewSummary';
import { PrChangeSummary } from '@/components/review/PrChangeSummary';
import { IssueList } from '@/components/review/IssueList';
import { ReviewActions } from '@/components/review/ReviewActions';
import { DiffViewer } from '@/components/review/DiffViewer';
import { DecorImage, FG_IMAGE_OPACITY } from '@/components/common/DecorImage';
import { DECOR_IMAGES } from '@/constants/decorImages';
import { Card, CardBody, CardHeader } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Spinner } from '@/components/ui/Spinner';
import { ErrorAlert } from '@/components/common/ErrorAlert';
import { formatDate } from '@/utils/format';

export function ReviewDetailPage() {
  const { id } = useParams<{ id: string }>();
  const reviewId = id ? Number.parseInt(id, 10) : null;
  const { data, loading, error, notFound } = useReviewDetail(reviewId);

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center gap-4 py-24">
        <Spinner size="lg" />
        <p className="text-sm text-white/40">正在加载审查详情...</p>
      </div>
    );
  }

  if (notFound) {
    return (
      <div className="flex flex-col items-center gap-6 py-16 text-center">
        <DecorImage
          src={DECOR_IMAGES.catgirl}
          variant="card"
          opacity={FG_IMAGE_OPACITY.card}
          objectPosition="center 15%"
          className="mx-auto h-56 w-56 sm:h-64 sm:w-64"
        />
        <div>
          <h1 className="text-xl font-bold gradient-text">记录不存在</h1>
          <p className="mt-2 text-sm text-white/40">
            该审查记录可能已被删除或 ID 无效
          </p>
        </div>
        <Link to="/reviews">
          <Button variant="secondary">← 返回列表</Button>
        </Link>
      </div>
    );
  }

  if (error || !data) {
    return <ErrorAlert message={error ?? '加载失败'} />;
  }

  const metaItems = [
    { label: '状态', value: data.review_status, icon: '⚙️' },
    { label: 'AI 模型', value: data.ai_model ?? '-', icon: '🤖' },
    { label: '创建时间', value: formatDate(data.created_at), icon: '🕐' },
    { label: '记录 ID', value: `#${data.id}`, icon: '🏷️' },
  ];

  return (
    <div className="space-y-6">
      <DecorImage
        src={DECOR_IMAGES.detailBanner}
        variant="banner"
        opacity={FG_IMAGE_OPACITY.banner}
        objectPosition="center center"
        className="h-32 sm:h-40"
      />

      <Card>
        <CardHeader>
          <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
            <div>
              <div className="mb-2 flex items-center gap-2">
                <span className="rounded-md bg-dream-500/20 px-2 py-0.5 text-xs text-dream-200">
                  PR #{data.pr_number}
                </span>
                {data.repo_name && (
                  <span className="text-xs text-white/40">{data.repo_name}</span>
                )}
              </div>
              <h1 className="text-xl font-bold gradient-text">{data.pr_title}</h1>
            </div>
            <a
              href={data.pr_url}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1.5 rounded-xl border border-white/15 bg-white/5 px-4 py-2 text-sm text-white/70 transition-all hover:bg-white/10 hover:text-white"
            >
              <span>🔗</span> 在 GitHub 查看
            </a>
          </div>
        </CardHeader>
        <CardBody>
          <dl className="grid grid-cols-2 gap-4 sm:grid-cols-4">
            {metaItems.map((item) => (
              <div
                key={item.label}
                className="rounded-xl border border-white/5 bg-white/[0.03] p-3"
              >
                <dt className="flex items-center gap-1 text-[11px] text-white/40">
                  <span>{item.icon}</span> {item.label}
                </dt>
                <dd className="mt-1.5 text-sm font-medium text-white/80">
                  {item.value}
                </dd>
              </div>
            ))}
          </dl>
        </CardBody>
      </Card>

      <ReviewSummary
        totalIssues={data.total_issues}
        highIssues={data.high_issues}
        mediumIssues={data.medium_issues}
        lowIssues={data.low_issues}
      />

      <PrChangeSummary summary={data.review_result.pr_change_summary} />

      <section>
        <h2 className="mb-5 flex items-center gap-2 text-lg font-semibold">
          <span>🔎</span>
          <span className="gradient-text">审查问题</span>
        </h2>
        <IssueList reviewResult={data.review_result} />
      </section>

      <ReviewActions
        reviewId={data.id}
        prChangeSummary={data.review_result.pr_change_summary}
      />

      <DiffViewer diff={data.raw_diff} reviewResult={data.review_result} />
    </div>
  );
}
