import { DecorImage, FG_IMAGE_OPACITY } from '@/components/common/DecorImage';
import { DECOR_IMAGES } from '@/constants/decorImages';
import { Spinner } from '@/components/ui/Spinner';

interface LoadingOverlayProps {
  message?: string;
}

export function LoadingOverlay({
  message = '审查进行中，请勿关闭页面...',
}: LoadingOverlayProps) {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-md">
      <div className="glass-strong relative w-full max-w-md overflow-hidden rounded-2xl">
        <DecorImage
          src={DECOR_IMAGES.warrior}
          variant="banner"
          opacity={FG_IMAGE_OPACITY.loading}
          objectPosition="center"
          className="absolute inset-0 h-full"
        />
        <div className="relative z-10 flex flex-col items-center gap-5 px-10 py-8">
          <Spinner size="lg" />
          <div className="text-center">
            <p className="text-base font-medium gradient-text">{message}</p>
            <p className="mt-2 text-xs text-white/40">AI 正在逐行分析代码变更...</p>
          </div>
          <div className="flex gap-1">
            {[0, 1, 2].map((i) => (
              <span
                key={i}
                className="h-1.5 w-1.5 animate-pulse rounded-full bg-sakura-400"
                style={{ animationDelay: `${i * 300}ms` }}
              />
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
