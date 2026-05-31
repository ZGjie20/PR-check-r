import { useEffect } from 'react';
import { checkHealth } from '@/api/health';
import { useAppStore } from '@/store/app.store';

export function useHealth() {
  const apiHealthy = useAppStore((s) => s.apiHealthy);
  const setApiHealthy = useAppStore((s) => s.setApiHealthy);

  useEffect(() => {
    let cancelled = false;

    async function run() {
      const healthy = await checkHealth();
      if (!cancelled) {
        setApiHealthy(healthy);
      }
    }

    void run();

    const interval = setInterval(() => {
      void run();
    }, 30_000);

    return () => {
      cancelled = true;
      clearInterval(interval);
    };
  }, [setApiHealthy]);

  return { apiHealthy };
}
