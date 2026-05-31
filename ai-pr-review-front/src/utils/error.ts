import type { ErrorResponse } from '@/types/api';

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

export function getFriendlyErrorMessage(error: unknown): string {
  if (error instanceof ApiError) {
    const msg = error.message;

    if (msg.includes('invalid request body')) {
      return '请求格式错误，请重试';
    }
    if (msg.includes('pr_url is required')) {
      return '请输入 PR 链接';
    }
    if (msg.includes('invalid GitHub PR URL')) {
      return '请输入有效的 GitHub PR 链接';
    }
    if (msg.includes('fetch PR')) {
      return '无法获取 PR 信息，请确认链接正确且 PR 可访问';
    }
    if (msg.includes('run review')) {
      return 'AI 审查失败，请稍后重试';
    }
    if (error.status === 404) {
      return '记录不存在';
    }
    if (error.status === 500) {
      return '服务暂时不可用，请稍后重试';
    }

    return msg;
  }

  if (error instanceof DOMException && error.name === 'TimeoutError') {
    return '审查耗时较长或网络中断，请稍后在历史记录中查看';
  }

  if (error instanceof Error) {
    return error.message;
  }

  return '未知错误';
}

export function isApiErrorBody(body: unknown): body is ErrorResponse {
  return (
    typeof body === 'object' &&
    body !== null &&
    'error' in body &&
    typeof (body as ErrorResponse).error === 'string'
  );
}
