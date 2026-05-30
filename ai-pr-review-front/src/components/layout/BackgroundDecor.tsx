import { useLocation } from 'react-router-dom';
import { FrostedImage, BG_IMAGE_OPACITY } from '@/components/common/FrostedImage';
import { DECOR_IMAGES } from '@/constants/decorImages';

interface PageBackground {
  src: string;
  objectPosition: string;
}

const DEFAULT_BACKGROUND: PageBackground = {
  src: DECOR_IMAGES.washitsu,
  objectPosition: 'center top',
};

function resolveBackground(pathname: string): PageBackground {
  if (/^\/reviews\/\d+$/.test(pathname)) {
    return {
      src: DECOR_IMAGES.detailBg,
      objectPosition: 'center center',
    };
  }

  if (pathname === '/reviews') {
    return {
      src: DECOR_IMAGES.historyBg,
      objectPosition: 'center center',
    };
  }

  return DEFAULT_BACKGROUND;
}

export function BackgroundDecor() {
  const { pathname } = useLocation();
  const bg = resolveBackground(pathname);

  return (
    <div className="pointer-events-none fixed inset-0 overflow-hidden" aria-hidden="true">
      <FrostedImage
        key={bg.src}
        src={bg.src}
        opacity={BG_IMAGE_OPACITY.main}
        frost="heavy"
        objectPosition={bg.objectPosition}
        className="scale-105"
      />

      <div className="absolute -left-32 -top-32 h-96 w-96 rounded-full bg-dream-500/15 blur-[100px]" />
      <div className="absolute -bottom-32 -right-32 h-96 w-96 rounded-full bg-sakura-500/15 blur-[100px]" />
      <div className="absolute left-1/2 top-1/3 h-64 w-64 -translate-x-1/2 rounded-full bg-cyan-400/8 blur-[80px]" />

      {PARTICLES.map((p) => (
        <span
          key={p.id}
          className={`absolute rounded-full ${p.size} ${p.color} ${p.animation}`}
          style={{ left: p.left, top: p.top }}
        />
      ))}
    </div>
  );
}

const PARTICLES = [
  { id: 1, left: '10%', top: '20%', size: 'h-1 w-1', color: 'bg-sakura-300/60', animation: 'animate-twinkle' },
  { id: 2, left: '25%', top: '60%', size: 'h-1.5 w-1.5', color: 'bg-dream-300/50', animation: 'animate-float' },
  { id: 3, left: '45%', top: '15%', size: 'h-1 w-1', color: 'bg-cyan-300/50', animation: 'animate-float-delayed' },
  { id: 4, left: '70%', top: '40%', size: 'h-2 w-2', color: 'bg-sakura-200/40', animation: 'animate-twinkle' },
  { id: 5, left: '85%', top: '70%', size: 'h-1 w-1', color: 'bg-dream-200/50', animation: 'animate-float-slow' },
  { id: 6, left: '55%', top: '80%', size: 'h-1.5 w-1.5', color: 'bg-white/30', animation: 'animate-float' },
];
