import { useCallback, useEffect, useState } from 'react';
import { reviewService } from '@/services/review.service';
import { useAppStore } from '@/store/app.store';
import type {
  CreateReviewResponse,
  ReviewListResponse,
  ReviewRecord,
} from '@/types/review';
import { ApiError, getFriendlyErrorMessage } from '@/utils/error';

export function useCreateReview() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const setIsSubmitting = useAppStore((s) => s.setIsSubmitting);

  const create = useCallback(
    async (prUrl: string): Promise<CreateReviewResponse | null> => {
      setLoading(true);
      setError(null);
      setIsSubmitting(true);

      try {
        const result = await reviewService.createReview(prUrl);
        return result;
      } catch (err) {
        setError(getFriendlyErrorMessage(err));
        return null;
      } finally {
        setLoading(false);
        setIsSubmitting(false);
      }
    },
    [setIsSubmitting],
  );

  return { create, loading, error, setError };
}

export function useReviewList(page: number, limit: number, prNumber?: number) {
  const [data, setData] = useState<ReviewListResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    async function fetchList() {
      setLoading(true);
      setError(null);

      try {
        const result = await reviewService.listReviews({ page, limit, prNumber });
        if (!cancelled) {
          setData(result);
        }
      } catch (err) {
        if (!cancelled) {
          setError(getFriendlyErrorMessage(err));
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    void fetchList();

    return () => {
      cancelled = true;
    };
  }, [page, limit, prNumber]);

  return { data, loading, error };
}

export function useReviewDetail(id: number | null) {
  const [data, setData] = useState<ReviewRecord | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [notFound, setNotFound] = useState(false);

  useEffect(() => {
    if (id === null || Number.isNaN(id) || id <= 0) {
      setLoading(false);
      setNotFound(true);
      return;
    }

    const reviewId = id;
    let cancelled = false;

    async function fetchDetail() {
      setLoading(true);
      setError(null);
      setNotFound(false);

      try {
        const result = await reviewService.getReview(reviewId);
        if (!cancelled) {
          setData(result);
        }
      } catch (err) {
        if (!cancelled) {
          const message = getFriendlyErrorMessage(err);
          setError(message);
          if (err instanceof ApiError && err.status === 404) {
            setNotFound(true);
          }
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    void fetchDetail();

    return () => {
      cancelled = true;
    };
  }, [id]);

  return { data, loading, error, notFound };
}
