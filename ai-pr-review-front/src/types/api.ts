/** 统一错误响应 */
export interface ErrorResponse {
  error: string;
}

/** GET /health 响应 */
export interface HealthResponse {
  status: string;
}
