interface SpinnerProps {
  size?: 'sm' | 'md' | 'lg';
}

const sizeClasses = {
  sm: 'h-4 w-4 border-2',
  md: 'h-6 w-6 border-2',
  lg: 'h-12 w-12 border-[3px]',
};

export function Spinner({ size = 'md' }: SpinnerProps) {
  return (
    <div className="relative flex items-center justify-center" role="status" aria-label="加载中">
      <div
        className={`animate-spin rounded-full border-sakura-400/30 border-t-sakura-400 ${sizeClasses[size]}`}
      />
      {size === 'lg' && (
        <div className="absolute h-6 w-6 animate-spin-slow rounded-full border border-dream-400/20 border-t-dream-400/60" />
      )}
    </div>
  );
}
