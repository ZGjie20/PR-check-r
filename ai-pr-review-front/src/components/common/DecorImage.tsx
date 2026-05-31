import type { CSSProperties } from 'react';

type DecorVariant = 'card' | 'banner' | 'corner';

interface DecorImageProps {
  src: string;
  alt?: string;
  /** 前景层透明度，宜高于背景 */
  opacity?: number;
  variant?: DecorVariant;
  className?: string;
  imageClassName?: string;
  objectPosition?: string;
}

const VARIANT_CLASSES: Record<DecorVariant, string> = {
  card: 'relative overflow-hidden rounded-2xl',
  banner: 'relative overflow-hidden rounded-2xl',
  corner: 'relative overflow-hidden rounded-2xl',
};

/** 前景装饰图：无磨砂，保持清晰 */
export function DecorImage({
  src,
  alt = '',
  opacity = 0.82,
  variant = 'card',
  className = '',
  imageClassName = '',
  objectPosition = 'center',
}: DecorImageProps) {
  const imageStyle: CSSProperties = { opacity, objectPosition };

  return (
    <div className={`${VARIANT_CLASSES[variant]} ${className}`} aria-hidden={!alt}>
      <img
        src={src}
        alt={alt}
        loading="lazy"
        decoding="async"
        className={`h-full w-full object-cover ${imageClassName}`}
        style={imageStyle}
      />
      {/* 仅底部轻渐变，避免遮挡主体 */}
      <div className="pointer-events-none absolute inset-x-0 bottom-0 h-1/3 bg-gradient-to-t from-[#0f0c29]/50 to-transparent" />
    </div>
  );
}

/** 前景层默认透明度（更清晰） */
export const FG_IMAGE_OPACITY = {
  banner: 0.78,
  aside: 0.85,
  card: 0.88,
  loading: 0.75,
} as const;
