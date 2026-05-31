import { DecorImage, FG_IMAGE_OPACITY } from '@/components/common/DecorImage';
import { DECOR_IMAGES } from '@/constants/decorImages';

interface PageDecorBannerProps {
  image: string;
  opacity?: number;
  objectPosition?: string;
  className?: string;
}

/** 页面顶部装饰横幅（前景，无磨砂） */
export function PageDecorBanner({
  image,
  opacity = FG_IMAGE_OPACITY.banner,
  objectPosition = 'center',
  className = '',
}: PageDecorBannerProps) {
  return (
    <DecorImage
      src={image}
      variant="banner"
      opacity={opacity}
      objectPosition={objectPosition}
      className={`mb-6 h-36 sm:h-44 ${className}`}
    />
  );
}

/** 侧边装饰插画（前景，无磨砂） */
export function PageDecorAside({
  image,
  opacity = FG_IMAGE_OPACITY.aside,
  objectPosition = 'center top',
}: {
  image: string;
  opacity?: number;
  objectPosition?: string;
}) {
  return (
    <DecorImage
      src={image}
      variant="card"
      opacity={opacity}
      objectPosition={objectPosition}
      className="h-full min-h-[280px] w-full"
    />
  );
}

export { DECOR_IMAGES };
