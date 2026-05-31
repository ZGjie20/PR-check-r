/** 装饰图片资源（由 png/ 原图压缩生成，见 scripts/optimize-images.py） */
export const DECOR_IMAGES = {
  washitsu: '/images/decor/washitsu.webp',
  historyBg: '/images/decor/history-bg.webp',
  detailBg: '/images/decor/detail-bg.webp',
  detailBanner: '/images/decor/detail-banner.webp',
  warrior: '/images/decor/warrior.webp',
  pool: '/images/decor/pool.webp',
  campfire: '/images/decor/campfire.webp',
  catgirl: '/images/decor/catgirl.webp',
} as const;

export type DecorImageKey = keyof typeof DECOR_IMAGES;
