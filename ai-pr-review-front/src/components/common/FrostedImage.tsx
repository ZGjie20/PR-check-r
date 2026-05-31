import type { CSSProperties, ReactNode } from 'react';

type FrostLevel = 'light' | 'medium' | 'heavy';

interface FrostedImageProps {
  src: string;
  alt?: string;
  /** 背景层透明度，宜低于前景图 */
  opacity?: number;
  frost?: FrostLevel;
  className?: string;
  imageClassName?: string;
  objectPosition?: string;
  children?: ReactNode;
}

const FROST_CLASSES: Record<FrostLevel, string> = {
  light: 'backdrop-blur-[2px] bg-white/[0.06]',
  medium: 'backdrop-blur-md bg-white/10',
  heavy: 'backdrop-blur-xl bg-white/[0.14]',
};

/** 仅用于最底层背景：磨砂 + 噪点 + 渐变遮罩 */
export function FrostedImage({
  src,
  alt = '',
  opacity = 0.28,
  frost = 'heavy',
  className = '',
  imageClassName = '',
  objectPosition = 'center',
  children,
}: FrostedImageProps) {
  const imageStyle: CSSProperties = { opacity, objectPosition };

  return (
    <div
      className={`pointer-events-none absolute inset-0 overflow-hidden ${className}`}
      aria-hidden={!alt}
    >
      <img
        src={src}
        alt={alt}
        loading="lazy"
        decoding="async"
        className={`h-full w-full object-cover ${imageClassName}`}
        style={imageStyle}
      />
      <div className={`absolute inset-0 ${FROST_CLASSES[frost]}`} />
      <div
        className="pointer-events-none absolute inset-0 opacity-[0.35] mix-blend-overlay"
        style={{
          backgroundImage:
            'url("data:image/svg+xml,%3Csvg viewBox=\'0 0 256 256\' xmlns=\'http://www.w3.org/2000/svg\'%3E%3Cfilter id=\'n\'%3E%3CfeTurbulence type=\'fractalNoise\' baseFrequency=\'0.85\' numOctaves=\'4\' stitchTiles=\'stitch\'/%3E%3C/filter%3E%3Crect width=\'100%25\' height=\'100%25\' filter=\'url(%23n)\' opacity=\'0.5\'/%3E%3C/svg%3E")',
        }}
      />
      <div className="absolute inset-0 bg-gradient-to-br from-[#0f0c29]/55 via-[#302b63]/35 to-[#0f0c29]/65" />
      {children}
    </div>
  );
}

/** 背景层默认透明度（更通透） */
export const BG_IMAGE_OPACITY = {
  main: 0.32,
} as const;
