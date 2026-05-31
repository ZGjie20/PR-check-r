import { apiRequest } from '@/api/client';
import type {
  CreateReviewResponse,
  PRActionResult,
  RejectCommentDraftResponse,
  ReviewListResponse,
  ReviewRecord,
} from '@/types/review';

export interface ListReviewsParams {
  page?: number;
  limit?: number;
  prNumber?: number;
}

export async function createReview(prUrl: string): Promise<CreateReviewResponse> {
  return apiRequest<CreateReviewResponse>('/api/v1/reviews', {
    method: 'POST',
    body: JSON.stringify({ pr_url: prUrl }),
    signal: AbortSignal.timeout(180_000),
  });
}

export async function listReviews(
  params: ListReviewsParams = {},
): Promise<ReviewListResponse> {
  const search = new URLSearchParams();
  search.set('page', String(params.page ?? 1));
  search.set('limit', String(params.limit ?? 20));
  if (params.prNumber && params.prNumber > 0) {
    search.set('pr_number', String(params.prNumber));
  }

  return apiRequest<ReviewListResponse>(`/api/v1/reviews?${search}`);
}

export async function getReview(id: number): Promise<ReviewRecord> {
  return apiRequest<ReviewRecord>(`/api/v1/reviews/${id}`);
}

export async function approveReview(
  id: number,
  comment?: string,
): Promise<PRActionResult> {
  return apiRequest<PRActionResult>(`/api/v1/reviews/${id}/approve`, {
    method: 'POST',
    body: comment ? JSON.stringify({ comment }) : undefined,
  });
}

export async function mergeReview(id: number): Promise<PRActionResult> {
  return apiRequest<PRActionResult>(`/api/v1/reviews/${id}/merge`, {
    method: 'POST',
  });
}

export async function getRejectCommentDraft(
  id: number,
): Promise<RejectCommentDraftResponse> {
  return apiRequest<RejectCommentDraftResponse>(
    `/api/v1/reviews/${id}/reject-comment-draft`,
  );
}

export async function rejectReview(
  id: number,
  comment: string,
): Promise<PRActionResult> {
  return apiRequest<PRActionResult>(`/api/v1/reviews/${id}/reject`, {
    method: 'POST',
    body: JSON.stringify({ comment }),
  });
}
