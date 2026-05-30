/** 单一审查问题 */
export interface ReviewIssueDetail {
  file: string;
  line: number;
  message: string;
  suggestion: string;
}

/** 依严重程度分组的审查结果 */
export interface ReviewResultBySeverity {
  high: ReviewIssueDetail[];
  medium: ReviewIssueDetail[];
  low: ReviewIssueDetail[];
}

/** POST /api/v1/reviews 请求体 */
export interface CreateReviewRequest {
  pr_url: string;
}

/** POST /api/v1/reviews 响应 */
export interface CreateReviewResponse {
  id: number;
  pr_title: string;
  pr_number: number;
  repo_name?: string;
  pr_url: string;
  ai_model?: string;
  review_status: string;
  total_issues: number;
  high_issues: number;
  medium_issues: number;
  low_issues: number;
  review_result: ReviewResultBySeverity;
  raw_diff: string;
  created_at: string;
  output_file: string;
}

/** GET /api/v1/reviews/:id 响应 */
export interface ReviewRecord {
  id: number;
  pr_title: string;
  pr_number: number;
  repo_name?: string;
  pr_url: string;
  ai_model?: string;
  review_status: string;
  total_issues: number;
  high_issues: number;
  medium_issues: number;
  low_issues: number;
  review_result: ReviewResultBySeverity;
  raw_diff: string;
  created_at: string;
}

/** 列表单笔摘要 */
export interface ReviewListItem {
  id: number;
  pr_title: string;
  pr_number: number;
  total_issues: number;
  created_at: string;
}

/** GET /api/v1/reviews 响应 */
export interface ReviewListResponse {
  items: ReviewListItem[];
  total: number;
  page: number;
  limit: number;
}

export type Severity = 'high' | 'medium' | 'low';
