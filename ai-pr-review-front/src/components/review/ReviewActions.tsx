import { useReviewActions } from '@/hooks/useReview';
import { RejectModal } from '@/components/review/RejectModal';
import { Card, CardBody, CardHeader } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { ErrorAlert } from '@/components/common/ErrorAlert';

interface ReviewActionsProps {
  reviewId: number;
  prChangeSummary?: string;
}

export function ReviewActions({ reviewId, prChangeSummary }: ReviewActionsProps) {
  const {
    actionState,
    loading,
    error,
    successMessage,
    rejectOpen,
    approve,
    merge,
    reject,
    openReject,
    closeReject,
  } = useReviewActions(reviewId);

  const isFinalState = actionState === 'merged' || actionState === 'rejected';
  const canApprove = actionState === 'idle' && !loading;
  const canReject = actionState === 'idle' && !loading;
  const showMerge = actionState === 'approved';

  return (
    <>
      <Card>
        <CardHeader>
          <h2 className="flex items-center gap-2 text-lg font-semibold">
            <span>⚡</span>
            <span className="gradient-text">GitHub 操作</span>
          </h2>
          <p className="mt-1 text-sm text-white/40">
            Approve 后可选择 Merge；Reject 将提交 REQUEST_CHANGES 并发布评论。
          </p>
        </CardHeader>
        <CardBody className="space-y-4">
          <div className="flex flex-wrap gap-3">
            <Button
              variant="primary"
              loading={loading && actionState === 'idle'}
              disabled={!canApprove}
              onClick={() => void approve(prChangeSummary)}
            >
              Approve
            </Button>

            <Button
              variant="secondary"
              disabled={!canReject}
              onClick={openReject}
            >
              Reject
            </Button>

            {showMerge && (
              <Button
                variant="primary"
                className="from-emerald-500 to-teal-500 shadow-emerald-500/25 hover:shadow-emerald-500/40"
                loading={loading && actionState === 'approved'}
                disabled={loading || isFinalState}
                onClick={() => void merge()}
              >
                Merge
              </Button>
            )}
          </div>

          {successMessage && (
            <div className="rounded-xl border border-emerald-400/20 bg-emerald-500/10 px-4 py-3 text-sm text-emerald-200">
              {successMessage}
            </div>
          )}

          {error && <ErrorAlert message={error} />}

          {isFinalState && (
            <p className="text-xs text-white/40">
              操作已完成。刷新页面后状态会重置。
            </p>
          )}
        </CardBody>
      </Card>

      <RejectModal
        open={rejectOpen}
        reviewId={reviewId}
        loading={loading}
        onClose={closeReject}
        onReject={reject}
      />
    </>
  );
}
