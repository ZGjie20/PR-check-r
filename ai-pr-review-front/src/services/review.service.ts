import * as reviewApi from '@/api/review';
import type { ListReviewsParams } from '@/api/review';

export const reviewService = {
  createReview: reviewApi.createReview,
  listReviews: (params: ListReviewsParams = {}) => reviewApi.listReviews(params),
  getReview: reviewApi.getReview,
};
