import { apiRequest } from '@/api/client';
import type { HealthResponse } from '@/types/api';

export async function checkHealth(): Promise<boolean> {
  try {
    const res = await apiRequest<HealthResponse>('/health');
    return res.status === 'ok';
  } catch {
    return false;
  }
}
