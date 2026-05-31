import { useNavigate } from 'react-router-dom';
import { PrUrlForm } from '@/components/review/PrUrlForm';
import { LoadingOverlay } from '@/components/common/LoadingOverlay';
import { PageHeader } from '@/components/layout/PageHeader';
import { PageDecorAside, DECOR_IMAGES } from '@/components/layout/PageDecor';
import { Card, CardBody } from '@/components/ui/Card';
import { useCreateReview } from '@/hooks/useReview';
import { useAppStore } from '@/store/app.store';

const FEATURES = [
  { icon: '🔍', title: '智能分析', desc: 'AI 逐块审查代码变更' },
  { icon: '⚡', title: '实时反馈', desc: '同步返回完整审查结果' },
  { icon: '📋', title: '分级报告', desc: '按高中低风险分类展示' },
  { icon: '📝', title: '变更总结', desc: '自动生成 PR 变更说明' },
];

export function ReviewCreatePage() {
  const navigate = useNavigate();
  const { create, loading, error, setError } = useCreateReview();
  const apiHealthy = useAppStore((s) => s.apiHealthy);

  async function handleSubmit(prUrl: string) {
    setError(null);
    const result = await create(prUrl);
    if (result) {
      navigate(`/reviews/${result.id}`);
    }
  }

  return (
    <>
      <PageHeader
        title="新建 PR 审查"
        subtitle="输入 GitHub Pull Request 链接，开启 AI 代码审查之旅"
        icon="🔮"
      />

      <div className="grid gap-6 lg:grid-cols-5">
        <Card className="lg:col-span-3">
          <CardBody>
            <PrUrlForm
              onSubmit={handleSubmit}
              loading={loading}
              disabled={apiHealthy === false}
              serverError={error}
            />
          </CardBody>
        </Card>

        <div className="flex flex-col gap-4 lg:col-span-2">
          <PageDecorAside
            image={DECOR_IMAGES.pool}
            objectPosition="center 20%"
          />
          {FEATURES.map((f) => (
            <div key={f.title} className="glass-card flex items-start gap-3 p-4">
              <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-white/5 text-base">
                {f.icon}
              </span>
              <div>
                <p className="text-sm font-medium text-white/80">{f.title}</p>
                <p className="mt-0.5 text-xs text-white/40">{f.desc}</p>
              </div>
            </div>
          ))}
        </div>
      </div>

      {loading && <LoadingOverlay />}
    </>
  );
}
