import { Outlet } from 'react-router-dom';
import { Header } from '@/components/layout/Header';
import { BackgroundDecor } from '@/components/layout/BackgroundDecor';
import { useAppStore } from '@/store/app.store';
import { ErrorAlert } from '@/components/common/ErrorAlert';

export function MainLayout() {
  const apiHealthy = useAppStore((s) => s.apiHealthy);

  return (
    <div className="relative min-h-screen">
      <BackgroundDecor />
      <div className="relative z-10 flex min-h-screen flex-col">
        <Header />
        {apiHealthy === false && (
          <div className="mx-auto w-full max-w-6xl px-4 pt-4">
            <ErrorAlert message="后端服务不可用，请确认 API 已启动后再试" />
          </div>
        )}
        <main className="mx-auto w-full max-w-6xl flex-1 px-4 py-8">
          <Outlet />
        </main>
        <footer className="relative z-10 border-t border-white/10 py-6 text-center text-xs text-white/30">
          <p>AI PR Review · 代码审查灵境</p>
          <p className="mt-1">writed by jie</p>
        </footer>
      </div>
    </div>
  );
}
