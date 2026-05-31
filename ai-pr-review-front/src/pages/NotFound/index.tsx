import { Link } from 'react-router-dom';
import { DecorImage, FG_IMAGE_OPACITY } from '@/components/common/DecorImage';
import { DECOR_IMAGES } from '@/constants/decorImages';
import { Button } from '@/components/ui/Button';

export function NotFoundPage() {
  return (
    <div className="flex flex-col items-center gap-6 py-16 text-center">
      <div className="relative">
        <DecorImage
          src={DECOR_IMAGES.catgirl}
          variant="card"
          opacity={FG_IMAGE_OPACITY.card}
          objectPosition="center 10%"
          className="mx-auto h-64 w-48 sm:h-72 sm:w-56"
        />
        <span className="absolute -right-2 -top-4 text-5xl font-black gradient-text opacity-90">
          404
        </span>
      </div>
      <div>
        <p className="text-lg font-medium text-white/70">页面迷失在异次元了...</p>
        <p className="mt-1 text-sm text-white/40">你所寻找的页面不存在</p>
      </div>
      <Link to="/">
        <Button variant="secondary">✨ 返回首页</Button>
      </Link>
    </div>
  );
}
