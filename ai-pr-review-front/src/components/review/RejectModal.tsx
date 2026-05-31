import { useCallback, useEffect, useState } from 'react';
import { reviewService } from '@/services/review.service';
import { Modal } from '@/components/ui/Modal';
import { Button } from '@/components/ui/Button';
import { Spinner } from '@/components/ui/Spinner';
import { ErrorAlert } from '@/components/common/ErrorAlert';
import { getFriendlyErrorMessage } from '@/utils/error';

interface RejectModalProps {
  open: boolean;
  reviewId: number;
  loading: boolean;
  onClose: () => void;
  onReject: (comment: string) => Promise<boolean>;
}

export function RejectModal({
  open,
  reviewId,
  loading,
  onClose,
  onReject,
}: RejectModalProps) {
  const [comment, setComment] = useState('');
  const [draftLoading, setDraftLoading] = useState(false);
  const [draftError, setDraftError] = useState<string | null>(null);

  const loadDraft = useCallback(async () => {
    setDraftLoading(true);
    setDraftError(null);

    try {
      const result = await reviewService.getRejectCommentDraft(reviewId);
      setComment(result.comment);
    } catch (err) {
      setDraftError(getFriendlyErrorMessage(err));
    } finally {
      setDraftLoading(false);
    }
  }, [reviewId]);

  useEffect(() => {
    if (open) {
      void loadDraft();
    } else {
      setComment('');
      setDraftError(null);
    }
  }, [open, loadDraft]);

  const handleSubmit = async () => {
    const trimmed = comment.trim();
    if (!trimmed) {
      return;
    }
    await onReject(trimmed);
  };

  const canSubmit = comment.trim().length > 0 && !loading && !draftLoading;

  return (
    <Modal open={open} onClose={onClose} title="打回 PR">
      {draftLoading ? (
        <div className="flex flex-col items-center gap-3 py-12">
          <Spinner />
          <p className="text-sm text-white/40">正在加载评论草稿...</p>
        </div>
      ) : draftError ? (
        <div className="space-y-4">
          <ErrorAlert message={draftError} />
          <Button variant="secondary" onClick={() => void loadDraft()}>
            重试
          </Button>
        </div>
      ) : (
        <div className="space-y-4">
          <p className="text-sm text-white/50">
            以下为 AI 生成的打回评论草稿，可直接编辑后提交。
          </p>
          <textarea
            value={comment}
            onChange={(event) => setComment(event.target.value)}
            rows={14}
            className="w-full resize-y rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white/80 outline-none transition-colors focus:border-dream-400/50 focus:ring-1 focus:ring-dream-400/30"
            placeholder="请输入打回原因..."
          />
          <div className="flex justify-end gap-3">
            <Button variant="ghost" onClick={onClose} disabled={loading}>
              取消
            </Button>
            <Button
              variant="secondary"
              loading={loading}
              disabled={!canSubmit}
              onClick={() => void handleSubmit()}
            >
              提交打回
            </Button>
          </div>
        </div>
      )}
    </Modal>
  );
}
